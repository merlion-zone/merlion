package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name
	ModuleName = "ve"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_ve"

	EmissionPoolName     = "ve_emission_pool"
	DistributionPoolName = "ve_distribution_pool"
)

const (
	prefixTotalLockedAmount = iota + 1
	prefixLockedAmountByUser
	prefixNextVeID
	prefixEpoch
	prefixPointHistoryByEpoch
	prefixUserEpoch
	prefixUserPointHistoryByUserEpoch
	prefixSlopeChange

	prefixAttached
	prefixVoted

	prefixTotalEmission
	prefixEmissionAtLastPeriod
	prefixEmissionLastTimestamp

	prefixDistributionAccruedLastTimestamp
	prefixDistributionTotalAmount
	prefixDistributionPerPeriod
	prefixDistributionClaimLastTimestampByUser
)

var (
	KeyPrefixTotalLockedAmount           = []byte{prefixTotalLockedAmount}
	KeyPrefixLockedAmountByUser          = []byte{prefixLockedAmountByUser}
	KeyPrefixNextVeID                    = []byte{prefixNextVeID}
	KeyPrefixEpoch                       = []byte{prefixEpoch}
	KeyPrefixPointHistoryByEpoch         = []byte{prefixPointHistoryByEpoch}
	KeyPrefixUserEpoch                   = []byte{prefixUserEpoch}
	KeyPrefixUserPointHistoryByUserEpoch = []byte{prefixUserPointHistoryByUserEpoch}
	KeyPrefixSlopeChange                 = []byte{prefixSlopeChange}

	KeyPrefixAttached = []byte{prefixAttached}
	KeyPrefixVoted    = []byte{prefixVoted}

	KeyPrefixTotalEmission         = []byte{prefixTotalEmission}
	KeyPrefixEmissionAtLastPeriod  = []byte{prefixEmissionAtLastPeriod}
	KeyPrefixEmissionLastTimestamp = []byte{prefixEmissionLastTimestamp}

	KeyPrefixDistributionAccruedLastTimestamp     = []byte{prefixDistributionAccruedLastTimestamp}
	KeyPrefixDistributionTotalAmount              = []byte{prefixDistributionTotalAmount}
	KeyPrefixDistributionPerPeriod                = []byte{prefixDistributionPerPeriod}
	KeyPrefixDistributionClaimLastTimestampByUser = []byte{prefixDistributionClaimLastTimestampByUser}
)

func TotalLockedAmountKey() []byte {
	return KeyPrefixTotalLockedAmount
}

func LockedAmountByUserKey(veID uint64) []byte {
	return append(KeyPrefixLockedAmountByUser, sdk.Uint64ToBigEndian(veID)...)
}

func NextVeIDKey() []byte {
	return KeyPrefixNextVeID
}

func EpochKey() []byte {
	return KeyPrefixEpoch
}

func PointKey(epoch uint64) []byte {
	return append(KeyPrefixPointHistoryByEpoch, sdk.Uint64ToBigEndian(epoch)...)
}

func UserEpochKey(veID uint64) []byte {
	return append(KeyPrefixUserEpoch, sdk.Uint64ToBigEndian(veID)...)
}

func UserPointKey(veID uint64, userEpoch uint64) []byte {
	return append(append(KeyPrefixUserPointHistoryByUserEpoch, sdk.Uint64ToBigEndian(veID)...), sdk.Uint64ToBigEndian(userEpoch)...)
}

func SlopeChangeKey(timestamp uint64) []byte {
	return append(KeyPrefixSlopeChange, sdk.Uint64ToBigEndian(timestamp)...)
}

func AttachedKey(veID uint64) []byte {
	return append(KeyPrefixAttached, sdk.Uint64ToBigEndian(veID)...)
}

func VotedKey(veID uint64) []byte {
	return append(KeyPrefixVoted, sdk.Uint64ToBigEndian(veID)...)
}

func TotalEmissionKey() []byte {
	return KeyPrefixTotalEmission
}

func EmissionAtLastPeriodKey() []byte {
	return KeyPrefixEmissionAtLastPeriod
}

func EmissionLastTimestampKey() []byte {
	return KeyPrefixEmissionLastTimestamp
}

func DistributionAccruedLastTimestampKey() []byte {
	return KeyPrefixDistributionAccruedLastTimestamp
}

func DistributionTotalAmountKey() []byte {
	return KeyPrefixDistributionTotalAmount
}

func DistributionPerPeriodKey(timestamp uint64) []byte {
	return append(KeyPrefixDistributionPerPeriod, sdk.Uint64ToBigEndian(timestamp)...)
}

func DistributionClaimLastTimestampByUserKey(veID uint64) []byte {
	return append(KeyPrefixDistributionClaimLastTimestampByUser, sdk.Uint64ToBigEndian(veID)...)
}
