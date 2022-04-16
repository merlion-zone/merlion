package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/merlion-zone/merlion/x/oracle/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdAggregateExchangeRatePrevote(),
		CmdAggregateExchangeRateVote(),
		CmdDelegateFeederPermission(),
	)
	// this line is used by starport scaffolding # 1

	return cmd
}

// CmdAggregateExchangeRatePrevote will create a MsgAggregateExchangeRatePrevote tx and sign it with the given key.
func CmdAggregateExchangeRatePrevote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-prevote [salt] [exchange-rates] [validator]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Submit an oracle aggregate prevote for the exchange rates of various assets",
		Long: strings.TrimSpace(`
Submit an oracle aggregate prevote for the exchange rates of various assets denominated in $uUSD.
The purpose of aggregate prevote is to hide aggregate exchange rate vote with hash which is formatted 
as hex string in SHA256("{salt};{denom}:{exchange_rate},...,{denom}:{exchange_rate};{voter}")

# Aggregate Prevote
$ merliond tx oracle aggregate-prevote 1234 alion:1.234,uusd:0.99

where "alion,uusd" is the denominating currencies, and "1.234,0.99" is the exchange rates of these currencies in $uUSD from the voter's point of view.

If voting from a voting delegate, set "validator" to the address of the validator to vote on behalf of:
$ merliond tx oracle aggregate-prevote 1234 alion:1.234,uusd:0.99 mervaloper1...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			salt := args[0]
			exchangeRatesStr := args[1]
			_, err = types.ParseExchangeRateTuples(exchangeRatesStr)
			if err != nil {
				return fmt.Errorf("given exchange_rates {%s} is not a valid format; exchange_rate should be formatted as DecCoins; %s", exchangeRatesStr, err.Error())
			}

			// Get from address
			voter := clientCtx.GetFromAddress()

			// By default, the voter is voting on behalf of itself
			validator := sdk.ValAddress(voter)

			// Override validator if validator is given
			if len(args) == 3 {
				parsedVal, err := sdk.ValAddressFromBech32(args[2])
				if err != nil {
					return errors.Wrap(err, "validator address is invalid")
				}
				validator = parsedVal
			}

			hash := types.GetAggregateVoteHash(salt, exchangeRatesStr, validator)
			msg := types.NewMsgAggregateExchangeRatePrevote(hash, voter, validator)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdAggregateExchangeRateVote will create a MsgAggregateExchangeRateVote tx and sign it with the given key.
func CmdAggregateExchangeRateVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-vote [salt] [exchange-rates] [validator]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Submit an oracle aggregate vote for the exchange_rates of various assets",
		Long: strings.TrimSpace(`
Submit a aggregate vote for the exchange_rates of various assets w.r.t $uUSD. Companion to a prevote submitted in the previous vote period. 

$ merliond tx oracle aggregate-vote 1234 alion:1.234,uusd:0.99

where "alion,uusd" is the denominating currencies, and "1.234,0.99" is the exchange rates of these currencies in $uUSD from the voter's point of view.

"salt" should match the salt used to generate the SHA256 hex in the aggregated pre-vote. 

If voting from a voting delegate, set "validator" to the address of the validator to vote on behalf of:
$ merliond tx oracle aggregate-vote 1234 alion:1.234,uusd:0.99 mervaloper1...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			salt := args[0]
			exchangeRatesStr := args[1]
			_, err = types.ParseExchangeRateTuples(exchangeRatesStr)
			if err != nil {
				return fmt.Errorf("given exchange_rate {%s} is not a valid format; exchange rate should be formatted as DecCoin; %s", exchangeRatesStr, err.Error())
			}

			// Get from address
			voter := clientCtx.GetFromAddress()

			// By default, the voter is voting on behalf of itself
			validator := sdk.ValAddress(voter)

			// Override validator if validator is given
			if len(args) == 3 {
				parsedVal, err := sdk.ValAddressFromBech32(args[2])
				if err != nil {
					return errors.Wrap(err, "validator address is invalid")
				}
				validator = parsedVal
			}

			msg := types.NewMsgAggregateExchangeRateVote(salt, exchangeRatesStr, voter, validator)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdDelegateFeederPermission will create a MsgDelegateFeedConsent tx and sign it with the given key.
func CmdDelegateFeederPermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-feeder [feeder]",
		Args:  cobra.ExactArgs(1),
		Short: "Delegate the permission to vote for the oracle to an address",
		Long: strings.TrimSpace(`
Delegate the permission to submit exchange rate votes for the oracle to an address.

Delegation can keep your validator operator key offline and use a separate replaceable key online.

$ merliond tx oracle set-feeder mer1...

where "mer1..." is the address you want to delegate your voting rights to.
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Get from address
			voter := clientCtx.GetFromAddress()

			// The address the right is being delegated from
			validator := sdk.ValAddress(voter)

			feederStr := args[0]
			feeder, err := sdk.AccAddressFromBech32(feederStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgDelegateFeedConsent(validator, feeder)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
