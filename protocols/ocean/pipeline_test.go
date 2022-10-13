package ocean

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readJson(t *testing.T, path string) (u []byte) {
	f, err := os.Open(path)
	assert.Nil(t, err)
	u, err = ioutil.ReadAll(f)
	assert.Nil(t, err)
	return u
}

func TestPipelineAll(t *testing.T) {
	logger := log.Default()
	result, err := PipelineAll(logger)
	fmt.Println("len(addresses)", len(result.Addresses))
	assert.Nil(t, err)

	a, err := json.MarshalIndent(result.Addresses, "", "\t")
	assert.Nil(t, err)
	f, err := os.OpenFile("addresses.json", os.O_CREATE|os.O_WRONLY, 0644)
	assert.Nil(t, err)
	_, err = f.Write(a)
	assert.Nil(t, err)

	a, err = json.MarshalIndent(result.Assets, "", "\t")
	assert.Nil(t, err)
	f, err = os.OpenFile("assets.json", os.O_CREATE|os.O_WRONLY, 0644)
	assert.Nil(t, err)
	_, err = f.Write(a)
	assert.Nil(t, err)
}
