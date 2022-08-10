package wallet

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
)

type Wallet struct {
	Network  string    `json:"json"`
	Address  string    `json:"address"`
	Balances []Balance `json:"balances"`
}

type Balance struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Network string `json:"network"`
}

type Response struct {
	APIResponse interface{}
}

type scanBalancesFunc func(apiURL string) (*Response, error)

type balanceMapperFunc func(response *Response, network string) []Balance

type BalanceScannerConf struct {
	apiURL            string
	balancesScanner   scanBalancesFunc
	balanceMapperFunc balanceMapperFunc
}

var (
	addressesToScan      = make(map[string]struct{})
	placeholder          struct{}
	networkConfigMapping = make(map[string]BalanceScannerConf)
	networks             = []string{"ethereum", "polygon", "gnosis"}
	scannerFunctions     = map[string]scanBalancesFunc{
		"ethereum": scanCovalentBalance,
		"polygon":  scanCovalentBalance,
		"gnosis":   scanBlockscoutBalance,
	}
	balanceMapperFunctions = map[string]balanceMapperFunc{
		"ethereum": mapCovalentResponse,
		"polygon":  mapCovalentResponse,
		"gnosis":   mapToBalance,
	}
)

func Ready(cfg config.Schema) {
	for _, network := range networks {
		networkConfigMapping[network] = BalanceScannerConf{
			apiURL:            cfg.BalanceAPI[network],
			balancesScanner:   scannerFunctions[network],
			balanceMapperFunc: balanceMapperFunctions[network],
		}
	}
	log.Infof("Initialized balance scanner config %v", networkConfigMapping)
}

func Scan(address string, ch chan<- *Wallet) {
	if _, ok := addressesToScan[address]; ok {
		return
	}
	addressesToScan[address] = placeholder
	scan(address, ch)
	log.Infof("Started scanning token balances for %s", address)
}

func ScanCached(ch chan<- *Wallet) {
	log.Info("Scanning cached addresses...")
	for address := range addressesToScan {
		scan(address, ch)
	}
}

func scan(address string, ch chan<- *Wallet) {
	for network, networkScanConfig := range networkConfigMapping {
		log.Infof("Scanning tokens balances for network=%s", network)
		balances := networkScanConfig.scanBalances(address, network)
		if len(balances) > 0 {
			wallet := &Wallet{
				Network:  network,
				Address:  address,
				Balances: balances,
			}
			ch <- wallet
		} else {
			log.Infof("No token balances found for %s", address)
		}
	}
}

func (w *BalanceScannerConf) scanBalances(address string, network string) []Balance {
	requestURL := fmt.Sprintf(w.apiURL, address)
	response, err := w.balancesScanner(requestURL)
	if err != nil {
		log.Errorf("Failed to scan balances for %s, err=%s", address, err.Error())
		return []Balance{}
	}
	return w.balanceMapperFunc(response, network)
}
