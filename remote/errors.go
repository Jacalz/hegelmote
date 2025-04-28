package remote

import "errors"

// The following error codes were reverse engineered by sending incorrect commands.
var (
	errMalformedCommand = errors.New("malformed command") // -e.1
	errUnknownCommand   = errors.New("unknown command")   // -e.2
	errInvalidParameter = errors.New("invalid parameter") // -e.3
	errUnknownErrorCode = errors.New("received unknown error code")
	errInvalidVolume    = errors.New("invalid volume")
)

// Mapping of error values. Index zero corresponds to error 1 and so on.
var errorCodes = [3]error{errMalformedCommand, errUnknownCommand, errInvalidParameter}

func (c *Control) parseErrorResponse() error {
	buf := [7]byte{}
	_, err := c.Conn.Read(buf[:])
	if err != nil {
		return err
	}

	return parseErrorFromBuffer(buf[:])
}

func parseErrorFromBuffer(buf []byte) error {
	if len(buf) < 4 || buf[1] != 'e' {
		return nil
	}

	code := buf[3] - '1'
	if code > 2 {
		return errUnknownErrorCode
	}

	return errorCodes[code]
}
