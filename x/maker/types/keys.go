package types

const (
	// ModuleName defines the module name
	ModuleName = "maker"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_maker"
)

const (
	prefixBackingRatio = iota + 1
	prefixBackingRatioLastBlock
	prefixBackingParams
	prefixCollateralParams
	prefixBackingTotal
	prefixCollateralTotal
	prefixBackingPool
	prefixCollateralPool
	prefixBackingAccount
	prefixCollateralAccount
)

var (
	KeyPrefixBackingRatio          = []byte{prefixBackingRatio}
	KeyPrefixBackingRatioLastBlock = []byte{prefixBackingRatioLastBlock}
	KeyPrefixBackingParams         = []byte{prefixBackingParams}
	KeyPrefixCollateralParams      = []byte{prefixCollateralParams}
	KeyPrefixBackingTotal          = []byte{prefixBackingTotal}
	KeyPrefixCollateralTotal       = []byte{prefixCollateralTotal}
	KeyPrefixBackingPool           = []byte{prefixBackingPool}
	KeyPrefixCollateralPool        = []byte{prefixCollateralPool}
	KeyPrefixBackingAccount        = []byte{prefixBackingAccount}
	KeyPrefixCollateralAccount     = []byte{prefixCollateralAccount}
)
