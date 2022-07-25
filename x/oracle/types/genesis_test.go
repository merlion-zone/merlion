package types_test

import (
	"testing"

	"github.com/merlion-zone/merlion/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisValidation(t *testing.T) {
	genState := types.DefaultGenesis()
	require.NoError(t, genState.Validate())

	genState.Params.VotePeriod = 0
	require.Error(t, genState.Validate())
}
