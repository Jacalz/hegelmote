package remote

import "errors"

// The following error codes were reverse engineered by sending incorrect commands.
var (
	errMalformedCommand = errors.New("malformed command") // -e.1
	errUnknownCommand   = errors.New("unknown command")   // -e.2
	errInvalidParameter = errors.New("invalid parameter") // -e.3
)

// Mapping of error values. Index zero corresponds to error 1 and so on.
var errorCodes = [3]error{errMalformedCommand, errUnknownCommand, errInvalidParameter}
