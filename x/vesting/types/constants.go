package types

import merlion "github.com/merlion-zone/merlion/types"

const (
	StakingRewardVestingName = "staking_reward_vesting"
	CommunityPoolVestingName = "community_pool_vesting"
	TeamVestingName          = "team_vesting"

	// Strate reserve pool controlled by governance.
	// Not used now, maybe future.
	StrategicReservePoolName = "strategic_reserve_pool"

	StakingRewardVestingTime = merlion.SecondsPer4Years
	CommunityPoolVestingTime = merlion.SecondsPer4Years
	TeamVestingTime          = merlion.SecondsPer4Years

	ClaimVestedPeriod = 10
)
