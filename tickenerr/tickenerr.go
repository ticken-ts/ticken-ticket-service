package tickenerr

import "ticken-ticket-service/tickenerr/usererr"

type TickenError struct {
	Message       string
	Code          uint8
	UnderlyingErr error
}

func New(errCode uint8) TickenError {
	return FromError(errCode, nil)
}

func FromError(errCode uint8, underlyingError error) TickenError {
	var message string

	if between(errCode, 100, 200) {
		message = usererr.GetErrMessage(errCode)
	}

	return TickenError{Message: message, Code: errCode, UnderlyingErr: underlyingError}
}

func between(code, min, max uint8) bool {
	return code >= min && code <= max
}

func (ticketErr TickenError) Error() string {
	return string(ticketErr.Code)
}
