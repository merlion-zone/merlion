package keeper

import (
	"context"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/merlion-zone/merlion/x/staking/types"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
)

type MsgServer struct {
	stakingtypes.MsgServer
	Keeper
}

var (
	_ stakingtypes.MsgServer = MsgServer{}
	_ types.MsgServer        = MsgServer{}
)

func NewMsgServerImpl(keeper Keeper) MsgServer {
	return MsgServer{
		MsgServer: stakingkeeper.NewMsgServerImpl(keeper.Keeper),
		Keeper:    keeper,
	}
}

func (k MsgServer) VeDelegate(goCtx context.Context, msg *types.MsgVeDelegate) (*types.MsgVeDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if valErr != nil {
		return nil, valErr
	}

	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	bondDenom := k.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}

	owner := k.Keeper.nftKeeper.GetOwner(ctx, vetypes.VeNftClass.Id, msg.VeId)
	if owner.Equals(delegatorAddress) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ve %s not owned by delegator", msg.VeId)
	}

	veID := vetypes.Uint64FromVeID(msg.VeId)

	newShares, err := k.Keeper.VeDelegate(ctx, delegatorAddress, veID, msg.Amount.Amount, stakingtypes.Unbonded, validator)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, stakingtypes.ModuleName, "delegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", msg.Type()},
				float32(msg.Amount.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeVeDelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, stakingtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	})

	return &types.MsgVeDelegateResponse{}, nil
}

func (k MsgServer) CreateValidator(goCtx context.Context, msg *stakingtypes.MsgCreateValidator) (*stakingtypes.MsgCreateValidatorResponse, error) {
	return k.MsgServer.CreateValidator(goCtx, msg)
}

func (k MsgServer) EditValidator(goCtx context.Context, msg *stakingtypes.MsgEditValidator) (*stakingtypes.MsgEditValidatorResponse, error) {
	return k.MsgServer.EditValidator(goCtx, msg)
}

func (k MsgServer) Delegate(goCtx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error) {
	return k.MsgServer.Delegate(goCtx, msg)
}

func (k MsgServer) BeginRedelegate(goCtx context.Context, msg *stakingtypes.MsgBeginRedelegate) (*stakingtypes.MsgBeginRedelegateResponse, error) {
	return k.MsgServer.BeginRedelegate(goCtx, msg)
}

func (k MsgServer) Undelegate(goCtx context.Context, msg *stakingtypes.MsgUndelegate) (*stakingtypes.MsgUndelegateResponse, error) {
	return k.MsgServer.Undelegate(goCtx, msg)
}
