package main

import (
	"encoding/json"
	"io/ioutil"
)

// WriteJSONFile marshals data and writes it with identation to the target filename
func WriteJSONFile(input interface{}, filename string) error {
	b, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		return err
	}
	return nil
}
