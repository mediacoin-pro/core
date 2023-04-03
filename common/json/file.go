package json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

func UnmarshalFile(filename string, v interface{}) error {
	if data, err := ioutil.ReadFile(filename); err != nil {
		return err
	} else {
		return json.Unmarshal(bytes.TrimSpace(data), v)
	}
}

func MarshalToFile(filename string, v interface{}) error {
	if data, err := json.Marshal(v); err != nil {
		return err
	} else {
		return ioutil.WriteFile(filename, data, 0644)
	}
}

func MarshalIndentToFile(filename string, v interface{}) error {
	if data, err := json.MarshalIndent(v, "", "  "); err != nil {
		return err
	} else {
		return ioutil.WriteFile(filename, data, 0644)
	}
}
