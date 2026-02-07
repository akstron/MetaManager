package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadJSON reads the file at path and unmarshals its contents into v.
// v must be a pointer to a value that json.Unmarshal can decode into.
// Returns an error if the file cannot be read or JSON is invalid (use os.IsNotExist to check for missing file).
func ReadJSON(path string, v interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("decode JSON from %s: %w", path, err)
	}
	return nil
}

// WriteJSON marshals v to JSON and writes it to path.
// If indent is true, output is pretty-printed with two-space indentation.
// The file is created with mode 0600.
func WriteJSON(path string, v interface{}, indent bool) error {
	var b []byte
	var err error
	if indent {
		b, err = json.MarshalIndent(v, "", "  ")
	} else {
		b, err = json.Marshal(v)
	}
	if err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}
	return os.WriteFile(path, b, 0600)
}
