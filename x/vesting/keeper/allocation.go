package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/vesting/types"
	ethermint "github.com/tharsis/ethermint/types"
)

func (k Keeper) AllocateAtGenesis(ctx sdk.Context, genState types.GenesisState) {
	startTime := ctx.BlockTime().Unix()
	alloc := genState.Params.Allocation

	k.createContinuousVestingAccount(ctx, types.StakingRewardVestingName, alloc.StakingRewardAmount, startTime, types.StakingRewardVestingTime)
	k.createContinuousVestingAccount(ctx, types.CommunityPoolVestingName, alloc.CommunityPoolAmount, startTime, types.CommunityPoolVestingTime)
	k.createContinuousVestingAccount(ctx, types.TeamVestingName, alloc.TeamVestingAmount, startTime, types.TeamVestingTime)

	k.veKeeper.AddTotalEmission(ctx, alloc.VeVestingAmount)

	srAmount := sdk.NewCoin(merlion.BaseDenom, alloc.StrategicReserveAmount)
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(srAmount))
	if err != nil {
		panic(err)
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, genState.AllocationAddresses.GetStrategicReserveCustodianAddr(), sdk.NewCoins(srAmount))
	if err != nil {
		panic(err)
	}
}

func (k Keeper) ClaimVested(ctx sdk.Context) {
	vestedClaim := []struct {
		vestingName     string
		destinationName string
		destinationAddr sdk.AccAddress
	}{
		{types.StakingRewardVestingName, k.feeCollectorName, sdk.AccAddress{}},
		{types.CommunityPoolVestingName, "", sdk.AccAddress{}},
		{types.TeamVestingName, "", k.GetAllocationAddresses(ctx).GetTeamVestingAddr()},
	}

	for _, claim := range vestedClaim {
		vestingAddr := k.getVestingAddress(claim.vestingName)
		spendable := k.bankKeeper.SpendableCoins(ctx, vestingAddr)
		if !spendable.IsAllPositive() {
			continue
		}

		var err error
		if len(claim.destinationName) != 0 {
			err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, vestingAddr, claim.destinationName, spendable)
		} else if !claim.destinationAddr.Empty() {
			err = k.bankKeeper.SendCoins(ctx, vestingAddr, claim.destinationAddr, spendable)
		} else {
			err = k.distrKeeper.FundCommunityPool(ctx, spendable, vestingAddr)
		}
		if err != nil {
			panic(err)
		}
	}
}

func (k Keeper) createContinuousVestingAccount(ctx sdk.Context, vestingName string, amount sdk.Int, startTime int64, duration int64) {
	baseAccount := k.accountKeeper.NewAccountWithAddress(ctx, k.getVestingAddress(vestingName))
	amt := sdk.NewCoin(merlion.BaseDenom, amount)
	vestingAccount := vestingtypes.NewContinuousVestingAccount(baseAccount.(*ethermint.EthAccount).BaseAccount, sdk.NewCoins(amt), startTime, startTime+duration)
	k.accountKeeper.SetAccount(ctx, vestingAccount)

	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(amt))
	if err != nil {
		panic(err)
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, vestingAccount.GetAddress(), sdk.NewCoins(amt))
	if err != nil {
		panic(err)
	}
}

func (k Keeper) getVestingAddress(vestingName string) sdk.AccAddress {
	return authtypes.NewModuleAddress(vestingName)
}

// SetAllocationAddresses sets allocation target addresses
func (k Keeper) SetAllocationAddresses(ctx sdk.Context, addresses types.AllocationAddresses) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&addresses)
	store.Set(types.AllocationAddrKey(), bz)
}

// GetAllocationAddresses gets allocation target addresses
func (k Keeper) GetAllocationAddresses(ctx sdk.Context) types.AllocationAddresses {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AllocationAddrKey())
	if bz == nil {
		return types.AllocationAddresses{}
	}
	var addresses types.AllocationAddresses
	k.cdc.MustUnmarshal(bz, &addresses)
	return addresses
}

// SetAirdropTotalAmount sets airdrop total amount
func (k Keeper) SetAirdropTotalAmount(ctx sdk.Context, total sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{Int: total})
	store.Set(types.AirdropsTotalAmountKey(), bz)
}

// GetAirdropTotalAmount gets airdrop total amount
func (k Keeper) GetAirdropTotalAmount(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AirdropsTotalAmountKey())
	if bz == nil {
		return sdk.ZeroInt()
	}
	var total sdk.IntProto
	k.cdc.MustUnmarshal(bz, &total)
	return total.Int
}

// SetAirdrop sets airdrop target
func (k Keeper) SetAirdrop(ctx sdk.Context, acc sdk.AccAddress, airdrop types.Airdrop) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&airdrop)
	store.Set(types.AirdropsKey(acc), bz)
}

// GetAirdrop gets airdrop target
func (k Keeper) GetAirdrop(ctx sdk.Context, acc sdk.AccAddress) types.Airdrop {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AirdropsKey(acc))
	if bz == nil {
		return types.Airdrop{}
	}
	var airdrop types.Airdrop
	k.cdc.MustUnmarshal(bz, &airdrop)
	return airdrop
}

// DeleteAirdrop deletes airdrop target
func (k Keeper) DeleteAirdrop(ctx sdk.Context, acc sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AirdropsKey(acc))
}

// IterateAirdrops iterates airdrop targets
func (k Keeper) IterateAirdrops(ctx sdk.Context, handler func(airdrop types.Airdrop) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixAirdrops)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var airdrop types.Airdrop
		k.cdc.MustUnmarshal(iter.Value(), &airdrop)
		if handler(airdrop) {
			break
		}
	}
}

// SetAirdropCompleted sets completed airdrop target
func (k Keeper) SetAirdropCompleted(ctx sdk.Context, acc sdk.AccAddress, airdrop types.Airdrop) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&airdrop)
	store.Set(types.AirdropsCompletedKey(acc), bz)
}

// GetAirdropCompleted gets completed airdrop target
func (k Keeper) GetAirdropCompleted(ctx sdk.Context, acc sdk.AccAddress) types.Airdrop {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AirdropsCompletedKey(acc))
	if bz == nil {
		return types.Airdrop{}
	}
	var airdrop types.Airdrop
	k.cdc.MustUnmarshal(bz, &airdrop)
	return airdrop
}

// IterateAirdropsCompleted iterates completed airdrop targets
func (k Keeper) IterateAirdropsCompleted(ctx sdk.Context, handler func(airdrop types.Airdrop) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixAirdropsCompleted)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var airdrop types.Airdrop
		k.cdc.MustUnmarshal(iter.Value(), &airdrop)
		if handler(airdrop) {
			break
		}
	}
}
