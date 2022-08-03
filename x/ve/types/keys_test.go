package types

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTotalLockedAmountKey(t *testing.T) {
	key := TotalLockedAmountKey()
	require.Equal(t, "01", hex.EncodeToString(key))
}

func TestLockedAmountByUserKey(t *testing.T) {
	key := LockedAmountByUserKey(uint64(10000))
	require.Equal(t, "020000000000002710", hex.EncodeToString(key))
}

func TestNextVeIDKey(t *testing.T) {
	key := NextVeIDKey()
	require.Equal(t, "03", hex.EncodeToString(key))
}

func TestEpochKey(t *testing.T) {
	key := EpochKey()
	require.Equal(t, "04", hex.EncodeToString(key))
}

func TestPointKey(t *testing.T) {
	key := PointKey(uint64(10000))
	require.Equal(t, "050000000000002710", hex.EncodeToString(key))
}

func TestUserEpochKey(t *testing.T) {
	key := UserEpochKey(uint64(10000))
	require.Equal(t, "060000000000002710", hex.EncodeToString(key))
}

func TestUserPointKey(t *testing.T) {
	key := UserPointKey(uint64(10000), uint64(10000))
	require.Equal(t, "0700000000000027100000000000002710", hex.EncodeToString(key))
}

func TestSlopeChangeKey(t *testing.T) {
	key := SlopeChangeKey(uint64(1659426807))
	require.Equal(t, "080000000062e8d7f7", hex.EncodeToString(key))
}

func TestAttachedKey(t *testing.T) {
	key := AttachedKey(uint64(10000))
	require.Equal(t, "090000000000002710", hex.EncodeToString(key))
}

func TestVotedKey(t *testing.T) {
	key := VotedKey(uint64(10000))
	require.Equal(t, "0a0000000000002710", hex.EncodeToString(key))
}

func TestTotalEmissionKey(t *testing.T) {
	key := TotalEmissionKey()
	require.Equal(t, "0b", hex.EncodeToString(key))
}

func TestEmissionAtLastPeriodKey(t *testing.T) {
	key := EmissionAtLastPeriodKey()
	require.Equal(t, "0c", hex.EncodeToString(key))
}

func TestEmissionLastTimestampKey(t *testing.T) {
	key := EmissionLastTimestampKey()
	require.Equal(t, "0d", hex.EncodeToString(key))
}

func TestDistributionAccruedLastTimestampKey(t *testing.T) {
	key := DistributionAccruedLastTimestampKey()
	require.Equal(t, "0e", hex.EncodeToString(key))
}

func TestDistributionTotalAmountKey(t *testing.T) {
	key := DistributionTotalAmountKey()
	require.Equal(t, "0f", hex.EncodeToString(key))
}

func TestDistributionPerPeriodKey(t *testing.T) {
	key := DistributionPerPeriodKey(uint64(1659426807))
	require.Equal(t, "100000000062e8d7f7", hex.EncodeToString(key))
}

func TestDistributionClaimLastTimestampByUserKey(t *testing.T) {
	key := DistributionClaimLastTimestampByUserKey(uint64(10000))
	require.Equal(t, "110000000000002710", hex.EncodeToString(key))
}
