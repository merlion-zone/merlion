package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	NegOneInt = sdk.OneInt().Neg()
	NegOneDec = sdk.OneDec().Neg()
)
