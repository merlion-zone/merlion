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

	EmissionPoolName = "ve_emission_pool"
)

const (
	prefixTotalLockedAmount = iota + 1
	prefixLockedAmountByUser
	prefixNextVeNftID
	prefixEpoch
	prefixPointHistoryByEpoch
	prefixUserEpoch
	prefixUserPointHistoryByUserEpoch
	prefixSlopeChange

	prefixTotalEmission
	prefixEmissionAtLastPeriod
	prefixEmissionLastTimestamp
)

var (
	KeyPrefixTotalLockedAmount           = []byte{prefixTotalLockedAmount}
	KeyPrefixLockedAmountByUser          = []byte{prefixLockedAmountByUser}
	KeyPrefixNextVeNftID                 = []byte{prefixNextVeNftID}
	KeyPrefixEpoch                       = []byte{prefixEpoch}
	KeyPrefixPointHistoryByEpoch         = []byte{prefixPointHistoryByEpoch}
	KeyPrefixUserEpoch                   = []byte{prefixUserEpoch}
	KeyPrefixUserPointHistoryByUserEpoch = []byte{prefixUserPointHistoryByUserEpoch}
	KeyPrefixSlopeChange                 = []byte{prefixSlopeChange}

	KeyPrefixTotalEmission         = []byte{prefixTotalEmission}
	KeyPrefixEmissionAtLastPeriod  = []byte{prefixEmissionAtLastPeriod}
	KeyPrefixEmissionLastTimestamp = []byte{prefixEmissionLastTimestamp}
)

func TotalLockedAmountKey() []byte {
	return KeyPrefixTotalLockedAmount
}

func LockedAmountByUserKey(veID uint64) []byte {
	return append(KeyPrefixLockedAmountByUser, sdk.Uint64ToBigEndian(veID)...)
}

func NextVeNftIDKey() []byte {
	return KeyPrefixNextVeNftID
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

func TotalEmissionKey() []byte {
	return KeyPrefixTotalEmission
}

func EmissionAtLastPeriodKey() []byte {
	return KeyPrefixEmissionAtLastPeriod
}

func EmissionLastTimestampKey() []byte {
	return KeyPrefixEmissionLastTimestamp
}
