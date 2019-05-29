package mapper

import (
	"fmt"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type testInput struct {
	id                 uuid.UUID
	replicas           int32
	nodes              int
	nodesZoneRedundant bool
}

type testExpected struct {
	isRedundant                        bool
	isRedundantAcrossNodes             bool
	isRedundantAcrossAvailabilityZones bool
}

func newNodes(numberOfNodes int, nodesZoneRedundant bool) *[]corev1.Node {
	nodes := make([]corev1.Node, numberOfNodes)
	for i := range nodes {
		name := fmt.Sprintf("kubelet-%d", i)
		zone := "0"
		if nodesZoneRedundant {
			zone = string(i)
		}
		nodes[i] = corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"name":                                   name,
					"failure-domain.beta.kubernetes.io/zone": zone,
				},
			}}
	}
	return &nodes
}
