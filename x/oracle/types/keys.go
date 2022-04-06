package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "oracle"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_oracle"
)

// Prefix keys for oracle module store
var (
	ExchangeRateKey                 = []byte{0x01} // prefix for each key to a rate
	FeederDelegationKey             = []byte{0x02} // prefix for each key to a feeder delegation
	MissCounterKey                  = []byte{0x03} // prefix for each key to a miss counter
	AggregateExchangeRatePrevoteKey = []byte{0x04} // prefix for each key to a aggregate prevote
	AggregateExchangeRateVoteKey    = []byte{0x05} // prefix for each key to a aggregate vote
	TobinTaxKey                     = []byte{0x06} // prefix for each key to a tobin tax
)

// GetExchangeRateKey - stored by *denom*
func GetExchangeRateKey(denom string) []byte {
	return append(ExchangeRateKey, []byte(denom)...)
}

// GetFeederDelegationKey - stored by *Validator* address
func GetFeederDelegationKey(v sdk.ValAddress) []byte {
	return append(FeederDelegationKey, address.MustLengthPrefix(v)...)
}

// GetMissCounterKey - stored by *Validator* address
func GetMissCounterKey(v sdk.ValAddress) []byte {
	return append(MissCounterKey, address.MustLengthPrefix(v)...)
}

// GetAggregateExchangeRatePrevoteKey - stored by *Validator* address
func GetAggregateExchangeRatePrevoteKey(v sdk.ValAddress) []byte {
	return append(AggregateExchangeRatePrevoteKey, address.MustLengthPrefix(v)...)
}

// GetAggregateExchangeRateVoteKey - stored by *Validator* address
func GetAggregateExchangeRateVoteKey(v sdk.ValAddress) []byte {
	return append(AggregateExchangeRateVoteKey, address.MustLengthPrefix(v)...)
}

// GetTobinTaxKey - stored by *denom* bytes
func GetTobinTaxKey(d string) []byte {
	return append(TobinTaxKey, []byte(d)...)
}

// ExtractDenomFromTobinTaxKey - split denom from the tobin tax key
func ExtractDenomFromTobinTaxKey(key []byte) (denom string) {
	denom = string(key[1:])
	return
}
