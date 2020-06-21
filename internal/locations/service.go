package locations

import (
	"context"

	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, location Location) error
	Update(ctx context.Context, location Location) error
	Get(ctx context.Context, id string) (*Location, error)
	Delete(ctx context.Context, id string) error
}

var _ Service = (*DefaultService)(nil)

type DefaultService struct {
	logger *zap.Logger
	store  Store
}

func NewDefaultService(logger *zap.Logger, store Store) *DefaultService {
	return &DefaultService{
		logger: logger,
		store:  store,
	}
}

func (d *DefaultService) Create(ctx context.Context, location Location) error {

	loc := LocationStoreModel{
		ID:           location.ID,
		Point:        toPoint(&location.GeoPoint),
		LocationName: location.MetaData.LocationName,
		LocationType: LocationType(location.MetaData.LocationType),
	}
	if err := d.store.Create(ctx, loc); err != nil {
		d.logger.Error("Create: failed to create location to store", zap.Any("location", location), zap.Error(err))
		return err
	}
	d.logger.Info("create", zap.Any("location", loc))
	return nil
}

func (d *DefaultService) Update(ctx context.Context, location Location) error {
	loc := LocationStoreModel{
		ID:           location.ID,
		Point:        toPoint(&location.GeoPoint),
		LocationName: location.MetaData.LocationName,
		LocationType: LocationType(location.MetaData.LocationType),
	}
	if err := d.store.Update(ctx, location.ID, loc); err != nil {
		d.logger.Error("Create: failed to update location to store", zap.Any("location", location), zap.Error(err))
		return err
	}
	d.logger.Info("create", zap.Any("location", loc))
	return nil

}

func (d *DefaultService) Get(ctx context.Context, id string) (*Location, error) {

	loc, err := d.store.Get(ctx, id)
	if err != nil {
		d.logger.Error("Get: failed to get location from store", zap.Any("id", id), zap.Error(err))
		return nil, err
	}
	l := &Location{
		ID: loc.ID,
		GeoPoint: GeoPoint{
			Longitude: loc.Point.Lon(),
			Latitude:  loc.Point.Lat(),
		},
		MetaData: MetaData{
			LocationName: loc.LocationName,
			LocationType: loc.LocationType.String(),
		},
	}
	return l, nil
}

func (d *DefaultService) Delete(ctx context.Context, id string) error {
	if err := d.store.Delete(ctx, id); err != nil {
		d.logger.Error("Delete: failed to delete location", zap.Any("id", id), zap.Error(err))
		return err
	}
	return nil
}
