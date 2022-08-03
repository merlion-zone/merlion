package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegulatedUnixTime(t *testing.T) {
	res := RegulatedUnixTime(1659426807)
	require.Equal(t, uint64(1658966400), res)
}

func TestNextRegulatedUnixTime(t *testing.T) {
	res := NextRegulatedUnixTime(1658966400)
	require.Equal(t, uint64(1659571200), res)
}

func TestPreviousRegulatedUnixTime(t *testing.T) {
	res := PreviousRegulatedUnixTime(1658966400)
	require.Equal(t, uint64(1658361600), res)
}
