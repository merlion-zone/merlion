package types

import (
	fmt "fmt"
	"testing"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	mertypes "github.com/merlion-zone/merlion/types"
	"github.com/stretchr/testify/require"
)

func Test_validateIBC(t *testing.T) {
	meta := banktypes.Metadata{
		Base: mertypes.MicroUSMDenom,
	}

	err := validateIBC(meta)
	require.Nil(t, err)

	meta.Base = "abc/xxx"
	err = validateIBC(meta)
	require.Error(t, err, fmt.Errorf("invalid metadata. %s denomination should be prefixed with the format 'ibc/", meta.Base))

	meta.Base = "ibc/xxx"
	err = validateIBC(meta)
	require.Error(t, err, fmt.Errorf("invalid metadata (Name) for ibc. %s should include channel", meta.Name))

	meta.Name = "channel-xxx"
	err = validateIBC(meta)
	require.Error(t, err, fmt.Errorf("invalid metadata (Symbol) for ibc. %s should include \"ibc\" prefix", meta.Symbol))

	meta.Symbol = "ibcxxx"
	err = validateIBC(meta)
	require.NoError(t, err)
}
