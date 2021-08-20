package ocean

import geth "github.com/ethereum/go-ethereum/common"

func checksumAddress(address string) string {
	return geth.HexToAddress(address).String()
}
