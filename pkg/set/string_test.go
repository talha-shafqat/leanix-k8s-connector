package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStringSet(t *testing.T) {
	s := NewStringSet()

	assert.NotNil(t, s.Map)
}

func TestStringSetAdd(t *testing.T) {
	s := NewStringSet()
	s.Add("foo")
	s.Add("bar")

	assert.Contains(t, s.Map, "foo")
	assert.Contains(t, s.Map, "bar")
}

func TestStringSetList(t *testing.T) {
	s := NewStringSet()
	s.Add("foo")
	s.Add("bar")

	list := s.Items()

	assert.Len(t, list, 2)
	assert.Contains(t, list, "foo")
	assert.Contains(t, list, "bar")
}

func TestStringSetContains_containsString(t *testing.T) {
	s := NewStringSet()
	s.Add("foo")

	found := s.Contains("foo")

	assert.Equal(t, true, found)
}

func TestStringSetContains_notContainsString(t *testing.T) {
	s := NewStringSet()
	s.Add("foo")

	found := s.Contains("test")

	assert.Equal(t, false, found)
}
