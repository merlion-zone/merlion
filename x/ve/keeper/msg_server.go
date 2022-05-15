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

	err := m.Keeper.checkLockDenom(ctx, msg.Amount)
	if err != nil {
		return nil, err
	}

	sender, receiver, err := getSenderReceiver(msg.Sender, msg.To)
	if err != nil {
		return nil, err
	}

	unlockTime := types.RegulatedUnixTimeFromNow(ctx, msg.LockDuration)
	if unlockTime <= uint64(ctx.BlockTime().Unix()) {
		return nil, sdkerrors.Wrapf(types.ErrPastLockTime, "past time: %s", time.Unix(int64(unlockTime), 0))
	}
	if unlockTime > uint64(ctx.BlockTime().Unix())+types.MaxLockTime {
		return nil, sdkerrors.Wrapf(types.ErrTooLongLockTime, "future time: %s", time.Unix(int64(unlockTime), 0))
	}

	// get new ve id
	veID := m.Keeper.GetNextVeID(ctx)
	if veID > types.MaxVeID || veID == types.EmptyVeID {
		return nil, sdkerrors.Wrap(types.ErrInvalidVeID, "no available ve id")
	}
	m.Keeper.SetNextVeID(ctx, veID+1)

	// mint nft for ve id
	err = m.Keeper.nftKeeper.Mint(ctx, nfttypes.NFT{
		ClassId: types.VeNftClass.Id,
		Id:      types.VeIDFromUint64(veID),
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

	err = ctx.EventManager().EmitTypedEvent(&types.EventCreate{
		Sender:     sender.String(),
		Receiver:   receiver.String(),
		VeId:       types.VeIDFromUint64(veID),
		Amount:     msg.Amount,
		UnlockTime: unlockTime,
	})
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	return &types.MsgCreateResponse{
		VeId:       types.VeIDFromUint64(veID),
		UnlockTime: unlockTime,
	}, nil
}

func (m msgServer) Deposit(c context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	err := m.Keeper.checkLockDenom(ctx, msg.Amount)
	if err != nil {
		return nil, err
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if !m.Keeper.nftKeeper.HasNFT(ctx, types.VeNftClass.Id, msg.VeId) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidVeID, "invalid ve id: %s", msg.VeId)
	}

	veID := types.Uint64FromVeID(msg.VeId)

	locked := m.Keeper.GetLockedAmountByUser(ctx, veID)
	if !locked.Amount.IsPositive() {
		// should not happen
		return nil, sdkerrors.Wrapf(types.ErrAmountNotPositive, "nothing is locked for ve %s", msg.VeId)
	}
	if locked.End <= uint64(ctx.BlockTime().Unix()) {
		return nil, sdkerrors.Wrapf(types.ErrLockExpired, "unlocking time %s but now %s", time.Unix(int64(locked.End), 0), ctx.BlockTime())
	}

	err = m.Keeper.DepositFor(ctx, sender, veID, msg.Amount.Amount, 0, locked, true)
	if err != nil {
		return nil, err
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventDeposit{
		Sender: sender.String(),
		VeId:   types.VeIDFromUint64(veID),
		Amount: msg.Amount,
	})
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	return &types.MsgDepositResponse{}, nil
}

func (m msgServer) ExtendTime(c context.Context, msg *types.MsgExtendTime) (*types.MsgExtendTimeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if !m.Keeper.nftKeeper.HasNFT(ctx, types.VeNftClass.Id, msg.VeId) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidVeID, "invalid ve id: %s", msg.VeId)
	}

	veID := types.Uint64FromVeID(msg.VeId)

	locked := m.Keeper.GetLockedAmountByUser(ctx, veID)
	if !locked.Amount.IsPositive() {
		// should not happen
		return nil, sdkerrors.Wrapf(types.ErrAmountNotPositive, "nothing is locked for ve %s", msg.VeId)
	}
	if locked.End <= uint64(ctx.BlockTime().Unix()) {
		return nil, sdkerrors.Wrapf(types.ErrLockExpired, "unlocking time %s but now %s", time.Unix(int64(locked.End), 0), ctx.BlockTime())
	}

	owner := m.Keeper.nftKeeper.GetOwner(ctx, types.VeNftClass.Id, msg.VeId)
	if !sender.Equals(owner) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "user %s do not own ve %s, so cannot extend its locking duration", sender, owner)
	}

	unlockTime := types.RegulatedUnixTimeFromNow(ctx, msg.LockDuration)
	if unlockTime <= locked.End {
		return nil, sdkerrors.Wrapf(types.ErrNotIncreasedLockTime, "unlocking time %s but existing %s", time.Unix(int64(unlockTime), 0), time.Unix(int64(locked.End), 0))
	}
	if unlockTime > uint64(ctx.BlockTime().Unix())+types.MaxLockTime {
		return nil, sdkerrors.Wrapf(types.ErrTooLongLockTime, "future time: %s", time.Unix(int64(unlockTime), 0))
	}

	err = m.Keeper.DepositFor(ctx, sender, veID, sdk.ZeroInt(), unlockTime, locked, false)
	if err != nil {
		return nil, err
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventExtendTime{
		Sender:     sender.String(),
		VeId:       types.VeIDFromUint64(veID),
		UnlockTime: unlockTime,
	})
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	return &types.MsgExtendTimeResponse{}, nil
}

func (m msgServer) Merge(c context.Context, msg *types.MsgMerge) (*types.MsgMergeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	fromOwner := m.Keeper.nftKeeper.GetOwner(ctx, types.VeNftClass.Id, msg.FromVeId)
	toOwner := m.Keeper.nftKeeper.GetOwner(ctx, types.VeNftClass.Id, msg.ToVeId)
	if !sender.Equals(fromOwner) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "user %s do not own ve %s", sender, msg.FromVeId)
	}
	if !sender.Equals(toOwner) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "user %s do not own ve %s", sender, msg.ToVeId)
	}

	fromVeID := types.Uint64FromVeID(msg.FromVeId)
	toVeID := types.Uint64FromVeID(msg.ToVeId)

	err = m.Keeper.CheckVeAttached(ctx, fromVeID)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "from ve id attached")
	}

	lockedFrom := m.Keeper.GetLockedAmountByUser(ctx, fromVeID)
	lockedTo := m.Keeper.GetLockedAmountByUser(ctx, toVeID)

	if m.Keeper.getDelegatedAmount != nil {
		delegatedAmt := m.Keeper.getDelegatedAmount(ctx, fromVeID)
		if delegatedAmt.IsPositive() {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "locked amount of from ve is delegated for staking")
		}
	}

	// NOTE: here do not check whether locks are expired

	// take the longest end time
	end := lockedFrom.End
	if lockedTo.End > end {
		end = lockedTo.End
	}

	// delete user locked of fromVeID
	m.Keeper.DeleteLockedAmountByUser(ctx, fromVeID)

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

	err = ctx.EventManager().EmitTypedEvent(&types.EventMerge{
		Sender:   sender.String(),
		FromVeId: msg.FromVeId,
		ToVeId:   msg.ToVeId,
	})
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	return &types.MsgMergeResponse{}, nil
}

func (m msgServer) Withdraw(c context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	owner := m.Keeper.nftKeeper.GetOwner(ctx, types.VeNftClass.Id, msg.VeId)
	if !sender.Equals(owner) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "user %s do not own ve %s", sender, owner)
	}

	veID := types.Uint64FromVeID(msg.VeId)

	err = m.Keeper.CheckVeAttached(ctx, veID)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "ve id attached")
	}

	locked := m.Keeper.GetLockedAmountByUser(ctx, veID)
	if locked.End > uint64(ctx.BlockTime().Unix()) {
		return nil, sdkerrors.Wrapf(types.ErrLockNotExpired, "unlocking time %s but now %s", time.Unix(int64(locked.End), 0), ctx.BlockTime())
	}

	if m.Keeper.getDelegatedAmount != nil {
		delegatedAmt := m.Keeper.getDelegatedAmount(ctx, veID)
		if delegatedAmt.IsPositive() {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "locked amount is delegated for staking")
		}
	}

	// delete user locked
	m.Keeper.DeleteLockedAmountByUser(ctx, veID)

	// update total locked
	totalLocked := m.Keeper.GetTotalLockedAmount(ctx)
	totalLocked = totalLocked.Sub(locked.Amount)
	if totalLocked.IsNegative() {
		// should never happen
		panic("total locked negative")
	}
	m.Keeper.SetTotalLockedAmount(ctx, totalLocked)

	// regulate checkpoint of veID
	m.Keeper.RegulateUserCheckpoint(ctx, veID, locked, types.NewLockedBalance())

	// burn nft of veID
	err = m.Keeper.nftKeeper.Burn(ctx, types.VeNftClass.Id, msg.VeId)
	if err != nil {
		return nil, err
	}

	// send amount to sender
	coin := sdk.NewCoin(m.Keeper.LockDenom(ctx), locked.Amount)
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(coin))
	if err != nil {
		return nil, err
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventWithdraw{
		Sender: sender.String(),
		VeId:   msg.VeId,
	})
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	return &types.MsgWithdrawResponse{}, nil
}

// DepositFor deposits some more amount and/or update locking end time for a veNFT.
// 	 veID: must be valid ve id
//   amount: locked amount to add; can be zero if no more amount to deposit
//   unlockTime: when unlocking; can be zero if no need to update
//   locked: existing locked; may be zero if no existing locked
//   sendCoins: false when extend time or merge
func (k Keeper) DepositFor(ctx sdk.Context, sender sdk.AccAddress, veID uint64, amount sdk.Int, unlockTime uint64, locked types.LockedBalance, sendCoins bool) error {
	if amount.IsPositive() {
		// take the amount from sender
		if sendCoins {
			coin := sdk.NewCoin(k.LockDenom(ctx), amount)
			err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(coin))
			if err != nil {
				return err
			}
		}

		// update total locked
		totalLocked := k.GetTotalLockedAmount(ctx)
		totalLocked = totalLocked.Add(amount)
		k.SetTotalLockedAmount(ctx, totalLocked)
	}

	// old locked
	lockedOld := locked
	// new locked
	locked.Amount = lockedOld.Amount.Add(amount)
	if unlockTime > 0 {
		locked.End = unlockTime
	}

	// update locked for veID
	k.SetLockedAmountByUser(ctx, veID, locked)

	// regulate user checkpoint for veID,
	// also regulate system checkpoint
	k.RegulateUserCheckpoint(ctx, veID, lockedOld, locked)

	return nil
}

func (k Keeper) checkLockDenom(ctx sdk.Context, amount sdk.Coin) error {
	lockDenom := k.LockDenom(ctx)
	if amount.Denom != lockDenom {
		return sdkerrors.Wrapf(types.ErrInvalidLockDenom, "amount denom %s but should be %s", amount.Denom, lockDenom)
	}
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
