package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestBlacklistNamespaces(t *testing.T) {
	// expected test result
	expectedResult := []string{
		"kube-system",
		"kube-public",
		"docker",
		"daemon-set",
	}

	// create a dummy blacklist
	dummyBlacklist := []string{
		"kube*",
		"docker",
		"*-set",
	}

	// create dummy namespaces
	dummyNamespaces := []runtime.Object{
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "docker"}},
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kube-system"}},
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "daemon-set"}},
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kube-public"}},
	}

	k := API{
		Client: fake.NewSimpleClientset(dummyNamespaces...),
	}

	blacklistedNamespaces, err := k.Namespaces(dummyBlacklist)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expectedResult, blacklistedNamespaces)
}
