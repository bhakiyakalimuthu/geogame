package players

import (
	"context"
	"geogame/internal/locations"
)

var _ Store = (*MockStore)(nil)

type MockStore struct {
	CreateClientFunc     func(model *ClientStoreModel) error
	UpdateNameFunc       func(clientID, name string) error
	UpdateLocationFunc   func(clientID string, point locations.LocationStoreModel) error
	GetClientByEmailFunc func(emailID string) (*ClientStoreModel, error)
	GetClientByIDFunc    func(id string) (*ClientStoreModel, error)
}

func NewMockStore() *MockStore {
	return &MockStore{
		CreateClientFunc: func(model *ClientStoreModel) error {
			return nil
		},
		UpdateNameFunc: func(clientID, name string) error {
			return nil
		},
		UpdateLocationFunc: func(clientID string, point locations.LocationStoreModel) error {
			return nil
		},
		GetClientByEmailFunc: func(emailID string) (model *ClientStoreModel, e error) {
			return &ClientStoreModel{}, nil
		},
		GetClientByIDFunc: func(id string) (model *ClientStoreModel, e error) {
			return &ClientStoreModel{}, nil
		},
	}
}

func (m *MockStore) CreateClient(ctx context.Context, model *ClientStoreModel) error {
	return m.CreateClientFunc(model)
}

func (m *MockStore) UpdateName(ctx context.Context, clientID, name string) error {
	return m.UpdateNameFunc(clientID, name)
}

func (m *MockStore) UpdateLocation(ctx context.Context, clientID string, point locations.LocationStoreModel) error {
	return m.UpdateLocationFunc(clientID, point)
}

func (m *MockStore) GetClientByEmail(ctx context.Context, emailID string) (*ClientStoreModel, error) {
	return m.GetClientByEmailFunc(emailID)
}

func (m *MockStore) GetClientByID(ctx context.Context, id string) (*ClientStoreModel, error) {
	return m.GetClientByIDFunc(id)
}
