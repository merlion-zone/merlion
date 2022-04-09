package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	TypeMsgMintBySwap        = "mint_by_swap"
	TypeMsgBurnBySwap        = "burn_by_swap"
	TypeMsgMintByCollateral  = "mint_by_collateral"
	TypeMsgBurnByCollateral  = "burn_by_collateral"
	TypeMsgDepositCollateral = "deposit_collateral"
	TypeMsgRedeemCollateral  = "redeem_collateral"
	TypeMsgBuyBack           = "buyback"
	TypeMsgReCollateralize   = "recollateralize"
)

var (
	_ sdk.Msg = &MsgMintBySwap{}
	_ sdk.Msg = &MsgBurnBySwap{}
	_ sdk.Msg = &MsgMintByCollateral{}
	_ sdk.Msg = &MsgBurnByCollateral{}
	_ sdk.Msg = &MsgDepositCollateral{}
	_ sdk.Msg = &MsgRedeemCollateral{}
	_ sdk.Msg = &MsgBuyBack{}
	_ sdk.Msg = &MsgReCollateralize{}
)

// Route Implements sdk.Msg
func (m *MsgMintBySwap) Route() string { return RouterKey }

// Type Implements sdk.Msg
func (m *MsgMintBySwap) Type() string { return TypeMsgMintBySwap }

// GetSignBytes implements sdk.Msg
func (m *MsgMintBySwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic sdk.Msg
func (m *MsgMintBySwap) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

// GetSigners sdk.Msg
func (m *MsgMintBySwap) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route Implements sdk.Msg
func (m *MsgBurnBySwap) Route() string { return RouterKey }

// Type Implements sdk.Msg
func (m *MsgBurnBySwap) Type() string { return TypeMsgBurnBySwap }

// GetSignBytes implements sdk.Msg
func (m *MsgBurnBySwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic sdk.Msg
func (m *MsgBurnBySwap) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

// GetSigners sdk.Msg
func (m *MsgBurnBySwap) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route Implements sdk.Msg
func (m *MsgMintByCollateral) Route() string { return RouterKey }

// Type Implements sdk.Msg
func (m *MsgMintByCollateral) Type() string { return TypeMsgMintByCollateral }

// GetSignBytes implements sdk.Msg
func (m *MsgMintByCollateral) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic sdk.Msg
func (m *MsgMintByCollateral) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

// GetSigners sdk.Msg
func (m *MsgMintByCollateral) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route Implements sdk.Msg
func (m *MsgBurnByCollateral) Route() string { return RouterKey }

// Type Implements sdk.Msg
func (m *MsgBurnByCollateral) Type() string { return TypeMsgBurnByCollateral }

// GetSignBytes implements sdk.Msg
func (m *MsgBurnByCollateral) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic sdk.Msg
func (m *MsgBurnByCollateral) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

// GetSigners sdk.Msg
func (m *MsgBurnByCollateral) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route Implements sdk.Msg
func (m *MsgDepositCollateral) Route() string { return RouterKey }

// Type Implements sdk.Msg
func (m *MsgDepositCollateral) Type() string { return TypeMsgDepositCollateral }

// GetSignBytes implements sdk.Msg
func (m *MsgDepositCollateral) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic sdk.Msg
func (m *MsgDepositCollateral) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

// GetSigners sdk.Msg
func (m *MsgDepositCollateral) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route Implements sdk.Msg
func (m *MsgRedeemCollateral) Route() string { return RouterKey }

// Type Implements sdk.Msg
func (m *MsgRedeemCollateral) Type() string { return TypeMsgRedeemCollateral }

// GetSignBytes implements sdk.Msg
func (m *MsgRedeemCollateral) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic sdk.Msg
func (m *MsgRedeemCollateral) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

// GetSigners sdk.Msg
func (m *MsgRedeemCollateral) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route Implements sdk.Msg
func (m *MsgBuyBack) Route() string { return RouterKey }

// Type Implements sdk.Msg
func (m *MsgBuyBack) Type() string { return TypeMsgBuyBack }

// GetSignBytes implements sdk.Msg
func (m *MsgBuyBack) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic sdk.Msg
func (m *MsgBuyBack) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

// GetSigners sdk.Msg
func (m *MsgBuyBack) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route Implements sdk.Msg
func (m *MsgReCollateralize) Route() string { return RouterKey }

// Type Implements sdk.Msg
func (m *MsgReCollateralize) Type() string { return TypeMsgReCollateralize }

// GetSignBytes implements sdk.Msg
func (m *MsgReCollateralize) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic sdk.Msg
func (m *MsgReCollateralize) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

// GetSigners sdk.Msg
func (m *MsgReCollateralize) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
