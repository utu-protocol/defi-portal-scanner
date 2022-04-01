package ocean

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestDecentralizedDataObjectGetNameDescription(t *testing.T) {
	data := []byte(`{"@context":"https://w3id.org/did/v1","id":"did:op:7Bce67697eD2858d0683c631DdE7Af823b7eea38","publicKey":[{"id":"did:op:7Bce67697eD2858d0683c631DdE7Af823b7eea38","type":"EthereumECDSAKey","owner":"0x655eFe6Eb2021b8CEfE22794d90293aeC37bb325"}],"authentication":[{"type":"RsaSignatureAuthentication2018","publicKey":"did:op:7Bce67697eD2858d0683c631DdE7Af823b7eea38"}],"service":[{"type":"metadata","attributes":{"curation":{"rating":0,"numVotes":0,"isListed":true},"main":{"type":"dataset","name":"🖼 DataUnion.app - Image & Annotation Vault 📸","dateCreated":"2020-11-15T12:27:05Z","author":"DataUnion.app","license":"https://market.oceanprotocol.com/terms","files":[{"contentLength":"61","contentType":"text/plain","index":0}],"datePublished":"2021-09-03T02:00:17Z"},"additionalInformation":{"description":"# Notice\n\nThis dataset and the software stack behind it are under constant development.\n[Check our Twitter for updates.](https://twitter.com/DataunionA)\nThis version is in alpha so staking is possible and contributions are now possible in the [alpha platform](https://alpha.dataunion.app).\nThe initial data contained in the dataset is open source data but we are working on the Data Portal, the window onto the data.\nIt has a long list of features which will be rolled out over time:\n* buying data assets via compute-to-data\n* initial algorithm offerings\n* data bounties\n* data organisation\n* inspection\n\n# Description\n\nThis dataset represents a collection of images and annotations as well as the software stack attached to it. The data is uploaded, annotated and verified by the dataset co-owners in exchange for datatoken rewards. The ownership will be similar to NFTs but one finished image and its annotation will be owned partially by all its contributors.\nParticipants invest their time to co-own the dataset and improve its quality.\nThe goal of this dataset is to be used for AI training via the Data Portal and compute-to-data.\n\nDataUnion.app, customers and the community will run Data Bounties that are looking for specific data and give extra rewards. The reward distributions will be in intervals and not directly for now but these distributions will be announced on our [Twitter](https://twitter.com/DataunionA).\n\nThe long term goal of this dataset is to contain billions of crowdsource verified and annotated images. DataUnion.app will also facilitate the training of algorithms. These algorithms will be available via the Data Portal and the algorithm contributors will become co-creators and co-owners of the algorithm. We are aiming to create the option of Universal Data Income. Especially for countries with low average salaries compared to e.g. Europe, at the moment we are looking for [interested communities to connect to](mailto:collaborate@dataunion.app).\n\nWe are working on several algorithms via data bounties:\n* Anonymization algorithm for e.g. the Data Portal\n* Traffic sign detection algorithm with [Evotegra](https://evotegra.de)\n* Litter detection algorithm with [Project.BB](https://project.bb)\n\n# Tokenomics\n\nTo enable rewards datatokens are bought from this dataset with $OCEAN and handed out as rewards to the contributors. This enables a growth circle of the data value in the dataset while keeping the supply of the datatoken the same. So the value of the dataset grows in the long term and early adopter as well as early contributors are rewarded more leading to the incentivised growth that is needed to scale the dataset. In the longer term we want to migrate to a DAO in the style of the OceanDAO to facilitate longest term growth. This will potentially also involve a governance token - we are fans of fair launches to further strengthen our community.\nThere also will be challenges and Data Bounties that reward for the creation of certain information, some have been started already.\n\n# Vision\n\nToday large corporations dominate the AI algorithm market. We want to change that by giving everyone the option to participate and become a part of this evergrowing market. The aim is to create the most diverse and comprehensive datasets to create the most acurate and versatile algorithms. We believe a global community of participants will be able to outpace companies without a problem. At the moment we are the product in their data generation machineries but now it is time to take action and use our data for our own profit.\n\n# Roadmap\n\nWe are planning our roadmap in many differnt directions:\n* New annotation tools and mechanisms\n* Release a mobile application that the mechanisms of our [web app](https://alpha.dataunion.app) and Swipe-AI\n* Release the Data Portal\n* Internationalize the platforms\n* Simulate the token value flow using simulation tools\n* Include NLP to translate annotations and cater to a worldwide audience\n* Increase decentralisation of our solution by moving the control to smart contracts\n* Allow addition of data while it resides on different storage\n* Algorithm training and sales via our Data Portal - have the data providers become co-owners of that as well\n* Recruit more people and facilitate development via OceanDAO grants\n* Potentially launch a governance token\n* Add other data vaults e.g. text, sound, and 3D objects\n\n# Community\n\nPlease join our community on [Telegram](https://t.me/dataunionapp), [Discord](https://discord.com/invite/Jm9C3yD8Sd) and follow us on [Twitter](https://twitter.com/DataunionA).","tags":["machine-learning","ai","computer-vision","image","object-recognition","crowdsourcing","growing-dataset","continuous","universal-data-income","social","worldwide","user-driven"],"links":[{"contentLength":"69","contentType":"text/plain","url":"https://dataunion.app/sample_data.txt"}],"termsAndConditions":true},"encryptedFiles":"0x0456fb2b0190ccd575361c7092c4a1ec9ff0194f8bd466519293529b7ffb9047020d585e1c32178318d224a0bf79fbd2d418bae6a9c42ee0efe03f9bc3c6e5993dc74f6ae70f59eb9f74c51a39401044ff0dd4d8f1d951c490e236398811b5834d41de7d1b659997afc1edbef8a4e5f3d348c78fc08b7d10e36cbb835078e45b5df7d9cb2b04c8ac378a1ec39b713cfbf29073265fdfe7c192192f6ac9379da0136faba25a9f69dbe8f61c4aabd7a09e42be63e6c1e16cffb8a6f62c8d82"},"index":0},{"type":"access","index":1,"serviceEndpoint":"https://provider.mainnet.oceanprotocol.com","attributes":{"main":{"creator":"0x655eFe6Eb2021b8CEfE22794d90293aeC37bb325","datePublished":"2020-11-15T12:27:05Z","cost":"1","timeout":0,"name":"dataAssetAccess"}}}],"dataToken":"0x7Bce67697eD2858d0683c631DdE7Af823b7eea38","created":"2020-11-15T12:27:48Z","proof":{"created":"2020-11-15T12:27:38Z","creator":"0x655eFe6Eb2021b8CEfE22794d90293aeC37bb325","type":"AddressHash","signatureValue":"0xe3484e8e5201c97ae747423b55f078fa31e7a5f955c85e91496bdc62f4992263"},"dataTokenInfo":{"address":"0x7Bce67697eD2858d0683c631DdE7Af823b7eea38","name":"Quiescent Crab Token","symbol":"QUICRA-0","decimals":18,"cap":1000},"updated":"2021-05-17T21:58:02Z","accessWhiteList":[],"price":{"datatoken":458.8974964665521,"ocean":570154.1494885942,"value":546.9805175263882,"type":"pool","exchange_id":"","address":"0xAAB9EaBa1AA2653c1Dda9846334700b9F5e14E44","pools":["0xAAB9EaBa1AA2653c1Dda9846334700b9F5e14E44","0xa81798FEc662C1A131029B28546Cc80ECd11CB61"],"isConsumable":"true"},"isInPurgatory":"false","event":{"txid":"0x4360656996dad6626e088f45969a84c7e46dbfb1e5a8482c9f933ef08d4fabef","blockNo":12454477,"from":"0x655eFe6Eb2021b8CEfE22794d90293aeC37bb325","contract":"0x1a4b70d8c9DcA47cD6D0Fb3c52BB8634CA1C0Fdf","update":true},"chainId":1}`)
	ddo := new(DecentralizedDataObject)
	err := json.Unmarshal(data, ddo)
	if err != nil {
		t.Fatal(err)
	}
	name, author, description, tags, categories, err := ddo.GetNameAuthorMetadata()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("name: %s\nauthor: %s\ndescription: %s\ntags: %s\ncategories: %s\n", name, author, description, strings.Join(tags, ", "), strings.Join(categories, ", "))
}
