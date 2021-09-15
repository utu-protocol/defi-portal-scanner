package ocean

import (
	"fmt"
	"log"

	"github.com/utu-crowdsale/defi-portal-scanner/collector"
)

// PostAssetsToUTU works like this: Post Asset, then Pool, then Datatoken, then
// post the relationships (Asset owns Pool and Datatokens) between them all
func PostAssetsToUTU(assets []*Asset, u *collector.UTUClient, log *log.Logger) (assetTes []*collector.TrustEntity, poolTes []*collector.TrustEntity, datatokenTes []*collector.TrustEntity) {
	for _, asset := range assets {
		assetTe := asset.toTrustEntity()
		err := u.PostEntity(assetTe)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("%s posted to UTU\n", asset.Identifier())
		}
		assetTes = append(assetTes, assetTe)

		poolTe := asset.Pool.toTrustEntity()
		err = u.PostEntity(poolTe)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("%s posted to UTU\n", asset.Pool.Identifier())
		}
		poolTes = append(poolTes, poolTe)

		datatokenTe := asset.Datatoken.toTrustEntity()
		err = u.PostEntity(datatokenTe)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("%s posted to UTU\n", asset.Datatoken.Identifier())
		}
		datatokenTes = append(datatokenTes, datatokenTe)

		assetPoolRelationship := asset.poolToTrustRelationship()
		err = u.PostRelationship(assetPoolRelationship)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Relationship between %s and %s posted to UTU\n", asset.Identifier(), asset.Pool.Identifier())
		}

		assetDatatokenRelationship := asset.datatokenToTrustRelationship()
		err = u.PostRelationship(assetDatatokenRelationship)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Relationship between %s and %s posted to UTU\n", asset.Identifier(), asset.Datatoken.Identifier())
		}
	}
	return
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
		poolsMap[x.Pool.Address] = x.Pool.toTrustEntity()
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
