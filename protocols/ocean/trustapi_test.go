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

func TestPostAddressesToUTU(t *testing.T) {
	logger := log.Default()
	as := readJson(t, "assets.json")
	var assets []*Asset
	err := json.Unmarshal(as, &assets)
	if err != nil {
		t.Fatal(err)
	}

	ad := readJson(t, "addresses.json")
	var addresses []*Address
	err = json.Unmarshal(ad, &addresses)
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
	PostAddressesToUTU(addresses, assets, utu, logger)
}
