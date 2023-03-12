package dto

type User struct {
	UserID        string   `json:"user_id"`
	WalletAddress string   `json:"wallet_address"`
	Profile       *Profile `json:"profile"`
}

type Profile struct {
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	EmailVerified bool   `json:"email_verified"`
}
