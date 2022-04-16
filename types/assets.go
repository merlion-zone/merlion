package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// nolint
const (
	AttoLionDenom = "alion" // 1e-18
	MicroUSDDenom = "uusd"  // 1e-6
)

var (
	MicroUSDTarget = sdk.OneDec()
)

func SetDenomMetaDataForStableTokens(ctx sdk.Context, k bankkeeper.Keeper) {
	for _, base := range []string{MicroUSDDenom} {
		if _, ok := k.GetDenomMetaData(ctx, base); ok {
			continue
		}

		display := base[1:] // e.g., usd
		// Register meta data to bank module
		k.SetDenomMetaData(ctx, banktypes.Metadata{
			Description: "The native stable token of the Merlion.",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: "u" + display, Exponent: uint32(0), Aliases: []string{"micro" + display}}, // e.g., uusd
				{Denom: "m" + display, Exponent: uint32(3), Aliases: []string{"milli" + display}}, // e.g., musd
				{Denom: display, Exponent: uint32(6), Aliases: []string{}},                        // e.g., usd
			},
			Base:    base,
			Display: display,
			Name:    fmt.Sprintf("%s MER", strings.ToUpper(display)),               // e.g., USD MER
			Symbol:  fmt.Sprintf("%sM", strings.ToUpper(display[:len(display)-1])), // e.g., USM
		})
	}
}
