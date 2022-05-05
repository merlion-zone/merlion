package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/gauge module sentinel errors
var (
	ErrInvalidDepositDenom  = sdkerrors.Register(ModuleName, 2, "pool denom cannot be deposited as reward")
	ErrInvalidAmount        = sdkerrors.Register(ModuleName, 3, "invalid amount")
	ErrTooSmallRewardAmount = sdkerrors.Register(ModuleName, 4, "too small reward amount")
	ErrTooLargeAmount       = sdkerrors.Register(ModuleName, 5, "too large amount")
)
