package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/merlion-zone/merlion/x/oracle/keeper"
	"github.com/merlion-zone/merlion/x/oracle/types"
)

func SimulateMsgAggregateExchangeRatePrevote(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgAggregateExchangeRatePrevote{}

		// TODO: Handling the AggregateExchangeRatePrevote simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "AggregateExchangeRatePrevote simulation not implemented"), nil, nil
	}
}
