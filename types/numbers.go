package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	ZeroInt   = sdk.ZeroInt()
	ZeroDec   = sdk.ZeroDec()
	OneInt    = sdk.OneInt()
	OneDec    = sdk.OneDec()
	NegOneInt = sdk.OneInt().Neg()
	NegOneDec = sdk.OneDec().Neg()
)
