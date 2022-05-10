package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/merlion-zone/merlion/x/staking/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(
	ctx sdk.Context, keeper keeper.Keeper, accountKeeper stakingtypes.AccountKeeper,
	bankKeeper stakingtypes.BankKeeper, data *stakingtypes.GenesisState,
) (res []abci.ValidatorUpdate) {
	// TODO: replace keeper.Keeper with keeper
	res = staking.InitGenesis(ctx, keeper.Keeper, accountKeeper, bankKeeper, data)

	keeper.CheckDenom(ctx)

	return
}
