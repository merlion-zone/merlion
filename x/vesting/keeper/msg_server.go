package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/merlion-zone/merlion/x/vesting/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) AddAirdrops(c context.Context, msg *types.MsgAddAirdrops) (*types.MsgAddAirdropsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Airdrops can only be added by team vesting address
	teamAddr := m.Keeper.GetAllocationAddresses(ctx).GetTeamVestingAddr()
	if !sender.Equals(teamAddr) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorized sender")
	}

	total := m.Keeper.GetAirdropTotalAmount(ctx)

	for _, airdrop := range msg.Airdrops {
		targetAddr, err := sdk.AccAddressFromBech32(airdrop.TargetAddr)
		if err != nil {
			return nil, err
		}
		amount, err := sdk.ParseCoinNormalized(airdrop.Amount.String())
		if err != nil {
			return nil, err
		}
		airdrop.Amount = amount

		m.Keeper.SetAirdrop(ctx, targetAddr, airdrop)
		total = total.Add(amount.Amount)
	}

	if total.GT(m.Keeper.GetParams(ctx).Allocation.AirdropAmount) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "total amount should not be greater than its cap")
	}

	m.Keeper.SetAirdropTotalAmount(ctx, total)

	return &types.MsgAddAirdropsResponse{}, nil
}

func (m msgServer) ExecuteAirdrops(c context.Context, msg *types.MsgExecuteAirdrops) (*types.MsgExecuteAirdropsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Airdrops can only be executed by team vesting address
	teamAddr := m.Keeper.GetAllocationAddresses(ctx).GetTeamVestingAddr()
	if !sender.Equals(teamAddr) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorized sender")
	}

	count := uint64(0)
	m.Keeper.IterateAirdrops(ctx, func(airdrop types.Airdrop) (stop bool) {
		// mint and send
		err = m.Keeper.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(airdrop.Amount))
		if err != nil {
			return true
		}
		err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, airdrop.GetTargetAddr(), sdk.NewCoins(airdrop.Amount))
		if err != nil {
			return true
		}

		// complete airdrop into store
		m.Keeper.DeleteAirdrop(ctx, airdrop.GetTargetAddr())
		m.Keeper.SetAirdropCompleted(ctx, airdrop.GetTargetAddr(), airdrop)

		count++
		return count >= msg.MaxCount
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgExecuteAirdropsResponse{}, nil
}

func (m msgServer) SetAllocationAddress(c context.Context, msg *types.MsgSetAllocationAddress) (*types.MsgSetAllocationAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if (len(msg.StrategicReserveCustodianAddr) != 0 && len(msg.TeamVestingAddr) != 0) || (len(msg.StrategicReserveCustodianAddr) == 0 && len(msg.TeamVestingAddr) == 0) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "only one allocation address must be given")
	}

	addr := m.Keeper.GetAllocationAddresses(ctx)

	if len(msg.StrategicReserveCustodianAddr) != 0 {
		newSrca, err := sdk.AccAddressFromBech32(msg.StrategicReserveCustodianAddr)
		if err != nil {
			return nil, err
		}

		if !sender.Equals(addr.GetStrategicReserveCustodianAddr()) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorized sender")
		}

		addr.StrategicReserveCustodianAddr = newSrca.String()

		// The initial strategic reserve custodian account should always be setup at genesis.
		// And subsequently it can only be changed by existing account.
	}

	if len(msg.TeamVestingAddr) != 0 {
		newTva, err := sdk.AccAddressFromBech32(msg.TeamVestingAddr)
		if err != nil {
			return nil, err
		}

		if !sender.Equals(addr.GetTeamVestingAddr()) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorized sender")
		}

		addr.TeamVestingAddr = newTva.String()
	}

	m.Keeper.SetAllocationAddresses(ctx, addr)

	return &types.MsgSetAllocationAddressResponse{}, nil
}
