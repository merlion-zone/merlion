package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/oracle/types"
)

// GetParams get all parameters as types.Params.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// VotePeriod returns the number of blocks during which voting takes place.
func (k Keeper) VotePeriod(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyVotePeriod, &res)
	return
}

// VoteThreshold returns the minimum percentage of votes that must be received for a ballot to pass.
func (k Keeper) VoteThreshold(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyVoteThreshold, &res)
	return
}

// RewardBand returns the ratio of allowable exchange rate error that a validator can be rewared.
func (k Keeper) RewardBand(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyRewardBand, &res)
	return
}

// RewardDistributionWindow returns the number of vote periods during which seigniorage reward comes in and then is distributed.
func (k Keeper) RewardDistributionWindow(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyRewardDistributionWindow, &res)
	return
}

// SlashFraction returns oracle voting penalty rate.
func (k Keeper) SlashFraction(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeySlashFraction, &res)
	return
}

// SlashWindow returns # of vote period for oracle slashing.
func (k Keeper) SlashWindow(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeySlashWindow, &res)
	return
}

// MinValidPerWindow returns oracle slashing threshold.
func (k Keeper) MinValidPerWindow(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyMinValidPerWindow, &res)
	return
}
