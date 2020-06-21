package players

import (
	"context"

	"geogame/internal/locations"
)

var _ Store = (*MemStore)(nil)

type MemStore struct {
	clientMap map[interface{}]*ClientStoreModel
}

func NewMemStore(clientMap map[interface{}]*ClientStoreModel) *MemStore {
	return &MemStore{
		clientMap: clientMap,
	}
}

func (m *MemStore) CreateClient(ctx context.Context, model *ClientStoreModel) error {
	m.clientMap[model.Email] = model
	m.clientMap[model.ID.String()] = model
	return nil
}

func (m *MemStore) UpdateName(ctx context.Context, userID, name string) error {
	client := m.clientMap[userID]
	client.Name = name
	return nil
}

func (m *MemStore) UpdateLocation(ctx context.Context, userID string, point locations.LocationStoreModel) error {
	client := m.clientMap[userID]
	client.LocationID = toNullString(point.ID)
	client.Point = point.Point
	client.LocationType = toNullString(point.LocationType.String())
	client.LocationName = toNullString(point.LocationName)
	return nil
}

func (m *MemStore) GetClientByEmail(ctx context.Context, emailID string) (*ClientStoreModel, error) {
	return m.clientMap[emailID], nil
}

func (m *MemStore) GetClientByID(ctx context.Context, id string) (*ClientStoreModel, error) {
	return m.clientMap[id], nil
}
