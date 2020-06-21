package players

import (
	"github.com/google/uuid"
)
type RegisterPayload struct {
	FullName string `json:"fullName"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type  LoginPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}


type UpdatePayload struct {
	FullName string `json:"fullName"`
}


type APIResponse struct {
	Token string `json:"token"`
}

type ClientStoreModel struct {
	ID uuid.UUID `db:"id"`
	FullName string `db:"full_name"`
	Email string `db:"email"`
	Password string `db:"password"`
	Point string `db:"point"`
	LocationName string `db:"location_name"`
	LocationType string `db:"location_type"`
}