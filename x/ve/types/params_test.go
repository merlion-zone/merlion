package types

import (
	"testing"

	merlion "github.com/merlion-zone/merlion/types"
	"github.com/stretchr/testify/require"
)

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()
	require.Equal(t, merlion.BaseDenom, params.LockDenom)
}
