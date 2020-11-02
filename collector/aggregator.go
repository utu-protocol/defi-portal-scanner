package collector

import (
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	txBufferM  sync.RWMutex
	txBuffer   map[string]chan EthEvent
	doneBuffer chan string
)

func init() {
	txBuffer = make(map[string]chan EthEvent)
	doneBuffer = make(chan string)
	go cleaner()
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
			sb.WriteString(fmt.Sprintf("Contract %s\n", evts[0].ContractName))
			sb.WriteString(fmt.Sprintf("Tx       %s\n", fmt.Sprint("https://etherscan.io/tx/", evts[0].TransactionHash)))

			for i, v := range evts {
				s := _name(v.FromName, v.FromAddress)
				r := _name(v.ToName, v.ToAddress)
				sb.WriteString(fmt.Sprintf("%d. %20s     %20s -> %20s\n", i+1, v.Action, s, r))
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
	log.Info("Cleaner started")
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
