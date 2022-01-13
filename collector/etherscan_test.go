package collector

import (
	"os"
	"testing"
)

func TestEtherscanClientGetPagedTransactions(t *testing.T) {
	apiKey := os.Getenv("ETHERSCAN_API_TOKEN")

	c := NewEtherscanClient(apiKey)
	c.PageSize = 100
	_, err := c.GetTransactions("0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a")
	if err != nil {
		t.Error(err)
	}
}
