package cli

import (
	"io/ioutil"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/proto"
	"github.com/merlion-zone/merlion/x/bank/types"
	"github.com/spf13/cobra"
)

func NewSetDenomMetaDataProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-denom-metadata [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a set denom metadata proposal",
		Long: strings.TrimSpace(
			`Submit a set denom metadata along with an initial deposit.
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

			var metadata banktypes.Metadata
			err = parseProposalContent(clientCtx.Codec, args[0], &metadata)
			if err != nil {
				return err
			}

			content := &types.SetDenomMetadataProposal{
				Title:       title,
				Description: description,
				Metadata:    metadata,
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
