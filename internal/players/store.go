package players

import (
	"context"

	"geogame/internal/locations"
)

type Store interface {
	CreateClient(ctx context.Context, model *ClientStoreModel) error
	UpdateName(ctx context.Context, clientID, name string) error
	UpdateLocation(ctx context.Context, clientID string, point locations.Location) error
	GetClientByEmail(ctx context.Context, emailID string) (*ClientStoreModel, error)
	GetClientByID(ctx context.Context, id string) (*ClientStoreModel, error)
}
