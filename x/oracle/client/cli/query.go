package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/merlion-zone/merlion/x/oracle/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group oracle queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryExchangeRates(),
		CmdQueryActives(),
		CmdQueryVoteTargets(),
		CmdQueryFeederDelegation(),
		CmdQueryMissCounter(),
		CmdQueryAggregatePrevote(),
		CmdQueryAggregateVote(),
		CmdQueryParams(),
	)
	// this line is used by starport scaffolding # 1

	return cmd
}

// CmdQueryExchangeRates implements the query rate command.
func CmdQueryExchangeRates() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exchange-rates [denom]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query the current exchange rate of an asset w.r.t $uUSD",
		Long: strings.TrimSpace(`
Query the current exchange rate of an asset with an $uUSD. 
You can find the current list of active denoms by running

$ merliond query oracle exchange-rates 

Or, can filter with denom

$ merliond query oracle exchange-rates uusd
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if len(args) == 0 {
				res, err := queryClient.ExchangeRates(context.Background(), &types.QueryExchangeRatesRequest{})
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			}

			denom := args[0]
			res, err := queryClient.ExchangeRate(
				context.Background(),
				&types.QueryExchangeRateRequest{Denom: denom},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryActives implements the query actives command.
func CmdQueryActives() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "actives",
		Args:  cobra.NoArgs,
		Short: "Query the active list of Mer assets recognized by the oracle",
		Long: strings.TrimSpace(`
Query the active list of Mer assets recognized by the types.

$ merliond query oracle actives
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Actives(context.Background(), &types.QueryActivesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryVoteTargets implements the query params command.
func CmdQueryVoteTargets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-targets",
		Args:  cobra.NoArgs,
		Short: "Query the current Oracle vote targets",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.VoteTargets(
				context.Background(),
				&types.QueryVoteTargetsRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryFeederDelegation implements the query feeder delegation command
func CmdQueryFeederDelegation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feeder [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the oracle feeder delegate account",
		Long: strings.TrimSpace(`
Query the account the validator's oracle voting right is delegated to.

$ merliond query oracle feeder mervaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			res, err := queryClient.FeederDelegation(
				context.Background(),
				&types.QueryFeederDelegationRequest{ValidatorAddr: validator.String()},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryMissCounter implements the query miss counter of the validator command
func CmdQueryMissCounter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "miss [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the # of the miss count",
		Long: strings.TrimSpace(`
Query the # of vote periods missed in this oracle slash window.

$ merliond query oracle miss mervaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			res, err := queryClient.MissCounter(
				context.Background(),
				&types.QueryMissCounterRequest{ValidatorAddr: validator.String()},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryAggregatePrevote implements the query aggregate prevote of the validator command
func CmdQueryAggregatePrevote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-prevotes [validator]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query outstanding oracle aggregate prevotes.",
		Long: strings.TrimSpace(`
Query outstanding oracle aggregate prevotes.

$ merliond query oracle aggregate-prevotes

Or, can filter with voter address

$ merliond query oracle aggregate-prevotes mervaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if len(args) == 0 {
				res, err := queryClient.AggregatePrevotes(
					context.Background(),
					&types.QueryAggregatePrevotesRequest{},
				)
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			}

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			res, err := queryClient.AggregatePrevote(
				context.Background(),
				&types.QueryAggregatePrevoteRequest{ValidatorAddr: validator.String()},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryAggregateVote implements the query aggregate prevote of the validator command
func CmdQueryAggregateVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-votes [validator]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query outstanding oracle aggregate votes.",
		Long: strings.TrimSpace(`
Query outstanding oracle aggregate vote.

$ merliond query oracle aggregate-votes 

Or, can filter with voter address

$ merliond query oracle aggregate-votes mervaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if len(args) == 0 {
				res, err := queryClient.AggregateVotes(
					context.Background(),
					&types.QueryAggregateVotesRequest{},
				)
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			}

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			res, err := queryClient.AggregateVote(
				context.Background(),
				&types.QueryAggregateVoteRequest{ValidatorAddr: validator.String()},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
