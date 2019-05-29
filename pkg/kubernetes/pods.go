package kubernetes

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pods gets a list of pods filtered by the given list of label matcher
func (k *API) Pods(matchLabels map[string]string) (*corev1.PodList, error) {
	labelMatcher := make([]string, 0)
	for label, val := range matchLabels {
		labelMatcher = append(labelMatcher, fmt.Sprintf("%s=%s", label, val))
	}
	selector := strings.Join(labelMatcher, ",")
	pods, err := k.Client.CoreV1().Pods("").List(metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return nil, err
	}
	return pods, nil
}
