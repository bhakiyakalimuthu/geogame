package players

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"

	"geogame/internal/locations"
)

func TestDefaultService_Register(t *testing.T) {

	tests := []struct {
		name    string
		payload RegisterPayload
		wantErr bool
	}{
		{
			name: "success",
			payload: RegisterPayload{
				Name:     "dummyname",
				Email:    "dummy@mail.com",
				Password: "dummypassword",
			},
			wantErr: false,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultService{
				logger:      zap.NewNop(),
				store:       NewMemStore(make(map[interface{}]*ClientStoreModel)),
				dbTimeOut:   time.Second * 10,
				tokenSecret: "",
			}
			err := d.Register(context.TODO(), tt.payload)
			assert.Nil(t, err)
		})
	}
}

func TestDefaultService_Login(t *testing.T) {

	tests := []struct {
		name    string
		payload LoginPayload
		want    *APIResponse
		wantErr bool
		err     error
	}{
		{
			name: "hashing error",
			payload: LoginPayload{
				Email:    "dummy@mail.com",
				Password: "dummypassword",
			},
			wantErr: true,
			err:     errors.New("crypto/bcrypt: hashedSecret too short to be a bcrypted password"),
		},
		// {
		// 	name: "success",
		// 	payload: LoginPayload{
		// 		Email:    "dummy@mail.com",
		// 		Password: "dummypassword",
		// 	},
		// 	wantErr: false,
		// 	err:     errors.New("crypto/bcrypt: hashedSecret too short to be a bcrypted password"),
		// },
	}
	m := &MockStore{
		GetClientByEmailFunc: func(emailID string) (model *ClientStoreModel, e error) {
			return &ClientStoreModel{
				ID:       uuid.New(),
				Name:     "dummy",
				Email:    "dummy@mail.com",
				Password: "dummypassword",
			}, nil
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultService{
				logger:      zap.NewNop(),
				store:       m,
				dbTimeOut:   time.Second * 10,
				tokenSecret: "",
			}
			got, err := d.Login(context.TODO(), tt.payload)
			if tt.wantErr {
				assert.Equal(t, tt.err, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}

		})
	}
}

func TestDefaultService_UpdateName(t *testing.T) {
	tests := []struct {
		name     string
		payload  UpdatePayload
		clientID string
		want     error
	}{
		{
			name: "empty client id",
			payload: UpdatePayload{
				Name: "dummyname",
			},
			want:     errors.New("invalid UUID length: 0"),
			clientID: "",
		},
		{
			name: "empty client id",
			payload: UpdatePayload{
				Name: "dummyname",
			},
			want:     nil,
			clientID: "5f5ec8c1-b900-48f9-bcc8-cb01dba0747d",
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultService{
				logger: zap.NewNop(),
				store: &MockStore{
					UpdateNameFunc: func(clientID, name string) error {
						return nil
					},
				},
				dbTimeOut:   time.Second * 10,
				tokenSecret: "",
			}
			err := d.UpdateName(context.TODO(), tt.payload, tt.clientID)
			assert.Equal(t, tt.want, err)

		})
	}
}

func TestDefaultService_UpdateLocation(t *testing.T) {
	tests := []struct {
		name     string
		payload  locations.Location
		clientID string
		want     error
	}{
		{
			name: "success",
			payload: locations.Location{
				ID: "1",
				GeoPoint: locations.GeoPoint{
					Longitude: 10.1,
					Latitude:  10.1,
				},
				MetaData: locations.MetaData{
					LocationName: "dummy name",
					LocationType: "dummy type",
				},
			},
			clientID: "5f5ec8c1-b900-48f9-bcc8-cb01dba0747d",
		},
		{
			name: "empty client id",
			payload: locations.Location{
				ID: "1",
				GeoPoint: locations.GeoPoint{
					Longitude: 10.1,
					Latitude:  10.1,
				},
				MetaData: locations.MetaData{
					LocationName: "dummy name",
					LocationType: "dummy type",
				},
			},
			want:     errors.New("invalid UUID length: 0"),
			clientID: "",
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultService{
				logger: zap.NewNop(),
				store: &MockStore{
					UpdateLocationFunc: func(clientID string, point locations.LocationStoreModel) error {
						return nil
					},
				},
				dbTimeOut:   time.Second * 10,
				tokenSecret: "",
			}
			err := d.UpdateLocation(context.TODO(), tt.payload, tt.clientID)
			assert.Equal(t, tt.want, err)

		})
	}
}
