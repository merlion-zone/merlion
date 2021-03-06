#!/usr/bin/env bash

rm -rf modules && mkdir -p modules

# Include the specs from Ethermint
if [ ! -d "ethermint" ]; then
    git clone https://github.com/tharsis/ethermint.git
fi
cp -r ethermint/x/evm/spec/ ./modules/evm
cp -r ethermint/x/feemarket/spec/ ./modules/feemarket

# Include the specs from Cosmos SDK
if [ ! -d "cosmos-sdk" ]; then
    git clone https://github.com/cosmos/cosmos-sdk.git
fi
cp -r cosmos-sdk/x/auth/spec/ ./modules/auth
cp -r cosmos-sdk/x/bank/spec/ ./modules/bank
cp -r cosmos-sdk/x/crisis/spec/ ./modules/crisis
cp -r cosmos-sdk/x/distribution/spec/ ./modules/distribution
cp -r cosmos-sdk/x/gov/spec/ ./modules/gov
cp -r cosmos-sdk/x/slashing/spec/ ./modules/slashing
cp -r cosmos-sdk/x/staking/spec/ ./modules/staking

# Include the specs from IBC go
if [ ! -d "ibc-go" ]; then
    git clone https://github.com/cosmos/ibc-go.git
fi
cp -r ibc-go/modules/apps/transfer/spec/ ./modules/transfer
