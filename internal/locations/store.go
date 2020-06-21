package locations

import "context"

type Store interface {
	Create(ctx context.Context, location LocationStoreModel) error
	Update(ctx context.Context, id string, location LocationStoreModel) error
	Get(ctx context.Context, id string) (*LocationStoreModel, error)
	Delete(ctx context.Context, id string) error
}
