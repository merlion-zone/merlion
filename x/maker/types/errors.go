package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/maker module sentinel errors
var (
	ErrBackingCoinAlreadyExists    = sdkerrors.Register(ModuleName, 2, "backing coin already exists")
	ErrCollateralCoinAlreadyExists = sdkerrors.Register(ModuleName, 3, "collateral coin already exists")
	ErrBackingCoinNotFound         = sdkerrors.Register(ModuleName, 4, "backing coin not found")
	ErrCollateralCoinNotFound      = sdkerrors.Register(ModuleName, 5, "collateral coin not found")
)
