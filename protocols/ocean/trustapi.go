package ocean

import (
	"fmt"
	"log"

	"github.com/remeh/sizedwaitgroup"
	"github.com/utu-crowdsale/defi-portal-scanner/collector"
)

// PostAssetsToUTU works like this: Post Asset, then Pool, then Datatoken, then
// post the relationships (Asset owns Pool and Datatokens) between them all
func PostAssetsToUTU(assets []*Asset, u *collector.UTUClient, log *log.Logger) {
	wg := sizedwaitgroup.New(20)
	for _, asset := range assets {
		wg.Add()
		go postAsset(asset, u, log, &wg)
	}
	wg.Wait()
}

func postAsset(asset *Asset, u *collector.UTUClient, log *log.Logger, wg *sizedwaitgroup.SizedWaitGroup) {
	defer wg.Done()
	assetTe := asset.toTrustEntity()
	err := u.PostEntity(assetTe)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%s posted to UTU\n", asset.Identifier())
	}

	datatokenTe := asset.Datatoken.toTrustEntity()
	err = u.PostEntity(datatokenTe)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%s posted to UTU\n", asset.Datatoken.Identifier())
	}

	assetDatatokenRelationship := asset.datatokenToTrustRelationship()
	err = u.PostRelationship(assetDatatokenRelationship)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Relationship between %s and %s posted to UTU\n", asset.Identifier(), asset.Datatoken.Identifier())
	}
}

func PostAddressesToUTU(addresses []*Address, assets []*Asset, u *collector.UTUClient, log *log.Logger) {
	var addressesMap = make(map[string]*collector.TrustEntity)
	var assetsMap = make(map[string]*collector.TrustEntity)
	var datatokensMap = make(map[string]*collector.TrustEntity)

	// I need to be able to look up things by their addresses later, so I
	// transform things into a map
	for _, x := range assets {
		assetsMap[x.Datatoken.Address] = x.toTrustEntity()
		datatokensMap[x.Datatoken.Address] = x.Datatoken.toTrustEntity()
	}

	// This particular piece of code used to be in the postUser() block, but it
	// was moved out here so that concurrent writes to this map (which is shared
	// between goroutines) don't happen.
	for _, address := range addresses {
		// Convert users to UTU Trust Entities
		userTe := address.toTrustEntity()
		addressesMap[address.Address] = userTe
	}

	// OKAY now we can start parallelized POSTING to UTU Trust API. Because we
	// only read from the maps, not write to them, the code doesn't have to be
	// rewritten so much.
	wg := sizedwaitgroup.New(20)
	for _, address := range addresses {
		wg.Add()
		go postAddress(address, addressesMap, datatokensMap, u, log, &wg)
	}
	wg.Wait()

}

func postAddress(address *Address, usersMap, datatokensMap map[string]*collector.TrustEntity, u *collector.UTUClient, log *log.Logger, wg *sizedwaitgroup.SizedWaitGroup) {
	defer wg.Done()

	// Now we can create the relationships between the Users and the
	// Datatokens.
	dtiTes := address.datatokenInteractionsToTrustRelationships(datatokensMap, log)

	if len(dtiTes) > 0 {
		fmt.Println("datatokenRelationships", len(dtiTes))
	}

	// POST Address to UTU API
	err := u.PostEntity(usersMap[address.Address])
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%s posted to UTU\n", address.Identifier())
	}

	// POST Relationships to UTU API
	for _, r := range dtiTes {
		err := u.PostRelationship(r)
		if err != nil {
			log.Println(err)
		}
	}
}
