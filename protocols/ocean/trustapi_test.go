package ocean

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/utu-crowdsale/defi-portal-scanner/collector"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
)

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
