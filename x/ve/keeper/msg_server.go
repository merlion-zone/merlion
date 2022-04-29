package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	nfttypes "github.com/cosmos/cosmos-sdk/x/nft"
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

func (m msgServer) Create(c context.Context, msg *types.MsgCreate) (*types.MsgCreateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, receiver, err := getSenderReceiver(msg.Sender, msg.To)
	if err != nil {
		return nil, err
	}

	unlockTime := types.RegulatedUnixTimeFromNow(ctx, msg.LockDuration)
	if unlockTime <= uint64(ctx.BlockTime().Unix()) {
		return nil, sdkerrors.Wrapf(types.ErrPastLockTime, "past time: %s", time.Unix(int64(unlockTime), 0))
	}
	if unlockTime > types.MaxLockTime {
		return nil, sdkerrors.Wrapf(types.ErrTooLongLockTime, "future time: %s", time.Unix(int64(unlockTime), 0))
	}

	// get new ve id
	veID := m.Keeper.GetNextVeID(ctx)
	if veID > types.MaxVeID || veID == types.EmptyVeID {
		return nil, sdkerrors.Wrapf(types.ErrNoValidVeID, "invalid ve id: %d", veID)
	}
	m.Keeper.SetNextVeID(ctx, veID+1)

	// mint nft for ve id
	err = m.Keeper.nftKeeper.Mint(ctx, nfttypes.NFT{
		ClassId: types.VeNftClass.Id,
		Id:      types.VeID(veID),
		Uri:     "", // TODO: implement Uri as method not field
	}, receiver)
	if err != nil {
		return nil, err
	}

	// deposit for ve id
	err = m.Keeper.DepositFor(ctx, sender, veID, msg.Amount.Amount, unlockTime, types.NewLockedBalance(), true)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateResponse{}, nil
}

func (m msgServer) Deposit(c context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	_ = sdk.UnwrapSDKContext(c)

	return &types.MsgDepositResponse{}, nil
}

func (m msgServer) Merge(c context.Context, msg *types.MsgMerge) (*types.MsgMergeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// TODO: check attachments and voted

	fromOwner := m.Keeper.nftKeeper.GetOwner(ctx, types.VeNftClass.Id, msg.FromVeId)
	toOwner := m.Keeper.nftKeeper.GetOwner(ctx, types.VeNftClass.Id, msg.ToVeId)
	if !sender.Equals(fromOwner) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrTooManySignatures, "user %s do not own ve %s", sender, msg.FromVeId)
	}
	if !sender.Equals(toOwner) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrTooManySignatures, "user %s do not own ve %s", sender, msg.ToVeId)
	}

	fromVeID := types.Uint64FromVeID(msg.FromVeId)
	toVeID := types.Uint64FromVeID(msg.ToVeId)

	lockedFrom := m.Keeper.GetLockedAmountByUser(ctx, fromVeID)
	lockedTo := m.Keeper.GetLockedAmountByUser(ctx, toVeID)

	// take the longest end time
	end := lockedFrom.End
	if lockedTo.End > end {
		end = lockedTo.End
	}

	// reset locked of fromVeID
	m.Keeper.SetLockedAmountByUser(ctx, fromVeID, types.NewLockedBalance())

	// regulate checkpoint of fromVeID
	m.Keeper.RegulateUserCheckpoint(ctx, fromVeID, lockedFrom, types.NewLockedBalance())

	// burn nft of fromVeID
	err = m.Keeper.nftKeeper.Burn(ctx, types.VeNftClass.Id, msg.FromVeId)
	if err != nil {
		return nil, err
	}

	// deposit for toVeID
	err = m.Keeper.DepositFor(ctx, sender, toVeID, lockedFrom.Amount, end, lockedTo, false)
	if err != nil {
		return nil, err
	}

	return &types.MsgMergeResponse{}, nil
}

func (m msgServer) Withdraw(ctx context.Context, withdraw *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	// TODO implement me
	panic("implement me")
}

// DepositFor deposits some more amount and/or update locking end time for a veNFT.
// veID: must be valid ve id
// amount: can be zero if no more amount to deposit
// end: can be zero if do not update locking end time
// locked: may be zero if no existing locked
func (k Keeper) DepositFor(ctx sdk.Context, sender sdk.AccAddress, veID uint64, amount sdk.Int, end uint64, locked types.LockedBalance, sendCoins bool) error {
	coin := sdk.NewCoin(k.LockDenom(ctx), amount)

	// take the amount from sender
	if sendCoins {
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
	}

	// update total locked
	if amount.IsPositive() {
		totalLocked := k.GetTotalLockedAmount(ctx)
		totalLocked = totalLocked.Add(amount)
		k.SetTotalLockedAmount(ctx, totalLocked)
	}

	lockedOld := locked
	locked.Amount = lockedOld.Amount.Add(amount)
	if end > 0 {
		locked.End = end
	}

	k.SetLockedAmountByUser(ctx, veID, locked)

	k.RegulateUserCheckpoint(ctx, veID, lockedOld, locked)

	return nil
}

func getSenderReceiver(senderStr, toStr string) (sender sdk.AccAddress, receiver sdk.AccAddress, err error) {
	sender, err = sdk.AccAddressFromBech32(senderStr)
	if err != nil {
		return
	}
	receiver = sender
	if len(toStr) > 0 {
		// user specifies receiver
		receiver, err = sdk.AccAddressFromBech32(toStr)
		if err != nil {
			return
		}
	}
	return
}
