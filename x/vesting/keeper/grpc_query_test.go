package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	mertypes "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/vesting/keeper"
	"github.com/merlion-zone/merlion/x/vesting/types"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *KeeperTestSuite) TestKeeper_Airdrops() {
	suite.SetupTest()
	k := suite.app.VestingKeeper
	ctx := sdk.WrapSDKContext(suite.ctx)
	res, err := k.Airdrops(ctx, nil)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "empty request"))

	// Prepare
	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	sender := sdk.AccAddress(priv.PubKey().Address())
	receiver := sdk.AccAddress(suite.address.Bytes())
	denom := mertypes.AttoLionDenom
	k.SetAllocationAddresses(suite.ctx, types.AllocationAddresses{
		TeamVestingAddr: receiver.String(),
	})
	teamAddr := receiver
	impl := keeper.NewMsgServerImpl(k)
	// AddAirdrops
	airdrops := []types.Airdrop{
		{TargetAddr: sender.String(), Amount: sdk.NewCoin(denom, sdk.NewInt(100))}}
	_, err = impl.AddAirdrops(ctx, &types.MsgAddAirdrops{
		Sender:   teamAddr.String(),
		Airdrops: airdrops,
	})
	suite.Require().NoError(err)

	// Query
	req := &types.QueryAirdropsRequest{
		Completed:  false,
		Pagination: &query.PageRequest{Key: nil, Limit: 10, CountTotal: true},
	}
	res, err = k.Airdrops(ctx, req)
	suite.Require().NoError(err)
	suite.Require().Equal(airdrops, res.Airdrops)
	suite.Require().Equal(uint64(len(airdrops)), res.Pagination.Total)

	// ExecuteAirdrops
	_, err = impl.ExecuteAirdrops(ctx, &types.MsgExecuteAirdrops{
		Sender:   teamAddr.String(),
		MaxCount: 100,
	})
	suite.Require().NoError(err)

	// Query again
	res, err = k.Airdrops(ctx, req)
	suite.Require().NoError(err)
	suite.Require().Equal(0, len(res.Airdrops))

	// Set Completed
	req.Completed = true
	res, err = k.Airdrops(ctx, req)
	suite.Require().NoError(err)
	suite.Require().Equal(airdrops, res.Airdrops)
	suite.Require().Equal(uint64(len(airdrops)), res.Pagination.Total)
}

func (suite *KeeperTestSuite) TestKeeper_Airdrop() {
	suite.SetupTest()
	k := suite.app.VestingKeeper
	ctx := sdk.WrapSDKContext(suite.ctx)
	res, err := k.Airdrop(ctx, nil)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "empty request"))

	// Prepare
	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	sender := sdk.AccAddress(priv.PubKey().Address())
	receiver := sdk.AccAddress(suite.address.Bytes())
	denom := mertypes.AttoLionDenom
	k.SetAllocationAddresses(suite.ctx, types.AllocationAddresses{
		TeamVestingAddr: receiver.String(),
	})
	teamAddr := receiver
	impl := keeper.NewMsgServerImpl(k)
	// AddAirdrops
	airdrops := []types.Airdrop{
		{TargetAddr: sender.String(), Amount: sdk.NewCoin(denom, sdk.NewInt(100))}}
	_, err = impl.AddAirdrops(ctx, &types.MsgAddAirdrops{
		Sender:   teamAddr.String(),
		Airdrops: airdrops,
	})
	suite.Require().NoError(err)

	// Query
	req := &types.QueryAirdropRequest{
		TargetAddr: sender.String(),
		Completed:  false,
	}
	res, err = k.Airdrop(ctx, req)
	suite.Require().NoError(err)
	suite.Require().Equal(airdrops[0], res.Airdrop)

	// ExecuteAirdrops
	_, err = impl.ExecuteAirdrops(ctx, &types.MsgExecuteAirdrops{
		Sender:   teamAddr.String(),
		MaxCount: 100,
	})
	suite.Require().NoError(err)

	// Query again
	res, err = k.Airdrop(ctx, req)
	suite.Require().Error(err, status.Error(codes.NotFound, "airdrop target not found"))
	suite.Require().Nil(res)

	// Set Completed
	req.Completed = true
	res, err = k.Airdrop(ctx, req)
	suite.Require().NoError(err)
	suite.Require().Equal(airdrops[0], res.Airdrop)
}

func (suite *KeeperTestSuite) TestKeeper_Params() {
	suite.SetupTest()
	k := suite.app.VestingKeeper
	ctx := sdk.WrapSDKContext(suite.ctx)
	res, err := k.Params(ctx, nil)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "invalid request"))
}
