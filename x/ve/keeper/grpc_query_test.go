package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/merlion-zone/merlion/x/ve/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *KeeperTestSuite) TestKeeper_TotalVotingPower() {
	suite.SetupTest()
	k := suite.app.VeKeeper

	res, err := k.TotalVotingPower(sdk.WrapSDKContext(suite.ctx), nil)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "invalid request"))

	msg := &types.QueryTotalVotingPowerRequest{
		AtTime:  uint64(1912495957),
		AtBlock: int64(100),
	}
	res, err = k.TotalVotingPower(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "at time and at block cannot both be specified"))

	msg = &types.QueryTotalVotingPowerRequest{
		AtBlock: int64(100),
	}
	res, err = k.TotalVotingPower(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "invalid block"))

	msg = &types.QueryTotalVotingPowerRequest{}
	res, err = k.TotalVotingPower(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "at time and at block cannot both be zero"))

	msg = &types.QueryTotalVotingPowerRequest{
		AtBlock: int64(1),
	}
	res, err = k.TotalVotingPower(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Equal(sdk.ZeroInt(), res.Power)
}

func (suite *KeeperTestSuite) TestKeeper_GetTotalVotingPower() {
	suite.SetupTest()
	k := suite.app.VeKeeper
	power := k.GetTotalVotingPower(suite.ctx, 0, int64(1))
	suite.Equal(sdk.ZeroInt(), power)
}

func (suite *KeeperTestSuite) TestKeeper_VotingPower() {
	suite.SetupTest()
	k := suite.app.VeKeeper

	res, err := k.VotingPower(sdk.WrapSDKContext(suite.ctx), nil)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "invalid request"))

	msg := &types.QueryVotingPowerRequest{
		VeId: "100",
	}
	res, err = k.VotingPower(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().Nil(res)
	suite.Require().Error(err, sdkerrors.Wrapf(types.ErrInvalidVeID, "invalid ve id: %s", msg.VeId))

	err = suite.app.NftKeeper.Mint(suite.ctx, nft.NFT{
		ClassId: types.VeNftClass.Id,
		Id:      msg.VeId,
	}, sdk.AccAddress(suite.address.Bytes()))
	suite.Require().NoError(err)

	msg = &types.QueryVotingPowerRequest{
		VeId:    "100",
		AtTime:  uint64(1912495957),
		AtBlock: int64(100),
	}
	res, err = k.VotingPower(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "at time and at block cannot both be specified"))

	msg = &types.QueryVotingPowerRequest{
		VeId:    "100",
		AtBlock: int64(100),
	}
	res, err = k.VotingPower(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "invalid block"))

	msg = &types.QueryVotingPowerRequest{
		VeId: "100",
	}
	res, err = k.VotingPower(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "at time and at block cannot both be zero"))

	msg = &types.QueryVotingPowerRequest{
		VeId:    "100",
		AtBlock: int64(1),
	}
	res, err = k.VotingPower(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Equal(sdk.ZeroInt(), res.Power)
}

func (suite *KeeperTestSuite) TestKeeper_GetVotingPower() {
	suite.SetupTest()
	k := suite.app.VeKeeper
	power := k.GetVotingPower(suite.ctx, uint64(100), 0, int64(1))
	suite.Equal(sdk.ZeroInt(), power)
}

func (suite *KeeperTestSuite) TestKeeper_VeNfts() {
	suite.SetupTest()
	k := suite.app.VeKeeper

	res, err := k.VeNfts(sdk.WrapSDKContext(suite.ctx), nil)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "invalid request"))

	addr := sdk.AccAddress(suite.address.Bytes())
	msg := &types.QueryVeNftsRequest{
		Owner:      addr.String(),
		Pagination: &query.PageRequest{Key: nil, Limit: 10, CountTotal: true},
	}

	// Mint
	oneNft := nft.NFT{
		ClassId: types.VeNftClass.Id,
		Id:      "abc100",
	}
	err = suite.app.NftKeeper.Mint(suite.ctx, oneNft, addr)
	suite.Require().NoError(err)

	res, err = k.VeNfts(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Require().Equal(uint64(1), res.Pagination.Total)
	suite.Require().Equal(res.Nfts[0].Id, oneNft.Id)
	suite.Require().Equal(res.Nfts[0].ClassId, types.VeNftClass.Id)
}

func (suite *KeeperTestSuite) TestKeeper_VeNft() {
	suite.SetupTest()
	k := suite.app.VeKeeper

	res, err := k.VeNft(sdk.WrapSDKContext(suite.ctx), nil)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "invalid request"))

	// Mint
	addr := sdk.AccAddress(suite.address.Bytes())
	oneNft := nft.NFT{
		ClassId: types.VeNftClass.Id,
		Id:      "abc100",
	}
	err = suite.app.NftKeeper.Mint(suite.ctx, oneNft, addr)
	suite.Require().NoError(err)

	msg := &types.QueryVeNftRequest{Id: oneNft.Id}
	res, err = k.VeNft(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Require().Equal(res.Nft.Id, oneNft.Id)
	suite.Require().Equal(res.Nft.ClassId, types.VeNftClass.Id)
}

func (suite *KeeperTestSuite) TestKeeper_Params() {
	suite.SetupTest()
	k := suite.app.VeKeeper

	res, err := k.Params(sdk.WrapSDKContext(suite.ctx), nil)
	suite.Require().Nil(res)
	suite.Require().Error(err, status.Error(codes.InvalidArgument, "invalid request"))

	res, err = k.Params(sdk.WrapSDKContext(suite.ctx), &types.QueryParamsRequest{})
	suite.Require().NoError(err)

	params := k.GetParams(suite.ctx)
	suite.Require().Equal(res.Params, params)
}
