package types

import "errors"

var (
	KeyManagerNotInitError = errors.New("Key manager is not initialized yet ")
	TokenNotSupportError   = errors.New("The Token is not supported ")
)
