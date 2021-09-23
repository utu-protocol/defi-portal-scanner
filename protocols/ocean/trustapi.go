package ocean

import (
	"fmt"
	"log"

	"github.com/utu-crowdsale/defi-portal-scanner/collector"
)

// PostAssetsToUTU works like this: Post Asset, then Pool, then Datatoken, then
// post the relationships (Asset owns Pool and Datatokens) between them all
func PostAssetsToUTU(assets []*Asset, u *collector.UTUClient, log *log.Logger) {
	for _, asset := range assets {
		assetTe := asset.toTrustEntity()
		err := u.PostEntity(assetTe)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("%s posted to UTU\n", asset.Identifier())
		}

		poolTes := asset.poolsToTrustEntities()
		for _, poolTe := range poolTes {
			err = u.PostEntity(poolTe)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("%s posted to UTU\n", poolTe.Name)
			}
		}

		datatokenTe := asset.Datatoken.toTrustEntity()
		err = u.PostEntity(datatokenTe)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("%s posted to UTU\n", asset.Datatoken.Identifier())
		}

		assetPoolRelationships := asset.poolsToTrustRelationships()
		for _, pr := range assetPoolRelationships {
			err = u.PostRelationship(pr)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Relationship between %s and %s (type: %s) posted to UTU\n", pr.SourceCriteria.Name, pr.TargetCriteria.Name, pr.Type)
			}
		}

		assetDatatokenRelationship := asset.datatokenToTrustRelationship()
		err = u.PostRelationship(assetDatatokenRelationship)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Relationship between %s and %s posted to UTU\n", asset.Identifier(), asset.Datatoken.Identifier())
		}
	}
}

func PostToUTU(users []*User, assets []*Asset, u *collector.UTUClient, log *log.Logger) {
	var usersMap = make(map[string]*collector.TrustEntity)
	var assetsMap = make(map[string]*collector.TrustEntity)
	var datatokensMap = make(map[string]*collector.TrustEntity)
	var poolsMap = make(map[string]*collector.TrustEntity)

	// I need to be able to look up things by their addresses later, so I
	// transform things into a map
	for _, x := range assets {
		assetsMap[x.Datatoken.Address] = x.toTrustEntity()
		datatokensMap[x.Datatoken.Address] = x.Datatoken.toTrustEntity()

		// For the pools, for each user, we need to check if we know about any
		// of the Pools that he has interacted with (through Datatokens/Assets),
		// and if so, get a TrustEntity for that pool. This means that poolsMap
		// is flat - it contains no information about which Asset holds the
		// Pool.
		for _, poolTe := range x.poolsToTrustEntities() {
			poolsMap[poolTe.Ids["address"]] = poolTe
		}
	}

	// PostAssetsToUTU(assets, u, log)
	// log.Println("Finished posting Assets to UTU")

	for _, user := range users {
		// Convert users to UTU Trust Entities
		userTe := user.toTrustEntity()
		usersMap[user.Address] = userTe

		// Now we can create the relationships between the Users and the
		// Datatokens.
		dtiTes := user.datatokenInteractionsToTrustRelationships(datatokensMap, log)
		poolTes := user.poolInteractionsToTrustRelationships(poolsMap, log)

		if len(dtiTes) > 0 {
			fmt.Println("datatokenRelationships", len(dtiTes))
		}
		if len(poolTes) > 0 {
			fmt.Println("poolRelationships", len(poolTes))
		}

		// POST User to UTU API
		err := u.PostEntity(userTe)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("%s posted to UTU\n", user.Identifier())
		}

		// POST Relationships to UTU API
		for _, r := range dtiTes {
			err := u.PostRelationship(r)
			if err != nil {
				log.Println(err)
			}
		}
		for _, r := range poolTes {
			err := u.PostRelationship(r)
			if err != nil {
				log.Println(err)
			}
		}

	}

}
