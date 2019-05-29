package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlacklistFieldSelector(t *testing.T) {
	blacklist := []string{"kube-system", "private"}

	fieldSelector := BlacklistFieldSelector(blacklist)

	assert.Equal(t, "metadata.namespace!=kube-system,metadata.namespace!=private", fieldSelector)
}

func TestPrefix(t *testing.T) {
	list := []string{"foo", "bar"}
	prefix := "new-"

	r := Prefix(list, prefix)

	assert.Equal(t, []string{"new-foo", "new-bar"}, r)
}
