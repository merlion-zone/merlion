package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/maker/types"
)

// AdjustCollateralRatio dynamically adjusts the collateral ratio, according to mer price change.
func (k Keeper) AdjustCollateralRatio(ctx sdk.Context) {
	// check cooldown period since last update
	if ctx.BlockHeight()-k.GetCollateralRatioLastBlock(ctx) < k.CollateralRatioCooldownPeriod(ctx) {
		return
	}

	collateralRatio := k.GetCollateralRatio(ctx)
	ratioStep := k.CollateralRatioStep(ctx)
	priceBand := merlion.MicroUSDTarget.Mul(k.CollateralRatioPriceBand(ctx))

	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		panic(err)
	}

	if merPrice.GT(merlion.MicroUSDTarget.Add(priceBand)) {
		// mer price is too high
		// decrease collateral ratio; min 0%
		collateralRatio = sdk.MaxDec(collateralRatio.Sub(ratioStep), sdk.ZeroDec())
	} else if merPrice.LT(merlion.MicroUSDTarget.Sub(priceBand)) {
		// mer price is too low
		// increase collateral ratio; max 100%
		collateralRatio = sdk.MinDec(collateralRatio.Add(ratioStep), sdk.OneDec())
	}

	// TODO: consider adjusting CR based on total minted Mer, even though Mer price is within the band

	k.SetCollateralRatio(ctx, collateralRatio)
	k.SetCollateralRatioLastBlock(ctx, ctx.BlockHeight())
}

func (k Keeper) SetCollateralRatio(ctx sdk.Context, cr sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.DecProto{Dec: cr})
	store.Set(types.KeyPrefixCollateralRatio, bz)
}

func (k Keeper) GetCollateralRatio(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixCollateralRatio)
	if bz == nil {
		return sdk.OneDec()
	}
	dp := sdk.DecProto{}
	k.cdc.MustUnmarshal(bz, &dp)
	return dp.Dec
}

func (k Keeper) SetCollateralRatioLastBlock(ctx sdk.Context, bh int64) {
	store := ctx.KVStore(k.storeKey)
	if bh < 0 {
		panic("invalid block height")
	}
	bz := sdk.Uint64ToBigEndian(uint64(bh))
	store.Set(types.KeyPrefixCollateralRatioLastBlock, bz)
}

func (k Keeper) GetCollateralRatioLastBlock(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixCollateralRatio)
	if bz == nil {
		return 0
	}
	return int64(sdk.BigEndianToUint64(bz))
}
