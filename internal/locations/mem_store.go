package locations

import "context"

var _ Store = (*MemStore)(nil)

type MemStore struct {
	locationMap map[interface{}]LocationStoreModel
}

func NewMemStore(locationMap map[interface{}]LocationStoreModel) *MemStore {
	return &MemStore{
		locationMap: locationMap,
	}
}

func (m *MemStore) Create(ctx context.Context, location LocationStoreModel) error {
	m.locationMap[location.ID] = location
	return nil
}

func (m *MemStore) Update(ctx context.Context, id string, location LocationStoreModel) error {
	m.locationMap[id] = location
	return nil
}

func (m *MemStore) Get(ctx context.Context, id string) (*LocationStoreModel, error) {
	l := m.locationMap[id]
	return &l, nil
}

func (m *MemStore) Delete(ctx context.Context, id string) error {
	delete(m.locationMap, id)
	return nil
}
