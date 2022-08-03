package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestVeIDFromUint64(t *testing.T) {
	id := uint64(10000)
	veID := VeIDFromUint64(id)
	require.Equal(t, "ve-10000", veID)
}

func TestUint64FromVeID(t *testing.T) {
	veID := "xxx"
	val := Uint64FromVeID(veID)
	require.Equal(t, uint64(EmptyVeID), val)

	veID = "ve-xxx"
	val = Uint64FromVeID(veID)
	require.Equal(t, uint64(EmptyVeID), val)

	veID = "ve-10000"
	val = Uint64FromVeID(veID)
	require.Equal(t, uint64(10000), val)
}

func TestVeNftUri(t *testing.T) {
	var (
		nftID     = "10000"
		balance   = sdk.NewInt(10000)
		lockedEnd = uint64(10000)
		value     = sdk.NewInt(10000)
		uri       = VeNftUri(nftID, balance, lockedEnd, value)
	)
	require.Equal(t, "data:application/json;base64,eyJuYW1lIjoibG9jayAjMTAwMDAiLCJkZXNjcmlwdGlvbiI6Ik1lcmxpb24gbG9ja3MsIGNhbiBiZSB1c2VkIHRvIGJvb3N0IGdhdWdlIHlpZWxkcywgdm90ZSBvbiB0b2tlbiBlbWlzc2lvbiwgYW5kIHJlY2VpdmUgYnJpYmVzIiwiaW1hZ2UiOiJkYXRhOmltYWdlL3N2Zyt4bWw7YmFzZTY0LFBITjJaeUI0Yld4dWN6MGlhSFIwY0RvdkwzZDNkeTUzTXk1dmNtY3ZNakF3TUM5emRtY2lJSEJ5WlhObGNuWmxRWE53WldOMFVtRjBhVzg5SW5oTmFXNVpUV2x1SUcxbFpYUWlJSFpwWlhkQ2IzZzlJakFnTUNBek5UQWdNelV3SWo0OGMzUjViR1UtTG1KaGMyVWdleUJtYVd4c09pQjNhR2wwWlRzZ1ptOXVkQzFtWVcxcGJIazZJSE5sY21sbU95Qm1iMjUwTFhOcGVtVTZJREUwY0hnN0lIMDhMM04wZVd4bFBqeHlaV04wSUhkcFpIUm9QU0l4TURBbElpQm9aV2xuYUhROUlqRXdNQ1VpSUdacGJHdzlJbUpzWVdOcklpQXZQangwWlhoMElIZzlJakV3SWlCNVBTSXlNQ0lnWTJ4aGMzTTlJbUpoYzJVaVBuUnZhMlZ1SURFd01EQXdQQzkwWlhoMFBqeDBaWGgwSUhnOUlqRXdJaUI1UFNJME1DSWdZMnhoYzNNOUltSmhjMlVpUG1KaGJHRnVZMlZQWmlBeE1EQXdNRHd2ZEdWNGRENDhkR1Y0ZENCNFBTSXhNQ0lnZVQwaU5qQWlJR05zWVhOelBTSmlZWE5sSWo1c2IyTnJaV1JmWlc1a0lERXdNREF3UEM5MFpYaDBQangwWlhoMElIZzlJakV3SWlCNVBTSTRNQ0lnWTJ4aGMzTTlJbUpoYzJVaVBuWmhiSFZsSURFd01EQXdQQzkwWlhoMFBqd3ZjM1puUGc9PSJ9", uri)
}
