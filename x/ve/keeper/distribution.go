package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/ve/types"
)

type Distributor struct {
	keeper Keeper
}

func NewDistributor(keeper Keeper) Distributor {
	return Distributor{keeper: keeper}
}

func (d Distributor) DistributePerPeriod(ctx sdk.Context) {
	totalAmount := d.keeper.bankKeeper.GetBalance(ctx, d.keeper.accountKeeper.GetModuleAddress(types.DistributionPoolName), d.keeper.LockDenom(ctx)).Amount

	totalAmountLast := d.keeper.GetDistributionTotalAmount(ctx)
	d.keeper.SetDistributionTotalAmount(ctx, totalAmount)

	amount := totalAmount.Sub(totalAmountLast)

	now := uint64(ctx.BlockTime().Unix())
	timeLast := d.keeper.GetDistributionAccruedLastTimestamp(ctx)
	d.keeper.SetDistributionAccruedLastTimestamp(ctx, now)

	duration := now - timeLast

	epochTime := types.RegulatedUnixTime(timeLast)
	for {
		amountExisting := d.keeper.GetDistributionPerPeriod(ctx, epochTime)
		if duration == 0 || epochTime >= now {
			d.keeper.SetDistributionPerPeriod(ctx, epochTime, amountExisting.Add(amount))
			break
		}

		nextEpochTime := types.NextRegulatedUnixTime(epochTime)
		endTime := merlion.Min(nextEpochTime, now)
		amountAdd := amount.MulRaw(int64(endTime - timeLast)).QuoRaw(int64(duration))
		d.keeper.SetDistributionPerPeriod(ctx, epochTime, amountExisting.Add(amountAdd))

		if nextEpochTime >= now {
			break
		}
		epochTime = nextEpochTime
	}
}

func (d Distributor) Claim(ctx sdk.Context, veID uint64) error {
	d.keeper.RegulateCheckpoint(ctx)

	owner := d.keeper.nftKeeper.GetOwner(ctx, types.VeNftClass.Id, types.VeID(veID))

	now := uint64(ctx.BlockTime().Unix())
	timeLast := d.keeper.GetDistributionClaimLastTimestampByUser(ctx, veID)
	if now-timeLast < types.RegulatedPeriod {
		return nil
	}
	d.keeper.SetDistributionClaimLastTimestampByUser(ctx, veID, now)

	amount := sdk.ZeroInt()
	epochTime := types.RegulatedUnixTime(timeLast)
	for {
		epochTime = types.NextRegulatedUnixTime(epochTime)
		if epochTime > types.RegulatedUnixTime(now) {
			break
		}

		amountOfPeriod := d.keeper.GetDistributionPerPeriod(ctx, types.PreviousRegulatedUnixTime(epochTime))
		votingPower := d.keeper.GetVotingPower(ctx, veID, epochTime, 0)
		totalVotingPower := d.keeper.GetTotalVotingPower(ctx, epochTime, 0)
		amount = amount.Add(amountOfPeriod.Mul(votingPower).Quo(totalVotingPower))
	}

	if !amount.IsPositive() {
		return nil
	}
	coin := sdk.NewCoin(d.keeper.LockDenom(ctx), amount)
	err := d.keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.DistributionPoolName, owner, sdk.NewCoins(coin))
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) SetDistributionAccruedLastTimestamp(ctx sdk.Context, timestamp uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DistributionAccruedLastTimestampKey(), sdk.Uint64ToBigEndian(timestamp))
}

func (k Keeper) GetDistributionAccruedLastTimestamp(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DistributionAccruedLastTimestampKey())
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetDistributionTotalAmount(ctx sdk.Context, total sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{total})
	store.Set(types.DistributionTotalAmountKey(), bz)
}

func (k Keeper) GetDistributionTotalAmount(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DistributionTotalAmountKey())
	if bz == nil {
		return sdk.ZeroInt()
	}
	var total sdk.IntProto
	k.cdc.MustUnmarshal(bz, &total)
	return total.Int
}

func (k Keeper) SetDistributionPerPeriod(ctx sdk.Context, timestamp uint64, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{amount})
	store.Set(types.DistributionPerPeriodKey(timestamp), bz)
}

func (k Keeper) GetDistributionPerPeriod(ctx sdk.Context, timestamp uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DistributionPerPeriodKey(timestamp))
	if bz == nil {
		return sdk.ZeroInt()
	}
	var amount sdk.IntProto
	k.cdc.MustUnmarshal(bz, &amount)
	return amount.Int
}

func (k Keeper) SetDistributionClaimLastTimestampByUser(ctx sdk.Context, veID uint64, timestamp uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DistributionClaimLastTimestampByUserKey(veID), sdk.Uint64ToBigEndian(timestamp))
}

func (k Keeper) GetDistributionClaimLastTimestampByUser(ctx sdk.Context, veID uint64) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DistributionClaimLastTimestampByUserKey(veID))
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}
