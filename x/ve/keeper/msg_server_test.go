package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/keeper"
	"github.com/merlion-zone/merlion/x/ve/types"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
)

func (suite *KeeperTestSuite) TestVeCreate() {
	require := suite.Require()

	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(err)
	receiver := sdk.AccAddress(priv.PubKey().Address())

	testCases := []struct {
		name         string
		pass         bool
		to           string
		amount       string
		lockDuration uint64
	}{
		{"receiver is different from sender", true, receiver.String(), "1alion", types.RegulatedPeriod},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)

			sender := sdk.AccAddress(suite.address.Bytes())

			amount, _ := sdk.ParseCoinNormalized(tc.amount)

			res, err := keeper.NewMsgServerImpl(suite.app.VeKeeper).Create(ctx, &types.MsgCreate{
				Sender:       sender.String(),
				To:           tc.to,
				Amount:       amount,
				LockDuration: tc.lockDuration,
			})
			if tc.pass {
				require.NoError(err, tc.name)

				expectedUnlockTime := types.RegulatedUnixTimeFromNow(suite.ctx, tc.lockDuration)
				require.Equal(expectedUnlockTime, res.UnlockTime)
			} else {
				require.Error(err, tc.name)
			}
		})
	}
}
