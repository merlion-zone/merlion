package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewLockedBalance() LockedBalance {
	return LockedBalance{
		Amount: sdk.ZeroInt(),
		End:    0,
	}
}
