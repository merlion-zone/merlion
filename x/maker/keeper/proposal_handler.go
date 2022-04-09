package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	k.SetBackingRiskParams(ctx, *params)
	// TODO: event
	return nil
}

func HandleRegisterCollateralProposal(ctx sdk.Context, k Keeper, p *types.RegisterCollateralProposal) error {
	params := p.RiskParams
	if k.IsBackingRegistered(ctx, params.CollateralDenom) {
		return sdkerrors.Wrapf(types.ErrCollateralCoinAlreadyExists, "collateral coin denomination already registered: %s", params.CollateralDenom)
	}
	if !k.bankKeeper.HasSupply(ctx, params.CollateralDenom) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidCoins,
			"denomination '%s' cannot have a supply of 0", params.CollateralDenom,
		)
	}
	k.SetCollateralRiskParams(ctx, *params)
	// TODO: event
	return nil
}

func updateDecimal(target *sdk.Dec, patch *sdk.Dec) uint8 {
	if patch == nil || patch.IsZero() {
		// no set
		return 0
	} else if patch.IsNegative() {
		// negative means the maximum
		*target = sdk.MaxSortableDec
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
	updated |= updateDecimal(params.MaxBacking, patch.MaxBacking)
	updated |= updateDecimal(params.MaxMerMint, patch.MaxMerMint)
	updated |= updateDecimal(params.MintFee, patch.MintFee)
	updated |= updateDecimal(params.BurnFee, patch.BurnFee)
	updated |= updateDecimal(params.BuybackFee, patch.BuybackFee)
	updated |= updateDecimal(params.RecollateralizeFee, patch.RecollateralizeFee)
	updated |= updateDecimal(params.MintPriceThreshold, patch.MintPriceThreshold)
	updated |= updateDecimal(params.BurnPriceThreshold, patch.BurnPriceThreshold)

	if updated > 0 {
		k.SetBackingRiskParams(ctx, params)
		// TODO: event
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
	updated |= updateDecimal(params.MaxCollateral, patch.MaxCollateral)
	updated |= updateDecimal(params.MaxMerMint, patch.MaxMerMint)
	updated |= updateDecimal(params.LiquidationThreshold, patch.LiquidationThreshold)
	updated |= updateDecimal(params.LoanToValue, patch.LoanToValue)
	updated |= updateDecimal(params.BasicLoanToValue, patch.BasicLoanToValue)
	updated |= updateDecimal(params.CatalyticLionRatio, patch.CatalyticLionRatio)
	updated |= updateDecimal(params.LiquidationFee, patch.LiquidationFee)
	updated |= updateDecimal(params.MintFee, patch.MintFee)
	updated |= updateDecimal(params.InterestFee, patch.InterestFee)

	if updated > 0 {
		k.SetCollateralRiskParams(ctx, params)
		// TODO: event
	}
	return nil
}

func HandleSetBackingRiskParamsProposal(ctx sdk.Context, k Keeper, p *types.SetBackingRiskParamsProposal) error {
	return setBackingRiskParamsProposal(ctx, k, p.RiskParams)
}

func HandleSetCollateralRiskParamsProposal(ctx sdk.Context, k Keeper, p *types.SetCollateralRiskParamsProposal) error {
	return setCollateralRiskParamsProposal(ctx, k, p.RiskParams)
}

func HandleBatchSetBackingRiskParamsProposal(ctx sdk.Context, k Keeper, p *types.BatchSetBackingRiskParamsProposal) error {
	for _, params := range p.RiskParams {
		if err := setBackingRiskParamsProposal(ctx, k, params); err != nil {
			return err
		}
	}
	return nil
}

func HandleBatchSetCollateralRiskParamsProposal(ctx sdk.Context, k Keeper, p *types.BatchSetCollateralRiskParamsProposal) error {
	for _, params := range p.RiskParams {
		if err := setCollateralRiskParamsProposal(ctx, k, params); err != nil {
			return err
		}
	}
	return nil
}
