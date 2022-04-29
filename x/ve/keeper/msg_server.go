package keeper

import (
	"context"

	"github.com/merlion-zone/merlion/x/ve/types"
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

func (m msgServer) Create(ctx context.Context, create *types.MsgCreate) (*types.MsgCreateResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) Deposit(ctx context.Context, deposit *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) Merge(ctx context.Context, merge *types.MsgMerge) (*types.MsgMergeResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) Withdraw(ctx context.Context, withdraw *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	// TODO implement me
	panic("implement me")
}
