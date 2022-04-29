package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgCreate   = "create"
	TypeMsgDeposit  = "deposit"
	TypeMsgMerge    = "merge"
	TypeMsgWithdraw = "withdraw"
)

var (
	_ sdk.Msg = &MsgCreate{}
	_ sdk.Msg = &MsgDeposit{}
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
	if !m.Amount.IsPositive() {
	}
	if m.LockDuration > MaxLockTime {
		return ErrTooLongLockTime
	}
	// TODO implement me
	panic("implement me")
}

// GetSigners implements sdk.Msg
func (m *MsgCreate) GetSigners() []sdk.AccAddress {
	// TODO implement me
	panic("implement me")
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
	// TODO implement me
	panic("implement me")
}

// GetSigners implements sdk.Msg
func (m *MsgDeposit) GetSigners() []sdk.AccAddress {
	// TODO implement me
	panic("implement me")
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
	if Uint64FromVeID(m.FromVeId) == EmptyVeID {
	}
	if Uint64FromVeID(m.ToVeId) == EmptyVeID {
	}
	if m.FromVeId == m.ToVeId {
	}

	// TODO implement me
	panic("implement me")
}

// GetSigners implements sdk.Msg
func (m *MsgMerge) GetSigners() []sdk.AccAddress {
	// TODO implement me
	panic("implement me")
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
	// TODO implement me
	panic("implement me")
}

// GetSigners implements sdk.Msg
func (m *MsgWithdraw) GetSigners() []sdk.AccAddress {
	// TODO implement me
	panic("implement me")
}
