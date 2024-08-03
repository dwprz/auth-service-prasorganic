package dto

import (
	"github.com/dwprz/prasorganic-auth-service/src/model/entity"
	"time"
)

type LoginWithGoogleRes struct {
	UserId       string    `json:"user_id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	PhotoProfile string    `json:"photo_profile"`
	Whatsapp     string    `json:"whatsapp"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LoginRes struct {
	Data   *entity.SanitizedUser `json:"data"`
	Tokens *entity.Tokens        `json:"tokens"`
}
