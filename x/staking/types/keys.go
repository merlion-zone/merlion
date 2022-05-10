package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

var (
	VeValidatorsKey          = []byte{0xA1}
	VeDelegationKey          = []byte{0xA2}
	VeUnbondingDelegationKey = []byte{0xA3}
	VeRedelegationKey        = []byte{0xA4}
	VeTokensKey              = []byte{0xA5}
)

func GetVeValidatorKey(operatorAddr sdk.ValAddress) []byte {
	return append(VeValidatorsKey, address.MustLengthPrefix(operatorAddr)...)
}

func GetVeDelegationKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetVeDelegationsKey(delAddr), address.MustLengthPrefix(valAddr)...)
}

func GetVeDelegationsKey(delAddr sdk.AccAddress) []byte {
	return append(VeDelegationKey, address.MustLengthPrefix(delAddr)...)
}

func GetVeUBDKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetVeUBDsKey(delAddr.Bytes()), address.MustLengthPrefix(valAddr)...)
}

func GetVeUBDsKey(delAddr sdk.AccAddress) []byte {
	return append(VeUnbondingDelegationKey, address.MustLengthPrefix(delAddr)...)
}

func GetVeREDKey(delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress) []byte {
	// key is of the form GetVeREDsKey || valSrcAddrLen (1 byte) || valSrcAddr || valDstAddrLen (1 byte) || valDstAddr
	key := make([]byte, 1+3+len(delAddr)+len(valSrcAddr)+len(valDstAddr))

	copy(key[0:2+len(delAddr)], GetVeREDsKey(delAddr.Bytes()))
	key[2+len(delAddr)] = byte(len(valSrcAddr))
	copy(key[3+len(delAddr):3+len(delAddr)+len(valSrcAddr)], valSrcAddr.Bytes())
	key[3+len(delAddr)+len(valSrcAddr)] = byte(len(valDstAddr))
	copy(key[4+len(delAddr)+len(valSrcAddr):], valDstAddr.Bytes())

	return key
}

func GetVeREDsKey(delAddr sdk.AccAddress) []byte {
	return append(VeRedelegationKey, address.MustLengthPrefix(delAddr)...)
}

func GetVeTokensKey(veID uint64) []byte {
	return append(VeTokensKey, sdk.Uint64ToBigEndian(veID)...)
}
