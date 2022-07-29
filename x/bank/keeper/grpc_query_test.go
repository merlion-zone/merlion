package keeper_test

import (
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	mertypes "github.com/merlion-zone/merlion/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *KeeperTestSuite) TestKeeper_Balance() {
	suite.SetupTest()
	var (
		t     = suite.T()
		k     = suite.app.BankKeeper
		denom = mertypes.AttoLionDenom
	)

	res, err := k.Balance(sdk.WrapSDKContext(suite.ctx), nil)
	require.Nil(t, res)
	require.Error(t, err, status.Error(codes.InvalidArgument, "empty request"))

	res, err = k.Balance(sdk.WrapSDKContext(suite.ctx), &types.QueryBalanceRequest{})
	require.Nil(t, res)
	require.Error(t, err, status.Error(codes.InvalidArgument, "address cannot be empty"))

	res, err = k.Balance(sdk.WrapSDKContext(suite.ctx),
		&types.QueryBalanceRequest{Address: "xxx"})
	require.Nil(t, res)
	require.Error(t, err, status.Error(codes.InvalidArgument, "invalid denom"))

	res, err = k.Balance(sdk.WrapSDKContext(suite.ctx),
		&types.QueryBalanceRequest{Address: "xxx", Denom: denom})
	require.Nil(t, res)
	require.Error(t, err, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error()))

	req := &types.QueryBalanceRequest{Address: suite.addrs[0].String(), Denom: denom}
	res, err = k.Balance(sdk.WrapSDKContext(suite.ctx), req)
	require.NoError(t, err)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1134)), *res.Balance)
}

func (suite *KeeperTestSuite) TestKeeper_AllBalances() {
	suite.SetupTest()
	var (
		t     = suite.T()
		k     = suite.app.BankKeeper
		denom = mertypes.AttoLionDenom
	)

	req := &types.QueryAllBalancesRequest{
		Address:    suite.addrs[0].String(),
		Pagination: &query.PageRequest{Key: nil, Limit: 10, CountTotal: true}}
	res, err := k.AllBalances(sdk.WrapSDKContext(suite.ctx), req)
	require.NoError(t, err)
	require.Equal(t, 1, len(res.Balances))
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1134)), res.Balances[0])
	require.Equal(t, uint64(1), res.Pagination.Total)
}

func (suite *KeeperTestSuite) TestKeeper_SpendableBalances() {
	suite.SetupTest()
	var (
		t     = suite.T()
		k     = suite.app.BankKeeper
		denom = mertypes.AttoLionDenom
	)

	req := &types.QuerySpendableBalancesRequest{
		Address:    suite.addrs[0].String(),
		Pagination: &query.PageRequest{Key: nil, Limit: 10, CountTotal: true}}
	res, err := k.SpendableBalances(sdk.WrapSDKContext(suite.ctx), req)
	require.NoError(t, err)
	require.Equal(t, 1, len(res.Balances))
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1134)), res.Balances[0])
	require.Equal(t, uint64(1), res.Pagination.Total)
}

func (suite *KeeperTestSuite) TestKeeper_TotalSupply() {
	suite.SetupTest()
	var (
		t     = suite.T()
		k     = suite.app.BankKeeper
		denom = mertypes.AttoLionDenom
	)

	req := &types.QueryTotalSupplyRequest{
		Pagination: &query.PageRequest{Key: nil, Limit: 10, CountTotal: true}}
	res, err := k.TotalSupply(sdk.WrapSDKContext(suite.ctx), req)
	require.NoError(t, err)
	require.Equal(t, 1, len(res.Supply))

	amt, _ := new(big.Int).SetString("500000000000000000000006170", 10)
	expectedSupply := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromBigInt(amt)))
	require.Equal(t, expectedSupply, res.Supply)
	require.Equal(t, uint64(1), res.Pagination.Total)
}

func (suite *KeeperTestSuite) TestKeeper_SupplyOf() {
	suite.SetupTest()
	var (
		t     = suite.T()
		k     = suite.app.BankKeeper
		denom = mertypes.AttoLionDenom
	)

	res, err := k.SupplyOf(sdk.WrapSDKContext(suite.ctx), nil)
	require.Nil(t, res)
	require.Error(t, err, status.Error(codes.InvalidArgument, "empty request"))

	res, err = k.SupplyOf(sdk.WrapSDKContext(suite.ctx), &types.QuerySupplyOfRequest{Denom: ""})
	require.Nil(t, res)
	require.Error(t, err, status.Error(codes.InvalidArgument, "invalid denom"))

	req := &types.QuerySupplyOfRequest{Denom: denom}
	res, err = k.SupplyOf(sdk.WrapSDKContext(suite.ctx), req)
	require.NoError(t, err)

	amt, _ := new(big.Int).SetString("500000000000000000000006170", 10)
	expectedSupply := sdk.NewCoin(denom, sdk.NewIntFromBigInt(amt))
	require.Equal(t, expectedSupply, res.Amount)
}

func (suite *KeeperTestSuite) TestKeeper_Params() {
	suite.SetupTest()
	var (
		t = suite.T()
		k = suite.app.BankKeeper
	)

	res, err := k.Params(sdk.WrapSDKContext(suite.ctx), nil)
	require.Nil(t, res)
	require.Error(t, err, status.Error(codes.InvalidArgument, "empty request"))

	res, err = k.Params(sdk.WrapSDKContext(suite.ctx), &types.QueryParamsRequest{})
	require.NoError(t, err)
	params := types.DefaultParams()
	require.Equal(t, params.DefaultSendEnabled, res.Params.DefaultSendEnabled)
}

func (suite *KeeperTestSuite) TestKeeper_DenomsMetadata() {
	suite.SetupTest()
	var (
		t        = suite.T()
		k        = suite.app.BankKeeper
		denom    = mertypes.MicroUSMDenom
		base     = denom
		display  = base[1:]
		uusmMeta = banktypes.Metadata{
			Description: "The native stable token of the Merlion.",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: "u" + display, Exponent: uint32(0), Aliases: []string{"micro" + display}}, // e.g., uusm
				{Denom: "m" + display, Exponent: uint32(3), Aliases: []string{"milli" + display}}, // e.g., musm
				{Denom: display, Exponent: uint32(6), Aliases: []string{""}},                      // e.g., usm
			},
			Base:    base,
			Display: display,
			Name:    strings.ToUpper(display), // e.g., USM
			Symbol:  strings.ToUpper(display), // e.g., USM
		}
	)

	req := &types.QueryDenomsMetadataRequest{
		Pagination: &query.PageRequest{Key: nil, Limit: 10, CountTotal: true},
	}
	res, err := k.DenomsMetadata(sdk.WrapSDKContext(suite.ctx), req)
	require.NoError(t, err)
	meta := res.Metadatas[0]
	require.Equal(t, 1, len(res.Metadatas))
	require.Equal(t, uusmMeta.Description, meta.Description)
	require.Equal(t, uusmMeta.Base, meta.Base)
	require.Equal(t, uusmMeta.Display, meta.Display)
	require.Equal(t, uusmMeta.Name, meta.Name)
	require.Equal(t, uusmMeta.Symbol, meta.Symbol)
	require.Equal(t, uusmMeta.DenomUnits[0], meta.DenomUnits[0])
	require.Equal(t, uusmMeta.DenomUnits[1], meta.DenomUnits[1])
	require.Equal(t, uint64(1), res.Pagination.Total)
}

func (suite *KeeperTestSuite) TestKeeper_DenomMetadata() {
	suite.SetupTest()
	var (
		t        = suite.T()
		k        = suite.app.BankKeeper
		denom    = mertypes.MicroUSMDenom
		base     = denom
		display  = base[1:]
		uusmMeta = banktypes.Metadata{
			Description: "The native stable token of the Merlion.",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: "u" + display, Exponent: uint32(0), Aliases: []string{"micro" + display}}, // e.g., uusm
				{Denom: "m" + display, Exponent: uint32(3), Aliases: []string{"milli" + display}}, // e.g., musm
				{Denom: display, Exponent: uint32(6), Aliases: []string{""}},                      // e.g., usm
			},
			Base:    base,
			Display: display,
			Name:    strings.ToUpper(display), // e.g., USM
			Symbol:  strings.ToUpper(display), // e.g., USM
		}
	)
	res, err := k.DenomMetadata(sdk.WrapSDKContext(suite.ctx), nil)
	require.Nil(t, res)
	require.Error(t, err, status.Error(codes.InvalidArgument, "empty request"))

	res, err = k.DenomMetadata(sdk.WrapSDKContext(suite.ctx), &types.QueryDenomMetadataRequest{})
	require.Nil(t, res)
	require.Error(t, err, status.Error(codes.InvalidArgument, "invalid denom"))

	res, err = k.DenomMetadata(sdk.WrapSDKContext(suite.ctx), &types.QueryDenomMetadataRequest{
		Denom: "xxx",
	})
	require.Nil(t, res)
	require.Error(t, err, status.Errorf(codes.NotFound, "client metadata for denom %s", "xxx"))

	req := &types.QueryDenomMetadataRequest{Denom: denom}

	res, err = k.DenomMetadata(sdk.WrapSDKContext(suite.ctx), req)
	require.NoError(t, err)
	meta := res.Metadata
	require.Equal(t, uusmMeta.Description, meta.Description)
	require.Equal(t, uusmMeta.Base, meta.Base)
	require.Equal(t, uusmMeta.Display, meta.Display)
	require.Equal(t, uusmMeta.Name, meta.Name)
	require.Equal(t, uusmMeta.Symbol, meta.Symbol)
	require.Equal(t, uusmMeta.DenomUnits[0], meta.DenomUnits[0])
	require.Equal(t, uusmMeta.DenomUnits[1], meta.DenomUnits[1])
}
