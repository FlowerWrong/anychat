package utils

import (
	"encoding/json"
)

// RawMsg ...
func RawMsg(data interface{}) ([]byte, error) {
	b, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	raw := json.RawMessage(b)
	return raw, nil
}
