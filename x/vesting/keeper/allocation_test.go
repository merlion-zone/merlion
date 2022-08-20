package keeper_test

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	mertypes "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/vesting/types"
)

func (suite *KeeperTestSuite) TestKeeper_AllocateAtGenesis() {
	suite.SetupTest()
	k := suite.app.VestingKeeper
	receiver := sdk.AccAddress(suite.address.Bytes())
	denom := mertypes.BaseDenom
	amt, _ := new(big.Int).SetString("50000000000000000000000000", 10)
	amt2, _ := new(big.Int).SetString("200000000000000000000000000", 10)
	amt3, _ := new(big.Int).SetString("450000000000000000000000000", 10)
	baseAmt := sdk.NewIntFromBigInt(amt)
	baseAmt2 := sdk.NewIntFromBigInt(amt2)
	baseAmt3 := sdk.NewIntFromBigInt(amt3)
	genState := types.GenesisState{
		Params: types.Params{
			Allocation: types.AllocationAmounts{
				TotalAmount:            sdk.NewInt(10000),
				AirdropAmount:          sdk.NewInt(10000),
				VeVestingAmount:        sdk.NewInt(10000),
				StakingRewardAmount:    sdk.NewInt(10000),
				CommunityPoolAmount:    sdk.NewInt(10000),
				StrategicReserveAmount: sdk.NewInt(10000),
				TeamVestingAmount:      sdk.NewInt(10000),
			},
		},
		AllocationAddresses: types.AllocationAddresses{
			StrategicReserveCustodianAddr: receiver.String(),
		},
	}
	k.AllocateAtGenesis(suite.ctx, genState)

	alloc := genState.Params.Allocation
	// Check Results
	suite.Require().Equal(
		alloc.StakingRewardAmount.Add(baseAmt),
		suite.app.BankKeeper.GetBalance(suite.ctx,
			authtypes.NewModuleAddress(types.StakingRewardVestingName), denom).Amount,
	)
	suite.Require().Equal(
		alloc.CommunityPoolAmount.Add(baseAmt),
		suite.app.BankKeeper.GetBalance(suite.ctx,
			authtypes.NewModuleAddress(types.CommunityPoolVestingName), denom).Amount,
	)
	suite.Require().Equal(
		alloc.TeamVestingAmount.Add(baseAmt2),
		suite.app.BankKeeper.GetBalance(suite.ctx,
			authtypes.NewModuleAddress(types.TeamVestingName), denom).Amount,
	)
	suite.Require().Equal(
		alloc.StrategicReserveAmount.Mul(sdk.NewInt(2)),
		suite.app.BankKeeper.GetBalance(suite.ctx, receiver, denom).Amount,
	)

	emission := suite.app.VeKeeper.GetTotalEmission(suite.ctx)
	suite.Require().Equal(alloc.VeVestingAmount.Add(baseAmt3), emission)
}

func (suite *KeeperTestSuite) TestKeeper_GetAllocationAddresses_SetAllocationAddresses() {
	suite.SetupTest()
	k := suite.app.VestingKeeper
	receiver := sdk.AccAddress(suite.address.Bytes())
	addrs := types.AllocationAddresses{
		TeamVestingAddr:               receiver.String(),
		StrategicReserveCustodianAddr: receiver.String(),
	}
	k.SetAllocationAddresses(suite.ctx, addrs)
	suite.Require().Equal(addrs, k.GetAllocationAddresses(suite.ctx))
}

func (suite *KeeperTestSuite) TestKeeper_GetAirdropTotalAmount_SetAirdropTotalAmount() {
	suite.SetupTest()
	k := suite.app.VestingKeeper
	total := sdk.NewInt(100)
	k.SetAirdropTotalAmount(suite.ctx, total)
	suite.Require().Equal(total, k.GetAirdropTotalAmount(suite.ctx))
}

func (suite *KeeperTestSuite) TestKeeper_GetAirdrop_SetAirdrop() {
	suite.SetupTest()
	k := suite.app.VestingKeeper
	receiver := sdk.AccAddress(suite.address.Bytes())
	airdrop := types.Airdrop{
		TargetAddr: receiver.String(),
		Amount:     sdk.NewCoin(mertypes.BaseDenom, sdk.NewInt(1)),
	}
	k.SetAirdrop(suite.ctx, airdrop.GetTargetAddr(), airdrop)
	suite.Require().Equal(airdrop, k.GetAirdrop(suite.ctx, airdrop.GetTargetAddr()))
}

func (suite *KeeperTestSuite) TestKeeper_DeleteAirdrop() {
	suite.SetupTest()
	k := suite.app.VestingKeeper
	receiver := sdk.AccAddress(suite.address.Bytes())
	airdrop := types.Airdrop{
		TargetAddr: receiver.String(),
		Amount:     sdk.NewCoin(mertypes.BaseDenom, sdk.NewInt(1)),
	}
	target := airdrop.GetTargetAddr()
	k.SetAirdrop(suite.ctx, target, airdrop)
	k.DeleteAirdrop(suite.ctx, target)
	suite.Require().Equal(types.Airdrop{}, k.GetAirdrop(suite.ctx, target))
}

func (suite *KeeperTestSuite) TestKeeper_GetAirdropCompleted_SetAirdropCompleted() {
	suite.SetupTest()
	k := suite.app.VestingKeeper
	receiver := sdk.AccAddress(suite.address.Bytes())
	airdrop := types.Airdrop{
		TargetAddr: receiver.String(),
		Amount:     sdk.NewCoin(mertypes.BaseDenom, sdk.NewInt(1)),
	}
	k.SetAirdropCompleted(suite.ctx, airdrop.GetTargetAddr(), airdrop)
	suite.Require().Equal(airdrop, k.GetAirdropCompleted(suite.ctx, airdrop.GetTargetAddr()))
}
