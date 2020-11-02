package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
	"github.com/utu-crowdsale/defi-portal-scanner/utils"
)

type TxStats struct {
	In  int
	Out int
}

var (
	tokenNames map[string]Token
)

func init() {
	// init vars
	tokenNames = make(map[string]Token)
}

// EthEvent an contract interaction
type EthEvent struct {
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
	Amount          uint64    `json:"amount,omitempty"`
}

// LogEvent print an event on stdout
func LogEvent(evt *EthEvent) []byte {

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
	log.Printf(format, "amount", evt.Amount)
	log.Printf(format, "topics", strings.Join(evt.Topics, "\n\t"))

	if evt.ContractAddress == evt.FromAddress {
		log.Warnf("ACTION %s: CONTRACT and SENDER match: %s (%s)", evt.Action, evt.ContractAddress, evt.ContractName)
	}
	if evt.ContractAddress == evt.ToAddress {
		log.Warnf("ACTION %s: CONTRACT and RECIPIENT match: %s (%s)", evt.Action, evt.ContractAddress, evt.ContractName)
	}
	if evt.FromAddress == evt.ToAddress {
		log.Warnf("ACTION %s: SENDER and RECIPIENT match: %s (%s)", evt.Action, evt.ContractAddress, evt.ContractName)
	}

	log.Println("-------- -------- -------- --------")
	data, _ := json.Marshal(evt)
	return data
}

func tokenName(address string) (name string) {
	t, found := tokenNames[strings.ToLower(address)]
	if !found {
		return
	}
	name = t.Name
	return
}

func ParseLog(vLog types.Log, client *ethclient.Client) (m *EthEvent, err error) {
	m = &EthEvent{
		BlockNumber:     vLog.BlockNumber,
		TransactionHash: vLog.TxHash.Hex(),
		ContractAddress: vLog.Address.Hex(),
		ContractName:    tokenName(vLog.Address.Hex()),
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
		return
	}
	if isPending {
		log.Warnf("transaction %s is pending, skipped", vLog.TxHash)
		return
	}
	//
	if r := tx.To(); r != nil {
		m.ToAddress = r.Hex()
		m.ToName = tokenName(r.Hex())
	}

	switch m.Action {
	case "Transfer":
		m.FromAddress = common.BytesToAddress(vLog.Topics[1].Bytes()).Hex()
		m.FromName = tokenName(m.FromAddress)
		m.ToAddress = common.BytesToAddress(vLog.Topics[2].Bytes()).Hex()
		m.ToName = tokenName(m.ToAddress)
		var a big.Int
		a.SetBytes(vLog.Data)
		m.Amount = a.Uint64()
	default:
		// sender
		sender, e := client.TransactionSender(context.Background(), tx, vLog.BlockHash, vLog.TxIndex)
		if e != nil {
			err = e
			log.Error(err)
			return
		}
		m.FromAddress = sender.Hex()
		m.FromName = tokenName(sender.Hex())
	}
	// timestamp
	block, err := client.BlockByHash(context.Background(), vLog.BlockHash)
	if err != nil {
		log.Error(err)
		return
	}
	m.BlockTime = time.Unix(int64(block.Time()), 0)
	return
}

// Start the service
func Start(cfg config.Schema) (err error) {
	client, err := ethclient.Dial(cfg.Ethereum.WssURL)
	if err != nil {
		return
	}
	// open the database
	store, err := OpenStore(cfg.DbFolder)
	if err != nil {
		return
	}
	// now get the etherescan api
	// escli := etherscan.New(etherscan.Mainnet, cfg.EtherscanAPIToken)
	// read the list of monitored tokens
	var tokens []Token
	err = getTokens(cfg.DefiSourcesFile, &tokens)
	if err != nil {
		return
	}
	// load the tokens
	tokenUpdated := 0 // wherever the token file was updated
	var addresses []common.Address
	for i, t := range tokens {
		log.Infof("registering %s with name %s", t.Address, t.Name)
		tokens[i] = t
		//lowercase the
		tokenNames[strings.ToLower(t.Address)] = t
		addresses = append(addresses, common.HexToAddress(t.Address))
	}
	if tokenUpdated > 0 {
		log.Info("Token file received %d updates", tokenUpdated)
		if err := utils.WriteJSON(cfg.DefiSourcesFile, tokens); err != nil {
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
	f, err := os.OpenFile(cfg.LogOutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
			m, err := ParseLog(vLog, client)
			if err != nil {
				continue
			}

			// does nothing
			store.Tx(m.FromAddress, m.ToAddress, m.BlockTime)
			// aggregate
			queue(*m)

			// write it
			data := LogEvent(m)
			f.Write(data)
			f.WriteString("\n")
		}
	}
}
