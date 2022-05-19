package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgAddAirdrops          = "add_airdrops"
	TypeMsgExecuteAirdrops      = "execute_airdrops"
	TypeMsgSetAllocationAddress = "set_allocation_address"
)

var (
	_ sdk.Msg = &MsgAddAirdrops{}
	_ sdk.Msg = &MsgExecuteAirdrops{}
	_ sdk.Msg = &MsgSetAllocationAddress{}
)

// Route implements sdk.Msg
func (m *MsgAddAirdrops) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgAddAirdrops) Type() string { return TypeMsgAddAirdrops }

// GetSignBytes implements sdk.Msg
func (m *MsgAddAirdrops) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgAddAirdrops) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	for _, airdrop := range m.Airdrops {
		_, err := sdk.AccAddressFromBech32(airdrop.TargetAddr)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid airdrop target address (%s)", err)
		}
		err = airdrop.Amount.Validate()
		if err != nil {
			return err
		}
		_, err = sdk.ParseCoinNormalized(airdrop.Amount.String())
		if err != nil {
			return err
		}
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgAddAirdrops) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgExecuteAirdrops) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgExecuteAirdrops) Type() string { return TypeMsgExecuteAirdrops }

// GetSignBytes implements sdk.Msg
func (m *MsgExecuteAirdrops) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgExecuteAirdrops) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if m.MaxCount == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "max count must be > 0")
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgExecuteAirdrops) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgSetAllocationAddress) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgSetAllocationAddress) Type() string { return TypeMsgSetAllocationAddress }

// GetSignBytes implements sdk.Msg
func (m *MsgSetAllocationAddress) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgSetAllocationAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if len(m.TeamVestingAddr) != 0 {
		_, err := sdk.AccAddressFromBech32(m.TeamVestingAddr)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid team vesting address (%s)", err)
		}
	}
	if len(m.StrategicReserveCustodianAddr) != 0 {
		_, err := sdk.AccAddressFromBech32(m.StrategicReserveCustodianAddr)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid strategic reserve custodian address (%s)", err)
		}
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgSetAllocationAddress) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
