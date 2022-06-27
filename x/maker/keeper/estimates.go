package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/maker/types"
)

func (k Keeper) estimateMintBySwapIn(
	ctx sdk.Context,
	mintOut sdk.Coin,
	backingDenom string,
	fullBacking bool,
) (
	backingIn sdk.Coin,
	lionIn sdk.Coin,
	mintFee sdk.Coin,
	err error,
) {
	backingIn = sdk.NewCoin(backingDenom, sdk.ZeroInt())
	lionIn = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	mintFee = sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())

	err = k.checkMerPriceLowerBound(ctx)
	if err != nil {
		return
	}

	backingParams, err := k.getEnabledBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	mintFee = computeFee(mintOut, backingParams.MintFee)
	mintTotal := mintOut.Add(mintFee)
	mintTotalInUSD := mintTotal.Amount.ToDec().Mul(merlion.MicroUSDTarget)

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}
	poolBacking.MerMinted = poolBacking.MerMinted.Add(mintTotal)
	if backingParams.MaxMerMint != nil && poolBacking.MerMinted.Amount.GT(*backingParams.MaxMerMint) {
		err = sdkerrors.Wrapf(types.ErrMerCeiling, "mer over ceiling")
		return
	}

	backingRatio := k.GetBackingRatio(ctx)
	if backingRatio.GTE(sdk.OneDec()) || fullBacking {
		// full/over backing, or user selects full backing
		backingIn.Amount = mintTotalInUSD.QuoRoundUp(backingPrice).RoundInt()
	} else if backingRatio.IsZero() {
		// full algorithmic
		lionIn.Amount = mintTotalInUSD.QuoRoundUp(lionPrice).RoundInt()
	} else {
		// fractional
		backingIn.Amount = mintTotalInUSD.Mul(backingRatio).QuoRoundUp(backingPrice).RoundInt()
		lionIn.Amount = mintTotalInUSD.Mul(sdk.OneDec().Sub(backingRatio)).QuoRoundUp(lionPrice).RoundInt()
	}

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		err = sdkerrors.Wrapf(types.ErrBackingCeiling, "backing over ceiling")
		return
	}

	return
}

func (k Keeper) estimateMintBySwapOut(
	ctx sdk.Context,
	backingInMax sdk.Coin,
	lionInMax sdk.Coin,
	fullBacking bool,
) (
	backingIn sdk.Coin,
	lionIn sdk.Coin,
	mintOut sdk.Coin,
	mintFee sdk.Coin,
	err error,
) {
	backingDenom := backingInMax.Denom

	backingIn = sdk.NewCoin(backingDenom, sdk.ZeroInt())
	lionIn = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	mintOut = sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())
	mintFee = sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())

	err = k.checkMerPriceLowerBound(ctx)
	if err != nil {
		return
	}

	backingParams, err := k.getEnabledBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	backingAvailInUSD := backingPrice.MulInt(backingInMax.Amount)
	lionAvailInUSD := lionPrice.MulInt(lionInMax.Amount)

	mintTotalInUSD := sdk.ZeroDec()
	backingRatio := k.GetBackingRatio(ctx)
	if backingRatio.GTE(sdk.OneDec()) || fullBacking {
		// full/over backing, or user selects full backing
		mintTotalInUSD = backingAvailInUSD
		backingIn.Amount = backingInMax.Amount
	} else if backingRatio.IsZero() {
		// full algorithmic
		mintTotalInUSD = lionAvailInUSD
		lionIn.Amount = lionInMax.Amount
	} else {
		// fractional
		mintTotalWithBackingInUSD := backingAvailInUSD.Quo(backingRatio)
		mintTotalWithLionInUSD := lionAvailInUSD.Quo(sdk.OneDec().Sub(backingRatio))
		if lionInMax.IsZero() || (backingInMax.IsPositive() && mintTotalWithBackingInUSD.LT(mintTotalWithLionInUSD)) {
			mintTotalInUSD = mintTotalWithBackingInUSD
			backingIn.Amount = backingInMax.Amount
			lionIn.Amount = mintTotalInUSD.Mul(sdk.OneDec().Sub(backingRatio)).QuoRoundUp(lionPrice).RoundInt()
		} else {
			mintTotalInUSD = mintTotalWithLionInUSD
			lionIn.Amount = lionInMax.Amount
			backingIn.Amount = mintTotalInUSD.Mul(backingRatio).QuoRoundUp(backingPrice).RoundInt()
		}
	}

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}

	poolBacking.MerMinted = poolBacking.MerMinted.AddAmount(mintTotalInUSD.RoundInt())
	if backingParams.MaxMerMint != nil && poolBacking.MerMinted.Amount.GT(*backingParams.MaxMerMint) {
		err = sdkerrors.Wrapf(types.ErrMerCeiling, "mer over ceiling")
		return
	}

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		err = sdkerrors.Wrapf(types.ErrBackingCeiling, "backing over ceiling")
		return
	}

	mintFeeRate := sdk.ZeroDec()
	if backingParams.MintFee != nil {
		mintFeeRate = *backingParams.MintFee
	}

	mintOutValue := mintTotalInUSD.Quo(merlion.MicroUSDTarget).Quo(sdk.OneDec().Add(mintFeeRate))
	mintFeeValue := mintOutValue.Mul(mintFeeRate)
	mintOut = sdk.NewCoin(merlion.MicroUSDDenom, mintOutValue.RoundInt())
	mintFee = sdk.NewCoin(merlion.MicroUSDDenom, mintFeeValue.RoundInt())
	return
}

func (k Keeper) estimateBurnBySwapIn(
	ctx sdk.Context,
	backingOutMax sdk.Coin,
	lionOutMax sdk.Coin,
) (
	burnIn sdk.Coin,
	backingOut sdk.Coin,
	lionOut sdk.Coin,
	burnFee sdk.Coin,
	err error,
) {
	backingDenom := backingOutMax.Denom

	burnIn = sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())
	backingOut = sdk.NewCoin(backingOutMax.Denom, sdk.ZeroInt())
	lionOut = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	burnFee = sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())

	err = k.checkMerPriceUpperBound(ctx)
	if err != nil {
		return
	}

	backingParams, err := k.getEnabledBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	backingOutMaxInUSD := backingPrice.MulInt(backingOutMax.Amount)
	lionOutMaxInUSD := lionPrice.MulInt(lionOutMax.Amount)

	burnActualInUSD := sdk.ZeroDec()
	backingRatio := k.GetBackingRatio(ctx)
	if backingRatio.GTE(sdk.OneDec()) {
		// full/over backing
		burnActualInUSD = backingOutMaxInUSD
		backingOut.Amount = backingOutMax.Amount
	} else if backingRatio.IsZero() {
		// full algorithmic
		burnActualInUSD = lionOutMaxInUSD
		lionOut.Amount = lionOutMax.Amount
	} else {
		// fractional
		burnActualWithBackingInUSD := backingOutMaxInUSD.Quo(backingRatio)
		burnActualWithLionInUSD := lionOutMaxInUSD.Quo(sdk.OneDec().Sub(backingRatio))
		if lionOutMax.IsZero() || (backingOutMax.IsPositive() && burnActualWithBackingInUSD.LT(burnActualWithLionInUSD)) {
			burnActualInUSD = burnActualWithBackingInUSD
			backingOut.Amount = backingOutMax.Amount
			lionOut.Amount = burnActualInUSD.Mul(sdk.OneDec().Sub(backingRatio)).QuoRoundUp(lionPrice).RoundInt()
		} else {
			burnActualInUSD = burnActualWithLionInUSD
			lionOut.Amount = lionOutMax.Amount
			backingOut.Amount = burnActualInUSD.Mul(backingRatio).QuoRoundUp(backingPrice).RoundInt()
		}
	}

	moduleOwnedBacking := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), backingDenom)
	if moduleOwnedBacking.IsLT(backingOut) {
		err = sdkerrors.Wrapf(types.ErrBackingCoinInsufficient, "backing coin out(%s) < balance(%s)", backingOut, moduleOwnedBacking)
		return
	}

	burnFeeRate := sdk.ZeroDec()
	if backingParams.BurnFee != nil {
		burnFeeRate = *backingParams.BurnFee
	}

	burnInValue := burnActualInUSD.Quo(merlion.MicroUSDTarget).Quo(sdk.OneDec().Sub(burnFeeRate))
	burnFeeValue := burnInValue.Mul(burnFeeRate)
	burnIn = sdk.NewCoin(merlion.MicroUSDDenom, burnInValue.RoundInt())
	burnFee = sdk.NewCoin(merlion.MicroUSDDenom, burnFeeValue.RoundInt())
	return
}

func (k Keeper) estimateBurnBySwapOut(
	ctx sdk.Context,
	burnIn sdk.Coin,
	backingDenom string,
) (
	backingOut sdk.Coin,
	lionOut sdk.Coin,
	burnFee sdk.Coin,
	err error,
) {
	backingOut = sdk.NewCoin(backingDenom, sdk.ZeroInt())
	lionOut = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	burnFee = sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())

	err = k.checkMerPriceUpperBound(ctx)
	if err != nil {
		return
	}

	backingParams, err := k.getEnabledBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	burnFee = computeFee(burnIn, backingParams.BurnFee)
	burnActual := burnIn.Sub(burnFee)
	burnActualInUSD := burnActual.Amount.ToDec().Mul(merlion.MicroUSDTarget)

	backingRatio := k.GetBackingRatio(ctx)
	if backingRatio.GTE(sdk.OneDec()) {
		// full/over backing
		backingOut.Amount = burnActualInUSD.QuoRoundUp(backingPrice).RoundInt()
	} else if backingRatio.IsZero() {
		// full algorithmic
		lionOut.Amount = burnActualInUSD.QuoRoundUp(lionPrice).RoundInt()
	} else {
		// fractional
		backingOut.Amount = burnActualInUSD.Mul(backingRatio).QuoRoundUp(backingPrice).RoundInt()
		lionOut.Amount = burnActualInUSD.Mul(sdk.OneDec().Sub(backingRatio)).QuoRoundUp(lionPrice).RoundInt()
	}

	moduleOwnedBacking := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), backingDenom)
	if moduleOwnedBacking.IsLT(backingOut) {
		err = sdkerrors.Wrapf(types.ErrBackingCoinInsufficient, "backing coin out(%s) < balance(%s)", backingOut, moduleOwnedBacking)
		return
	}

	return
}

func (k Keeper) estimateBuyBackingIn(
	ctx sdk.Context,
	backingOut sdk.Coin,
) (
	lionIn sdk.Coin,
	buybackFee sdk.Coin,
	err error,
) {
	backingDenom := backingOut.Denom

	lionIn = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	buybackFee = sdk.NewCoin(backingDenom, sdk.ZeroInt())

	backingParams, err := k.getEnabledBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	totalBacking, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}

	if !totalBacking.MerMinted.IsPositive() {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
		return
	}

	excessBackingValue, err := k.getExcessBackingValue(ctx)
	if err != nil {
		return
	}

	buybackFeeRate := sdk.ZeroDec()
	if backingParams.MintFee != nil {
		buybackFeeRate = *backingParams.BuybackFee
	}
	backingOutTotalValue := backingOut.Amount.ToDec().Mul(backingPrice).Quo(sdk.OneDec().Sub(buybackFeeRate))

	if backingOutTotalValue.GT(excessBackingValue.Amount.ToDec()) {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
		return
	}
	if backingOutTotalValue.GT(poolBacking.Backing.Amount.ToDec()) {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
		return
	}

	lionIn = sdk.NewCoin(merlion.AttoLionDenom, backingOutTotalValue.Quo(lionPrice).RoundInt())
	buybackFee = sdk.NewCoin(backingDenom, backingOutTotalValue.Mul(buybackFeeRate).RoundInt())
	return
}

func (k Keeper) estimateBuyBackingOut(
	ctx sdk.Context,
	lionIn sdk.Coin,
	backingDenom string,
) (
	backingOut sdk.Coin,
	buybackFee sdk.Coin,
	err error,
) {
	backingOut = sdk.NewCoin(backingDenom, sdk.ZeroInt())
	buybackFee = sdk.NewCoin(backingDenom, sdk.ZeroInt())

	backingParams, err := k.getEnabledBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	totalBacking, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}
	if !totalBacking.MerMinted.IsPositive() {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
		return
	}

	excessBackingValue, err := k.getExcessBackingValue(ctx)
	if err != nil {
		return
	}

	lionInValue := lionIn.Amount.ToDec().Mul(lionPrice)
	if lionInValue.GT(excessBackingValue.Amount.ToDec()) {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
		return
	}

	backingOutTotal := sdk.NewCoin(backingDenom, lionInValue.Quo(backingPrice).RoundInt())
	if poolBacking.Backing.IsLT(backingOutTotal) {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
		return
	}

	buybackFee = computeFee(backingOutTotal, backingParams.BuybackFee)
	backingOut = backingOutTotal.Sub(buybackFee)
	return
}

func (k Keeper) estimateSellBackingIn(
	ctx sdk.Context,
	lionOut sdk.Coin,
	backingDenom string,
) (
	backingIn sdk.Coin,
	sellbackFee sdk.Coin,
	err error,
) {
	backingIn = sdk.NewCoin(backingDenom, sdk.ZeroInt())
	sellbackFee = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())

	backingParams, err := k.getEnabledBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}

	excessBackingValue, err := k.getExcessBackingValue(ctx)
	if err != nil {
		return
	}
	missingBackingValue := excessBackingValue.Amount.Neg()
	availableLionMint := missingBackingValue.ToDec().Quo(lionPrice)

	bonusRatio := k.RebackBonus(ctx)
	sellbackFeeRate := sdk.ZeroDec()
	if backingParams.RebackFee != nil {
		sellbackFeeRate = *backingParams.RebackFee
	}

	lionMint := lionOut.Amount.ToDec().Quo(sdk.OneDec().Add(bonusRatio).Sub(sellbackFeeRate))
	lionMintWithBonus := lionMint.Mul(sdk.OneDec().Add(bonusRatio))

	backingIn = sdk.NewCoin(backingDenom, lionMint.Mul(lionPrice).Quo(backingPrice).RoundInt())
	sellbackFee = sdk.NewCoin(merlion.AttoLionDenom, lionMint.Mul(sellbackFeeRate).RoundInt())

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		err = sdkerrors.Wrap(types.ErrBackingCeiling, "backing over ceiling")
		return
	}
	if lionMintWithBonus.GT(availableLionMint) {
		err = sdkerrors.Wrap(types.ErrLionCoinInsufficient, "insufficient available lion coin")
		return
	}

	return
}

func (k Keeper) estimateSellBackingOut(
	ctx sdk.Context,
	backingIn sdk.Coin,
) (
	lionOut sdk.Coin,
	sellbackFee sdk.Coin,
	err error,
) {
	backingDenom := backingIn.Denom

	lionOut = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	sellbackFee = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())

	backingParams, err := k.getEnabledBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		err = sdkerrors.Wrap(types.ErrBackingCeiling, "backing over ceiling")
		return
	}

	excessBackingValue, err := k.getExcessBackingValue(ctx)
	if err != nil {
		return
	}
	missingBackingValue := excessBackingValue.Amount.Neg()
	availableLionMint := missingBackingValue.ToDec().Quo(lionPrice)

	bonusRatio := k.RebackBonus(ctx)
	lionMint := sdk.NewCoin(merlion.AttoLionDenom, backingIn.Amount.ToDec().Mul(backingPrice).Quo(lionPrice).RoundInt())
	bonus := computeFee(lionMint, &bonusRatio)
	sellbackFee = computeFee(lionMint, backingParams.RebackFee)

	lionMintWithBonus := lionMint.Add(bonus)
	if lionMintWithBonus.Amount.ToDec().GT(availableLionMint) {
		err = sdkerrors.Wrap(types.ErrLionCoinInsufficient, "insufficient available lion coin")
		return
	}

	lionOut = lionMintWithBonus.Sub(sellbackFee)
	return
}

func (k Keeper) estimateMintByCollateralIn(
	ctx sdk.Context,
	account sdk.AccAddress,
	mintOut sdk.Coin,
	collateralDenom string,
	ltv sdk.Dec,
) (
	collateralIn sdk.Coin,
	lionIn sdk.Coin,
	mintFee sdk.Coin,
	totalColl types.TotalCollateral,
	poolColl types.PoolCollateral,
	accColl types.AccountCollateral,
	err error,
) {
	collateralIn = sdk.NewCoin(collateralDenom, sdk.ZeroInt())
	lionIn = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	mintFee = sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())

	// check price lower bound ?
	//err = k.checkMerPriceLowerBound(ctx)
	//if err != nil {
	//	return
	//}

	collateralParams, err := k.getEnabledCollateralParams(ctx, collateralDenom)
	if err != nil {
		return
	}

	err = checkLTV(ltv, &collateralParams)
	if err != nil {
		return
	}

	// get prices in usd
	collateralPrice, err := k.oracleKeeper.GetExchangeRate(ctx, collateralDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	totalColl, poolColl, accColl, err = k.getCollateral(ctx, account, collateralDenom)
	if err != nil {
		return
	}

	// settle interest fee
	settleInterestFee(ctx, &accColl, &poolColl, &totalColl, collateralParams.InterestFee)

	// compute mint total amount
	mintFee = computeFee(mintOut, collateralParams.MintFee)
	mintTotal := mintOut.Add(mintFee)

	// check mer ceiling
	if collateralParams.MaxMerMint != nil && poolColl.MerDebt.Amount.Add(poolColl.MerByLion.Amount).Add(mintTotal.Amount).GT(*collateralParams.MaxMerMint) {
		err = sdkerrors.Wrapf(types.ErrMerCeiling, "mer over ceiling")
		return
	}

	// compute MerByLion and lion in
	merDue := accColl.MerDebt.Add(accColl.MerByLion).Add(mintTotal)
	catalyticLionRatio := ltv.Sub(*collateralParams.BasicLoanToValue).Quo(collateralParams.LoanToValue.Sub(*collateralParams.BasicLoanToValue)).Mul(*collateralParams.CatalyticLionRatio)
	minAccMerByLion := merDue.Amount.ToDec().Mul(catalyticLionRatio).RoundInt()

	addedMerByLion := sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())
	if accColl.MerByLion.Amount.LT(minAccMerByLion) {
		addedMerByLion.Amount = minAccMerByLion.Sub(accColl.MerByLion.Amount)
	}

	accColl.MerByLion = accColl.MerByLion.Add(addedMerByLion)
	poolColl.MerByLion = poolColl.MerByLion.Add(addedMerByLion)
	totalColl.MerByLion = totalColl.MerByLion.Add(addedMerByLion)

	lionIn = sdk.NewCoin(merlion.AttoLionDenom, addedMerByLion.Amount.ToDec().Quo(lionPrice).RoundInt())
	accColl.LionBurned = accColl.LionBurned.Add(lionIn)
	poolColl.LionBurned = poolColl.LionBurned.Add(lionIn)
	totalColl.LionBurned = totalColl.LionBurned.Add(lionIn)

	// compute MerDebt and collateral in
	addedMerDebt := mintTotal.Sub(addedMerByLion)

	accColl.MerDebt = accColl.MerDebt.Add(addedMerDebt)
	poolColl.MerDebt = poolColl.MerDebt.Add(addedMerDebt)
	totalColl.MerDebt = totalColl.MerDebt.Add(addedMerDebt)

	minAccCollateral := accColl.MerDebt.Amount.ToDec().Quo(ltv).Quo(collateralPrice).RoundInt()
	if accColl.Collateral.Amount.LT(minAccCollateral) {
		collateralIn.Amount = minAccCollateral.Sub(accColl.Collateral.Amount)
	}

	accColl.Collateral = accColl.Collateral.Add(collateralIn)
	poolColl.Collateral = poolColl.Collateral.Add(collateralIn)

	if collateralParams.MaxCollateral != nil && poolColl.Collateral.Amount.GT(*collateralParams.MaxCollateral) {
		err = sdkerrors.Wrap(types.ErrCollateralCeiling, "collateral over ceiling")
		return
	}

	return
}

func (k Keeper) estimateMintByCollateralOut(
	ctx sdk.Context,
	account sdk.AccAddress,
	collateralIn sdk.Coin,
	ltv sdk.Dec,
) (
	lionIn sdk.Coin,
	mintOut sdk.Coin,
	mintFee sdk.Coin,
	totalColl types.TotalCollateral,
	poolColl types.PoolCollateral,
	accColl types.AccountCollateral,
	err error,
) {
	collateralDenom := collateralIn.Denom

	lionIn = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	mintOut = sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())
	mintFee = sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())

	// check price lower bound ?
	//err = k.checkMerPriceLowerBound(ctx)
	//if err != nil {
	//	return
	//}

	collateralParams, err := k.getEnabledCollateralParams(ctx, collateralDenom)
	if err != nil {
		return
	}

	err = checkLTV(ltv, &collateralParams)
	if err != nil {
		return
	}

	// get prices in usd
	collateralPrice, err := k.oracleKeeper.GetExchangeRate(ctx, collateralDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	totalColl, poolColl, accColl, err = k.getCollateral(ctx, account, collateralDenom)
	if err != nil {
		return
	}

	// update Collateral
	accColl.Collateral = accColl.Collateral.Add(collateralIn)
	poolColl.Collateral = poolColl.Collateral.Add(collateralIn)

	if collateralParams.MaxCollateral != nil && poolColl.Collateral.Amount.GT(*collateralParams.MaxCollateral) {
		err = sdkerrors.Wrap(types.ErrCollateralCeiling, "collateral over ceiling")
		return
	}

	// settle interest fee
	settleInterestFee(ctx, &accColl, &poolColl, &totalColl, collateralParams.InterestFee)

	// compute MerDebt
	maxAccMerDebt := accColl.Collateral.Amount.ToDec().Mul(collateralPrice).Mul(ltv).RoundInt()
	addedMerDebt := sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())
	if accColl.MerDebt.Amount.LT(maxAccMerDebt) {
		addedMerDebt.Amount = maxAccMerDebt.Sub(accColl.MerDebt.Amount)
	}

	accColl.MerDebt = accColl.MerDebt.Add(addedMerDebt)
	poolColl.MerDebt = poolColl.MerDebt.Add(addedMerDebt)
	totalColl.MerDebt = totalColl.MerDebt.Add(addedMerDebt)

	// compute MerByLion and lion in
	catalyticLionRatio := ltv.Sub(*collateralParams.BasicLoanToValue).Quo(collateralParams.LoanToValue.Sub(*collateralParams.BasicLoanToValue)).Mul(*collateralParams.CatalyticLionRatio)
	minAccMerByLion := catalyticLionRatio.MulInt(accColl.MerDebt.Amount).Quo(sdk.OneDec().Sub(catalyticLionRatio)).RoundInt()

	addedMerByLion := sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())
	if accColl.MerByLion.Amount.LT(minAccMerByLion) {
		addedMerByLion.Amount = minAccMerByLion.Sub(accColl.MerByLion.Amount)
	}

	accColl.MerByLion = accColl.MerByLion.Add(addedMerByLion)
	poolColl.MerByLion = poolColl.MerByLion.Add(addedMerByLion)
	totalColl.MerByLion = totalColl.MerByLion.Add(addedMerByLion)

	lionIn = sdk.NewCoin(merlion.AttoLionDenom, addedMerByLion.Amount.ToDec().Quo(lionPrice).RoundInt())
	accColl.LionBurned = accColl.LionBurned.Add(lionIn)
	poolColl.LionBurned = poolColl.LionBurned.Add(lionIn)
	totalColl.LionBurned = totalColl.LionBurned.Add(lionIn)

	// check mer ceiling
	if collateralParams.MaxMerMint != nil && poolColl.MerDebt.Amount.Add(poolColl.MerByLion.Amount).GT(*collateralParams.MaxMerMint) {
		err = sdkerrors.Wrapf(types.ErrMerCeiling, "mer over ceiling")
		return
	}

	// compute mint out
	mintTotal := addedMerDebt.Add(addedMerByLion)
	mintOut = sdk.NewCoin(merlion.MicroUSDDenom, mintTotal.Amount.ToDec().Quo(sdk.OneDec().Add(*collateralParams.MintFee)).RoundInt())
	mintFee = sdk.NewCoin(merlion.MicroUSDDenom, mintOut.Amount.ToDec().Mul(*collateralParams.MintFee).RoundInt())

	return
}

func computeFee(coin sdk.Coin, rate *sdk.Dec) sdk.Coin {
	amt := sdk.ZeroInt()
	if rate != nil {
		amt = coin.Amount.ToDec().Mul(*rate).RoundInt()
	}
	return sdk.NewCoin(coin.Denom, amt)
}

func (k Keeper) checkMerPriceLowerBound(ctx sdk.Context) error {
	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return err
	}
	merPriceLowerBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Sub(k.MintPriceBias(ctx)))
	if merPrice.LT(merPriceLowerBound) {
		return sdkerrors.Wrapf(types.ErrMerPriceTooLow, "%s price too low: %s", merlion.MicroUSDDenom, merPrice)
	}
	return nil
}

func (k Keeper) checkMerPriceUpperBound(ctx sdk.Context) error {
	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return err
	}
	merPriceUpperBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Add(k.BurnPriceBias(ctx)))
	if merPrice.GT(merPriceUpperBound) {
		return sdkerrors.Wrapf(types.ErrMerPriceTooHigh, "%s price too high: %s", merlion.MicroUSDDenom, merPrice)
	}
	return nil
}

func (k Keeper) getEnabledBackingParams(ctx sdk.Context, backingDenom string) (backingParams types.BackingRiskParams, err error) {
	backingParams, found := k.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
		return
	}
	if !backingParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
		return
	}
	return
}

func (k Keeper) getEnabledCollateralParams(ctx sdk.Context, collateralDenom string) (collateralParams types.CollateralRiskParams, err error) {
	collateralParams, found := k.GetCollateralRiskParams(ctx, collateralDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", collateralDenom)
		return
	}
	if !collateralParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrCollateralCoinDisabled, "collateral coin disabled: %s", collateralDenom)
		return
	}
	return
}

func (k Keeper) getExcessBackingValue(ctx sdk.Context) (excessBackingValue sdk.Coin, err error) {
	excessBackingValue = sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())

	totalBacking, found := k.GetTotalBacking(ctx)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "total backing not found")
		return
	}

	backingRatio := k.GetBackingRatio(ctx)
	requiredBackingValue := totalBacking.MerMinted.Amount.ToDec().Mul(backingRatio).RoundInt()

	totalBackingValue, err := k.totalBackingInUSD(ctx)
	if err != nil {
		return
	}

	// may be negative
	excessBackingValue.Amount = totalBackingValue.Amount.Sub(requiredBackingValue)
	return
}

func (k Keeper) totalBackingInUSD(ctx sdk.Context) (sdk.Coin, error) {
	totalBackingValue := sdk.ZeroDec()
	for _, pool := range k.GetAllPoolBacking(ctx) {
		// get price in usd
		backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, pool.Backing.Denom)
		if err != nil {
			return sdk.Coin{}, err
		}
		totalBackingValue = totalBackingValue.Add(pool.Backing.Amount.ToDec().Mul(backingPrice))
	}
	return sdk.NewCoin(merlion.MicroUSDDenom, totalBackingValue.RoundInt()), nil
}

func checkLTV(ltv sdk.Dec, collateralParams *types.CollateralRiskParams) error {
	if collateralParams.BasicLoanToValue != nil && ltv.LT(*collateralParams.BasicLoanToValue) {
		return sdkerrors.Wrapf(types.ErrLTVOutOfRange, "%s < %s", ltv, *collateralParams.BasicLoanToValue)
	}
	if collateralParams.LoanToValue != nil && ltv.GT(*collateralParams.LoanToValue) {
		return sdkerrors.Wrapf(types.ErrLTVOutOfRange, "%s > %s", ltv, *collateralParams.LoanToValue)
	}
	return nil
}
