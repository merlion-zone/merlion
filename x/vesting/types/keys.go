package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "vesting"

	// StoreKey defines the primary module store key
	// Here use "vs" to avoid potential key collision with the ve module
	StoreKey = "vs"

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

const (
	prefixAllocationAddress = iota + 1
	prefixAirdropsTotalAmount
	prefixAirdrops
	prefixAirdropsCompleted
)

var (
	KeyPrefixAllocationAddress   = []byte{prefixAllocationAddress}
	KeyPrefixAirdropsTotalAmount = []byte{prefixAirdropsTotalAmount}
	KeyPrefixAirdrops            = []byte{prefixAirdrops}
	KeyPrefixAirdropsCompleted   = []byte{prefixAirdropsCompleted}
)

func AllocationAddrKey() []byte {
	return KeyPrefixAllocationAddress
}

func AirdropsTotalAmountKey() []byte {
	return KeyPrefixAirdropsTotalAmount
}

func AirdropsKey(acc sdk.AccAddress) []byte {
	return append(KeyPrefixAirdrops, address.MustLengthPrefix(acc)...)
}

func AirdropsCompletedKey(acc sdk.AccAddress) []byte {
	return append(KeyPrefixAirdropsCompleted, address.MustLengthPrefix(acc)...)
}
