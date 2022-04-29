package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	merlion "github.com/merlion-zone/merlion/types"
)

func RegulatedUnixTime(timestamp uint64) uint64 {
	return timestamp / merlion.SecondsPerWeek * merlion.SecondsPerWeek
}

func RegulatedUnixTimeFromNow(ctx sdk.Context, seconds uint64) uint64 {
	timestamp := ctx.BlockTime().Add(time.Duration(seconds) * time.Second).Unix()
	return RegulatedUnixTime(uint64(timestamp))
}

func NextRegulatedUnixTime(timestamp uint64) uint64 {
	if timestamp > MaxUnixTime-merlion.SecondsPerWeek {
		panic("too large unix time")
	}
	next := timestamp + merlion.SecondsPerWeek
	CheckRegulatedUnixTime(next)
	return next
}

func CheckRegulatedUnixTime(timestamp uint64) {
	if RegulatedUnixTime(timestamp) != timestamp {
		panic("invalid regulated unix time")
	}
}
