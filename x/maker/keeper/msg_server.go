package keeper

import (
	"context"

	"github.com/merlion-zone/merlion/x/maker/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) MintBySwap(ctx context.Context, swap *types.MsgMintBySwap) (*types.MsgMintBySwapResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) BurnBySwap(ctx context.Context, swap *types.MsgBurnBySwap) (*types.MsgBurnBySwapResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) MintByCollateral(ctx context.Context, collateral *types.MsgMintByCollateral) (*types.MsgMintByCollateralResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) BurnByCollateral(ctx context.Context, collateral *types.MsgBurnByCollateral) (*types.MsgBurnByCollateralResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) DepositCollateral(ctx context.Context, collateral *types.MsgDepositCollateral) (*types.MsgDepositCollateralResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) RedeemCollateral(ctx context.Context, collateral *types.MsgRedeemCollateral) (*types.MsgRedeemCollateralResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) BuyBack(ctx context.Context, back *types.MsgBuyBack) (*types.MsgBuyBackResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) ReCollateralize(ctx context.Context, collateralize *types.MsgReCollateralize) (*types.MsgReCollateralizeResponse, error) {
	// TODO implement me
	panic("implement me")
}
