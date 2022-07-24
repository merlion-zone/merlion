package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/merlion-zone/merlion/app"
	"github.com/merlion-zone/merlion/types"
	merlion "github.com/merlion-zone/merlion/types"
	erc20types "github.com/merlion-zone/merlion/x/erc20/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"
)

func TestKeeper_RegisterCoin(t *testing.T) {
	var (
		coinMetadata = banktypes.Metadata{
			Description: "USDT",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: "u" + "usd", Exponent: uint32(0), Aliases: []string{"micro" + "usd"}}, // e.g., uusd
				{Denom: "m" + "usd", Exponent: uint32(3), Aliases: []string{"milli" + "usd"}}, // e.g., musd
				{Denom: "usd", Exponent: uint32(6), Aliases: []string{}},                      // e.g., usd
			},
			Base:    "uusd",
			Display: "usd",
			Name:    "USDT",
			Symbol:  "USDT",
		}
		PKS   = simapp.CreateTestPubKeys(5)
		addrs = []sdk.AccAddress{
			sdk.AccAddress(PKS[0].Address()),
			sdk.AccAddress(PKS[1].Address()),
			sdk.AccAddress(PKS[2].Address()),
			sdk.AccAddress(PKS[3].Address()),
			sdk.AccAddress(PKS[4].Address()),
		}
		valConsPk1 = PKS[0]
		myapp      = app.Setup(false)
		ctx        = myapp.BaseApp.NewContext(false, tmproto.Header{})
		err        error
	)

	app.FundTestAddrs(myapp, ctx, addrs, sdk.NewInt(1234))
	ctx = myapp.BaseApp.NewContext(false, tmproto.Header{
		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		ChainID:         "merlion_5000-101",
		Height:          1,
		Time:            time.Now().UTC(),
		ProposerAddress: addrs[0].Bytes(),
	})

	tstaking := teststaking.NewHelper(t, ctx, myapp.StakingKeeper.Keeper)
	tstaking.Denom = types.AttoLionDenom
	// create validator with 50% commission
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(sdk.ValAddress(valConsPk1.Address()), valConsPk1, sdk.NewInt(100), true)

	_, err = myapp.Erc20Keeper.RegisterCoin(ctx, merlion.DisplayDenom)
	require.Error(t, err, sdkerrors.Wrapf(erc20types.ErrEVMDenom, "cannot register the EVM denomination %s", merlion.DisplayDenom))

	_, err = myapp.Erc20Keeper.RegisterCoin(ctx, "USDT")
	require.Error(t, err, sdkerrors.Wrapf(erc20types.ErrEVMDenom, "cannot get metadata of denom %s", "USDT"))

	myapp.BankKeeper.SetDenomMetaData(ctx, coinMetadata)
	_, err = myapp.Erc20Keeper.RegisterCoin(ctx, "uusd")
	require.Error(t, err, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "base denomination 'uusd' cannot have a supply of 0"))

	totalSupply := sdk.NewCoins(sdk.NewCoin("uusd", sdk.NewInt(10000000)))
	myapp.BankKeeper.MintCoins(ctx, erc20types.ModuleName, totalSupply)

	_, err = myapp.Erc20Keeper.RegisterCoin(ctx, "uusd")
	require.Error(t, err, sdkerrors.Wrapf(erc20types.ErrTokenPairAlreadyExists, "coin denomination already registered: %s", coinMetadata.Base))

	addr, err := myapp.Erc20Keeper.DeployERC20Contract(ctx, coinMetadata)
	require.NoError(t, err)
	require.Equal(t, "0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75", addr.String())

	tokenPair, err := myapp.Erc20Keeper.RegisterERC20(ctx, addr)
	require.NoError(t, err)

	require.Equal(t, "0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75", tokenPair.Erc20Address)
	require.Equal(t, "erc20/0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75", tokenPair.Denom)
	require.Equal(t, erc20types.Owner(2), tokenPair.ContractOwner)

	// QueryERC20
	erc20Data, err := myapp.Erc20Keeper.QueryERC20(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, "USDT", erc20Data.Name)
	require.Equal(t, "USDT", erc20Data.Symbol)
	require.Equal(t, uint8(0), erc20Data.Decimals)
}
