package locations

import (
	"context"
	"errors"
	"testing"

	"gopkg.in/go-playground/assert.v1"

	"go.uber.org/zap"
)

func TestDefaultService_Create(t *testing.T) {

	tests := []struct {
		name     string
		store    Store
		location Location
		want     error
	}{
		{
			name: "success",
			store: &MockStore{
				CreateFunc: func(location LocationStoreModel) error {
					return nil
				},
			},
			location: Location{
				ID: "1",
				GeoPoint: GeoPoint{
					Longitude: 10.1,
					Latitude:  10.1,
				},
				MetaData: MetaData{
					LocationName: "locationName",
					LocationType: "locationType",
				},
			},
			want: nil,
		},
		{
			name: "success",
			store: &MockStore{
				CreateFunc: func(location LocationStoreModel) error {
					return errors.New("failed to insert")
				},
			},
			location: Location{
				ID: "1",
				GeoPoint: GeoPoint{
					Longitude: 10.1,
					Latitude:  10.1,
				},
				MetaData: MetaData{
					LocationName: "locationName",
					LocationType: "locationType",
				},
			},
			want: errors.New("failed to insert"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultService{
				logger: zap.NewNop(),
				store:  tt.store,
			}
			err := d.Create(context.TODO(), tt.location)
			assert.Equal(t, tt.want, err)
		})
	}
}

func TestDefaultService_Update(t *testing.T) {

	tests := []struct {
		name     string
		store    Store
		location Location
		want     error
	}{
		{
			name: "success",
			store: &MockStore{
				UpdateFunc: func(id string, location LocationStoreModel) error {
					return nil
				},
			},
			location: Location{
				ID: "1",
				GeoPoint: GeoPoint{
					Longitude: 10.1,
					Latitude:  10.1,
				},
				MetaData: MetaData{
					LocationName: "locationName",
					LocationType: "locationType",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultService{
				logger: zap.NewNop(),
				store:  tt.store,
			}
			err := d.Update(context.TODO(), tt.location)
			assert.Equal(t, tt.want, err)
		})
	}
}

func TestDefaultService_Delete(t *testing.T) {
	tests := []struct {
		name  string
		store Store
		id    string
		want  error
	}{
		{
			name: "success",
			store: &MockStore{
				DeleteFunc: func(id string) error {
					return nil
				},
			},
			id:   "1",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultService{
				logger: zap.NewNop(),
				store:  tt.store,
			}
			err := d.Delete(context.TODO(), tt.id)
			assert.Equal(t, tt.want, err)
		})
	}
}
