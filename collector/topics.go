package collector

var eventNames map[string]string

func init() {
	// action names
	eventNames = map[string]string{
		"0x623b3804fa71d67900d064613da8f94b9617215ee90799290593e1745087ad18": "TokenPurchase",
		"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef": "Transfer",
		"0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822": "Swap",
		"0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496": "Burn",
		"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925": "Approval",
		"0x875352fb3fadeb8c0be7cbbe8ff761b308fa7033470cd0287f02f3436fd76cb9": "AccrueInterest",
		"0x4dec04e750ca11537cabcd8a9eab06494de08da3735bc8871cd41250e190bc04": "AccrueInterest",
		"0xe5b754fb1abb7f01b499791d0b820ae3b6af3424ac1c59768edb53f4ec31a929": "Redeem",
		"0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f": "Mint",
		"0x3c67396e9c55d2fc8ad68875fc5beca1d96ad2a2f23b210ccc1d986551ab6fdf": "TokensTransferred",
		"0x1a2a22cb034d26d1854bdc6666a5b91fe25efbbb5dcad3b0355478d6f5c362a1": "RepayBorrow",
		"0x45b96fe442630264581b197e84bbada861235052c5a1aadfff9ea4e40a969aa0": "Failure",
		"0xbd5034ffbd47e4e72a94baa2cdb74c6fad73cb3bcdc13036b72ec8306f5a7646": "Redeem",
		"0x5e3cad45b1fe24159d1cb39788d82d0f69cc15770aa96fba1d3d1a7348735594": "InterestStreamRedirected",
		"0x9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d": "Staked",
		"0x34fcbac0073d7c3d388e51312faf357774904998eeb8fca628b9e6f65ee1cbf7": "Claim",
		"0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885": "Mint",
		"0xdec2bacdd2f05b59de34da9b523dff8be42e5e38e818c82fdb0bae774387a724": "DelegateVotesChanged", // Compound
	}

}
