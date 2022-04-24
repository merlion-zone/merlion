#!/usr/bin/env bash

rm -rf modules && mkdir -p modules

# Include the specs from Ethermint
if [ ! -d "ethermint" ]; then
    git clone https://github.com/tharsis/ethermint.git
fi
mv ethermint/x/evm/spec/ ./modules/evm
mv ethermint/x/feemarket/spec/ ./modules/feemarket

# Include the specs from Cosmos SDK
if [ ! -d "cosmos-sdk" ]; then
    git clone https://github.com/cosmos/cosmos-sdk.git
fi
mv cosmos-sdk/x/auth/spec/ ./modules/auth
mv cosmos-sdk/x/bank/spec/ ./modules/bank
mv cosmos-sdk/x/crisis/spec/ ./modules/crisis
mv cosmos-sdk/x/distribution/spec/ ./modules/distribution
mv cosmos-sdk/x/gov/spec/ ./modules/gov
mv cosmos-sdk/x/slashing/spec/ ./modules/slashing
mv cosmos-sdk/x/staking/spec/ ./modules/staking

# Include the specs from IBC go
if [ ! -d "ibc-go" ]; then
    git clone https://github.com/cosmos/ibc-go.git
fi
mv ibc-go/modules/apps/transfer/spec/ ./modules/transfer
