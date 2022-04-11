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
	ErrMerPriceTooLow              = sdkerrors.Register(ModuleName, 6, "mer stablecoin price too low")
	ErrBackingCoinSlippage         = sdkerrors.Register(ModuleName, 7, "backing coin over slippage")
	ErrCollateralCoinSlippage      = sdkerrors.Register(ModuleName, 8, "collateral coin over slippage")
	ErrLionCoinSlippage            = sdkerrors.Register(ModuleName, 9, "lion coin over slippage")
	ErrBackingParamsInvalid        = sdkerrors.Register(ModuleName, 10, "backing params invalid")
	ErrCollateralParamsInvalid     = sdkerrors.Register(ModuleName, 11, "collateral params invalid")
	ErrBackingCoinDisabled         = sdkerrors.Register(ModuleName, 12, "backing coin disabled")
	ErrCollateralCoinDisabled      = sdkerrors.Register(ModuleName, 13, "collateral coin disabled")
)
