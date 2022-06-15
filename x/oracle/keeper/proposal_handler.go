package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/oracle/types"
)

func HandleRegisterTargetProposal(ctx sdk.Context, k Keeper, p *types.RegisterTargetProposal) error {
	params := p.TargetParams

	if k.IsTarget(ctx, params.Denom) {
		return sdkerrors.Wrapf(types.ErrExistingTarget, "existing target denom '%s'", params.Denom)
	}

	// Check if the coin exists by ensuring the supply is set
	if !k.bankKeeper.HasSupply(ctx, params.Denom) && params.Denom != merlion.MicroUSDDenom {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidCoins,
			"target denom '%s' cannot have a supply of 0", params.Denom,
		)
	}

	k.SetTarget(ctx, params.Denom)

	switch params.Source {
	case types.TARGET_SOURCE_VALIDATORS:
		k.SetVoteTarget(ctx, params.Denom)
	default:
		// TODO
	}

	return nil
}
