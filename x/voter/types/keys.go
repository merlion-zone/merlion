package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name
	ModuleName = "voter"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_voter"
)

const (
	prefixTotalVotes = iota + 1
	prefixTotalVotesByUser
	prefixPoolWeightedVotes
	prefixPoolWeightedVotesByUser
	prefixIndex
	prefixIndexAtLastUpdatedByGauge
	prefixClaimableRewardByGauge
)

var (
	KeyPrefixTotalVotes                = []byte{prefixTotalVotes}
	KeyPrefixTotalVotesByUser          = []byte{prefixTotalVotesByUser}
	KeyPrefixPoolWeightedVotes         = []byte{prefixPoolWeightedVotes}
	KeyPrefixPoolWeightedVotesByUser   = []byte{prefixPoolWeightedVotesByUser}
	KeyPrefixIndex                     = []byte{prefixIndex}
	KeyPrefixIndexAtLastUpdatedByGauge = []byte{prefixIndexAtLastUpdatedByGauge}
	KeyPrefixClaimableRewardByGauge    = []byte{prefixClaimableRewardByGauge}
)

func TotalVotesKey() []byte {
	return KeyPrefixTotalVotes
}

func TotalVotesByUserKey(veID uint64) []byte {
	return append(KeyPrefixTotalVotesByUser, sdk.Uint64ToBigEndian(veID)...)
}

func PoolWeightedVotesKey(poolDenom string) []byte {
	return append(KeyPrefixPoolWeightedVotes, poolDenom...)
}

func PoolWeightedVotesByUserKey(veID uint64, poolDenom string) []byte {
	return append(append(KeyPrefixPoolWeightedVotes, sdk.Uint64ToBigEndian(veID)...), poolDenom...)
}

func IndexKey() []byte {
	return KeyPrefixIndex
}

func IndexAtLastUpdatedByGaugeKey(poolDenom string) []byte {
	return append(KeyPrefixIndexAtLastUpdatedByGauge, poolDenom...)
}

func ClaimableRewardByGaugeKey(poolDenom string) []byte {
	return append(KeyPrefixClaimableRewardByGauge, poolDenom...)
}
