package types

// DONTCOVER

import (
	erc20types "github.com/tharsis/evmos/v4/x/erc20/types"
)

// x/erc20 module sentinel errors
var (
	ErrInternalTokenPair      = erc20types.ErrInternalTokenPair
	ErrTokenPairNotFound      = erc20types.ErrTokenPairNotFound
	ErrTokenPairAlreadyExists = erc20types.ErrTokenPairAlreadyExists
	ErrABIPack                = erc20types.ErrABIPack
	ErrABIUnpack              = erc20types.ErrABIUnpack
	ErrEVMDenom               = erc20types.ErrEVMDenom
)
