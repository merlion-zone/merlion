package app

import (
	"strings"
	"sync"

	mgravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/multigravity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethermint "github.com/tharsis/ethermint/types"

	merlion "github.com/merlion-zone/merlion/types"
)

const (
	AccountAddressPrefix = "mer"
)

// SetBech32Prefixes sets the global prefixes to be used when serializing addresses and public keys to Bech32 strings.
func SetBech32Prefixes(config *sdk.Config, accountAddressPrefix string) {
	// Set prefixes
	accountPubKeyPrefix := accountAddressPrefix + "pub"
	validatorAddressPrefix := accountAddressPrefix + "valoper"
	validatorPubKeyPrefix := accountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := accountAddressPrefix + "valcons"
	consNodePubKeyPrefix := accountAddressPrefix + "valconspub"

	config.SetBech32PrefixForAccount(accountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdk.Config) {
	config.SetCoinType(ethermint.Bip44CoinType)
	config.SetPurpose(sdk.Purpose)                      // Shared
	config.SetFullFundraiserPath(ethermint.BIP44HDPath) // nolint: staticcheck
}

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms() {
	if err := sdk.RegisterDenom(merlion.DisplayDenom, sdk.OneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(merlion.BaseDenom, sdk.NewDecWithPrec(1, ethermint.BaseDenomUnit)); err != nil {
		panic(err)
	}

	mgravitytypes.SetGasCoinMetata(banktypes.Metadata{
		Description: "The native gas token of the Merlion.",
		DenomUnits: []*banktypes.DenomUnit{{
			Denom:    merlion.DisplayDenom,
			Exponent: ethermint.BaseDenomUnit,
			Aliases:  []string{},
		}, {
			Denom:    merlion.BaseDenom,
			Exponent: 0,
			Aliases:  []string{},
		}},
		Base:    merlion.BaseDenom,
		Display: merlion.DisplayDenom,
		Name:    strings.ToUpper(merlion.DisplayDenom),
		Symbol:  strings.ToUpper(merlion.DisplayDenom),
	})
}

var setup sync.Once

func SetupConfig() {
	setup.Do(func() {
		config := sdk.GetConfig()
		SetBech32Prefixes(config, AccountAddressPrefix)
		SetBip44CoinType(config)
		RegisterDenoms()
		config.Seal()
	})
}
