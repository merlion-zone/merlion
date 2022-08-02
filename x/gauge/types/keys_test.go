package types

import (
	"encoding/hex"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGaugeKey(t *testing.T) {
	key := GaugeKey("alion")
	require.Equal(t, "01616c696f6e", hex.EncodeToString(key))
}

func TestBribeKey(t *testing.T) {
	key := BribeKey("alion")
	require.Equal(t, "02616c696f6e", hex.EncodeToString(key))
}

func TestTotalDepositedAmountKey(t *testing.T) {
	key := TotalDepositedAmountKey([]byte("gauge"))
	require.Equal(t, "036761756765", hex.EncodeToString(key))
}

func TestDepositedAmountByUserKey(t *testing.T) {
	key := DepositedAmountByUserKey([]byte("gauge"), uint64(1000))
	require.Equal(t, "04676175676500000000000003e8", hex.EncodeToString(key))
}

func TestTotalDerivedAmountKey(t *testing.T) {
	key := TotalDerivedAmountKey([]byte("gaugeKey"))
	require.Equal(t, "0567617567654b6579", hex.EncodeToString(key))
}

func TestDerivedAmountByUserKey(t *testing.T) {
	key := DerivedAmountByUserKey([]byte("gaugeKey"), uint64(1000))
	require.Equal(t, "0667617567654b657900000000000003e8", hex.EncodeToString(key))
}

func TestRewardKeyPrefix(t *testing.T) {
	key := RewardKeyPrefix([]byte("gauge"))
	require.Equal(t, "076761756765", hex.EncodeToString(key))
}

func TestRewardKey(t *testing.T) {
	key := RewardKey([]byte("gauge"), "alion")
	require.Equal(t, "076761756765616c696f6e", hex.EncodeToString(key))
}

func TestUserRewardKey(t *testing.T) {
	key := UserRewardKey([]byte("gauge"), "alion", uint64(1000))
	require.Equal(t, "086761756765616c696f6e00000000000003e8", hex.EncodeToString(key))
}

func TestUserVeIDByAddressKey(t *testing.T) {
	PKS := simapp.CreateTestPubKeys(5)
	addr := sdk.AccAddress(PKS[0].Address())
	key := UserVeIDByAddressKey([]byte("gauge"), addr)
	require.Equal(t, "096761756765dcd3b2e3d86a013b5b5a823b30f8fb791bbc0ea1", hex.EncodeToString(key))
}

func TestEpochKey(t *testing.T) {
	key := EpochKey([]byte("gauge"))
	require.Equal(t, "0a6761756765", hex.EncodeToString(key))
}

func TestPointKey(t *testing.T) {
	key := PointKey([]byte("gauge"), uint64(1000))
	require.Equal(t, "0b676175676500000000000003e8", hex.EncodeToString(key))
}

func TestUserEpochKey(t *testing.T) {
	key := UserEpochKey([]byte("gauge"), uint64(1000))
	require.Equal(t, "0c676175676500000000000003e8", hex.EncodeToString(key))
}

func TestUserPointKey(t *testing.T) {
	key := UserPointKey([]byte("gauge"), uint64(1000), uint64(1000))
	require.Equal(t, "0d676175676500000000000003e800000000000003e8", hex.EncodeToString(key))
}

func TestRewardEpochKey(t *testing.T) {
	key := RewardEpochKey([]byte("gauge"), "alion")
	require.Equal(t, "0e6761756765616c696f6e", hex.EncodeToString(key))
}

func TestRewardPointKey(t *testing.T) {
	key := RewardPointKey([]byte("gauge"), "alion", uint64(1000))
	require.Equal(t, "0f6761756765616c696f6e00000000000003e8", hex.EncodeToString(key))
}
