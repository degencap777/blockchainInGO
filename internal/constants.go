package internal

const (
	// Blockchain version
	Version = byte(0x00)

	// Wallet address checksum length
	AddressChecksumLen = 4

	// Wallet filename
	WalletFile = "wallets.dat"

	// Rewards to mining a new block
	Subsidy = 10

	// Blockchain database filename
	DbFile = "blockchain_%s.db"

	// Blockchain bucket for blocks
	BlockBucket = "blocks"

	// Blockchain bucket for utxo sets
	UtxoBucket = "chainstate"

	// Genesis data for genesis block
	GenesisData = "Bitcoin was created in 2009 by a person or group of people using the pseudonym Satoshi Nakamoto"
)
