package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestStatefulSets(t *testing.T) {
	// create a dummy statefulsets
	dummyStatefulSets := []runtime.Object{
		&appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "myapp"}},
		&appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Namespace: "prod", Name: "wordpress"}},
	}
	k := API{
		Client: fake.NewSimpleClientset(dummyStatefulSets...),
	}
	// Since the Field Selector functionality is done server side, the fake client
	// does not support it. So we can not test with blacklisting here.
	statefulSetList, err := k.StatefulSets([]string{})
	if err != nil {
		t.Error(err)
	}
	statefulsets := statefulSetList.Items

	assert.Len(t, statefulsets, 2)
}
