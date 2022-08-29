package env

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/macro/enigma"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	dctv1 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v1"
	dctv2 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
)

func init() {
	println("set enigma")
	enigmaConfigFilename := "../env/enigma.yml"
	if err := newEnigmav1(enigmaConfigFilename); err != nil {
		panic(err)
	}
	if err := newEnigmav2(enigmaConfigFilename); err != nil {
		panic(err)
	}
}

func newEnigmav1(configFilename string) error {
	config := enigma.Config{}
	if err := configor.Load(&config, configFilename); err != nil {
		return errors.Wrapf(err, "read enigma config file %v",
			logs.KVL(
				"filename", configFilename,
			))
	}

	if err := enigma.LoadConfig(config); err != nil {
		b, _ := ioutil.ReadFile(configFilename)

		return errors.Wrapf(err, "load enigma config %v",
			logs.KVL(
				"filename", configFilename,
				"config", string(b),
			))
	}

	if len(config.CryptoAlgorithmSet) == 0 {
		return errors.New("'enigma cripto alg set' is empty")
	}

	for _, k := range dctv1.CiperKeyNames() {
		if _, ok := config.CryptoAlgorithmSet[k]; !ok {
			return errors.Errorf("not found enigma machine name%s",
				logs.KVL(
					"key", k,
				))
		}
	}

	enigma.PrintConfig(os.Stdout, config)

	for key := range config.CryptoAlgorithmSet {
		const quickbrownfox = `the quick brown fox jumps over the lazy dog`
		encripted, err := enigma.CipherSet(key).Encode([]byte(quickbrownfox))
		if err != nil {
			return errors.Wrapf(err, "enigma test: encode %v",
				logs.KVL("config-name", key))
		}
		plain, err := enigma.CipherSet(key).Decode(encripted)
		if err != nil {
			return errors.Wrapf(err, "enigma test: decode %v",
				logs.KVL("config-name", key))
		}

		if strings.Compare(quickbrownfox, string(plain)) != 0 {
			return errors.Errorf("enigma test: diff result %v",
				logs.KVL("config-name", key))
		}
	}

	return nil
}

func newEnigmav2(configFilename string) error {
	config := enigma.Config{}
	if err := configor.Load(&config, configFilename); err != nil {
		return errors.Wrapf(err, "read enigma config file %v",
			logs.KVL(
				"filename", configFilename,
			))
	}

	if err := enigma.LoadConfig(config); err != nil {
		b, _ := ioutil.ReadFile(configFilename)

		return errors.Wrapf(err, "load enigma config %v",
			logs.KVL(
				"filename", configFilename,
				"config", string(b),
			))
	}

	if len(config.CryptoAlgorithmSet) == 0 {
		return errors.New("'enigma cripto alg set' is empty")
	}

	for _, k := range dctv2.CiperKeyNames() {
		if _, ok := config.CryptoAlgorithmSet[k]; !ok {
			return errors.Errorf("not found enigma machine name%s",
				logs.KVL(
					"key", k,
				))
		}
	}

	enigma.PrintConfig(os.Stdout, config)

	for key := range config.CryptoAlgorithmSet {
		const quickbrownfox = `the quick brown fox jumps over the lazy dog`
		encripted, err := enigma.CipherSet(key).Encode([]byte(quickbrownfox))
		if err != nil {
			return errors.Wrapf(err, "enigma test: encode %v",
				logs.KVL("config-name", key))
		}
		plain, err := enigma.CipherSet(key).Decode(encripted)
		if err != nil {
			return errors.Wrapf(err, "enigma test: decode %v",
				logs.KVL("config-name", key))
		}

		if strings.Compare(quickbrownfox, string(plain)) != 0 {
			return errors.Errorf("enigma test: diff result %v",
				logs.KVL("config-name", key))
		}
	}

	return nil
}
