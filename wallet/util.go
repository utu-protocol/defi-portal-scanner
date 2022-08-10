package wallet

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
)

func formatBalance(balance *big.Int, decimals int) string {
	dec := big.NewInt(int64(math.Pow(10, float64(decimals))))
	integral, modulo := big.NewInt(0).DivMod(balance, dec, big.NewInt(0))
	if modulo.BitLen() == 0 {
		return integral.String()
	} else {
		return fmt.Sprintf("%s.%s", integral.String(), modulo.String())
	}
}

func httpGet(apiUrl string) ([]byte, error) {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
