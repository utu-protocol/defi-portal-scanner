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

func TestPipelineAssets(t *testing.T) {
	logger := log.Default()
	assets, err := PipelineAssets(logger)
	assert.Nil(t, err)

	a, err := json.MarshalIndent(assets, "", "\t")
	assert.Nil(t, err)
	f, err := os.OpenFile("assets.json", os.O_CREATE|os.O_WRONLY, 0644)
	assert.Nil(t, err)
	_, err = f.Write(a)
	assert.Nil(t, err)
}

func TestPipelineUsers(t *testing.T) {
	logger := log.Default()
	addresses, err := PipelineUsers(logger)
	fmt.Println("len(addresses)", len(addresses))
	assert.Nil(t, err)

	a, err := json.MarshalIndent(addresses, "", "\t")
	assert.Nil(t, err)
	f, err := os.OpenFile("addresses.json", os.O_CREATE|os.O_WRONLY, 0644)
	assert.Nil(t, err)
	_, err = f.Write(a)
	assert.Nil(t, err)
}
