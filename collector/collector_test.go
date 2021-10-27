package collector

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
)

func checkEntityHasLowercaseAddress(t *testing.T, te *TrustEntity) {
	if te.Name != "" {
		assert.Equal(t, strings.ToLower(te.Name), te.Name)
	}
	assert.Equal(t, strings.ToLower(te.Ids["address"]), te.Ids["address"])
}

func TestAddressProcessorEmitsChangesetsWithLowercaseAddresses(t *testing.T) {
	cfg := new(config.Schema)
	go func(t *testing.T) {
		for {
			cs, more := <-csQueue
			if !more {
				break
			}
			for _, e := range cs.Entities {
				checkEntityHasLowercaseAddress(t, e)
			}
			for _, r := range cs.Relationship {
				checkEntityHasLowercaseAddress(t, r.SourceCriteria)
				checkEntityHasLowercaseAddress(t, r.TargetCriteria)
			}
		}
	}(t)
	go addressProcessor(*cfg)

	// If you put in a lowercase string here into the channel, yes of course it
	// will break. However, addrQueue is only added to from Scan(), which is
	// started by server.go:Serve(). As long as that converts any user input
	// from a string into a Address, we are safe.
	addrQueue <- NewAddressFromString("0x0000000000007f150bd6f54c40a34d7c3d5e9f56")
	addrQueue <- NewAddressFromString("0x0000000000007F150Bd6f54c40A34d7C3d5e9f56")
	// here the cache should detect that they are the same address and skip scanning
	addrQueue <- NewAddressFromString("0xDe5CAf81E2446BA4BAf9A35E1DB1ecF247f1eF89")
}
