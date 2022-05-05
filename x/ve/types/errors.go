package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/ve module sentinel errors
var (
	ErrInvalidLockDenom     = sdkerrors.Register(ModuleName, 2, "invalid lock denom")
	ErrInvalidVeID          = sdkerrors.Register(ModuleName, 3, "invalid ve id")
	ErrPastLockTime         = sdkerrors.Register(ModuleName, 4, "cannot lock until time in the past")
	ErrTooLongLockTime      = sdkerrors.Register(ModuleName, 5, "too long lock time")
	ErrNotIncreasedLockTime = sdkerrors.Register(ModuleName, 6, "lock time can only be increased")
	ErrLockNotExpired       = sdkerrors.Register(ModuleName, 7, "lock didn't expire")
	ErrLockExpired          = sdkerrors.Register(ModuleName, 8, "lock expired")
	ErrAmountNotPositive    = sdkerrors.Register(ModuleName, 9, "amount must be positive")
	ErrSameVeID             = sdkerrors.Register(ModuleName, 10, "from ve id and to ve id must be different")
	ErrVeAttached           = sdkerrors.Register(ModuleName, 11, "ve owner deposited into gauge or ve voted")
)
