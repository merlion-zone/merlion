package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

var (
	VeValidatorsKey = []byte{0xA1}
	VeDelegationKey = []byte{0xA2}
)

func GetVeValidatorKey(operatorAddr sdk.ValAddress) []byte {
	return append(VeValidatorsKey, address.MustLengthPrefix(operatorAddr)...)
}

func GetDelegationKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetDelegationsKey(delAddr), address.MustLengthPrefix(valAddr)...)
}

func GetDelegationsKey(delAddr sdk.AccAddress) []byte {
	return append(VeDelegationKey, address.MustLengthPrefix(delAddr)...)
}
