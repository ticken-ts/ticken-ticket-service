package eventerr

const (
	EventNotFoundErrorCode = iota + 200
)

func GetErrMessage(code uint32) string {
	switch code {
	case EventNotFoundErrorCode:
		return "event not found"
	default:
		return "an error has occurred"
	}
}
