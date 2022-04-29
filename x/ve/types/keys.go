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
)

const (
	prefixTotalLockedAmount = iota + 1
	prefixLockedAmountByUser
	prefixNextVeNftID
	prefixEpoch
	prefixPointHistoryByEpoch
	prefixUserEpochByVeNftID
	prefixUserPointHistoryByUserEpochByVeNftID
	prefixSlopeChange
)

var (
	KeyPrefixTotalLockedAmount                    = []byte{prefixTotalLockedAmount}
	KeyPrefixLockedAmountByUser                   = []byte{prefixLockedAmountByUser}
	KeyPrefixNextVeNftID                          = []byte{prefixNextVeNftID}
	KeyPrefixEpoch                                = []byte{prefixEpoch}
	KeyPrefixPointHistoryByEpoch                  = []byte{prefixPointHistoryByEpoch}
	KeyPrefixUserEpochByVeNftID                   = []byte{prefixUserEpochByVeNftID}
	KeyPrefixUserPointHistoryByUserEpochByVeNftID = []byte{prefixUserPointHistoryByUserEpochByVeNftID}
	KeyPrefixSlopeChange                          = []byte{prefixSlopeChange}
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
	return append(KeyPrefixUserEpochByVeNftID, sdk.Uint64ToBigEndian(veID)...)
}

func UserPointKey(veID uint64, userEpoch uint64) []byte {
	return append(append(KeyPrefixUserPointHistoryByUserEpochByVeNftID, sdk.Uint64ToBigEndian(veID)...), sdk.Uint64ToBigEndian(userEpoch)...)
}

func SlopeChangeKey(timestamp uint64) []byte {
	return append(KeyPrefixSlopeChange, sdk.Uint64ToBigEndian(timestamp)...)
}
