package commonerr

const (
	FailedToEstablishConnectionWithPVTBCErrorCode = iota
	ElementNotFoundInDatabase
	FailedToUpdateElement
)

func GetErrMessage(code uint32) string {
	switch code {
	case FailedToEstablishConnectionWithPVTBCErrorCode:
		return "failed to establish connection with the private blockchain"
	case ElementNotFoundInDatabase:
		return "element not found in database"
	case FailedToUpdateElement:
		return "failed to update element in database"
	default:
		return "an error has occurred"
	}
}
