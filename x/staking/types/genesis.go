package types

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	merlion "github.com/merlion-zone/merlion/types"
)

// DefaultGenesis gets the raw genesis raw message for testing
func DefaultGenesis() *stakingtypes.GenesisState {
	params := stakingtypes.DefaultParams()
	params.BondDenom = merlion.BaseDenom
	return &stakingtypes.GenesisState{
		Params: params,
	}
}
