package types_test

import (
	"encoding/hex"
	"testing"

	"github.com/merlion-zone/merlion/x/voter/types"
	"github.com/stretchr/testify/require"
)

func TestTotalVotesKey(t *testing.T) {
	key := types.TotalVotesKey()
	require.Equal(t, "01", hex.EncodeToString(key))
}

func TestTotalVotesByUserKey(t *testing.T) {
	key := types.TotalVotesByUserKey(uint64(100))
	require.Equal(t, "020000000000000064", hex.EncodeToString(key))
}

func TestPoolWeightedVotesKey(t *testing.T) {
	key := types.PoolWeightedVotesKey("alion")
	require.Equal(t, "03616c696f6e", hex.EncodeToString(key))
}

func TestPoolWeightedVotesByUserKey(t *testing.T) {
	key := types.PoolWeightedVotesByUserKey(uint64(100), "alion")
	require.Equal(t, "030000000000000064616c696f6e", hex.EncodeToString(key))
}

func TestIndexKey(t *testing.T) {
	key := types.IndexKey()
	require.Equal(t, "05", hex.EncodeToString(key))
}

func TestIndexAtLastUpdatedByGaugeKey(t *testing.T) {
	key := types.IndexAtLastUpdatedByGaugeKey("alion")
	require.Equal(t, "06616c696f6e", hex.EncodeToString(key))
}

func TestClaimableRewardByGaugeKey(t *testing.T) {
	key := types.ClaimableRewardByGaugeKey("alion")
	require.Equal(t, "07616c696f6e", hex.EncodeToString(key))
}
