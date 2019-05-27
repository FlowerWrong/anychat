package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NextFlakeID ...
func NextFlakeID() (string, error) {
	resp, err := http.Get("http://localhost:8090")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var m map[string]interface{}
	json.Unmarshal(body, &m)
	return fmt.Sprintf("%d", int64(m["id"].(float64))), nil
}
