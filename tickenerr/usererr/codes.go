package usererr

const (
	UserAlreadyExistErrorCode = 100
	PrivateKeyStoreErrorCode  = 101
	CreateWallerErrorCode     = 102
	StoreUserInDatabase       = 103
)

func GetErrMessage(code uint8) string {
	switch code {
	case UserAlreadyExistErrorCode:
		return "user already exists"
	case PrivateKeyStoreErrorCode:
		return "failed to store user private key"
	case CreateWallerErrorCode:
		return "failed to create user wallet in public blockchain"
	case StoreUserInDatabase:
		return "failed to store new user in database"
	default:
		return "an error has occurred"
	}
}
