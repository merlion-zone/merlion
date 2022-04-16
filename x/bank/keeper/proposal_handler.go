package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/merlion-zone/merlion/x/bank/types"
)

func HandleSetDenomMetaDataProposal(ctx sdk.Context, k bankkeeper.Keeper, p *types.SetDenomMetadataProposal) error {
	k.SetDenomMetaData(ctx, p.Metadata)
	return nil
}
