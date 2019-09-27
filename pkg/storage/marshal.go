package storage

import (
	"encoding/json"
)

// Marshal marshals data with indentation
func Marshal(input interface{}) ([]byte, error) {
	return json.MarshalIndent(input, "", "  ")
}
