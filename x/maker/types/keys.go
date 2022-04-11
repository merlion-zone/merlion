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
	prefixCollateralRatio = iota + 1
	prefixBacking
	prefixCollateral
)

var (
	KeyPrefixCollateralRatio = []byte{prefixCollateralRatio}
	KeyPrefixBacking         = []byte{prefixBacking}
	KeyPrefixCollateral      = []byte{prefixCollateral}
)
