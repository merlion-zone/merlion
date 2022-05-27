package cli

import (
	"context"
	"fmt"
	// "strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/merlion-zone/merlion/x/maker/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group maker queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetAllBackingRiskParamsCmd(),
		GetAllCollateralRiskParamsCmd(),
		GetAllBackingPoolsCmd(),
		GetAllCollateralPoolsCmd(),
		GetBackingPoolCmd(),
		GetCollateralPoolCmd(),
		GetCollateralOfAccountCmd(),
		GetTotalBackingCmd(),
		GetTotalCollateralCmd(),
		GetBackingRatioCmd(),
		GetParamsCmd(),
	)

	return cmd
}

func GetAllBackingRiskParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-backing-risk-params",
		Short: "Gets risk params of all the backing pools",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryAllBackingRiskParamsRequest{}

			res, err := queryClient.AllBackingRiskParams(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetAllCollateralRiskParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-collateral-risk-params",
		Short: "Gets risk params of all the collateral pools",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryAllCollateralRiskParamsRequest{}

			res, err := queryClient.AllCollateralRiskParams(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetAllBackingPoolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-backing-pools",
		Short: "Gets all the backing pools",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryAllBackingPoolsRequest{}

			res, err := queryClient.AllBackingPools(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetAllCollateralPoolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-collateral-pools",
		Short: "Gets all the collateral pools",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryAllCollateralPoolsRequest{}

			res, err := queryClient.AllCollateralPools(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetBackingPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backing-pool [backing_denom]",
		Short: "Gets a backing pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryBackingPoolRequest{
				BackingDenom: args[0],
			}

			res, err := queryClient.BackingPool(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCollateralPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collateral-pool [collateral_denom]",
		Short: "Gets a collateral pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryCollateralPoolRequest{
				CollateralDenom: args[0],
			}

			res, err := queryClient.CollateralPool(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCollateralOfAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collateral-account [account] [collateral_denom]",
		Short: "Gets an account's collateral",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryCollateralOfAccountRequest{
				Account:         args[0],
				CollateralDenom: args[1],
			}

			res, err := queryClient.CollateralOfAccount(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetTotalBackingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backing-total",
		Short: "Gets the total backing",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryTotalBackingRequest{}

			res, err := queryClient.TotalBacking(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetTotalCollateralCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collateral-total",
		Short: "Gets the total collateral",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryTotalCollateralRequest{}

			res, err := queryClient.TotalCollateral(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetBackingRatioCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collateral-ratio",
		Short: "Gets the backing ratio",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryBackingRatioRequest{}

			res, err := queryClient.BackingRatio(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetParamsCmd() *cobra.Command {
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
