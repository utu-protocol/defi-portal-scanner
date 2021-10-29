package utils

import (
	"encoding/json"
	"io/ioutil"

	geth "github.com/ethereum/go-ethereum/common"
)

// WriteJSON write data to a json file
func WriteJSON(file string, pretty bool, data interface{}) (err error) {
	var bin []byte
	if pretty {
		bin, err = json.MarshalIndent(data, "", "\t")
	} else {
		bin, err = json.Marshal(data)
	}
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

func ChecksumAddress(address string) string {
	return geth.HexToAddress(address).String()
}
