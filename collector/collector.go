package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
)

var (
	tokenNames map[string]Token
	eventNames map[string]string
)

func init() {
	// init vars
	tokenNames = make(map[string]Token)
	// action names
	eventNames = map[string]string{
		"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef": "Transfer",
		"0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822": "Swap",
		"0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496": "Burn",
		"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925": "Approval",
		"0x875352fb3fadeb8c0be7cbbe8ff761b308fa7033470cd0287f02f3436fd76cb9": "AccrueInterest",
		"0x4dec04e750ca11537cabcd8a9eab06494de08da3735bc8871cd41250e190bc04": "AccrueInterest",
		"0xe5b754fb1abb7f01b499791d0b820ae3b6af3424ac1c59768edb53f4ec31a929": "Redeem",
		"0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f": "Mint",
		"0x3c67396e9c55d2fc8ad68875fc5beca1d96ad2a2f23b210ccc1d986551ab6fdf": "TokensTransferred",
		"0x1a2a22cb034d26d1854bdc6666a5b91fe25efbbb5dcad3b0355478d6f5c362a1": "RepayBorrow",
		"0x45b96fe442630264581b197e84bbada861235052c5a1aadfff9ea4e40a969aa0": "Failure",
		"0xbd5034ffbd47e4e72a94baa2cdb74c6fad73cb3bcdc13036b72ec8306f5a7646": "Redeem",
		"0x5e3cad45b1fe24159d1cb39788d82d0f69cc15770aa96fba1d3d1a7348735594": "InterestStreamRedirected",
		"0x9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d": "Staked",
		"0x34fcbac0073d7c3d388e51312faf357774904998eeb8fca628b9e6f65ee1cbf7": "Claim",
		"0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885": "Mint",
	}
}

// Token is a source token
type Token struct {
	Name    string  `json:"name,omitempty"`
	Address string  `json:"address,omitempty"`
	AbiJSON string  `json:"abi,omitempty"`
	Abi     abi.ABI `json:"-"`
}

func readTokens(file string) (tokens []Token, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &tokens)
	return
}

func writeJSON(file string, data interface{}) (err error) {
	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer f.Close()
	bin, err := json.Marshal(data)
	if err != nil {
		return
	}
	_, err = f.Write(bin)
	return
}

// DefiInteraction an contract interaction
type DefiInteraction struct {
	BlockNumber     uint64    `json:"block_number,omitempty"`
	BlockTime       time.Time `json:"block_time,omitempty"`
	ContractAddress string    `json:"contract_address"`
	ContractName    string    `json:"contract_name"`
	TransactionHash string    `json:"tx_hash,omitempty"`
	FromAddress     string    `json:"from_address,omitempty"`
	FromName        string    `json:"from_name,omitempty"`
	ToAddress       string    `json:"to_address,omitempty"`
	ToName          string    `json:"to_name,omitempty"`
	Topics          []string  `json:"topics,omitempty"`
	Action          string    `json:"action,omitempty"`
}

// LogEvent print an event on stdout
func LogEvent(evt DefiInteraction) []byte {
	log.Println("-------- -------- -------- --------")
	format := "%20s: %v"
	log.Printf(format, "block_number", fmt.Sprint("https://etherscan.io/blocks/", evt.BlockNumber))
	log.Printf(format, "block_time", evt.BlockTime)
	log.Printf(format, "contract_address", fmt.Sprint("https://etherscan.io/address/", evt.ContractAddress))
	log.Printf(format, "contract_name", evt.ContractName)
	log.Printf(format, "tx_hash", fmt.Sprint("https://etherscan.io/tx/", evt.TransactionHash))
	log.Printf(format, "from_address", fmt.Sprint("https://etherscan.io/address/", evt.FromAddress))
	log.Printf(format, "from_name", evt.FromName)
	log.Printf(format, "to_address", fmt.Sprint("https://etherscan.io/address/", evt.ToAddress))
	log.Printf(format, "to_name", evt.ToName)
	log.Printf(format, "action", evt.Action)
	log.Printf(format, "topics", strings.Join(evt.Topics, "\n\t"))
	data, _ := json.Marshal(evt)
	return data
}

// Start the service
func Start(cfg config.Schema) (err error) {
	client, err := ethclient.Dial(cfg.EthWssURL)
	if err != nil {
		return
	}
	// now get the etherescan api
	//escli := etherscan.New(etherscan.Mainnet, cfg.EtherscanAPIToken)
	// read the list of monitored tokens
	tokens, err := readTokens(cfg.DefiSourcesFile)
	if err != nil {
		log.Fatal(err)
		log.Info()
	}
	// load the tokens
	tokenUpdated := 0 // wherever the token file was updated
	var addresses []common.Address
	for i, t := range tokens {
		log.Infof("registering %s with name %s", t.Address, t.Name)
		// if t.AbiJSON == "" {
		// 	abiJSON, err := escli.ContractABI(t.Address)
		// 	if err != nil {
		// 		log.Warnf("can't retrieve the ABI for %s: %v", t.Address, err)
		// 		continue
		// 	}
		// 	t.AbiJSON = abiJSON
		// 	tokenUpdated++
		// }
		// now parse the json abi
		// abi, err := abi.JSON(strings.NewReader(t.AbiJSON))
		// if err != nil {
		// 	log.Warnf("cannot parse the json ABI for %s: %v", t.Address, err)
		// }
		// t.Abi = abi
		// all good
		tokens[i] = t
		tokenNames[t.Address] = t
		addresses = append(addresses, common.HexToAddress(t.Address))
	}
	if tokenUpdated > 0 {
		log.Info("Token file received %d updates", tokenUpdated)
		if err := writeJSON(cfg.DefiSourcesFile, tokens); err != nil {
			log.Warnf("failed to update the tokens file %s: %v", cfg.DefiSourcesFile, err)
		}
	}
	log.Infof("registered %d names", len(tokenNames))
	// prepare query
	query := ethereum.FilterQuery{
		Addresses: addresses,
	}
	// prepare the channel for subscrition
	logs := make(chan types.Log)
	// make the query
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	// propare output
	f, err := os.Create(cfg.LogOutputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// get them
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			m := DefiInteraction{
				BlockNumber:     vLog.BlockNumber,
				TransactionHash: vLog.TxHash.Hex(),
				ContractAddress: vLog.Address.Hex(),
				ContractName:    tokenNames[vLog.Address.Hex()].Name,
			}
			// action
			for _, t := range vLog.Topics {
				m.Topics = append(m.Topics, t.Hex())
			}
			if action, found := eventNames[m.Topics[0]]; !found {
				log.Errorf("undefined name for action signature %s", m.Topics[0])
			} else {
				m.Action = action
			}

			// recipient
			tx, isPending, err := client.TransactionByHash(context.Background(), vLog.TxHash)
			if err != nil {
				log.Error(err)
				continue
			}
			if isPending {
				log.Warnf("transaction %s is pending, skipped", vLog.TxHash)
				continue
			}
			//
			if r := tx.To(); r != nil {
				m.ToAddress = r.Hex()
				m.ToName = tokenNames[r.Hex()].Name
			}

			// sender
			sender, err := client.TransactionSender(context.Background(), tx, vLog.BlockHash, vLog.TxIndex)
			if err != nil {
				log.Error(err)
				continue
			}
			m.FromAddress = sender.Hex()
			m.FromName = tokenNames[sender.Hex()].Name
			// timestamp
			block, err := client.BlockByHash(context.Background(), vLog.BlockHash)
			if err != nil {
				log.Error(err)
				continue
			}
			m.BlockTime = time.Unix(int64(block.Time()), 0)

			// write it
			data := LogEvent(m)
			f.Write(data)
			f.WriteString("\n")
		}
	}
}
