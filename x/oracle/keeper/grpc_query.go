package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/oracle/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) ExchangeRate(c context.Context, req *types.QueryExchangeRateRequest) (*types.QueryExchangeRateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if len(req.Denom) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty denom")
	}

	ctx := sdk.UnwrapSDKContext(c)
	exchangeRate, err := k.GetLionExchangeRate(ctx, req.Denom)
	if err != nil {
		return nil, err
	}

	return &types.QueryExchangeRateResponse{ExchangeRate: exchangeRate}, nil
}

func (k Keeper) ExchangeRates(c context.Context, req *types.QueryExchangeRatesRequest) (*types.QueryExchangeRatesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var exchangeRates sdk.DecCoins
	k.IterateLionExchangeRates(ctx, func(denom string, rate sdk.Dec) (stop bool) {
		exchangeRates = append(exchangeRates, sdk.NewDecCoinFromDec(denom, rate))
		return false
	})

	return &types.QueryExchangeRatesResponse{ExchangeRates: exchangeRates}, nil
}

func (k Keeper) TobinTax(c context.Context, req *types.QueryTobinTaxRequest) (*types.QueryTobinTaxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if len(req.Denom) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty denom")
	}

	ctx := sdk.UnwrapSDKContext(c)
	tobinTax, err := k.GetTobinTax(ctx, req.Denom)
	if err != nil {
		return nil, err
	}

	return &types.QueryTobinTaxResponse{TobinTax: tobinTax}, nil
}

func (k Keeper) TobinTaxes(c context.Context, req *types.QueryTobinTaxesRequest) (*types.QueryTobinTaxesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var tobinTaxes types.DenomList
	k.IterateTobinTaxes(ctx, func(denom string, rate sdk.Dec) (stop bool) {
		tobinTaxes = append(tobinTaxes, types.Denom{
			Name:     denom,
			TobinTax: rate,
		})
		return false
	})

	return &types.QueryTobinTaxesResponse{TobinTaxes: tobinTaxes}, nil
}

func (k Keeper) Actives(c context.Context, req *types.QueryActivesRequest) (*types.QueryActivesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var denoms []string
	k.IterateLionExchangeRates(ctx, func(denom string, rate sdk.Dec) (stop bool) {
		denoms = append(denoms, denom)
		return false
	})

	return &types.QueryActivesResponse{Actives: denoms}, nil
}

func (k Keeper) VoteTargets(c context.Context, req *types.QueryVoteTargetsRequest) (*types.QueryVoteTargetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryVoteTargetsResponse{VoteTargets: k.GetVoteTargets(ctx)}, nil
}

func (k Keeper) FeederDelegation(c context.Context, req *types.QueryFeederDelegationRequest) (*types.QueryFeederDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryFeederDelegationResponse{
		FeederAddr: k.GetFeederDelegation(ctx, valAddr).String(),
	}, nil
}

func (k Keeper) MissCounter(c context.Context, req *types.QueryMissCounterRequest) (*types.QueryMissCounterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryMissCounterResponse{
		MissCounter: k.GetMissCounter(ctx, valAddr),
	}, nil
}

func (k Keeper) AggregatePrevote(c context.Context, req *types.QueryAggregatePrevoteRequest) (*types.QueryAggregatePrevoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	prevote, err := k.GetAggregateExchangeRatePrevote(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	return &types.QueryAggregatePrevoteResponse{
		AggregatePrevote: prevote,
	}, nil
}

func (k Keeper) AggregatePrevotes(c context.Context, req *types.QueryAggregatePrevotesRequest) (*types.QueryAggregatePrevotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var prevotes []types.AggregateExchangeRatePrevote
	k.IterateAggregateExchangeRatePrevotes(ctx, func(_ sdk.ValAddress, prevote types.AggregateExchangeRatePrevote) bool {
		prevotes = append(prevotes, prevote)
		return false
	})

	return &types.QueryAggregatePrevotesResponse{
		AggregatePrevotes: prevotes,
	}, nil
}

func (k Keeper) AggregateVote(c context.Context, req *types.QueryAggregateVoteRequest) (*types.QueryAggregateVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	vote, err := k.GetAggregateExchangeRateVote(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	return &types.QueryAggregateVoteResponse{
		AggregateVote: vote,
	}, nil
}

func (k Keeper) AggregateVotes(c context.Context, req *types.QueryAggregateVotesRequest) (*types.QueryAggregateVotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var votes []types.AggregateExchangeRateVote
	k.IterateAggregateExchangeRateVotes(ctx, func(_ sdk.ValAddress, vote types.AggregateExchangeRateVote) bool {
		votes = append(votes, vote)
		return false
	})

	return &types.QueryAggregateVotesResponse{
		AggregateVotes: votes,
	}, nil
}
