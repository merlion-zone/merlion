package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgMintBySwap{}, "merlion/MsgMintBySwap", nil)
	cdc.RegisterConcrete(&MsgBurnBySwap{}, "merlion/MsgBurnBySwap", nil)
	cdc.RegisterConcrete(&MsgBuyBacking{}, "merlion/MsgBuyBacking", nil)
	cdc.RegisterConcrete(&MsgSellBacking{}, "merlion/MsgSellBacking", nil)
	cdc.RegisterConcrete(&MsgMintByCollateral{}, "merlion/MsgMintByCollateral", nil)
	cdc.RegisterConcrete(&MsgBurnByCollateral{}, "merlion/MsgBurnByCollateral", nil)
	cdc.RegisterConcrete(&MsgDepositCollateral{}, "merlion/MsgDepositCollateral", nil)
	cdc.RegisterConcrete(&MsgRedeemCollateral{}, "merlion/MsgRedeemCollateral", nil)
	cdc.RegisterConcrete(&MsgLiquidateCollateral{}, "merlion/MsgLiquidateCollateral", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&RegisterBackingProposal{},
		&RegisterCollateralProposal{},
		&SetBackingRiskParamsProposal{},
		&SetCollateralRiskParamsProposal{},
		&BatchSetBackingRiskParamsProposal{},
		&BatchSetCollateralRiskParamsProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(Amino)
)

func init() {
	RegisterCodec(Amino)
}
