package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/maker/types"
)

func HandleRegisterBackingProposal(ctx sdk.Context, k Keeper, p *types.RegisterBackingProposal) error {
	params := p.RiskParams

	if k.IsBackingRegistered(ctx, params.BackingDenom) {
		return sdkerrors.Wrapf(types.ErrBackingCoinAlreadyExists, "backing coin denomination already registered: %s", params.BackingDenom)
	}
	if !k.bankKeeper.HasSupply(ctx, params.BackingDenom) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidCoins,
			"denomination '%s' cannot have a supply of 0", params.BackingDenom,
		)
	}

	// assign missing fields with default value
	if params.MintFee == nil {
		dec := sdk.ZeroDec()
		params.MintFee = &dec
	}
	if params.BurnFee == nil {
		dec := sdk.ZeroDec()
		params.BurnFee = &dec
	}
	if params.BuybackFee == nil {
		dec := sdk.ZeroDec()
		params.BuybackFee = &dec
	}
	if params.RebackFee == nil {
		dec := sdk.ZeroDec()
		params.RebackFee = &dec
	}

	if err := validateBackingRiskParams(ctx, k, &params); err != nil {
		return err
	}

	k.SetBackingRiskParams(ctx, params)

	_, found := k.GetTotalBacking(ctx)
	if !found {
		k.SetTotalBacking(ctx, types.TotalBacking{
			MerMinted:  sdk.NewCoin(merlion.MicroUSMDenom, sdk.ZeroInt()),
			LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt()),
		})
	}

	k.SetPoolBacking(ctx, types.PoolBacking{
		MerMinted:  sdk.NewCoin(merlion.MicroUSMDenom, sdk.ZeroInt()),
		Backing:    sdk.NewCoin(params.BackingDenom, sdk.ZeroInt()),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt()),
	})

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeRegisterBacking,
			sdk.NewAttribute(types.AttributeKeyRiskParams, params.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return nil
}

func HandleRegisterCollateralProposal(ctx sdk.Context, k Keeper, p *types.RegisterCollateralProposal) error {
	params := p.RiskParams

	if k.IsCollateralRegistered(ctx, params.CollateralDenom) {
		return sdkerrors.Wrapf(types.ErrCollateralCoinAlreadyExists, "collateral coin denomination already registered: %s", params.CollateralDenom)
	}
	if !k.bankKeeper.HasSupply(ctx, params.CollateralDenom) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidCoins,
			"denomination '%s' cannot have a supply of 0", params.CollateralDenom,
		)
	}

	// assign missing fields with default value
	if params.LiquidationThreshold == nil {
		dec := sdk.NewDecWithPrec(70, 2)
		params.LiquidationThreshold = &dec
	}
	if params.LoanToValue == nil {
		dec := sdk.NewDecWithPrec(60, 2)
		params.LoanToValue = &dec
	}
	if params.BasicLoanToValue == nil {
		dec := sdk.NewDecWithPrec(40, 2)
		params.BasicLoanToValue = &dec
	}
	if params.CatalyticLionRatio == nil {
		dec := sdk.NewDecWithPrec(10, 2)
		params.CatalyticLionRatio = &dec
	}
	if params.LiquidationFee == nil {
		dec := sdk.NewDecWithPrec(10, 2)
		params.LiquidationFee = &dec
	}
	if params.MintFee == nil {
		dec := sdk.NewDecWithPrec(5, 3)
		params.MintFee = &dec
	}
	if params.InterestFee == nil {
		dec := sdk.NewDecWithPrec(1, 2)
		params.InterestFee = &dec
	}

	if err := validateCollateralRiskParams(&params); err != nil {
		return err
	}

	k.SetCollateralRiskParams(ctx, params)

	_, found := k.GetTotalCollateral(ctx)
	if !found {
		k.SetTotalCollateral(ctx, types.TotalCollateral{
			MerDebt:            sdk.NewCoin(merlion.MicroUSMDenom, sdk.ZeroInt()),
			LionCollateralized: sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt()),
		})
	}

	k.SetPoolCollateral(ctx, types.PoolCollateral{
		Collateral:         sdk.NewCoin(params.CollateralDenom, sdk.ZeroInt()),
		MerDebt:            sdk.NewCoin(merlion.MicroUSMDenom, sdk.ZeroInt()),
		LionCollateralized: sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt()),
	})

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeRegisterCollateral,
			sdk.NewAttribute(types.AttributeKeyRiskParams, params.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return nil
}

func updateInt(target *sdk.Int, patch *sdk.Int) uint8 {
	if patch == nil {
		// no set
		return 0
	} else {
		*target = *patch
	}
	return 1
}

func updateDecimal(target *sdk.Dec, patch *sdk.Dec) uint8 {
	if patch == nil {
		// no set
		return 0
	} else {
		*target = *patch
	}
	return 1
}

func setBackingRiskParamsProposal(ctx sdk.Context, k Keeper, patch *types.BackingRiskParams) error {
	params, found := k.GetBackingRiskParams(ctx, patch.BackingDenom)
	if !found {
		return sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", params.BackingDenom)
	}

	var updated uint8
	if params.Enabled != patch.Enabled {
		params.Enabled = patch.Enabled
		updated |= 1
	}
	updated |= updateInt(params.MaxBacking, patch.MaxBacking)
	updated |= updateInt(params.MaxMerMint, patch.MaxMerMint)
	updated |= updateDecimal(params.MintFee, patch.MintFee)
	updated |= updateDecimal(params.BurnFee, patch.BurnFee)
	updated |= updateDecimal(params.BuybackFee, patch.BuybackFee)
	updated |= updateDecimal(params.RebackFee, patch.RebackFee)

	if updated > 0 {
		if err := validateBackingRiskParams(ctx, k, &params); err != nil {
			return err
		}

		k.SetBackingRiskParams(ctx, params)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(types.EventTypeSetBackingRiskParams,
				sdk.NewAttribute(types.AttributeKeyRiskParams, params.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		})
	}
	return nil
}

func setCollateralRiskParamsProposal(ctx sdk.Context, k Keeper, patch *types.CollateralRiskParams) error {
	params, found := k.GetCollateralRiskParams(ctx, patch.CollateralDenom)
	if !found {
		return sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", params.CollateralDenom)
	}

	var updated uint8
	if params.Enabled != patch.Enabled {
		params.Enabled = patch.Enabled
		updated |= 1
	}
	updated |= updateInt(params.MaxCollateral, patch.MaxCollateral)
	updated |= updateInt(params.MaxMerMint, patch.MaxMerMint)
	updated |= updateDecimal(params.LiquidationThreshold, patch.LiquidationThreshold)
	updated |= updateDecimal(params.LoanToValue, patch.LoanToValue)
	updated |= updateDecimal(params.BasicLoanToValue, patch.BasicLoanToValue)
	updated |= updateDecimal(params.CatalyticLionRatio, patch.CatalyticLionRatio)
	updated |= updateDecimal(params.LiquidationFee, patch.LiquidationFee)
	updated |= updateDecimal(params.MintFee, patch.MintFee)
	updated |= updateDecimal(params.InterestFee, patch.InterestFee)

	if updated > 0 {
		if err := validateCollateralRiskParams(&params); err != nil {
			return err
		}

		k.SetCollateralRiskParams(ctx, params)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(types.EventTypeSetCollateralRiskParams,
				sdk.NewAttribute(types.AttributeKeyRiskParams, params.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		})
	}
	return nil
}

func validateBackingRiskParams(ctx sdk.Context, keeper Keeper, params *types.BackingRiskParams) error {
	if params.RebackFee.GT(keeper.RebackBonus(ctx)) {
		return sdkerrors.Wrap(types.ErrBackingParamsInvalid, "reback fee ratio should not be greater than reback bonus ratio")
	}
	return nil
}

func validateCollateralRiskParams(params *types.CollateralRiskParams) error {
	if params.LoanToValue.GTE(*params.LiquidationThreshold) {
		return sdkerrors.Wrap(types.ErrCollateralParamsInvalid, "loan-to-value must be < liquidation threshold")
	}
	if params.BasicLoanToValue.GT(*params.LoanToValue) {
		return sdkerrors.Wrap(types.ErrCollateralParamsInvalid, "basic loan-to-value must be <= loan-to-value")
	}
	return nil
}

func HandleSetBackingRiskParamsProposal(ctx sdk.Context, k Keeper, p *types.SetBackingRiskParamsProposal) error {
	return setBackingRiskParamsProposal(ctx, k, &p.RiskParams)
}

func HandleSetCollateralRiskParamsProposal(ctx sdk.Context, k Keeper, p *types.SetCollateralRiskParamsProposal) error {
	return setCollateralRiskParamsProposal(ctx, k, &p.RiskParams)
}

func HandleBatchSetBackingRiskParamsProposal(ctx sdk.Context, k Keeper, p *types.BatchSetBackingRiskParamsProposal) error {
	for _, params := range p.RiskParams {
		if err := setBackingRiskParamsProposal(ctx, k, &params); err != nil {
			return err
		}
	}
	return nil
}

func HandleBatchSetCollateralRiskParamsProposal(ctx sdk.Context, k Keeper, p *types.BatchSetCollateralRiskParamsProposal) error {
	for _, params := range p.RiskParams {
		if err := setCollateralRiskParamsProposal(ctx, k, &params); err != nil {
			return err
		}
	}
	return nil
}
