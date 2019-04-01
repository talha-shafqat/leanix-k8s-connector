package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelations(t *testing.T) {
	source := "daf29ef3-8a55-415d-aa11-848a7a9e815b"
	targets := []interface{}{
		"95df59d1-b0a3-4523-aaa5-e8db58139d7d",
		"b4aa4eb2-79e1-4a5b-b1ec-e248fdc53621",
	}

	r := Relations(source, targets)

	assert.Contains(t, r, source)
	assert.Len(t, r[source], 2)
}

func TestNewRelation(t *testing.T) {
	target := "102b9da4-549d-11e9-8647-d663bd873d93"

	r := NewRelation(target)

	assert.Equal(t, r["uid"], target)
	assert.Equal(t, r["relName"], "relToRequiredBy")
}
