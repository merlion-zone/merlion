package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/merlion-zone/merlion/x/staking/keeper"
	"github.com/merlion-zone/merlion/x/staking/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

type AppModuleBasic struct {
	staking.AppModuleBasic
}

type AppModule struct {
	staking.AppModule
	keeper keeper.Keeper
}

func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper) AppModule {
	return AppModule{
		AppModule: staking.NewAppModule(cdc, keeper.Keeper, ak, bk),
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	stakingtypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))

	querier := stakingkeeper.Querier{Keeper: am.keeper.Keeper}
	stakingtypes.RegisterQueryServer(cfg.QueryServer(), querier)

	m := stakingkeeper.NewMigrator(am.keeper.Keeper)
	cfg.RegisterMigration(stakingtypes.ModuleName, 1, m.Migrate1to2)
}
