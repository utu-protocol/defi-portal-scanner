package collector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/utu-crowdsale/defi-portal-scanner/utils"
)

// EtherscanReply a reply from etherscan
type EtherscanReply struct {
	Status  string           `json:"status,omitempty"`
	Message string           `json:"message,omitempty"`
	Result  []EthTransaction `json:"result,omitempty"`
}

// ChecksumAddress is a custom type that guarantees that the Ethereum address (a string) is checksummed
type ChecksumAddress string

func NewAddressFromString(s string) ChecksumAddress {
	return ChecksumAddress(utils.ChecksumAddress(s))
}

func (a *ChecksumAddress) UnmarshalJSON(data []byte) error {
	var addr string
	if err := json.Unmarshal(data, &addr); err != nil {
		return err
	}

	addr = utils.ChecksumAddress(addr)
	*a = ChecksumAddress(addr)
	return nil
}

// EthTransaction a transaction from etherescan
type EthTransaction struct {
	BlockNumber       string          `json:"blockNumber,omitempty"`
	TimeStamp         string          `json:"timeStamp,omitempty"`
	Hash              string          `json:"hash,omitempty"`
	Nonce             string          `json:"nonce,omitempty"`
	BlockHash         string          `json:"blockHash,omitempty"`
	TransactionIndex  string          `json:"transactionIndex,omitempty"`
	From              ChecksumAddress `json:"from"`
	To                ChecksumAddress `json:"to"`
	Value             string          `json:"value,omitempty"`
	Gas               string          `json:"gas,omitempty"`
	GasPrice          string          `json:"gasPrice,omitempty"`
	IsError           string          `json:"isError,omitempty"`
	TxReceiptStatus   string          `json:"txreceipt_status,omitempty"`
	Input             string          `json:"input,omitempty"`
	ContractAddress   string          `json:"contractAddress,omitempty"`
	CumulativeGasUsed string          `json:"cumulativeGasUsed,omitempty"`
	GasUsed           string          `json:"gasUsed,omitempty"`
	Confirmations     string          `json:"confirmations,omitempty"`
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
}

// NewEtherscanClient create a new utu client
func NewEtherscanClient(apiToken string) *EtherscanClient {
	return &EtherscanClient{
		APIEndpoint: "https://api.etherscan.io/api",
		APIToken:    apiToken,
		HTTPCli: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// GetTransactions get the list of transactions
func (c EtherscanClient) GetTransactions(address string) (txs []EthTransaction, err error) {
	// create request
	req, err := http.NewRequest("GET", c.APIEndpoint, nil)
	if err != nil {
		return
	}
	// build parameters
	q := req.URL.Query()
	q.Add("module", "account")
	q.Add("action", "txlist")
	q.Add("address", address)
	q.Add("apikey", c.APIToken)
	//TODO make it parametrized
	q.Add("page", "1")     // which page
	q.Add("offset", "100") // how many items
	req.URL.RawQuery = q.Encode()
	// execute the request
	res, err := c.HTTPCli.Do(req)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	//log.Debugf("etherscan reply: %s", data)
	var r EtherscanReply
	// now parse the result
	err = json.Unmarshal(data, &r)
	txs = r.Result
	return

}
