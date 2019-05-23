package main

import (
	"encoding/json"
)

// Marshal marshals data with identation
func Marshal(input interface{}) ([]byte, error) {
	return json.MarshalIndent(input, "", "  ")
}
