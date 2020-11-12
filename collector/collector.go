package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
	"github.com/utu-crowdsale/defi-portal-scanner/utils"
)

// some constants
const (
	ZeroAddress = "0x0000000000000000000000000000000000000000"
)

var (
	csQueue   chan *TrustAPIChangeSet
	addrQueue chan string
)

func init() {
	csQueue = make(chan *TrustAPIChangeSet)
	addrQueue = make(chan string)
}

func topic2Addr(l *types.Log, index int) string {
	return common.BytesToAddress(l.Topics[index].Bytes()).Hex()
}

func criteria(address string) (entity *TrustEntity, isNew bool) {
	// cache lookup
	label, typ, found := cacheGet(address)
	if !found {
		// here is a user, we store 0x123, address, address
		typ = TypeAddress
		label = address
		cachePush(address, label, typ)
		isNew = true
	}
	// create the entity to be used as criteria
	entity = NewTrustEntity()
	entity.Type = typ
	entity.Ids = map[string]string{"address": label}
	return
}

// ParseLog take a log and return an Event
func ParseLog(vLog *types.Log, client *ethclient.Client) (cs TrustAPIChangeSet, err error) {

	// action
	action, found := eventNames[vLog.Topics[0].Hex()]
	if !found {
		err = fmt.Errorf("undefined name for action signature %s", vLog.Topics[0].Hex())
		log.Error(err)
		return
	}

	// recipient
	_, isPending, err := client.TransactionByHash(context.Background(), vLog.TxHash)
	if err != nil {
		log.Error(err)
		return
	}
	if isPending {
		err = fmt.Errorf("transaction %s is pending, skipped", vLog.TxHash)
		log.Warn(err)
		return
	}

	// timestamp
	block, err := client.BlockByHash(context.Background(), vLog.BlockHash)
	if err != nil {
		log.Error(err)
		return
	}

	// now parse the types
	var rel *TrustRelationship
	switch action {
	case "Transfer":
		// process entities
		contractAddress := vLog.Address.Hex()
		senderAddress := topic2Addr(vLog, 1)
		recipientAddress := topic2Addr(vLog, 2)

		// skip 0x0 address
		if senderAddress == ZeroAddress || recipientAddress == ZeroAddress {
			err = fmt.Errorf("skip tx %s event log: zero-address detected", vLog.TxHash.Hex())
			return
		}

		// case sender is a defi-ptocol
		c, _ := criteria(contractAddress)
		s, sIsNew := criteria(senderAddress)
		r, rIsNew := criteria(recipientAddress)

		if s.Type == r.Type {
			// if they are both defi-portal then skip
			if s.Type == TypeDefiProtocol {
				err = fmt.Errorf("skip tx %s event log:  both sender and recipient are defi-protocols", vLog.TxHash.Hex())
				return
			}
			// if they are both address
			// then create 2 relationships to the contract
			rel = NewTrustRelationship()
			rel.Type = TypeInteraction
			rel.Properties = map[string]interface{}{
				"txId":      vLog.TxHash.Hex(),
				"action":    action,
				"timestamp": time.Unix(int64(block.Time()), 0),
			}
			rel.SourceCriteria = s // the sender is the source
			rel.TargetCriteria = c
			cs.AddRel(rel)
			// second one
			rel := NewTrustRelationship()
			rel.Type = TypeInteraction
			rel.Properties = map[string]interface{}{
				"txId":      vLog.TxHash.Hex(),
				"action":    action,
				"timestamp": time.Unix(int64(block.Time()), 0),
			}
			rel.SourceCriteria = r // the recipient is the source
			rel.TargetCriteria = c
			cs.AddRel(rel)
		} else {
			rel = NewTrustRelationship()
			rel.Type = TypeInteraction
			rel.Properties = map[string]interface{}{
				"txId":      vLog.TxHash.Hex(),
				"action":    action,
				"timestamp": time.Unix(int64(block.Time()), 0),
			}
			if s.Type == TypeAddress {
				// if the sender is type address and recipient defi-portal
				// then best case scenario
				rel.SourceCriteria = s // the sender is the source
				rel.TargetCriteria = r
				cs.AddRel(rel)
			} else {
				// if the sender is type defi-portal and sender address
				// then swap them around
				rel.SourceCriteria = r // the sender is the source
				rel.TargetCriteria = s
				cs.AddRel(rel)
			}

		}

		// now add missing stuff
		if sIsNew {
			// TODO copying here is ugly
			entity := NewTrustEntity()
			entity.Ids = s.Ids
			entity.Type = s.Type
			entity.Name = senderAddress
			entity.Image = fmt.Sprintf("https://via.placeholder.com/150/FFFF00/000000/?text=%s", senderAddress)
			cs.AddEntity(entity)
		}
		if rIsNew {
			// TODO copying here is ugly
			entity := NewTrustEntity()
			entity.Ids = r.Ids
			entity.Type = r.Type
			entity.Name = recipientAddress
			entity.Image = fmt.Sprintf("https://via.placeholder.com/150/FFFF00/000000/?text=%s", recipientAddress)
			cs.AddEntity(entity)
		}

	default:
		err = fmt.Errorf("action %s not supported", action)
	}
	return
}

func changesetsProcessor(cfg config.TrustEngineSchema) {
	utuCli := NewUTUClient(cfg)
	if cfg.DryRun {
		log.Info("Utu client is in dry run mode, CHANGES WILL NOT BE SUBMITTED!")
	}
	// listen on the queue
	for {
		cs, more := <-csQueue
		if !more {
			log.Info("changeset queue is closed, exiting")
			break
		}
		// if dryrun just print the outcome
		if cfg.DryRun {
			v, _ := json.MarshalIndent(cs, "", "  ")
			log.Infof("%s", v)
			continue
		}

		for _, e := range cs.Entities {
			// cache addresses
			for a, n := range e.Ids {
				// push to the address cache
				cachePush(a, n, e.Type)
			}
			// execute the request
			if err := utuCli.PostEntity(e); err != nil {
				log.Error("error posting entity:", err)
			}
		}

		for _, r := range cs.Relationship {
			// execute the request
			if err := utuCli.PostRelationship(r); err != nil {
				log.Error("error posting relationship:", err)
			}
		}
	}
}

// Ready setup the processing queue
func Ready(cfg config.Schema) {
	// start the processor
	go changesetsProcessor(cfg.UTUTrustAPI)
	go addressProcessor(cfg)

}

// Start the service
func Start(cfg config.Schema) (err error) {
	log.Info("starting collector for protocols at", cfg.DefiSourcesFile)
	client, err := ethclient.Dial(cfg.Ethereum.WssURL)
	if err != nil {
		return
	}
	// // open the database
	// store, err := OpenStore(cfg.DbFolder)
	// if err != nil {
	// 	return
	// }
	// log.Debug(store)
	// prepare the entities cache
	var addresFilters []common.Address

	// now get the etherescan api
	// escli := etherscan.New(etherscan.Mainnet, cfg.EtherscanAPIToken)
	// read the list of monitored protocols
	var protocols []Protocol
	err = utils.ReadJSON(cfg.DefiSourcesFile, &protocols)
	if err != nil {
		log.Errorf("cannot retrieve the defi protocols from %s: %v", cfg.DefiSourcesFile, err)
		return
	}

	for _, p := range protocols {
		// if there are no filters skip
		if len(p.Filters) == 0 {
			log.Warnf("skip protocol %s: empty filters", p.Name)
			continue
		}
		// build the entity
		protocolID := strcase.ToCamel(p.Name)
		log.Infof("protocol %s added with %d addresses", p.Name, len(p.Filters))
		//
		e := NewTrustEntity()
		e.Name = p.Name
		e.Type = TypeDefiProtocol
		e.Image = p.IconURL
		e.Ids = map[string]string{"address": protocolID}
		e.Properties = map[string]interface{}{
			"url":         p.URL,
			"description": p.Description,
			"category":    p.Category,
		}
		// queue it to the processor
		csQueue <- NewChangeset(e)
		// cache addresses
		for a := range p.Filters {
			// push to the address cache
			cachePush(a, protocolID, TypeDefiProtocol)
			// add to the list of filter for ethereum
			addresFilters = append(addresFilters, common.HexToAddress(a))
			log.Debugf("registered protocol %s filter %s at %s", p.Name, protocolID, a)
		}
	}

	log.Infof("registered %d filters", len(addresFilters))
	// prepare query
	query := ethereum.FilterQuery{
		Addresses: addresFilters,
	}
	// prepare the channel for subscrition
	logs := make(chan types.Log)
	// make the query
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	// propare output
	// f, err := os.OpenFile(cfg.LogOutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()

	// get them
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			// check if the log is for an address we know
			_, _, found := cacheGet(vLog.Address.Hex())
			if !found {
				err = fmt.Errorf("skip unknown contract address: %s ", vLog.Address.Hex())
				continue
			}
			// this should return a relationship
			changeset, err := ParseLog(&vLog, client)
			if err != nil {
				log.Warn("error parsing log: ", err)
				continue
			}
			csQueue <- &changeset
			// aggregate
			//queue(changeset)
		}
	}
}

func addressProcessor(cfg config.Schema) {
	cache := make(map[string]bool)
	// get the etherscan client
	client := NewEtherscanClient(cfg.Ethereum.EtherscanAPIToken)
	for {
		addr, more := <-addrQueue
		log.Info("received request to scan address ", addr)
		if !more {
			log.Info("changeset queue is closed, exiting")
			break
		}
		if _, found := cache[addr]; found {
			log.Infof("skip address %s, already scanned", addr)
			continue
		}
		cache[addr] = true
		// recursion levels
		currentLevel := 0
		maxLevel := 1
		// now we go through all transactions an we search for:
		processedAddress := make(map[string]bool)
		scan(client, processedAddress, addr, currentLevel, maxLevel)
	}
}

// Scan scan the relationships of a new address
func Scan(cfg config.Schema, address string) (err error) {
	addrQueue <- address
	return
}

// actually process the addresses
func scan(client *EtherscanClient, processedAddress map[string]bool, a string, level, maxLevel int) {
	// don't go too deep
	if level > maxLevel {
		return
	}
	// if already visited don't do it again
	if _, doneAlready := processedAddress[a]; doneAlready {
		return
	}
	processedAddress[a] = true
	// retrieve the address transactions
	txs, err := client.GetTransactions(a)
	if err != nil {
		log.Error("error retrieving transactions: ", err)
		return
	}
	// create the source criteria
	sc, isNew := criteria(a)
	// it's a contract
	if sc.Type == TypeDefiProtocol {
		return
	}
	// create the changeset queue and start the processor
	if isNew {
		sc.Name = a
		csQueue <- NewChangeset(sc)
	}
	// process the relationships
	var src, dst string
	var cs *TrustAPIChangeSet
	for _, y := range txs {
		src, dst = y.From, y.To
		// if source is eq destination skip it
		if src == dst {
			continue
		}
		// the subject address is always the sender
		if src != a {
			src, dst = y.To, y.From
		}
		dc, isNew := criteria(dst)
		cs = NewChangeset()
		if isNew {
			dc.Name = dst
			cs.AddEntity(dc)
		}
		//
		rel := NewTrustRelationship()
		rel.Type = TypeInteraction
		rel.Properties = map[string]interface{}{
			"txId":      y.Hash,
			"action":    "interaction",
			"timestamp": y.GetTime(),
		}
		rel.SourceCriteria = sc // the sender is the source
		rel.TargetCriteria = dc
		cs.AddRel(rel)
		// add to the processed list
		csQueue <- cs
		// recursively call on the destination
		scan(client, processedAddress, dst, level+1, maxLevel)
	}
}
