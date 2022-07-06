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
	// DisplayDenom defines the denomination displayed to users in client applications.
	DisplayDenom = "lion"
	// BaseDenom defines to the default denomination used in Merlion (staking, EVM, governance, etc.)
	BaseDenom = AttoLionDenom

	AttoLionDenom = "alion" // 1e-18
	MicroUSMDenom = "uusm"  // 1e-6
)

var (
	// MicroUSMTarget defines the target exchange rate of uusm denominated in uUSD.
	MicroUSMTarget = sdk.OneDec()
)

func SetDenomMetaDataForStableCoins(ctx sdk.Context, k bankkeeper.Keeper) {
	for _, base := range []string{MicroUSMDenom} {
		if _, ok := k.GetDenomMetaData(ctx, base); ok {
			continue
		}

		display := base[1:] // e.g., usm
		// Register meta data to bank module
		k.SetDenomMetaData(ctx, banktypes.Metadata{
			Description: "The native stable token of the Merlion.",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: "u" + display, Exponent: uint32(0), Aliases: []string{"micro" + display}}, // e.g., uusm
				{Denom: "m" + display, Exponent: uint32(3), Aliases: []string{"milli" + display}}, // e.g., musm
				{Denom: display, Exponent: uint32(6), Aliases: []string{}},                        // e.g., usm
			},
			Base:    base,
			Display: display,
			Name:    fmt.Sprintf("%s", strings.ToUpper(display)), // e.g., USM
			Symbol:  fmt.Sprintf("%s", strings.ToUpper(display)), // e.g., USM
		})
	}
}
