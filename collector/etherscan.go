package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// EtherscanReply a reply from etherscan
type EtherscanReply struct {
	Status  string           `json:"status,omitempty"`
	Message string           `json:"message,omitempty"`
	Result  []EthTransaction `json:"result,omitempty"`
}

// Address is a custom type that guarantees that the Ethereum address (a string)
// is lowercased/checksummed
type Address string

func NewAddressFromString(s string) Address {
	return Address(strings.ToLower(s))
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var addr string
	if err := json.Unmarshal(data, &addr); err != nil {
		return err
	}

	// this class was originally written to ensure that Addresses were
	// checksummed everywhere. Now we ensure that it is lowercase everywhere,
	// but the following line is kept, commented out to explain the context of
	// why this class exists

	// addr = utils.ChecksumAddress(addr)

	// We decided for lowercase everywhere
	addr = strings.ToLower(addr)
	*a = Address(addr)
	return nil
}

type NoMoreTransactionsError struct {
	address Address
	page    int
	offset  int
}

func (n *NoMoreTransactionsError) Error() string {
	return fmt.Sprintf("No more transactions for %s: page %d offset %d", n.address, n.page, n.offset)
}

// EthTransaction a transaction from etherescan
type EthTransaction struct {
	BlockNumber       string  `json:"blockNumber,omitempty"`
	TimeStamp         string  `json:"timeStamp,omitempty"`
	Hash              string  `json:"hash,omitempty"`
	Nonce             string  `json:"nonce,omitempty"`
	BlockHash         string  `json:"blockHash,omitempty"`
	TransactionIndex  string  `json:"transactionIndex,omitempty"`
	From              Address `json:"from"`
	To                Address `json:"to"`
	Value             string  `json:"value,omitempty"`
	Gas               string  `json:"gas,omitempty"`
	GasPrice          string  `json:"gasPrice,omitempty"`
	IsError           string  `json:"isError,omitempty"`
	TxReceiptStatus   string  `json:"txreceipt_status,omitempty"`
	Input             string  `json:"input,omitempty"`
	ContractAddress   string  `json:"contractAddress,omitempty"`
	CumulativeGasUsed string  `json:"cumulativeGasUsed,omitempty"`
	GasUsed           string  `json:"gasUsed,omitempty"`
	Confirmations     string  `json:"confirmations,omitempty"`
}

// GetTime return the tx time
func (et EthTransaction) GetTime() (t time.Time) {
	v, err := strconv.ParseInt(et.TimeStamp, 10, 64)
	if err != nil {
		return
	}
	t = time.Unix(int64(v), 0)
	return
}

// EtherscanClient trust api client
type EtherscanClient struct {
	APIEndpoint string
	APIToken    string
	HTTPCli     *http.Client
	PageSize    int
}

// NewEtherscanClient create a new utu client
func NewEtherscanClient(apiToken string) *EtherscanClient {
	return &EtherscanClient{
		APIEndpoint: "https://api.etherscan.io/api",
		APIToken:    apiToken,
		HTTPCli: &http.Client{
			Timeout: time.Second * 10,
		},
		PageSize: 100,
	}
}

// GetTransactions gets normal (not internal) Transactions from Etherscan.
// Etherscan will only return 10000 records maximum, regardless of how many
// pages you request/size of those pages.
func (c EtherscanClient) GetTransactions(address Address) (txs []EthTransaction, err error) {
	var pagedTxs []EthTransaction
	page := 1
	for err == nil {
		pagedTxs, err = c.getPagedTransactions(address, page, c.PageSize)
		txs = append(txs, pagedTxs...)
		page++
	}
	_, ok := err.(*NoMoreTransactionsError)
	if ok {
		return txs, nil
	}
	return

}

// getPagedTransactions execute GET query and parse possible responses
func (c EtherscanClient) getPagedTransactions(address Address, page, offset int) (txs []EthTransaction, err error) {

	req, err := http.NewRequest("GET", c.APIEndpoint, nil)
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add("module", "account")
	q.Add("action", "txlist")
	q.Add("address", string(address))
	q.Add("apikey", c.APIToken)
	q.Add("page", fmt.Sprint(page))     // which page
	q.Add("offset", fmt.Sprint(offset)) // how many items
	req.URL.RawQuery = q.Encode()

	res, err := c.HTTPCli.Do(req)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var r EtherscanReply
	err = json.Unmarshal(data, &r)
	txs = r.Result

	// if no more transactions are found, r.status will also be 0
	if r.Message == "No transactions found" && len(txs) == 0 {
		return txs, &NoMoreTransactionsError{
			address: address,
			page:    page,
			offset:  offset,
		}
	}

	if r.Status == "0" {
		return txs, fmt.Errorf(r.Message)
	}
	return

}
