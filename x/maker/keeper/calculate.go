package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/maker/types"
)

func (k Keeper) calculateMintBySwapIn(
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
	mintFee = sdk.NewCoin(merlion.MicroUSMDenom, sdk.ZeroInt())

	err = k.checkMintPriceLowerBound(ctx)
	if err != nil {
		return
	}

	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
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
	mintTotalInUSD := mintTotal.Amount.ToDec().Mul(merlion.MicroUSMTarget)

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

func (k Keeper) calculateMintBySwapOut(
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

	err = k.checkMintPriceLowerBound(ctx)
	if err != nil {
		return
	}

	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
	if err != nil {
		return
	}

	// get prices in uusd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	backingRatio := k.GetBackingRatio(ctx)

	backingInMaxInUSD := backingPrice.MulInt(backingInMax.Amount)
	lionInMaxInUSD := lionPrice.MulInt(lionInMax.Amount)

	mintTotalInUSD := sdk.ZeroDec()
	backingIn = sdk.NewCoin(backingDenom, sdk.ZeroInt())
	lionIn = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())

	if backingRatio.GTE(sdk.OneDec()) || fullBacking {
		// full/over backing, or user selects full backing
		mintTotalInUSD = backingInMaxInUSD
		backingIn.Amount = backingInMax.Amount
	} else if backingRatio.IsZero() {
		// full algorithmic
		mintTotalInUSD = lionInMaxInUSD
		lionIn.Amount = lionInMax.Amount
	} else {
		// fractional
		max1 := backingInMaxInUSD.Quo(backingRatio)
		max2 := lionInMaxInUSD.Quo(sdk.OneDec().Sub(backingRatio))
		if backingInMax.IsPositive() && (lionInMax.IsZero() || max1.LTE(max2)) {
			mintTotalInUSD = max1
			backingIn.Amount = backingInMax.Amount
			lionIn.Amount = mintTotalInUSD.Mul(sdk.OneDec().Sub(backingRatio)).QuoRoundUp(lionPrice).RoundInt()
			if lionInMax.IsPositive() && lionInMax.IsLT(lionIn) {
				lionIn.Amount = lionInMax.Amount
			}
		} else {
			mintTotalInUSD = max2
			lionIn.Amount = lionInMax.Amount
			backingIn.Amount = mintTotalInUSD.Mul(backingRatio).QuoRoundUp(backingPrice).RoundInt()
			if backingInMax.IsPositive() && backingInMax.IsLT(backingIn) {
				backingIn.Amount = backingInMax.Amount
			}
		}
	}

	mintTotal := sdk.NewCoin(merlion.MicroUSMDenom, mintTotalInUSD.Quo(merlion.MicroUSMTarget).TruncateInt())

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}

	poolBacking.MerMinted = poolBacking.MerMinted.AddAmount(mintTotal.Amount)
	if backingParams.MaxMerMint != nil && poolBacking.MerMinted.Amount.GT(*backingParams.MaxMerMint) {
		err = sdkerrors.Wrap(types.ErrMerCeiling, "")
		return
	}

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		err = sdkerrors.Wrap(types.ErrBackingCeiling, "")
		return
	}

	mintFee = computeFee(mintTotal, backingParams.MintFee)
	mintOut = mintTotal.Sub(mintFee)
	return
}

func (k Keeper) calculateBurnBySwapIn(
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

	burnIn = sdk.NewCoin(merlion.MicroUSMDenom, sdk.ZeroInt())
	backingOut = sdk.NewCoin(backingOutMax.Denom, sdk.ZeroInt())
	lionOut = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	burnFee = sdk.NewCoin(merlion.MicroUSMDenom, sdk.ZeroInt())

	err = k.checkBurnPriceUpperBound(ctx)
	if err != nil {
		return
	}

	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
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

	burnInValue := burnActualInUSD.Quo(merlion.MicroUSMTarget).Quo(sdk.OneDec().Sub(burnFeeRate))
	burnFeeValue := burnInValue.Mul(burnFeeRate)
	burnIn = sdk.NewCoin(merlion.MicroUSMDenom, burnInValue.RoundInt())
	burnFee = sdk.NewCoin(merlion.MicroUSMDenom, burnFeeValue.RoundInt())
	return
}

func (k Keeper) calculateBurnBySwapOut(
	ctx sdk.Context,
	burnIn sdk.Coin,
	backingDenom string,
) (
	backingOut sdk.Coin,
	lionOut sdk.Coin,
	burnFee sdk.Coin,
	err error,
) {
	err = k.checkBurnPriceUpperBound(ctx)
	if err != nil {
		return
	}

	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
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

	backingRatio := k.GetBackingRatio(ctx)

	burnFee = computeFee(burnIn, backingParams.BurnFee)
	burnActual := burnIn.Sub(burnFee)
	burnActualInUSD := burnActual.Amount.ToDec().Mul(merlion.MicroUSMTarget)

	backingOut = sdk.NewCoin(backingDenom, sdk.ZeroInt())
	lionOut = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())

	if backingRatio.GTE(sdk.OneDec()) {
		// full/over backing
		backingOut.Amount = burnActualInUSD.QuoTruncate(backingPrice).TruncateInt()
	} else if backingRatio.IsZero() {
		// full algorithmic
		lionOut.Amount = burnActualInUSD.QuoTruncate(lionPrice).TruncateInt()
	} else {
		// fractional
		backingOut.Amount = burnActualInUSD.Mul(backingRatio).QuoTruncate(backingPrice).TruncateInt()
		lionOut.Amount = burnActualInUSD.Mul(sdk.OneDec().Sub(backingRatio)).QuoTruncate(lionPrice).TruncateInt()
	}

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}
	moduleOwnedBacking := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), backingDenom)

	poolBackingBalance := sdk.NewCoin(backingDenom, sdk.MinInt(poolBacking.Backing.Amount, moduleOwnedBacking.Amount))
	if poolBackingBalance.IsLT(backingOut) {
		err = sdkerrors.Wrapf(types.ErrBackingCoinInsufficient, "backing coin out(%s) > balance(%s)", backingOut, poolBackingBalance)
		return
	}

	return
}

func (k Keeper) calculateBuyBackingIn(
	ctx sdk.Context,
	backingOut sdk.Coin,
) (
	lionIn sdk.Coin,
	buybackFee sdk.Coin,
	err error,
) {
	backingDenom := backingOut.Denom

	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
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

	excessBackingValue, err := k.getExcessBackingValue(ctx)
	if err != nil {
		return
	}

	backingOutTotal := sdk.NewCoin(backingDenom, backingOut.Amount.ToDec().Quo(sdk.OneDec().Sub(*backingParams.BuybackFee)).TruncateInt())
	lionInValue := backingOutTotal.Amount.ToDec().Mul(backingPrice)

	if lionInValue.GT(excessBackingValue.ToDec()) {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "")
		return
	}

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}
	moduleOwnedBacking := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), backingDenom)

	poolBackingBalance := sdk.NewCoin(backingDenom, sdk.MinInt(poolBacking.Backing.Amount, moduleOwnedBacking.Amount))
	if poolBackingBalance.IsLT(backingOutTotal) {
		err = sdkerrors.Wrapf(types.ErrBackingCoinInsufficient, "backing coin out(%s) > balance(%s)", backingOutTotal, poolBackingBalance)
		return
	}

	lionIn = sdk.NewCoin(merlion.AttoLionDenom, lionInValue.Quo(lionPrice).RoundInt())
	buybackFee = sdk.NewCoin(backingDenom, backingOutTotal.Amount.ToDec().Mul(*backingParams.BuybackFee).RoundInt())
	return
}

func (k Keeper) calculateBuyBackingOut(
	ctx sdk.Context,
	lionIn sdk.Coin,
	backingDenom string,
) (
	backingOut sdk.Coin,
	buybackFee sdk.Coin,
	err error,
) {
	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
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

	excessBackingValue, err := k.getExcessBackingValue(ctx)
	if err != nil {
		return
	}

	lionInValue := lionIn.Amount.ToDec().Mul(lionPrice)
	if lionInValue.GT(excessBackingValue.ToDec()) {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "")
		return
	}

	backingOutTotal := sdk.NewCoin(backingDenom, lionInValue.Quo(backingPrice).TruncateInt())

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}
	moduleOwnedBacking := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), backingDenom)

	poolBackingBalance := sdk.NewCoin(backingDenom, sdk.MinInt(poolBacking.Backing.Amount, moduleOwnedBacking.Amount))
	if poolBackingBalance.IsLT(backingOutTotal) {
		err = sdkerrors.Wrapf(types.ErrBackingCoinInsufficient, "backing coin out(%s) > balance(%s)", backingOutTotal, poolBackingBalance)
		return
	}

	buybackFee = computeFee(backingOutTotal, backingParams.BuybackFee)
	backingOut = backingOutTotal.Sub(buybackFee)
	return
}

func (k Keeper) calculateSellBackingIn(
	ctx sdk.Context,
	lionOut sdk.Coin,
	backingDenom string,
) (
	backingIn sdk.Coin,
	rebackFee sdk.Coin,
	err error,
) {
	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
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
	missingBackingValue := excessBackingValue.Neg()
	availableLionMint := missingBackingValue.ToDec().Quo(lionPrice)

	bonusRatio := k.RebackBonus(ctx)

	lionMint := lionOut.Amount.ToDec().Quo(sdk.OneDec().Add(bonusRatio).Sub(*backingParams.RebackFee))

	backingIn = sdk.NewCoin(backingDenom, lionMint.Mul(lionPrice).Quo(backingPrice).RoundInt())
	rebackFee = sdk.NewCoin(merlion.AttoLionDenom, lionMint.Mul(*backingParams.RebackFee).RoundInt())

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		err = sdkerrors.Wrap(types.ErrBackingCeiling, "")
		return
	}
	if lionMint.GT(availableLionMint) {
		err = sdkerrors.Wrap(types.ErrLionCoinInsufficient, "")
		return
	}

	return
}

func (k Keeper) calculateSellBackingOut(
	ctx sdk.Context,
	backingIn sdk.Coin,
) (
	lionOut sdk.Coin,
	rebackFee sdk.Coin,
	err error,
) {
	backingDenom := backingIn.Denom

	backingParams, err := k.getAvailableBackingParams(ctx, backingDenom)
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
		err = sdkerrors.Wrap(types.ErrBackingCeiling, "")
		return
	}

	excessBackingValue, err := k.getExcessBackingValue(ctx)
	if err != nil {
		return
	}
	missingBackingValue := excessBackingValue.Neg()
	availableLionMint := missingBackingValue.ToDec().Quo(lionPrice)

	bonusRatio := k.RebackBonus(ctx)
	lionMint := sdk.NewCoin(merlion.AttoLionDenom, backingIn.Amount.ToDec().Mul(backingPrice).Quo(lionPrice).TruncateInt())
	bonus := computeFee(lionMint, &bonusRatio)
	rebackFee = computeFee(lionMint, backingParams.RebackFee)

	if lionMint.Amount.ToDec().GT(availableLionMint) {
		err = sdkerrors.Wrap(types.ErrLionCoinInsufficient, "")
		return
	}

	lionOut = lionMint.Add(bonus).Sub(rebackFee)
	return
}

func (k Keeper) calculateMintByCollateral(
	ctx sdk.Context,
	account sdk.AccAddress,
	collateralDenom string,
	mintOut sdk.Coin,
) (
	mintFee sdk.Coin,
	totalColl types.TotalCollateral,
	poolColl types.PoolCollateral,
	accColl types.AccountCollateral,
	err error,
) {
	collateralParams, err := k.getAvailableCollateralParams(ctx, collateralDenom)
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
	settleInterestFee(ctx, &accColl, &poolColl, &totalColl, *collateralParams.InterestFee)

	// compute mint total
	mintFee = computeFee(mintOut, collateralParams.MintFee)
	mintTotal := mintOut.Add(mintFee)

	// update mer debt
	accColl.MerDebt = accColl.MerDebt.Add(mintTotal)
	poolColl.MerDebt = poolColl.MerDebt.Add(mintTotal)
	totalColl.MerDebt = totalColl.MerDebt.Add(mintTotal)

	if collateralParams.MaxMerMint != nil && poolColl.MerDebt.Amount.GT(*collateralParams.MaxMerMint) {
		err = sdkerrors.Wrapf(types.ErrMerCeiling, "")
		return
	}

	collateralValue := accColl.Collateral.Amount.ToDec().Mul(collateralPrice)
	lionCollateralizedValue := accColl.LionCollateralized.Amount.ToDec().Mul(lionPrice)
	if !collateralValue.IsPositive() {
		err = sdkerrors.Wrapf(types.ErrAccountInsufficientCollateral, "")
		return
	}

	actualCatalyticRatio := sdk.MinDec(lionCollateralizedValue.Quo(collateralValue), *collateralParams.CatalyticLionRatio)

	// actualCatalyticRatio / catalyticRatio = (availableLTV - basicLTV) / (maxLTV - basicLTV)
	availableLTV := *collateralParams.BasicLoanToValue
	if collateralParams.CatalyticLionRatio.IsPositive() {
		availableLTV = availableLTV.Add(actualCatalyticRatio.Mul(collateralParams.LoanToValue.Sub(*collateralParams.BasicLoanToValue)).Quo(*collateralParams.CatalyticLionRatio))
	}
	availableDebtMax := collateralValue.Mul(availableLTV).Quo(merlion.MicroUSMTarget).TruncateInt()

	if availableDebtMax.LT(accColl.MerDebt.Amount) {
		err = sdkerrors.Wrapf(types.ErrAccountInsufficientCollateral, "")
		return
	}

	return
}

func computeFee(coin sdk.Coin, rate *sdk.Dec) sdk.Coin {
	amt := sdk.ZeroInt()
	if rate != nil {
		amt = coin.Amount.ToDec().Mul(*rate).RoundInt()
	}
	return sdk.NewCoin(coin.Denom, amt)
}

func (k Keeper) checkMintPriceLowerBound(ctx sdk.Context) error {
	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSMDenom)
	if err != nil {
		return err
	}
	// market price must be >= target price + mint bias
	mintPriceLowerBound := merlion.MicroUSMTarget.Mul(sdk.OneDec().Add(k.MintPriceBias(ctx)))
	if merPrice.LT(mintPriceLowerBound) {
		return sdkerrors.Wrapf(types.ErrMerPriceTooLow, "%s price too low: %s", merlion.MicroUSMDenom, merPrice)
	}
	return nil
}

func (k Keeper) checkBurnPriceUpperBound(ctx sdk.Context) error {
	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSMDenom)
	if err != nil {
		return err
	}
	// market price must be <= target price - burn bias
	burnPriceUpperBound := merlion.MicroUSMTarget.Mul(sdk.OneDec().Sub(k.BurnPriceBias(ctx)))
	if merPrice.GT(burnPriceUpperBound) {
		return sdkerrors.Wrapf(types.ErrMerPriceTooHigh, "%s price too high: %s", merlion.MicroUSMDenom, merPrice)
	}
	return nil
}

func (k Keeper) getAvailableBackingParams(ctx sdk.Context, backingDenom string) (backingParams types.BackingRiskParams, err error) {
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

func (k Keeper) getAvailableCollateralParams(ctx sdk.Context, collateralDenom string) (collateralParams types.CollateralRiskParams, err error) {
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

func (k Keeper) getExcessBackingValue(ctx sdk.Context) (excessBackingValue sdk.Int, err error) {
	totalBacking, found := k.GetTotalBacking(ctx)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "total backing not found")
		return
	}

	backingRatio := k.GetBackingRatio(ctx)
	requiredBackingValue := totalBacking.MerMinted.Amount.ToDec().Mul(backingRatio).Ceil().TruncateInt()
	if requiredBackingValue.IsNegative() {
		requiredBackingValue = sdk.ZeroInt()
	}

	totalBackingValue, err := k.totalBackingInUSD(ctx)
	if err != nil {
		return
	}

	// may be negative
	excessBackingValue = totalBackingValue.Sub(requiredBackingValue)
	return
}

func (k Keeper) totalBackingInUSD(ctx sdk.Context) (sdk.Int, error) {
	totalBackingValue := sdk.ZeroDec()
	for _, pool := range k.GetAllPoolBacking(ctx) {
		// get price in usd
		backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, pool.Backing.Denom)
		if err != nil {
			return sdk.Int{}, err
		}
		totalBackingValue = totalBackingValue.Add(pool.Backing.Amount.ToDec().Mul(backingPrice))
	}
	return totalBackingValue.TruncateInt(), nil
}
