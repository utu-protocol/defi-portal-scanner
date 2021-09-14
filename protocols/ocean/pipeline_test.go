package ocean

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utu-crowdsale/defi-portal-scanner/collector"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
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
	_, err := pipelineAssets(logger)
	assert.Nil(t, err)

	// a, err := json.MarshalIndent(assets, "", "\t")
	// assert.Nil(t, err)
	// f, err := os.OpenFile("assets.json", os.O_CREATE|os.O_WRONLY, 0644)
	// assert.Nil(t, err)
	// _, err = f.Write(a)
	// assert.Nil(t, err)
}

func TestPipelineUsers(t *testing.T) {
	logger := log.Default()
	users, err := pipelineUsers(logger)
	fmt.Println("users", users)
	assert.Nil(t, err)

	u, err := json.MarshalIndent(users, "", "\t")
	assert.Nil(t, err)
	f, err := os.OpenFile("users.json", os.O_CREATE|os.O_WRONLY, 0644)
	assert.Nil(t, err)
	_, err = f.Write(u)
	assert.Nil(t, err)
}

func TestPostAssetsToUTU(t *testing.T) {
	logger := log.Default()
	a := readJson(t, "assets.json")
	var assets []*Asset
	err := json.Unmarshal(a, &assets)
	if err != nil {
		t.Fatal(err)
	}

	apiKey := os.Getenv("APIKEY")
	s := &config.TrustEngineSchema{
		URL:           "https://stage-api.ututrust.com/core-api",
		Authorization: apiKey,
		DryRun:        false,
	}
	utu := collector.NewUTUClient(*s)
	PostAssetsToUTU(assets, utu, logger)
}

func TestPostToUTU(t *testing.T) {
	logger := log.Default()
	a := readJson(t, "assets.json")
	var assets []*Asset
	err := json.Unmarshal(a, &assets)
	if err != nil {
		t.Fatal(err)
	}

	u := readJson(t, "users.json")
	var users []*User
	err = json.Unmarshal(u, &users)
	if err != nil {
		t.Fatal(err)
	}

	apiKey := os.Getenv("APIKEY")
	s := &config.TrustEngineSchema{
		URL:           "https://stage-api.ututrust.com/core-api",
		Authorization: apiKey,
		DryRun:        false,
	}
	utu := collector.NewUTUClient(*s)
	PostToUTU(users, assets, utu, logger)
}
