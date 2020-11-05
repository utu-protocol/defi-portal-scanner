package utils

import (
	"encoding/json"
	"io/ioutil"
)

// WriteJSON write data to a json file
func WriteJSON(file string, data interface{}) (err error) {
	bin, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(file, bin, 0600)
	return
}

// ReadJSON write data to a json file
func ReadJSON(file string, target interface{}) (err error) {
	bin, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(bin, target)
	return
}
