package ocean

import (
	"fmt"
	"testing"
)

func TestPools(t *testing.T) {
	q := `{
		pools{
			controller,
			totalSwapVolume,
			datatokenReserve,
			datatokenAddress,
			spotPrice,
			consumePrice,
			tokens{denormWeight, tokenAddress, balance}
		  }
	}`
	respData, err := grQuery(q)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(respData)
}
