package sudoryclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"
	"sigs.k8s.io/yaml"
)

var (
	SudoryclientNamespace  = "default"
	SudoryclientSecretName = "sudoryclient-credential"
)

func init() {
	// get namespace
	namespace, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err == nil {
		SudoryclientNamespace = string(namespace)
	}
}

func (c *Client) Credential(verb string, params map[string]interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	found := true
	secret, err := c.k8sClient.GetK8sClientset().CoreV1().Secrets(SudoryclientNamespace).Get(ctx, SudoryclientSecretName, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return "", err
		}
		found = false
		secret = newSecretForSudoryclient()
	}

	var result []string
	switch verb {
	case "add":
		result, err = addCredential(secret, params)
	case "get":
		result, err = getCredential(secret)
	case "update":
		result, err = updateCredential(secret, params)
	case "remove":
		result, err = removeCredential(secret, params)
	default:
		return "", fmt.Errorf("unsupported verb: %s", verb)
	}

	if err != nil {
		return "", err
	}

	resBytes, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	if !found {
		// k8s create secret
		_, err := c.k8sClient.GetK8sClientset().CoreV1().Secrets(SudoryclientNamespace).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			return "", err
		}
	} else {
		// k8s update secret
		_, err := c.k8sClient.GetK8sClientset().CoreV1().Secrets(SudoryclientNamespace).Update(ctx, secret, metav1.UpdateOptions{})
		if err != nil {
			return "", err
		}
	}

	return string(resBytes), nil
}

func addCredential(secret *corev1.Secret, params map[string]interface{}) ([]string, error) {
	if params == nil && len(params) <= 0 {
		return nil, fmt.Errorf("params is empty")
	}

	credentialsInf, ok := params["credentials"]
	if !ok {
		return nil, fmt.Errorf("credentials argument is empty")
	}

	credentials, ok := credentialsInf.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed type assertion for credentials argument")
	}

	if len(credentials) <= 0 {
		return nil, fmt.Errorf("credentials is empty")
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	var result []string
	for key, data := range credentials {
		if errs := validation.IsConfigMapKey(key); len(errs) != 0 {
			return nil, fmt.Errorf("credential %q's key is invalid: %s", key, strings.Join(errs, ";"))
		}
		if _, ok := secret.Data[key]; ok {
			return nil, fmt.Errorf("credential %q already exists", key)
		}

		cred, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed type assertion for credential %q's data: want be map[string]interface{}, not %T", key, data)
		}

		jsonb, err := json.Marshal(cred)
		if err != nil {
			return nil, err
		}

		b, err := yaml.JSONToYAML(jsonb)
		if err != nil {
			return nil, err
		}

		secret.Data[key] = b
		result = append(result, key)
	}

	return result, nil
}

func getCredential(secret *corev1.Secret) ([]string, error) {
	var result []string

	if secret != nil && secret.Data != nil {
		for k := range secret.Data {
			result = append(result, k)
		}
	}

	return result, nil
}

func updateCredential(secret *corev1.Secret, params map[string]interface{}) ([]string, error) {
	if secret == nil || secret.Data == nil || len(secret.Data) <= 0 {
		return nil, fmt.Errorf("no credentials exist")
	}

	if params == nil && len(params) <= 0 {
		return nil, fmt.Errorf("params is empty")
	}

	credentialsInf, ok := params["credentials"]
	if !ok {
		return nil, fmt.Errorf("credentials argument is empty")
	}

	credentials, ok := credentialsInf.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed type assertion for credentials argument")
	}

	if len(credentials) <= 0 {
		return nil, fmt.Errorf("credentials is empty")
	}

	var result []string
	for key, data := range credentials {
		if _, ok := secret.Data[key]; !ok {
			return nil, fmt.Errorf("credential %q does not exist", key)
		}

		cred, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed type assertion for credential %q's data: want be map[string]interface{}, not %T", key, data)
		}

		jsonb, err := json.Marshal(cred)
		if err != nil {
			return nil, err
		}

		b, err := yaml.JSONToYAML(jsonb)
		if err != nil {
			return nil, err
		}

		secret.Data[key] = b
		result = append(result, key)
	}

	return result, nil
}

func removeCredential(secret *corev1.Secret, params map[string]interface{}) ([]string, error) {
	if secret == nil || secret.Data == nil || len(secret.Data) <= 0 {
		return nil, fmt.Errorf("no credentials exist")
	}

	if params == nil && len(params) <= 0 {
		return nil, fmt.Errorf("params is empty")
	}

	keysInf, ok := params["credential_keys"]
	if !ok {
		return nil, fmt.Errorf("credential_keys argument is empty")
	}

	keysInfs, ok := keysInf.([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed type assertion for credential_keys argument: want be []interface{}, not %T", keysInf)
	}

	if len(keysInfs) <= 0 {
		return nil, fmt.Errorf("credential_keys does not have an item")
	}

	var keys []string
	for _, ki := range keysInfs {
		key, ok := ki.(string)
		if !ok {
			return nil, fmt.Errorf("failed type assertion for credential_keys argument: want be string, not %T", ki)
		}

		if key == "" {
			continue
		}

		keys = append(keys, key)
	}

	for _, key := range keys {
		if _, ok = secret.Data[key]; !ok {
			return nil, fmt.Errorf("credential %q does not exist", key)
		}
		delete(secret.Data, key)
	}

	return keys, nil
}

func newSecretForSudoryclient() *corev1.Secret {
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      SudoryclientSecretName,
			Namespace: SudoryclientNamespace,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{},
	}
	return secret
}
