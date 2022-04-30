package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreate     = "create"
	TypeMsgDeposit    = "deposit"
	TypeMsgExtendTime = "extend_time"
	TypeMsgMerge      = "merge"
	TypeMsgWithdraw   = "withdraw"
)

var (
	_ sdk.Msg = &MsgCreate{}
	_ sdk.Msg = &MsgDeposit{}
	_ sdk.Msg = &MsgExtendTime{}
	_ sdk.Msg = &MsgMerge{}
	_ sdk.Msg = &MsgWithdraw{}
)

// Route implements sdk.Msg
func (m *MsgCreate) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgCreate) Type() string { return TypeMsgCreate }

// GetSignBytes implements sdk.Msg
func (m *MsgCreate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgCreate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if len(m.To) > 0 {
		_, err = sdk.AccAddressFromBech32(m.To)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
	}
	if !m.Amount.IsPositive() {
		return ErrAmountNotPositive
	}
	if m.LockDuration == 0 {
		return ErrPastLockTime
	}
	if m.LockDuration > MaxLockTime {
		return ErrTooLongLockTime
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgCreate) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgDeposit) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgDeposit) Type() string { return TypeMsgDeposit }

// GetSignBytes implements sdk.Msg
func (m *MsgDeposit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if Uint64FromVeID(m.VeId) == EmptyVeID {
		return ErrInvalidVeID
	}
	if !m.Amount.IsPositive() {
		return ErrAmountNotPositive
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgDeposit) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgExtendTime) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgExtendTime) Type() string { return TypeMsgExtendTime }

// GetSignBytes implements sdk.Msg
func (m *MsgExtendTime) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgExtendTime) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if Uint64FromVeID(m.VeId) == EmptyVeID {
		return ErrInvalidVeID
	}
	if m.LockDuration == 0 {
		return ErrPastLockTime
	}
	if m.LockDuration > MaxLockTime {
		return ErrTooLongLockTime
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgExtendTime) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgMerge) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgMerge) Type() string { return TypeMsgMerge }

// GetSignBytes implements sdk.Msg
func (m *MsgMerge) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgMerge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if Uint64FromVeID(m.FromVeId) == EmptyVeID {
		return ErrInvalidVeID
	}
	if Uint64FromVeID(m.ToVeId) == EmptyVeID {
		return ErrInvalidVeID
	}
	if m.FromVeId == m.ToVeId {
		return ErrSameVeID
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgMerge) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements sdk.Msg
func (m *MsgWithdraw) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m *MsgWithdraw) Type() string { return TypeMsgWithdraw }

// GetSignBytes implements sdk.Msg
func (m *MsgWithdraw) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic implements sdk.Msg
func (m *MsgWithdraw) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if Uint64FromVeID(m.VeId) == EmptyVeID {
		return ErrInvalidVeID
	}
	return nil
}

// GetSigners implements sdk.Msg
func (m *MsgWithdraw) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
