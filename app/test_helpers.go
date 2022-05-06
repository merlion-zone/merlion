package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/simapp"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"github.com/tharsis/ethermint/encoding"
)

// Setup initializes a new Merlion
func Setup(isCheckTx bool) *Merlion {
	SetupConfig()

	db := dbm.NewMemDB()
	app := NewMerlion(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, encoding.MakeConfig(ModuleBasics), simapp.EmptyAppOptions{})

	if !isCheckTx {
		genesisState := NewDefaultGenesisState()

		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		app.InitChain(
			abci.RequestInitChain{
				ChainId:         "merlion_5000-101",
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app.(*Merlion)
}
