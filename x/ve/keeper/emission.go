package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/types"
)

type Emitter struct {
	keeper Keeper
}

func NewEmitter(keeper Keeper) Emitter {
	return Emitter{keeper: keeper}
}

func (e Emitter) AddTotalEmission(ctx sdk.Context, emission sdk.Int) {
	if !emission.IsPositive() {
		panic("emission must be nonzero")
	}

	total := e.keeper.GetTotalEmission(ctx)
	e.keeper.SetTotalEmission(ctx, total.Add(emission))

	// for geometric sequence of weeks of every 4 years,
	// a * (1 - r^n) / (1 - r) = <total emission>
	emissionInitial := emission.ToDec().Mul(sdk.OneDec().Sub(types.EmissionRatio)).Quo(sdk.OneDec().Sub(types.EmissionRatio.Power(types.MaxLockTimeWeeks))).TruncateInt()

	emissionLast := e.keeper.GetEmissionAtLastPeriod(ctx)
	e.keeper.SetEmissionAtLastPeriod(ctx, emissionLast.Add(emissionInitial))
}

func (e Emitter) CirculationSupply(ctx sdk.Context) sdk.Int {
	totalSupply := e.keeper.bankKeeper.GetSupply(ctx, e.keeper.LockDenom(ctx)).Amount
	// actually voting power is degenerative locked amount by ve
	veLocked := e.keeper.GetTotalVotingPower(ctx, 0, ctx.BlockHeight())
	return totalSupply.Sub(veLocked)
}

func (e Emitter) CirculationRate(ctx sdk.Context) sdk.Dec {
	totalSupply := e.keeper.bankKeeper.GetSupply(ctx, e.keeper.LockDenom(ctx)).Amount
	return e.CirculationSupply(ctx).ToDec().QuoInt(totalSupply)
}

func (e Emitter) Emission(ctx sdk.Context) sdk.Int {
	emissionLast := e.keeper.GetEmissionAtLastPeriod(ctx)
	emission := emissionLast.ToDec().Mul(types.EmissionRatio)

	circulationRate := e.CirculationRate(ctx)
	if circulationRate.LT(types.MinEmissionCirculating) {
		circulationRate = types.MinEmissionCirculating
	}

	return emission.Mul(circulationRate).TruncateInt()
}

func (e Emitter) EmissionCompensation(ctx sdk.Context, emission sdk.Int) sdk.Int {
	return emission.ToDec().Mul(sdk.OneDec().Sub(e.CirculationRate(ctx))).TruncateInt()
}

// Emit emits coin rewards of every period, on the basis of predefined emission policy.
// The part of compensation for ve holders will be sent into the distribution pool.
// The remaining will be deposited as rewards by the voter module.
func (e Emitter) Emit(ctx sdk.Context) sdk.Int {
	timestamp := types.RegulatedUnixTimeFromNow(ctx, 0)
	timeLast := e.keeper.GetEmissionLastTimestamp(ctx)
	// only allow one emission per period
	if timestamp-timeLast < types.RegulatedPeriod {
		return sdk.ZeroInt()
	}
	// TODO: add initialization
	if true {
		return sdk.ZeroInt()
	}

	emission := e.Emission(ctx)

	// mint emission amount
	emissionAmt := sdk.NewCoin(e.keeper.LockDenom(ctx), emission)
	err := e.keeper.bankKeeper.MintCoins(ctx, types.EmissionPoolName, sdk.NewCoins(emissionAmt))
	if err != nil {
		panic(err)
	}

	e.keeper.SetEmissionLastTimestamp(ctx, timestamp)
	e.keeper.SetEmissionAtLastPeriod(ctx, emission)

	// calculate compensation for ve holders due to inflation loss
	compensation := e.EmissionCompensation(ctx, emission)
	compensationAmt := sdk.NewCoin(e.keeper.LockDenom(ctx), compensation)
	// send compensation into distribution pool
	err = e.keeper.bankKeeper.SendCoinsFromModuleToModule(ctx, types.EmissionPoolName, types.DistributionPoolName, sdk.NewCoins(compensationAmt))
	if err != nil {
		panic(err)
	}
	NewDistributor(e.keeper).DistributePerPeriod(ctx)

	e.keeper.RegulateCheckpoint(ctx)

	emission = emission.Sub(compensation)
	return emission
}

func (k Keeper) AddTotalEmission(ctx sdk.Context, emission sdk.Int) {
	NewEmitter(k).AddTotalEmission(ctx, emission)
}

func (k Keeper) SetTotalEmission(ctx sdk.Context, total sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{total})
	store.Set(types.TotalEmissionKey(), bz)
}

func (k Keeper) GetTotalEmission(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TotalEmissionKey())
	if bz == nil {
		return sdk.ZeroInt()
	}
	var total sdk.IntProto
	k.cdc.MustUnmarshal(bz, &total)
	return total.Int
}

func (k Keeper) SetEmissionAtLastPeriod(ctx sdk.Context, emission sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{emission})
	store.Set(types.EmissionAtLastPeriodKey(), bz)
}

func (k Keeper) GetEmissionAtLastPeriod(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EmissionAtLastPeriodKey())
	if bz == nil {
		return sdk.ZeroInt()
	}
	var emission sdk.IntProto
	k.cdc.MustUnmarshal(bz, &emission)
	return emission.Int
}

func (k Keeper) SetEmissionLastTimestamp(ctx sdk.Context, timestamp uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.EmissionLastTimestampKey(), sdk.Uint64ToBigEndian(timestamp))
}

func (k Keeper) GetEmissionLastTimestamp(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EmissionLastTimestampKey())
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}
