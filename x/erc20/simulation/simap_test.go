package simulation

import (
	"math/rand"
	"testing"
	"time"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/stretchr/testify/require"
)

func TestFindAccount(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	accs := simtypes.RandomAccounts(r, 5)
	addr := accs[0].Address
	res, ok := FindAccount(accs, addr.String())
	require.Equal(t, true, res.Equals(accs[0]))
	require.Equal(t, true, ok)
}
