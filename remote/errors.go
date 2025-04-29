package remote

import "errors"

var (
	errUnexpectedResponse = errors.New("unexpected response")
	errInvalidVolume      = errors.New("invalid volume")
	errInputIsZero        = errors.New("source indexing starts at 1")
)

// Mapping of error values. Index zero corresponds to error 1 and so on.
// The following error codes were reverse engineered by sending incorrect commands.
var errorCodes = [...]error{
	errors.New("malformed command"), // -e.1
	errors.New("unknown command"),   // -e.2
	errors.New("invalid parameter"), // -e.3
}

func errorFromCode(code byte) error {
	if code < '1' || code > '3' {
		return errUnexpectedResponse
	}

	return errorCodes[code-'1']
}
