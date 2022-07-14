package cli

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/merlion-zone/merlion/x/maker/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
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
		NewMintBySwapCmd(),
		NewBurnBySwapCmd(),
		NewBuyBackingCmd(),
		NewSellBackingCmd(),
		NewMintByCollateralCmd(),
		NewBurnByCollateralCmd(),
		NewDepositCollateralCmd(),
		NewRedeemCollateralCmd(),
		NewLiquidateCollateralCmd(),
	)

	return cmd
}

func NewMintBySwapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-by-swap [mint_out] [receiver]",
		Short: "Mint by swapping in backing asset and lion coin",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress().String()

			mintOut, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			var receiver string
			if len(args) == 2 {
				receiver = args[1]
				if _, err := sdk.AccAddressFromBech32(receiver); err != nil {
					return fmt.Errorf("invalid receiver bech32 address %w", err)
				}
			} else {
				receiver = sender
			}

			backingInMaxStr, err := cmd.Flags().GetString(FlagBackingInMax)
			if err != nil {
				return err
			}
			backingInMax, err := sdk.ParseCoinNormalized(backingInMaxStr)
			if err != nil {
				return fmt.Errorf("--%s: %w", FlagBackingInMax, err)
			}

			lionInMaxStr, err := cmd.Flags().GetString(FlagLionInMax)
			if err != nil {
				return err
			}
			lionInMax, err := sdk.ParseCoinNormalized(lionInMaxStr)
			if err != nil {
				return fmt.Errorf("--%s: %w", FlagLionInMax, err)
			}

			msg := &types.MsgMintBySwap{
				Sender:       sender,
				To:           receiver,
				MintOutMin:   mintOut,
				BackingInMax: backingInMax,
				LionInMax:    lionInMax,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagBackingInMax, "", "Maximum backing-in coin")
	cmd.Flags().String(FlagLionInMax, "0ulion", "Maximum lion-in coin")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewBurnBySwapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-by-swap [burn_in] [receiver]",
		Short: "Burn by swapping out backing asset and lion coin",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress().String()

			burnIn, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			var receiver string
			if len(args) == 2 {
				receiver = args[1]
				if _, err := sdk.AccAddressFromBech32(receiver); err != nil {
					return fmt.Errorf("invalid receiver bech32 address %w", err)
				}
			} else {
				receiver = sender
			}

			backingOutMinStr, err := cmd.Flags().GetString(FlagBackingOutMin)
			if err != nil {
				return err
			}
			backingOutMin, err := sdk.ParseCoinNormalized(backingOutMinStr)
			if err != nil {
				return fmt.Errorf("--%s: %w", FlagBackingOutMin, err)
			}

			lionOutMinStr, err := cmd.Flags().GetString(FlagLionOutMin)
			if err != nil {
				return err
			}
			lionOutMin, err := sdk.ParseCoinNormalized(lionOutMinStr)
			if err != nil {
				return fmt.Errorf("--%s: %w", FlagLionOutMin, err)
			}

			msg := &types.MsgBurnBySwap{
				Sender:        sender,
				To:            receiver,
				BurnIn:        burnIn,
				BackingOutMin: backingOutMin,
				LionOutMin:    lionOutMin,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagBackingOutMin, "", "Minimum backing-out coin")
	cmd.Flags().String(FlagLionOutMin, "", "Minimum lion-out coin")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewBuyBackingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy-backing [lion_in] [receiver]",
		Short: "Buy backing asset by spending lion coin",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress().String()

			lionIn, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			var receiver string
			if len(args) == 2 {
				receiver = args[1]
				if _, err := sdk.AccAddressFromBech32(receiver); err != nil {
					return fmt.Errorf("invalid receiver bech32 address %w", err)
				}
			} else {
				receiver = sender
			}

			backingOutMinStr, err := cmd.Flags().GetString(FlagBackingOutMin)
			if err != nil {
				return err
			}
			backingOutMin, err := sdk.ParseCoinNormalized(backingOutMinStr)
			if err != nil {
				return fmt.Errorf("--%s: %w", FlagBackingOutMin, err)
			}

			msg := &types.MsgBuyBacking{
				Sender:        sender,
				To:            receiver,
				LionIn:        lionIn,
				BackingOutMin: backingOutMin,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagBackingOutMin, "", "Minimum backing-out coin")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewSellBackingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sell-backing [backing_in] [receiver]",
		Short: "Sell backing asset by earning lion coin",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress().String()

			backingIn, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			var receiver string
			if len(args) == 2 {
				receiver = args[1]
				if _, err := sdk.AccAddressFromBech32(receiver); err != nil {
					return fmt.Errorf("invalid receiver bech32 address %w", err)
				}
			} else {
				receiver = sender
			}

			lionOutMinStr, err := cmd.Flags().GetString(FlagLionOutMin)
			if err != nil {
				return err
			}
			lionOutMin, err := sdk.ParseCoinNormalized(lionOutMinStr)
			if err != nil {
				return fmt.Errorf("--%s: %w", FlagLionOutMin, err)
			}

			msg := &types.MsgSellBacking{
				Sender:     sender,
				To:         receiver,
				BackingIn:  backingIn,
				LionOutMin: lionOutMin,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagLionOutMin, "", "Minimum lion-out coin")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewMintByCollateralCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-by-collateral [collateral_denom] [mint_out] [receiver]",
		Short: "Mint by locking collateral asset and catalytic lion coin",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress().String()

			collateralDenom := args[0]

			mintOut, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			var receiver string
			if len(args) == 2 {
				receiver = args[2]
				if _, err := sdk.AccAddressFromBech32(receiver); err != nil {
					return fmt.Errorf("invalid receiver bech32 address %w", err)
				}
			} else {
				receiver = sender
			}

			msg := &types.MsgMintByCollateral{
				Sender:          sender,
				To:              receiver,
				CollateralDenom: collateralDenom,
				MintOut:         mintOut,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewBurnByCollateralCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-by-collateral [collateral_denom] [repay_in_max]",
		Short: "Burn to repay debt and unlock collateral asset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress().String()

			collateralDenom := args[0]

			repayInMax, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := &types.MsgBurnByCollateral{
				Sender:          sender,
				CollateralDenom: collateralDenom,
				RepayInMax:      repayInMax,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewDepositCollateralCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-collateral [collateral] [receiver]",
		Short: "Deposit collateral asset",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress().String()

			collateral, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			var receiver string
			if len(args) == 2 {
				receiver = args[1]
				if _, err := sdk.AccAddressFromBech32(receiver); err != nil {
					return fmt.Errorf("invalid receiver bech32 address %w", err)
				}
			} else {
				receiver = sender
			}

			msg := &types.MsgDepositCollateral{
				Sender:       sender,
				To:           receiver,
				CollateralIn: collateral,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewRedeemCollateralCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem-collateral [collateral] [lion] [receiver]",
		Short: "Redeem collateral asset",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress().String()

			collateral, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			lion, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			var receiver string
			if len(args) == 3 {
				receiver = args[2]
				if _, err := sdk.AccAddressFromBech32(receiver); err != nil {
					return fmt.Errorf("invalid receiver bech32 address %w", err)
				}
			} else {
				receiver = sender
			}

			msg := &types.MsgRedeemCollateral{
				Sender:        sender,
				To:            receiver,
				CollateralOut: collateral,
				LionOut:       lion,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewLiquidateCollateralCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquidate-collateral [debtor] [collateral] [receiver]",
		Short: "Liquidate debtor's collateral asset",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress().String()

			debtor := args[0]
			if _, err := sdk.AccAddressFromBech32(debtor); err != nil {
				return fmt.Errorf("invalid debtor bech32 address %w", err)
			}

			collateral, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			var receiver string
			if len(args) == 2 {
				receiver = args[2]
				if _, err := sdk.AccAddressFromBech32(receiver); err != nil {
					return fmt.Errorf("invalid receiver bech32 address %w", err)
				}
			} else {
				receiver = sender
			}

			msg := &types.MsgLiquidateCollateral{
				Sender:     sender,
				To:         receiver,
				Debtor:     debtor,
				Collateral: collateral,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewRegisterBackingProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-backing [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a register backing proposal",
		Long: strings.TrimSpace(
			`Submit a register backing proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.`,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			title, description, deposit, err := getProposalArgs(cmd)
			if err != nil {
				return err
			}

			var backingParams types.BackingRiskParams
			err = parseProposalContent(clientCtx.Codec, args[0], &backingParams)
			if err != nil {
				return err
			}

			content := &types.RegisterBackingProposal{
				Title:       title,
				Description: description,
				RiskParams:  backingParams,
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addProposalTxFlagsToCmd(cmd)

	return cmd
}

func NewRegisterCollateralProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-collateral [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a register collateral proposal",
		Long: strings.TrimSpace(
			`Submit a register collateral proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.`,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			title, description, deposit, err := getProposalArgs(cmd)
			if err != nil {
				return err
			}

			var collateralParams types.CollateralRiskParams
			err = parseProposalContent(clientCtx.Codec, args[0], &collateralParams)
			if err != nil {
				return err
			}

			content := &types.RegisterCollateralProposal{
				Title:       title,
				Description: description,
				RiskParams:  collateralParams,
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addProposalTxFlagsToCmd(cmd)

	return cmd
}

func NewSetBackingProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-backing [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a set backing proposal",
		Long: strings.TrimSpace(
			`Submit a set backing proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.`,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			title, description, deposit, err := getProposalArgs(cmd)
			if err != nil {
				return err
			}

			var backingParams types.BackingRiskParams
			err = parseProposalContent(clientCtx.Codec, args[0], &backingParams)
			if err != nil {
				return err
			}

			content := &types.SetBackingRiskParamsProposal{
				Title:       title,
				Description: description,
				RiskParams:  backingParams,
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addProposalTxFlagsToCmd(cmd)

	return cmd
}

func NewSetCollateralProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-collateral [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a set collateral proposal",
		Long: strings.TrimSpace(
			`Submit a set collateral proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.`,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			title, description, deposit, err := getProposalArgs(cmd)
			if err != nil {
				return err
			}

			var collateralParams types.CollateralRiskParams
			err = parseProposalContent(clientCtx.Codec, args[0], &collateralParams)
			if err != nil {
				return err
			}

			content := &types.SetCollateralRiskParamsProposal{
				Title:       title,
				Description: description,
				RiskParams:  collateralParams,
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addProposalTxFlagsToCmd(cmd)

	return cmd
}

func NewBatchSetBackingProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-set-backing [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a batch set backing proposal",
		Long: strings.TrimSpace(
			`Submit a batch set backing proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.`,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			title, description, deposit, err := getProposalArgs(cmd)
			if err != nil {
				return err
			}

			var backingParams types.BatchBackingRiskParams
			err = parseProposalContent(clientCtx.Codec, args[0], &backingParams)
			if err != nil {
				return err
			}

			content := &types.BatchSetBackingRiskParamsProposal{
				Title:       title,
				Description: description,
				RiskParams:  backingParams.RiskParams,
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addProposalTxFlagsToCmd(cmd)

	return cmd
}

func NewBatchSetCollateralProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-set-collateral [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a batch set collateral proposal",
		Long: strings.TrimSpace(
			`Submit a batch set collateral proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.`,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			title, description, deposit, err := getProposalArgs(cmd)
			if err != nil {
				return err
			}

			var collateralParams types.BatchCollateralRiskParams
			err = parseProposalContent(clientCtx.Codec, args[0], &collateralParams)
			if err != nil {
				return err
			}

			content := &types.BatchSetCollateralRiskParamsProposal{
				Title:       title,
				Description: description,
				RiskParams:  collateralParams.RiskParams,
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addProposalTxFlagsToCmd(cmd)

	return cmd
}

func parseProposalContent(cdc codec.JSONCodec, proposalFile string, proposal proto.Message) error {
	content, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return err
	}
	if err = cdc.UnmarshalJSON(content, proposal); err != nil {
		return err
	}
	return nil
}

func getProposalArgs(cmd *cobra.Command) (title, description string, deposit sdk.Coins, err error) {
	title, err = cmd.Flags().GetString(cli.FlagTitle)
	if err != nil {
		return
	}

	description, err = cmd.Flags().GetString(cli.FlagDescription)
	if err != nil {
		return
	}

	depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
	if err != nil {
		return
	}

	deposit, err = sdk.ParseCoinsNormalized(depositStr)
	if err != nil {
		return
	}

	return
}

func addProposalTxFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "1ulion", "deposit of proposal")
	if err := cmd.MarkFlagRequired(cli.FlagTitle); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired(cli.FlagDescription); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired(cli.FlagDeposit); err != nil {
		panic(err)
	}
}

const (
	FlagBackingInMax  = "backing-in-max"
	FlagLionInMax     = "lion-in-max"
	FlagBackingOutMin = "backing-out-min"
	FlagLionOutMin    = "lion-out-min"
)
