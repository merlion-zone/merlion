package bank

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// AppModule implements an application module for the bank module.
type AppModule struct {
	bank.AppModule
	keeper bankkeeper.Keeper
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	banktypes.RegisterMsgServer(cfg.MsgServer(), bankkeeper.NewMsgServerImpl(am.keeper))
	banktypes.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	m := bankkeeper.NewMigrator(am.keeper)
	cfg.RegisterMigration(banktypes.ModuleName, 1, m.Migrate1to2)
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper bankkeeper.Keeper, accountKeeper banktypes.AccountKeeper) AppModule {
	return AppModule{
		AppModule: bank.NewAppModule(cdc, keeper, accountKeeper),
		keeper:    keeper,
	}
}
