package mappers

import (
	"ticken-ticket-service/api/dto"
	"ticken-ticket-service/models"
)

func MapUserToDTO(user *models.User) *dto.User {
	return &dto.User{
		UserID: user.UUID.String(),
	}
}
