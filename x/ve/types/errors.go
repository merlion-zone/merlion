package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/ve module sentinel errors
var (
	ErrNoValidVeID     = sdkerrors.Register(ModuleName, 2, "no valid ve id")
	ErrPastLockTime    = sdkerrors.Register(ModuleName, 3, "cannot lock until time in the past")
	ErrTooLongLockTime = sdkerrors.Register(ModuleName, 4, "too long lock time")
)
