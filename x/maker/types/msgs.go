package types

import sdk "github.com/cosmos/cosmos-sdk/types"

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

func (m *MsgMintBySwap) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *MsgMintBySwap) GetSigners() []sdk.AccAddress {
	// TODO implement me
	panic("implement me")
}

func (m *MsgBurnBySwap) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *MsgBurnBySwap) GetSigners() []sdk.AccAddress {
	// TODO implement me
	panic("implement me")
}

func (m *MsgMintByCollateral) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *MsgMintByCollateral) GetSigners() []sdk.AccAddress {
	// TODO implement me
	panic("implement me")
}

func (m *MsgBurnByCollateral) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *MsgBurnByCollateral) GetSigners() []sdk.AccAddress {
	// TODO implement me
	panic("implement me")
}

func (m *MsgDepositCollateral) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *MsgDepositCollateral) GetSigners() []sdk.AccAddress {
	// TODO implement me
	panic("implement me")
}

func (m *MsgRedeemCollateral) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *MsgRedeemCollateral) GetSigners() []sdk.AccAddress {
	// TODO implement me
	panic("implement me")
}

func (m *MsgBuyBack) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *MsgBuyBack) GetSigners() []sdk.AccAddress {
	// TODO implement me
	panic("implement me")
}

func (m *MsgReCollateralize) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *MsgReCollateralize) GetSigners() []sdk.AccAddress {
	// TODO implement me
	panic("implement me")
}
