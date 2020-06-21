package players

import (
	"database/sql"

	"github.com/google/uuid"

	"geogame/internal/locations"
)

type RegisterPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatePayload struct {
	Name string `json:"name"`
}

type APIResponse struct {
	Token string `json:"token"`
}

type ClientStoreModel struct {
	ID           uuid.UUID       `db:"id"`
	Name         string          `db:"name"`
	Email        string          `db:"email"`
	Password     string          `db:"password"`
	LocationID   sql.NullString  `db:"loc_id""`
	Point        locations.Point `db:"point"`
	LocationName sql.NullString  `db:"loc_name"`
	LocationType sql.NullString  `db:"loc_type"`
}
