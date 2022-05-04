package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name
	ModuleName = "gauge"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_gauge"

	GaugePoolName = ModuleName
	BribePoolName = "bribe"
)

const (
	prefixGaugeDenom = iota + 1
	prefixBribeDenom
	prefixTotalDepositedAmount
	prefixDepositedAmountByUser
	prefixTotalDerivedAmount
	prefixDerivedAmountByUser
	prefixReward
	prefixUserReward
	prefixUserVeIDByAddress
	prefixEpoch
	prefixPointHistoryByEpoch
	prefixUserEpoch
	prefixUserPointHistoryByUserEpoch
	prefixRewardEpoch
	prefixRewardPointHistoryByRewardEpoch
)

var (
	KeyPrefixGaugeDenom = []byte{prefixGaugeDenom}
	KeyPrefixBribeDenom = []byte{prefixBribeDenom}

	KeyPrefixTotalDepositedAmount  = []byte{prefixTotalDepositedAmount}
	KeyPrefixDepositedAmountByUser = []byte{prefixDepositedAmountByUser}

	KeyPrefixTotalDerivedAmount  = []byte{prefixTotalDerivedAmount}
	KeyPrefixDerivedAmountByUser = []byte{prefixDerivedAmountByUser}

	KeyPrefixReward     = []byte{prefixReward}
	KeyPrefixUserReward = []byte{prefixUserReward}

	KeyPrefixUserVeIDByAddress = []byte{prefixUserVeIDByAddress}

	KeyPrefixEpoch               = []byte{prefixEpoch}
	KeyPrefixPointHistoryByEpoch = []byte{prefixPointHistoryByEpoch}

	KeyPrefixUserEpoch                   = []byte{prefixUserEpoch}
	KeyPrefixUserPointHistoryByUserEpoch = []byte{prefixUserPointHistoryByUserEpoch}

	KeyPrefixRewardEpoch                     = []byte{prefixRewardEpoch}
	KeyPrefixRewardPointHistoryByRewardEpoch = []byte{prefixRewardPointHistoryByRewardEpoch}
)

func GaugeKey(denom string) []byte {
	return append(KeyPrefixGaugeDenom, denom...)
}

func BribeKey(denom string) []byte {
	return append(KeyPrefixBribeDenom, denom...)
}

func TotalDepositedAmountKey(gaugeOrBribe []byte) []byte {
	return append(KeyPrefixTotalDepositedAmount, gaugeOrBribe...)
}

func DepositedAmountByUserKey(gaugeOrBribe []byte, veID uint64) []byte {
	prefix := append(KeyPrefixDepositedAmountByUser, gaugeOrBribe...)
	return append(prefix, sdk.Uint64ToBigEndian(veID)...)
}

func TotalDerivedAmountKey(gaugeKey []byte) []byte {
	return append(KeyPrefixTotalDerivedAmount, gaugeKey...)
}

func DerivedAmountByUserKey(gaugeKey []byte, veID uint64) []byte {
	prefix := append(KeyPrefixDerivedAmountByUser, gaugeKey...)
	return append(prefix, sdk.Uint64ToBigEndian(veID)...)
}

func RewardKeyPrefix(gaugeOrBribe []byte) []byte {
	return append(KeyPrefixReward, gaugeOrBribe...)
}

func RewardKey(gaugeOrBribe []byte, rewardDenom string) []byte {
	return append(RewardKeyPrefix(gaugeOrBribe), rewardDenom...)
}

func UserRewardKey(gaugeOrBribe []byte, rewardDenom string, veID uint64) []byte {
	prefix := append(KeyPrefixUserReward, gaugeOrBribe...)
	prefix = append(prefix, rewardDenom...)
	return append(prefix, sdk.Uint64ToBigEndian(veID)...)
}

func UserVeIDByAddressKey(gaugeKey []byte, acc sdk.AccAddress) []byte {
	prefix := append(KeyPrefixUserVeIDByAddress, gaugeKey...)
	return append(prefix, acc...)
}

func EpochKey(gaugeOrBribe []byte) []byte {
	return append(KeyPrefixEpoch, gaugeOrBribe...)
}

func PointKey(gaugeOrBribe []byte, epoch uint64) []byte {
	prefix := append(KeyPrefixPointHistoryByEpoch, gaugeOrBribe...)
	return append(prefix, sdk.Uint64ToBigEndian(epoch)...)
}

func UserEpochKey(gaugeOrBribe []byte, veID uint64) []byte {
	return append(append(KeyPrefixUserEpoch, gaugeOrBribe...), sdk.Uint64ToBigEndian(veID)...)
}

func UserPointKey(gaugeOrBribe []byte, veID uint64, epoch uint64) []byte {
	prefix := append(KeyPrefixUserPointHistoryByUserEpoch, gaugeOrBribe...)
	return append(append(prefix, sdk.Uint64ToBigEndian(veID)...), sdk.Uint64ToBigEndian(epoch)...)
}

func RewardEpochKey(gaugeOrBribe []byte, rewardDenom string) []byte {
	return append(append(KeyPrefixRewardEpoch, gaugeOrBribe...), rewardDenom...)
}

func RewardPointKey(gaugeOrBribe []byte, rewardDenom string, epoch uint64) []byte {
	prefix := append(KeyPrefixRewardPointHistoryByRewardEpoch, gaugeOrBribe...)
	return append(append(prefix, rewardDenom...), sdk.Uint64ToBigEndian(epoch)...)
}
