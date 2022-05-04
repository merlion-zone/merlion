package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/gauge/types"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
)

const (
	EmptyEpoch = 0
	FirstEpoch = 1
)

func (k Keeper) SetGauge(ctx sdk.Context, depositDenom string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GaugeKey(depositDenom), []byte(depositDenom))
}

func (k Keeper) HasGauge(ctx sdk.Context, depositDenom string) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GaugeKey(depositDenom))
	return bz != nil
}

func (k Keeper) GetGauges(ctx sdk.Context) (denoms []string) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixGaugeDenom)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		denom := string(iter.Key()[len(types.KeyPrefixGaugeDenom):])
		denoms = append(denoms, denom)
	}
	return denoms
}

func (b *Base) SetTotalDepositedAmount(ctx sdk.Context, amount sdk.Int) {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := b.keeper.cdc.MustMarshal(&sdk.IntProto{amount})
	store.Set(types.TotalDepositedAmountKey(b.prefixKey), bz)
}

func (b *Base) GetTotalDepositedAmount(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := store.Get(types.TotalDepositedAmountKey(b.prefixKey))
	if bz == nil {
		return sdk.ZeroInt()
	}
	var amount sdk.IntProto
	b.keeper.cdc.MustUnmarshal(bz, &amount)
	return amount.Int
}

func (b *Base) SetDepositedAmountByUser(ctx sdk.Context, veID uint64, amount sdk.Int) {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := b.keeper.cdc.MustMarshal(&sdk.IntProto{amount})
	store.Set(types.DepositedAmountByUserKey(b.prefixKey, veID), bz)
}

func (b *Base) GetDepositedAmountByUser(ctx sdk.Context, veID uint64) sdk.Int {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := store.Get(types.DepositedAmountByUserKey(b.prefixKey, veID))
	if bz == nil {
		return sdk.ZeroInt()
	}
	var amount sdk.IntProto
	b.keeper.cdc.MustUnmarshal(bz, &amount)
	return amount.Int
}

func (b *Base) DeleteDepositedAmountByUser(ctx sdk.Context, veID uint64) {
	store := ctx.KVStore(b.keeper.storeKey)
	store.Delete(types.DepositedAmountByUserKey(b.prefixKey, veID))
}

func (b *Base) SetTotalDerivedAmount(ctx sdk.Context, amount sdk.Int) {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := b.keeper.cdc.MustMarshal(&sdk.IntProto{amount})
	store.Set(types.TotalDerivedAmountKey(b.prefixKey), bz)
}

func (b *Base) GetTotalDerivedAmount(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := store.Get(types.TotalDerivedAmountKey(b.prefixKey))
	if bz == nil {
		return sdk.ZeroInt()
	}
	var amount sdk.IntProto
	b.keeper.cdc.MustUnmarshal(bz, &amount)
	return amount.Int
}

func (b *Base) SetDerivedAmountByUser(ctx sdk.Context, veID uint64, amount sdk.Int) {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := b.keeper.cdc.MustMarshal(&sdk.IntProto{amount})
	store.Set(types.DerivedAmountByUserKey(b.prefixKey, veID), bz)
}

func (b *Base) GetDerivedAmountByUser(ctx sdk.Context, veID uint64) sdk.Int {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := store.Get(types.DerivedAmountByUserKey(b.prefixKey, veID))
	if bz == nil {
		return sdk.ZeroInt()
	}
	var amount sdk.IntProto
	b.keeper.cdc.MustUnmarshal(bz, &amount)
	return amount.Int
}

func (b *Base) SetReward(ctx sdk.Context, rewardDenom string, reward types.Reward) {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := b.keeper.cdc.MustMarshal(&reward)
	store.Set(types.RewardKey(b.prefixKey, rewardDenom), bz)
}

func (b *Base) GetReward(ctx sdk.Context, rewardDenom string) types.Reward {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := store.Get(types.RewardKey(b.prefixKey, rewardDenom))
	if bz == nil {
		return types.Reward{
			Denom:               rewardDenom,
			Rate:                sdk.ZeroInt(),
			FinishTime:          0,
			LastUpdateTime:      0,
			CumulativePerTicket: sdk.ZeroInt(),
			AccruedAmount:       sdk.ZeroInt(),
		}
	}
	var reward types.Reward
	b.keeper.cdc.MustUnmarshal(bz, &reward)
	return reward
}

func (b *Base) IterateRewards(ctx sdk.Context, handler func(reward types.Reward) (stop bool)) {
	store := ctx.KVStore(b.keeper.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.RewardKeyPrefix(b.prefixKey))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var reward types.Reward
		b.keeper.cdc.MustUnmarshal(iter.Value(), &reward)
		if handler(reward) {
			break
		}
	}
}

func (b *Base) SetUserReward(ctx sdk.Context, rewardDenom string, veID uint64, reward types.UserReward) {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := b.keeper.cdc.MustMarshal(&reward)
	store.Set(types.UserRewardKey(b.prefixKey, rewardDenom, veID), bz)
}

func (b *Base) GetUserReward(ctx sdk.Context, rewardDenom string, veID uint64) types.UserReward {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := store.Get(types.UserRewardKey(b.prefixKey, rewardDenom, veID))
	if bz == nil {
		return types.UserReward{
			Denom:               rewardDenom,
			VeId:                veID,
			LastClaimTime:       0,
			CumulativePerTicket: sdk.ZeroInt(),
		}
	}
	var reward types.UserReward
	b.keeper.cdc.MustUnmarshal(bz, &reward)
	return reward
}

func (b *Base) SetUserVeIDByAddress(ctx sdk.Context, acc sdk.AccAddress, veID uint64) {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := sdk.Uint64ToBigEndian(veID)
	store.Set(types.UserVeIDByAddressKey(b.prefixKey, acc), bz)
}

func (b *Base) GetUserVeIDByAddress(ctx sdk.Context, acc sdk.AccAddress) uint64 {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := store.Get(types.UserVeIDByAddressKey(b.prefixKey, acc))
	if bz == nil {
		return vetypes.EmptyVeID
	}
	return sdk.BigEndianToUint64(bz)
}

func (b *Base) DeleteUserVeIDByAddress(ctx sdk.Context, acc sdk.AccAddress) {
	store := ctx.KVStore(b.keeper.storeKey)
	store.Delete(types.UserVeIDByAddressKey(b.prefixKey, acc))
}

func (b *Base) SetEpoch(ctx sdk.Context, epoch uint64) {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := sdk.Uint64ToBigEndian(epoch)
	store.Set(types.EpochKey(b.prefixKey), bz)
}

func (b *Base) GetEpoch(ctx sdk.Context) uint64 {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := store.Get(types.EpochKey(b.prefixKey))
	if bz == nil {
		return EmptyEpoch
	}
	return sdk.BigEndianToUint64(bz)
}

func (b *Base) SetCheckpoint(ctx sdk.Context, epoch uint64, point types.Checkpoint) {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := b.keeper.cdc.MustMarshal(&point)
	store.Set(types.PointKey(b.prefixKey, epoch), bz)
}

func (b *Base) GetCheckpoint(ctx sdk.Context, epoch uint64) types.Checkpoint {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := store.Get(types.PointKey(b.prefixKey, epoch))
	if bz == nil {
		return types.Checkpoint{}
	}
	var point types.Checkpoint
	b.keeper.cdc.MustUnmarshal(bz, &point)
	return point
}

func (b *Base) SetUserEpoch(ctx sdk.Context, veID uint64, epoch uint64) {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := sdk.Uint64ToBigEndian(epoch)
	store.Set(types.UserEpochKey(b.prefixKey, veID), bz)
}

func (b *Base) GetUserEpoch(ctx sdk.Context, veID uint64) uint64 {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := store.Get(types.UserEpochKey(b.prefixKey, veID))
	if bz == nil {
		return EmptyEpoch
	}
	return sdk.BigEndianToUint64(bz)
}

func (b *Base) SetUserCheckpoint(ctx sdk.Context, veID uint64, epoch uint64, point types.Checkpoint) {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := b.keeper.cdc.MustMarshal(&point)
	store.Set(types.UserPointKey(b.prefixKey, veID, epoch), bz)
}

func (b *Base) GetUserCheckpoint(ctx sdk.Context, veID uint64, epoch uint64) types.Checkpoint {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := store.Get(types.UserPointKey(b.prefixKey, veID, epoch))
	if bz == nil {
		return types.Checkpoint{}
	}
	var point types.Checkpoint
	b.keeper.cdc.MustUnmarshal(bz, &point)
	return point
}

func (b *Base) SetRewardEpoch(ctx sdk.Context, rewardDenom string, epoch uint64) {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := sdk.Uint64ToBigEndian(epoch)
	store.Set(types.RewardEpochKey(b.prefixKey, rewardDenom), bz)
}

func (b *Base) GetRewardEpoch(ctx sdk.Context, rewardDenom string) uint64 {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := store.Get(types.RewardEpochKey(b.prefixKey, rewardDenom))
	if bz == nil {
		return EmptyEpoch
	}
	return sdk.BigEndianToUint64(bz)
}

func (b *Base) SetRewardCheckpoint(ctx sdk.Context, rewardDenom string, epoch uint64, point types.Checkpoint) {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := b.keeper.cdc.MustMarshal(&point)
	store.Set(types.RewardPointKey(b.prefixKey, rewardDenom, epoch), bz)
}

func (b *Base) GetRewardCheckpoint(ctx sdk.Context, rewardDenom string, epoch uint64) types.Checkpoint {
	store := ctx.KVStore(b.keeper.storeKey)
	bz := store.Get(types.RewardPointKey(b.prefixKey, rewardDenom, epoch))
	if bz == nil {
		return types.Checkpoint{}
	}
	var point types.Checkpoint
	b.keeper.cdc.MustUnmarshal(bz, &point)
	return point
}
