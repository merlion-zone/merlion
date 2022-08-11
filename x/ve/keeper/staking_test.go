package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/keeper"
	"github.com/merlion-zone/merlion/x/ve/types"
)

func (suite *KeeperTestSuite) TestKeeper_SlashLockedAmountByUser() {
	require := suite.Require()
	ctx := sdk.WrapSDKContext(suite.ctx)
	sender := sdk.AccAddress(suite.address.Bytes())
	k := suite.app.VeKeeper
	impl := keeper.NewMsgServerImpl(k)
	denom := "alion"
	veID := uint64(1)
	lockAmt := sdk.NewCoin(denom, sdk.NewInt(100))
	res, err := impl.Create(ctx, &types.MsgCreate{
		Sender:       sender.String(),
		To:           sender.String(),
		Amount:       lockAmt,
		LockDuration: types.RegulatedPeriod,
	})
	require.NoError(err)
	require.Equal("ve-1", res.VeId)

	locked := k.GetLockedAmountByUser(suite.ctx, veID)
	totalLocked := k.GetTotalLockedAmount(suite.ctx)
	require.Equal(lockAmt.Amount, locked.Amount)
	require.Equal(lockAmt.Amount, totalLocked)

	slashed := sdk.NewInt(50)
	k.SlashLockedAmountByUser(suite.ctx, veID, slashed)
	locked = k.GetLockedAmountByUser(suite.ctx, veID)
	totalLocked = k.GetTotalLockedAmount(suite.ctx)
	require.Equal(lockAmt.Amount.Sub(slashed), locked.Amount)
	require.Equal(lockAmt.Amount.Sub(slashed), totalLocked)
}
