package usererr

const (
	UserAlreadyExistErrorCode = iota + 100
	PrivateKeyStoreErrorCode
	CreateWallerErrorCode
	StoreUserInDatabaseErrorCode
	UserNotFoundInDatabaseErrorCode
)

func GetErrMessage(code uint8) string {
	switch code {
	case UserAlreadyExistErrorCode:
		return "user already exists"
	case PrivateKeyStoreErrorCode:
		return "failed to store user private key"
	case CreateWallerErrorCode:
		return "failed to create user wallet in public blockchain"
	case StoreUserInDatabaseErrorCode:
		return "failed to store new user in database"
	case UserNotFoundInDatabaseErrorCode:
		return "user not found in database"
	default:
		return "an error has occurred"
	}
}
