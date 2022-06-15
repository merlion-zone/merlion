package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAggregateExchangeRatePrevote{}, "oracle/MsgAggregateExchangeRatePrevote", nil)
	cdc.RegisterConcrete(&MsgAggregateExchangeRateVote{}, "oracle/MsgAggregateExchangeRateVote", nil)
	cdc.RegisterConcrete(&MsgDelegateFeedConsent{}, "oracle/MsgDelegateFeedConsent", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&RegisterTargetProposal{},
	)

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAggregateExchangeRatePrevote{},
		&MsgAggregateExchangeRateVote{},
		&MsgDelegateFeedConsent{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
