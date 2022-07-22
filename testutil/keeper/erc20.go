package keeper

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	nfttypes "github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	merlion "github.com/merlion-zone/merlion/types"
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
	"github.com/tharsis/ethermint/x/evm"
	evmkeeper "github.com/tharsis/ethermint/x/evm/keeper"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

const faucetAccountName = "faucet"

var (
	// module account permissions
	maccPerms = map[string][]string{
		faucetAccountName:              {authtypes.Minter},
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
	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		oracletypes.ModuleName: true,
	}
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		evm.AppModuleBasic{},
	)
	InitTokens = sdk.TokensFromConsensusPower(200, sdk.DefaultPowerReduction)
	InitCoins  = sdk.NewCoins(sdk.NewCoin("uusd", InitTokens))
)

// MakeEncodingConfig nolint
func MakeEncodingConfig(_ *testing.T) simparams.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	std.RegisterInterfaces(interfaceRegistry)
	std.RegisterLegacyAminoCodec(amino)

	ModuleBasics.RegisterLegacyAminoCodec(amino)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	types.RegisterCodec(amino)
	types.RegisterInterfaces(interfaceRegistry)
	return simparams.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

func Erc20KeeperWithExtra(t *testing.T) (*keeper.Keeper, *bankkeeper.BaseKeeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	keyBank := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyEvm := sdk.NewKVStoreKey(evmtypes.StoreKey)
	keyStaking := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(paramstypes.TStoreKey)

	db := tmdb.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, tmproto.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())
	encodingConfig := MakeEncodingConfig(t)
	appCodec, legacyAmino := encodingConfig.Marshaler, encodingConfig.Amino

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyEvm, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)

	require.NoError(t, ms.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		nil,
		"Erc20Params",
	)
	var erc20Keeper keeper.Keeper
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, keyParams, tKeyParams)
	accountKeeper := authkeeper.NewAccountKeeper(appCodec, keyAcc, paramsKeeper.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms)
	bankKeeper := bankkeeper.NewBaseKeeper(appCodec, keyBank, accountKeeper, paramsKeeper.Subspace(banktypes.ModuleName), blockedAddrs())

	// StakingKeeper
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keyStaking,
		accountKeeper,
		bankKeeper,
		paramsKeeper.Subspace(stakingtypes.ModuleName),
	)

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = merlion.AttoLionDenom
	stakingKeeper.SetParams(ctx, stakingParams)

	// EVMKeeper
	evmKeeper := evmkeeper.NewKeeper(
		appCodec, keyEvm,
		sdk.NewTransientStoreKey(evmtypes.TransientKey),
		paramsKeeper.Subspace(evmtypes.ModuleName),
		accountKeeper, bankKeeper, stakingKeeper, nil,
		"callTracer",
	)
	evmKeeper.SetParams(ctx, evmtypes.DefaultParams())

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		paramsSubspace,
		accountKeeper,
		bankKeeper,
		evmKeeper,
	)
	erc20Keeper = *k

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	evmKeeper = evmKeeper.SetHooks(
		evmkeeper.NewMultiEvmHooks(
			erc20Keeper.EvmHooks(),
		),
	)
	erc20Acc := authtypes.NewEmptyModuleAccount(types.ModuleName)
	accountKeeper.SetModuleAccount(ctx, erc20Acc)
	return k, &bankKeeper, ctx
}

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

func blockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}
