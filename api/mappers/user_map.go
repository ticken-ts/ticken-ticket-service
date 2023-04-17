package mappers

import (
	"ticken-ticket-service/api/dto"
	"ticken-ticket-service/models"
	"ticken-ticket-service/security/jwt"
)

func MapUserToDTO(user *models.Attendant, email string, profile *jwt.Profile) *dto.User {
	return &dto.User{
		UserID:        user.UUID.String(),
		WalletAddress: user.WalletAddress,
		Profile: &dto.Profile{
			Email:         email,
			FirstName:     profile.FirstName,
			LastName:      profile.LastName,
			EmailVerified: profile.EmailVerified,
		},
	}
}
