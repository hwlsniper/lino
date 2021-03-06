package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewError - create an error
func NewError(code sdk.CodeType, msg string) sdk.Error {
	return sdk.NewError(LinoErrorCodeSpace, code, msg)
}

// ErrInvalidCoins - error if convert LNO to Coin failed
func ErrInvalidCoins(msg string) sdk.Error {
	return NewError(CodeInvalidCoins, msg)
}
