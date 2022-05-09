package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
)

const (
	TypeMsgVeDelegate = "ve_delegate"
)

// Route implements the sdk.Msg interface.
func (m MsgVeDelegate) Route() string { return stakingtypes.RouterKey }

// Type implements the sdk.Msg interface.
func (m MsgVeDelegate) Type() string { return TypeMsgVeDelegate }

// GetSigners implements the sdk.Msg interface.
func (m MsgVeDelegate) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delAddr}
}

// GetSignBytes implements the sdk.Msg interface.
func (m MsgVeDelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// ValidateBasic implements the sdk.Msg interface.
func (m MsgVeDelegate) ValidateBasic() error {
	if m.DelegatorAddress == "" {
		return stakingtypes.ErrEmptyDelegatorAddr
	}

	if m.ValidatorAddress == "" {
		return stakingtypes.ErrEmptyValidatorAddr
	}

	if vetypes.Uint64FromVeID(m.VeId) == vetypes.EmptyVeID {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ve id")
	}

	if !m.Amount.IsValid() || !m.Amount.Amount.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid delegation amount",
		)
	}

	return nil
}
