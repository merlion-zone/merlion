package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/exp/constraints"
)

var (
	ZeroInt   = sdk.ZeroInt()
	ZeroDec   = sdk.ZeroDec()
	OneInt    = sdk.OneInt()
	OneDec    = sdk.OneDec()
	NegOneInt = sdk.OneInt().Neg()
	NegOneDec = sdk.OneDec().Neg()
)

func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}
