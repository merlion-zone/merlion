package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	nfttypes "github.com/cosmos/cosmos-sdk/x/nft"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	"github.com/merlion-zone/merlion/x/erc20/keeper"
	"github.com/merlion-zone/merlion/x/erc20/types"
	erc20types "github.com/merlion-zone/merlion/x/erc20/types"
	gaugetypes "github.com/merlion-zone/merlion/x/gauge/types"
	makertypes "github.com/merlion-zone/merlion/x/maker/types"
	oracletypes "github.com/merlion-zone/merlion/x/oracle/types"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
	customvestingtypes "github.com/merlion-zone/merlion/x/vesting/types"
	votertypes "github.com/merlion-zone/merlion/x/voter/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
	ethermint "github.com/tharsis/ethermint/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

func Erc20Keeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		nil,
		"Erc20Params",
	)
	// AccountKeeper
	akKey := sdk.NewKVStoreKey(authtypes.StoreKey)
	akSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		akKey,
		nil,
		authtypes.ModuleName,
	)
	// module account permissions
	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		evmtypes.ModuleName:            {authtypes.Minter, authtypes.Burner}, // used for secure addition and subtraction of balance using module account
		erc20types.ModuleName:          {authtypes.Minter, authtypes.Burner},
		oracletypes.ModuleName:         nil,
		makertypes.ModuleName:          {authtypes.Minter, authtypes.Burner, authtypes.Staking},
		nfttypes.ModuleName:            nil,
		vetypes.ModuleName:             {authtypes.Burner},
		vetypes.EmissionPoolName:       {authtypes.Minter},
		vetypes.DistributionPoolName:   nil,
		gaugetypes.ModuleName:          nil,
		votertypes.ModuleName:          nil,
		customvestingtypes.ModuleName:  {authtypes.Minter},
	}
	accountKeeper := authkeeper.NewAccountKeeper(
		cdc, akKey, akSubspace, ethermint.ProtoAccount, maccPerms,
	)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		paramsSubspace,
		accountKeeper,
		nil,
		nil,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
