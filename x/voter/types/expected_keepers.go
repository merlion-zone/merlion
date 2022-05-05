package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	gaugekeeper "github.com/merlion-zone/merlion/x/gauge/keeper"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

type Vekeeper interface {
	LockDenom(ctx sdk.Context) string
	GetVotingPower(ctx sdk.Context, veID uint64, atTime uint64, atBlock int64) sdk.Int
	SetVeVoted(ctx sdk.Context, veID uint64, voted bool)
}

type GaugeKeeper interface {
	CreateGauge(ctx sdk.Context, depoistDenom string)
	HasGauge(ctx sdk.Context, depositDenom string) bool
	GetGauges(ctx sdk.Context) (denoms []string)
	Gauge(ctx sdk.Context, depoistDenom string) gaugekeeper.Gauge
	Bribe(ctx sdk.Context, depoistDenom string) gaugekeeper.Bribe
}
