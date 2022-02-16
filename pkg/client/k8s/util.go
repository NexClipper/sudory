package k8s

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func convertMapToLabelSelector(m map[string]string) (string, error) {
	if len(m) == 0 {
		return "", nil
	}

	labelSelector := &metav1.LabelSelector{MatchLabels: make(map[string]string)}
	labelSelector.MatchLabels = m

	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return "", err
	}

	return selector.String(), nil
}
