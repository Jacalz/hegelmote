package remote

import (
	"errors"
	"fmt"
)

var errInputIsZero = errors.New("input 0 is not a valid source")

// Mapping of error values. Index zero corresponds to error 1 and so on.
// The following error codes were reverse engineered by sending incorrect commands.
var errorCodes = [...]error{
	errors.New("malformed command"), // -e.1
	errors.New("unknown command"),   // -e.2
	errors.New("invalid parameter"), // -e.3
}

func errorFromCode(code byte) error {
	if code < '1' || code > '3' {
		return fmt.Errorf("unexpected error code: %d", code)
	}

	return errorCodes[code-'1']
}
