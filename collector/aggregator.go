package collector

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
)

var (
	txBufferM  sync.RWMutex
	txBuffer   map[string]chan EthEvent
	doneBuffer chan string
)

func init() {
	// txBuffer = make(map[string]chan EthEvent)
	// doneBuffer = make(chan string)
	// go cleaner()
}

func _name(name, address string) string {
	if name == "" {
		return address
	}
	return name
}

func aggregate(events chan EthEvent) {

	evts := []EthEvent{}

	time.AfterFunc(time.Duration(5)*time.Second, func() {
		close(events)
	})

	for {
		evt, more := <-events
		if !more {
			var sb strings.Builder
			sb.WriteString("---------------TX SUMMARY-------------\n")
			sb.WriteString(fmt.Sprintf("Protocol %s\n", evts[0].Context.Name))
			sb.WriteString(fmt.Sprintf("Tx       %s\n", fmt.Sprint("https://etherscan.io/tx/", evts[0].TransactionHash)))

			for i, v := range evts {
				sb.WriteString(fmt.Sprintf("%d. %20s     %20v -> %20v\n", i+1, v.Action, v.Recipients, v.Senders))
			}
			sb.WriteString("---------------//////////-------------\n")
			log.Info(sb.String())

			doneBuffer <- evts[0].TransactionHash
			break
		}
		evts = append(evts, evt)
	}

}

func cleaner() {
	for {
		txh, more := <-doneBuffer
		if !more {
			break
		}
		txBufferM.Lock()
		delete(txBuffer, txh)
		txBufferM.Unlock()
	}
	log.Info("cleaner done")
}

func queue(e EthEvent) {
	txBufferM.Lock()

	evtChan, found := txBuffer[e.TransactionHash]
	if !found {
		evtChan = make(chan EthEvent)
		go aggregate(evtChan)
		txBuffer[e.TransactionHash] = evtChan
	}
	evtChan <- e
	txBufferM.Unlock()
}

// EthEvent an contract interaction
type EthEvent struct {
	BlockNumber     uint64          `json:"block_number,omitempty"`
	BlockTime       time.Time       `json:"block_time,omitempty"`
	Context         config.Protocol `json:"actor,omitempty"`
	Action          string          `json:"action,omitempty"`
	TransactionHash string          `json:"tx_hash,omitempty"`
	Recipients      []string        `json:"recipients,omitempty"`
	Senders         []string        `json:"senders,omitempty"`
}

// LogEvent print an event on stdout
func LogEvent(evt *EthEvent) []byte {

	format := "%20s: %v"
	log.Printf(format, "block_number", fmt.Sprint("https://etherscan.io/blocks/", evt.BlockNumber))
	log.Printf(format, "block_time", evt.BlockTime)
	// log.Printf(format, "contract_address", fmt.Sprint("https://etherscan.io/address/", evt.Context.Address))
	log.Printf(format, "contract_name", evt.Context.Name)
	log.Printf(format, "protocol", evt.Context.Name)
	log.Printf(format, "tx_hash", fmt.Sprint("https://etherscan.io/tx/", evt.TransactionHash))
	log.Printf(format, "recipients:", "")
	for _, v := range evt.Recipients {
		log.Printf(format, "from_address", fmt.Sprint("https://etherscan.io/address/", v))
	}
	for _, v := range evt.Senders {
		log.Printf(format, "to_address", fmt.Sprint("https://etherscan.io/address/", v))
	}
	log.Printf(format, "action", evt.Action)
	//log.Printf(format, "amount", evt.Amount)
	//log.Printf(format, "topics", strings.Join(evt.Topics, "\n\t"))

	log.Println("-------- -------- -------- --------")
	data, _ := json.Marshal(evt)
	return data
}
