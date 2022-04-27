package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/go-bip39"
	"github.com/gogo/protobuf/proto"
	merlion "github.com/merlion-zone/merlion/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
)

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage) printInfo {
	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(info printInfo) error {
	out, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stderr, "%s\n", string(sdk.MustSortJSON(out)))

	return err
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd(mbm module.BasicManager, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			config.P2P.MaxNumInboundPeers = 100
			config.P2P.MaxNumOutboundPeers = 30
			config.Mempool.Size = 10000
			config.StateSync.TrustPeriod = 112 * time.Hour

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("merlion_5000-%v", tmrand.Str(6))
			}

			// Get bip39 mnemonic
			var mnemonic string
			recover, _ := cmd.Flags().GetBool(genutilcli.FlagRecover)
			if recover {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				value, err := input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}

				mnemonic = value
				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
			}

			nodeID, _, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, mnemonic)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			genFile := config.GenesisFile()
			overwrite, _ := cmd.Flags().GetBool(genutilcli.FlagOverwrite)

			if !overwrite && tmos.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}

			appState, err := overwriteDefaultGenState(cdc, mbm.DefaultGenesis(cdc))
			if err != nil {
				return errors.Wrap(err, "Failed to marshall default genesis state")
			}

			genDoc := &types.GenesisDoc{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				genDoc, err = types.GenesisDocFromFile(genFile)
				if err != nil {
					return errors.Wrap(err, "Failed to read genesis doc from file")
				}
			}

			genDoc.ChainID = chainID
			genDoc.Validators = nil
			genDoc.AppState = appState

			if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
				return errors.Wrap(err, "Failed to export genesis file")
			}

			toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return displayInfo(toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(genutilcli.FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().Bool(genutilcli.FlagRecover, false, "provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")

	return cmd
}

func overwriteDefaultGenState(cdc codec.JSONCodec, appState map[string]json.RawMessage) ([]byte, error) {
	var stakingGenState stakingtypes.GenesisState
	cdc.MustUnmarshalJSON(appState[stakingtypes.ModuleName], &stakingGenState)
	stakingGenState.Params.BondDenom = merlion.AttoLionDenom
	appState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingGenState)

	var govGenState govtypes.GenesisState
	cdc.MustUnmarshalJSON(appState[govtypes.ModuleName], &govGenState)
	govGenState.DepositParams.MinDeposit[0].Denom = merlion.AttoLionDenom
	appState[govtypes.ModuleName] = cdc.MustMarshalJSON(&govGenState)

	var crisisGenState crisistypes.GenesisState
	cdc.MustUnmarshalJSON(appState[crisistypes.ModuleName], &crisisGenState)
	crisisGenState.ConstantFee.Denom = merlion.AttoLionDenom
	appState[crisistypes.ModuleName] = cdc.MustMarshalJSON(&crisisGenState)

	var evmGenState evmtypes.GenesisState
	cdc.MustUnmarshalJSON(appState[evmtypes.ModuleName], &evmGenState)
	evmGenState.Params.EvmDenom = merlion.AttoLionDenom
	appState[evmtypes.ModuleName] = cdc.MustMarshalJSON(&evmGenState)

	var mintGenState minttypes.GenesisState
	cdc.MustUnmarshalJSON(appState[minttypes.ModuleName], &mintGenState)
	mintGenState.Params.MintDenom = merlion.AttoLionDenom
	appState[minttypes.ModuleName] = cdc.MustMarshalJSON(&mintGenState)

	return json.MarshalIndent(appState, "", " ")
}

type AppStateParts struct {
	Bank    banktypes.GenesisState    `json:"bank"`
	Staking stakingtypes.GenesisState `json:"staking"`
	Crisis  crisistypes.GenesisState  `json:"crisis"`
	Gov     govtypes.GenesisState     `json:"gov"`
	Evm     evmtypes.GenesisState     `json:"evm"`
}

func (m *AppStateParts) Reset()         { *m = AppStateParts{} }
func (m *AppStateParts) String() string { return proto.CompactTextString(m) }
func (*AppStateParts) ProtoMessage()    {}
