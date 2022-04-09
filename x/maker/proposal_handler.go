package maker

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/merlion-zone/merlion/x/maker/keeper"
	"github.com/merlion-zone/merlion/x/maker/types"
)

func NewMakerProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.RegisterBackingProposal:
			return nil
		case *types.RegisterCollateralProposal:
			return nil
		case *types.SetBackingRiskParamsProposal:
			return nil
		case *types.SetCollateralRiskParamsProposal:
			return nil
		case *types.BatchSetBackingRiskParamsProposal:
			return nil
		case *types.BatchSetCollateralRiskParamsProposal:
			return nil
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s proposal content type: %T", types.ModuleName, c)
		}
	}
}
