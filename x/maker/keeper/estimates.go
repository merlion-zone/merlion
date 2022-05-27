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

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}
	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return
	}

	// check price lower bound
	merPriceLowerBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Sub(k.MintPriceBias(ctx)))
	if merPrice.LT(merPriceLowerBound) {
		err = sdkerrors.Wrapf(types.ErrMerPriceTooLow, "%s price too low: %s", merlion.MicroUSDDenom, merPrice)
		return
	}

	backingParams, found := k.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
		return
	}
	if !backingParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
		return
	}

	mintFee = computeFee(mintOut, backingParams.MintFee)
	mintTotal := mintOut.Add(mintFee)

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}
	poolBacking.MerMinted = poolBacking.MerMinted.Add(mintTotal)
	if backingParams.MaxMerMint != nil && poolBacking.MerMinted.Amount.GT(*backingParams.MaxMerMint) {
		err = sdkerrors.Wrapf(types.ErrMerCeiling, "mer over ceiling")
		return
	}

	mintTotalInUSD := mintTotal.Amount.ToDec().Mul(merlion.MicroUSDTarget)

	backingRatio := k.GetBackingRatio(ctx)
	if backingRatio.GTE(sdk.OneDec()) || fullBacking {
		// full/over backed, or user selects full backing
		backingIn.Amount = mintTotalInUSD.QuoRoundUp(backingPrice).RoundInt()
	} else if backingRatio.IsZero() {
		// algorithmic
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

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}
	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return
	}

	// check price lower bound
	merPriceLowerBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Sub(k.MintPriceBias(ctx)))
	if merPrice.LT(merPriceLowerBound) {
		err = sdkerrors.Wrapf(types.ErrMerPriceTooLow, "%s price too low: %s", merlion.MicroUSDDenom, merPrice)
		return
	}

	backingParams, found := k.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
		return
	}
	if !backingParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
		return
	}

	backingAvailInUSD := backingPrice.MulInt(backingInMax.Amount)
	lionAvailInUSD := lionPrice.MulInt(lionInMax.Amount)

	mintTotalInUSD := sdk.ZeroDec()

	backingRatio := k.GetBackingRatio(ctx)
	if backingRatio.GTE(sdk.OneDec()) || lionInMax.IsZero() {
		// full/over backed, or user selects full backing
		mintTotalInUSD = backingAvailInUSD
		backingIn.Amount = backingInMax.Amount
	} else if backingRatio.IsZero() {
		// algorithmic
		mintTotalInUSD = lionAvailInUSD
		lionIn.Amount = lionInMax.Amount
	} else {
		// fractional
		mintTotalWithBackingInUSD := backingAvailInUSD.QuoRoundUp(backingRatio)
		mintTotalWithLionInUSD := lionAvailInUSD.QuoRoundUp(sdk.OneDec().Sub(backingRatio))
		if mintTotalWithBackingInUSD.LT(mintTotalWithLionInUSD) {
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
	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		err = sdkerrors.Wrapf(types.ErrBackingCeiling, "backing over ceiling")
		return
	}

	mintTotalValue := mintTotalInUSD.QuoRoundUp(merPrice)
	poolBacking.MerMinted = poolBacking.MerMinted.AddAmount(mintTotalValue.RoundInt())
	if backingParams.MaxMerMint != nil && poolBacking.MerMinted.Amount.GT(*backingParams.MaxMerMint) {
		err = sdkerrors.Wrapf(types.ErrMerCeiling, "mer over ceiling")
		return
	}

	mintFeeRate := sdk.ZeroDec()
	if backingParams.MintFee != nil {
		mintFeeRate = *backingParams.MintFee
	}
	mintOutValue := mintTotalValue.QuoRoundUp(sdk.OneDec().Add(mintFeeRate))
	mintFeeValue := mintOutValue.Mul(mintFeeRate)

	mintOut = sdk.NewCoin(merlion.MicroUSDDenom, mintOutValue.RoundInt())
	mintFee = sdk.NewCoin(merlion.MicroUSDDenom, mintFeeValue.RoundInt())
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

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}
	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return
	}

	// check price upper bound
	merPriceUpperBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Add(k.BurnPriceBias(ctx)))
	if merPrice.GT(merPriceUpperBound) {
		err = sdkerrors.Wrapf(types.ErrMerPriceTooHigh, "%s price too high: %s", merlion.MicroUSDDenom, merPrice)
		return
	}

	backingParams, found := k.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
		return
	}
	if !backingParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
		return
	}

	burnFee = computeFee(burnIn, backingParams.BurnFee)
	burnActual := burnIn.Sub(burnFee)
	burnActualInUSD := burnActual.Amount.ToDec().Mul(merlion.MicroUSDTarget)

	backingRatio := k.GetBackingRatio(ctx)
	if backingRatio.GTE(sdk.OneDec()) {
		// full/over collateralized
		backingOut.Amount = burnActualInUSD.QuoRoundUp(backingPrice).RoundInt()
	} else if backingRatio.IsZero() {
		// algorithmic
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

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	backingParams, found := k.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
		return
	}
	if !backingParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
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

	backingRatio := k.GetBackingRatio(ctx)
	requiredBackingValue := totalBacking.MerMinted.Amount.ToDec().Mul(backingRatio).TruncateInt()

	totalBackingValue, err := k.totalBackingInUSD(ctx)
	if err != nil {
		return
	}

	availableExcessBackingValue := sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())
	if requiredBackingValue.LT(totalBackingValue.Amount) {
		availableExcessBackingValue.Amount = totalBackingValue.Amount.Sub(requiredBackingValue)
	}

	lionInValue := lionIn.Amount.ToDec().Mul(lionPrice)
	if lionInValue.TruncateInt().GT(availableExcessBackingValue.Amount) {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
		return
	}

	backingOut = sdk.NewCoin(backingDenom, lionInValue.Quo(backingPrice).TruncateInt())
	if poolBacking.Backing.IsLT(backingOut) {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
		return
	}

	buybackFee = computeFee(backingOut, backingParams.BuybackFee)
	backingOut = backingOut.Sub(buybackFee)

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
	lionOut = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	sellbackFee = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())

	backingDenom := backingIn.Denom

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	backingParams, found := k.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
		return
	}
	if !backingParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
		return
	}

	totalBacking, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		err = sdkerrors.Wrap(types.ErrBackingCeiling, "over ceiling")
		return
	}

	backingRatio := k.GetBackingRatio(ctx)
	requiredBackingValue := totalBacking.MerMinted.Amount.ToDec().Mul(backingRatio).TruncateInt()

	totalBackingValue, err := k.totalBackingInUSD(ctx)
	if err != nil {
		return
	}

	availableMissingBackingValue := sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())
	if requiredBackingValue.GT(totalBackingValue.Amount) {
		availableMissingBackingValue.Amount = requiredBackingValue.Sub(totalBackingValue.Amount)
	}
	availableLionOut := availableMissingBackingValue.Amount.ToDec().Quo(lionPrice)

	bonusRatio := k.RecollateralizeBonus(ctx)
	lionMint := sdk.NewCoin(merlion.AttoLionDenom, backingIn.Amount.ToDec().Mul(backingPrice).Quo(lionPrice).TruncateInt())
	bonus := computeFee(lionMint, &bonusRatio)
	sellbackFee = computeFee(lionMint, backingParams.RecollateralizeFee)

	lionMint = lionMint.Add(bonus)
	if lionMint.Amount.ToDec().GT(availableLionOut) {
		err = sdkerrors.Wrap(types.ErrLionCoinInsufficient, "insufficient available lion coin")
		return
	}

	lionOut = lionMint.Sub(sellbackFee)
	return
}

func (k Keeper) estimateMintByCollateralIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	mintOut sdk.Coin,
	collateralDenom string,
	lionInMax sdk.Coin,
) (
	lionIn sdk.Coin,
	mintFee sdk.Coin,
	totalColl types.TotalCollateral,
	poolColl types.PoolCollateral,
	accColl types.AccountCollateral,
	err error,
) {
	lionIn = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	mintFee = sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())

	// get prices in usd
	collateralPrice, err := k.oracleKeeper.GetExchangeRate(ctx, collateralDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}
	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return
	}

	// check price lower bound
	merPriceLowerBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Sub(k.MintPriceBias(ctx)))
	if merPrice.LT(merPriceLowerBound) {
		err = sdkerrors.Wrapf(types.ErrMerPriceTooLow, "%s price too low: %s", merlion.MicroUSDDenom, merPrice)
		return
	}

	collateralParams, found := k.GetCollateralRiskParams(ctx, collateralDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", collateralDenom)
		return
	}
	if !collateralParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrCollateralCoinDisabled, "collateral coin disabled: %s", collateralDenom)
		return
	}

	totalColl, poolColl, accColl, err = k.getCollateral(ctx, sender, collateralDenom)
	if err != nil {
		return
	}

	// settle interestFee fee
	settleInterestFee(ctx, &accColl, &poolColl, &totalColl, collateralParams.InterestFee)

	// compute mint amount
	mintFee = computeFee(mintOut, collateralParams.MintFee)
	mint := mintOut.Add(mintFee)

	// update debt
	accColl.MerDebt = accColl.MerDebt.Add(mint)
	poolColl.MerDebt = poolColl.MerDebt.Add(mint)
	totalColl.MerDebt = totalColl.MerDebt.Add(mint)

	if collateralParams.MaxMerMint != nil && poolColl.MerDebt.Amount.GT(*collateralParams.MaxMerMint) {
		err = sdkerrors.Wrapf(types.ErrMerCeiling, "mer over ceiling")
		return
	}

	// compute actual catalytic lion
	merDue := accColl.MerDebt.Add(accColl.MerByLion)
	bestCatalyticLionInUSD := merDue.Amount.ToDec().Mul(*collateralParams.CatalyticLionRatio)
	lionInMaxInUSD := lionInMax.Amount.ToDec().Mul(lionPrice).TruncateInt()
	catalyticLionInUSD := sdk.MinDec(bestCatalyticLionInUSD, accColl.MerByLion.Amount.Add(lionInMaxInUSD).ToDec()).TruncateInt()

	// compute actual lion-in
	lionInInUSD := catalyticLionInUSD.Sub(accColl.MerByLion.Amount)
	if !lionInInUSD.IsPositive() {
		lionInInUSD = sdk.ZeroInt()
	} else {
		accColl.MerByLion = accColl.MerByLion.AddAmount(lionInInUSD)
		poolColl.MerByLion = poolColl.MerByLion.AddAmount(lionInInUSD)
		totalColl.MerByLion = totalColl.MerByLion.AddAmount(lionInInUSD)
		accColl.MerDebt = accColl.MerDebt.SubAmount(lionInInUSD)
		poolColl.MerDebt = poolColl.MerDebt.SubAmount(lionInInUSD)
		totalColl.MerDebt = totalColl.MerDebt.SubAmount(lionInInUSD)
	}
	lionIn = sdk.NewCoin(merlion.AttoLionDenom, lionInInUSD.ToDec().Quo(lionPrice).TruncateInt())

	accColl.LionBurned = accColl.LionBurned.Add(lionIn)
	poolColl.LionBurned = poolColl.LionBurned.Add(lionIn)
	totalColl.LionBurned = totalColl.LionBurned.Add(lionIn)

	// compute actual catalytic ratio and max loan-to-value
	maxLoanToValue := maxLoanToValueForAccount(&accColl, &collateralParams)

	// check max mintable mer
	collateralInUSD := accColl.Collateral.Amount.ToDec().Mul(collateralPrice)
	maxMerMintable := collateralInUSD.Mul(maxLoanToValue)
	if maxMerMintable.LT(accColl.MerDebt.Amount.ToDec()) {
		err = sdkerrors.Wrapf(types.ErrAccountInsufficientCollateral, "account has insufficient collateral: %s", collateralDenom)
		return
	}

	return
}

func (k Keeper) estimateBurnByCollateralIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	collateralDenom string,
	repayInMax sdk.Coin,
) (
	repayIn sdk.Coin,
	interestFee sdk.Coin,
	totalColl types.TotalCollateral,
	poolColl types.PoolCollateral,
	accColl types.AccountCollateral,
	err error,
) {
	// get prices in usd
	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return
	}

	// check price upper bound
	merPriceUpperBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Add(k.BurnPriceBias(ctx)))
	if merPrice.GT(merPriceUpperBound) {
		err = sdkerrors.Wrapf(types.ErrMerPriceTooHigh, "%s price too high: %s", merlion.MicroUSDDenom, merPrice)
		return
	}

	collateralParams, found := k.GetCollateralRiskParams(ctx, collateralDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", collateralDenom)
		return
	}
	if !collateralParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrCollateralCoinDisabled, "collateral coin disabled: %s", collateralDenom)
		return
	}

	totalColl, poolColl, accColl, err = k.getCollateral(ctx, sender, collateralDenom)
	if err != nil {
		return
	}

	// settle interestFee fee
	settleInterestFee(ctx, &accColl, &poolColl, &totalColl, collateralParams.InterestFee)

	// compute burn-in
	if !accColl.MerDebt.IsPositive() {
		err = sdkerrors.Wrapf(types.ErrAccountNoDebt, "account has no debt for %s collateral", collateralDenom)
		return
	}
	repayIn = sdk.NewCoin(repayInMax.Denom, sdk.MinInt(accColl.MerDebt.Amount, repayInMax.Amount))
	interestFee = sdk.NewCoin(repayInMax.Denom, sdk.MinInt(accColl.LastInterest.Amount, repayIn.Amount))

	// update debt
	accColl.LastInterest = accColl.LastInterest.Sub(interestFee)
	accColl.MerDebt = accColl.MerDebt.Sub(repayIn)
	poolColl.MerDebt = poolColl.MerDebt.Sub(repayIn)
	totalColl.MerDebt = totalColl.MerDebt.Sub(repayIn)

	return
}
