package locations

import "context"

var _ Store = (*MockStore)(nil)

type MockStore struct {
	CreateFunc func(location LocationStoreModel) error
	UpdateFunc func(id string, location LocationStoreModel) error
	GetFunc    func(id string) (*LocationStoreModel, error)
	DeleteFunc func(id string) error
}

func NewMockStore() *MockStore {
	return &MockStore{
		CreateFunc: func(location LocationStoreModel) error {
			return nil
		},
		UpdateFunc: func(id string, location LocationStoreModel) error {
			return nil
		},
		GetFunc: func(id string) (model *LocationStoreModel, e error) {
			return &LocationStoreModel{}, nil
		},
		DeleteFunc: func(id string) error {
			return nil
		},
	}
}

func (m *MockStore) Create(ctx context.Context, location LocationStoreModel) error {
	return m.CreateFunc(location)
}

func (m *MockStore) Update(ctx context.Context, id string, location LocationStoreModel) error {
	return m.UpdateFunc(id, location)
}

func (m *MockStore) Get(ctx context.Context, id string) (*LocationStoreModel, error) {
	return m.GetFunc(id)
}

func (m *MockStore) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(id)
}
