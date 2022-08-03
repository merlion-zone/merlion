package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestNewLockedBalance(t *testing.T) {
	bal := NewLockedBalance()
	require.Equal(t, sdk.ZeroInt(), bal.Amount)
	require.Equal(t, uint64(0), bal.End)
}
