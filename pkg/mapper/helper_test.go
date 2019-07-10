package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplacerKubernetesLabel(t *testing.T) {
	input := "app.kubernetes.io/managed-by"
	expectedOutput := "app_kubernetes_io_managed_by"

	output := replacer.Replace(input)

	assert.Equal(t, expectedOutput, output)
}
func TestReplacer(t *testing.T) {
	input := "/.+-*\\"
	expectedOutput := "______"

	output := replacer.Replace(input)

	assert.Equal(t, expectedOutput, output)
}
