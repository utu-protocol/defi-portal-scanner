package ocean

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestPools(t *testing.T) {
// 	q := `{
// 		pools{
// 			controller,
// 			totalSwapVolume,
// 			datatokenReserve,
// 			datatokenAddress,
// 			spotPrice,
// 			consumePrice,
// 			tokens{denormWeight, tokenAddress, balance}
// 		  }
// 	}`
// 	err = graphQuery(q, nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Println(respData)
// }

func TestPipeline(t *testing.T) {
	logger := log.Default()
	err := pipeline(logger)
	assert.Nil(t, err)
}
