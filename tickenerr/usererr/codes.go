package usererr

const (
	UserAlreadyExistErrorCode = iota + 100
	RegisterAttendantErrorCode
	PrivateKeyStoreErrorCode
	PrivateKeyRetrieveErrorCode
	CreateWallerErrorCode
	StoreUserInDatabaseErrorCode
	UserNotFoundErrorCode
)

func GetErrMessage(code uint32) string {
	switch code {
	case UserAlreadyExistErrorCode:
		return "user already exists"
	case RegisterAttendantErrorCode:
		return "failed to register attendant"
	case PrivateKeyStoreErrorCode:
		return "failed to store user private key"
	case CreateWallerErrorCode:
		return "failed to create user wallet in public blockchain"
	case StoreUserInDatabaseErrorCode:
		return "failed to store new user in database"
	case UserNotFoundErrorCode:
		return "user not found"
	case PrivateKeyRetrieveErrorCode:
		return "could not retrieve user private key"
	default:
		return "an error has occurred"
	}
}
