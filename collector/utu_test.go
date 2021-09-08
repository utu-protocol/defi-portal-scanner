package collector

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
)

func TestUTUClient(t *testing.T) {
	apiKey := os.Getenv("APIKEY")
	s := &config.TrustEngineSchema{
		URL:           "https://stage-api.ututrust.com/core-api",
		Authorization: apiKey,
		DryRun:        false,
	}
	utu := NewUTUClient(*s)
	testEntity := NewTrustEntity("Pool 1")
	testEntity.Type = "pool"
	testEntity.Ids = map[string]string{"address": "0xdeadbeef"}
	err := utu.PostEntity(testEntity)
	assert.Nil(t, err)
}
