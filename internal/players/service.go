package players

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"geogame/internal/locations"
	"geogame/internal/middleware"
)

type Service interface {
	Register(ctx context.Context, payload RegisterPayload) error
	Login(ctx context.Context, payload LoginPayload) (*APIResponse, error)
	UpdateName(ctx context.Context, payload UpdatePayload, clientID string) error
	UpdateLocation(ctx context.Context, payload locations.Location, clientID string) error
	GetLocation(ctx context.Context, clientID string) (*locations.Location, error)
}

var _ Service = (*DefaultService)(nil)

type DefaultService struct {
	logger      *zap.Logger
	store       Store
	dbTimeOut   time.Duration
	tokenSecret string
}

func NewDefaultService(logger *zap.Logger, store Store, dbTimeOut time.Duration, tokenSecret string) *DefaultService {
	return &DefaultService{
		logger:      logger,
		store:       store,
		dbTimeOut:   dbTimeOut,
		tokenSecret: tokenSecret,
	}
}

func (d *DefaultService) Register(ctx context.Context, payload RegisterPayload) error {
	d.logger.Named("Register").With(zap.String("email", payload.Email))
	// generate hash from password
	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	if err != nil {
		d.logger.Error("Register: failed to generate hash from password", zap.String("email", payload.Email), zap.Error(err))
		return err
	}

	// unique clientID
	clientID := uuid.New()

	client := ClientStoreModel{
		ID:           clientID,
		Name:         payload.Name,
		Email:        payload.Email,
		Password:     string(hash),
		LocationID:   toNullString(""),
		LocationName: toNullString(""),
		LocationType: toNullString(""),
	}
	if err := d.store.CreateClient(ctx, &client); err != nil {
		d.logger.Error("Register: failed to create client", zap.String("email", payload.Email), zap.Error(err))
		return err
	}
	d.logger.Info("client store model : ", zap.Any("client", client))
	return nil
}

func (d *DefaultService) Login(ctx context.Context, payload LoginPayload) (*APIResponse, error) {
	d.logger.Named("Login").With(zap.String("email", payload.Email))
	dbCtx, cancel := context.WithTimeout(ctx, d.dbTimeOut)
	defer cancel()
	client, err := d.store.GetClientByEmail(dbCtx, payload.Email)
	if err != nil {
		d.logger.Error("Login: failed to get client from db", zap.String("email", payload.Email), zap.Error(err))
		return nil, err
	}
	// un-hash the hashed password by using password
	if err := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(payload.Password)); err != nil {
		d.logger.Error("Login: failed to un-hash the password", zap.String("email", payload.Email), zap.Error(err))
		return nil, err
	}

	// Generate jwt access token using token secret
	key := middleware.NewJwtKey(d.tokenSecret)
	token, err := key.GenerateToken(client.ID.String())
	if err != nil {
		d.logger.Error("Login: failed to generate token", zap.String("email", payload.Email), zap.Error(err))
		return nil, err
	}
	return &APIResponse{Token: token}, nil

}

func (d *DefaultService) UpdateName(ctx context.Context, payload UpdatePayload, clientID string) error {
	d.logger.Named("Register").With(zap.String("fullName", payload.Name), zap.String("clientID", clientID))
	_, err := uuid.Parse(clientID)
	if err != nil {
		d.logger.Error("UpdateName: failed to parse clientId", zap.String("clientID", clientID), zap.Error(err))
		return err
	}
	dbCtx, cancel := context.WithTimeout(ctx, d.dbTimeOut)
	defer cancel()
	if err := d.store.UpdateName(dbCtx, clientID, payload.Name); err != nil {
		d.logger.Error("UpdateName: failed to updateName to db", zap.String("clientID", clientID), zap.Error(err))
		return err
	}

	return nil
}

func (d *DefaultService) UpdateLocation(ctx context.Context, payload locations.Location, clientID string) error {
	_, err := uuid.Parse(clientID)
	if err != nil {
		d.logger.Error("UpdateLocation: failed to parse clientId", zap.String("clientID", clientID), zap.Error(err))
		return err
	}
	point := locations.LocationStoreModel{
		ID:           payload.ID,
		Point:        locations.NewPoint(payload.GeoPoint.Longitude, payload.GeoPoint.Latitude),
		LocationName: payload.MetaData.LocationName,
		LocationType: locations.LocationType(payload.MetaData.LocationType),
	}

	dbCtx, cancel := context.WithTimeout(ctx, d.dbTimeOut)
	defer cancel()
	if err := d.store.UpdateLocation(dbCtx, clientID, point); err != nil {
		d.logger.Error("UpdateLocation: failed to update location to db", zap.String("clientID", clientID), zap.Error(err))
		return err
	}

	return nil
}

func (d *DefaultService) GetLocation(ctx context.Context, clientID string) (*locations.Location, error) {
	_, err := uuid.Parse(clientID)
	if err != nil {
		d.logger.Error("UpdateLocation: failed to parse clientId", zap.String("clientID", clientID), zap.Error(err))
		return nil, err
	}
	dbCtx, cancel := context.WithTimeout(ctx, d.dbTimeOut)
	defer cancel()
	model, err := d.store.GetClientByID(dbCtx, clientID)
	if err != nil {
		d.logger.Error("UpdateLocation: failed to update location to db", zap.String("clientID", clientID), zap.Error(err))
		return nil, err
	}
	loc := &locations.Location{
		ID: model.LocationID.String,
		GeoPoint: locations.GeoPoint{
			Longitude: model.Point.Lon(),
			Latitude:  model.Point.Lat(),
		},
		MetaData: locations.MetaData{
			LocationName: model.LocationName.String,
			LocationType: model.LocationType.String,
		},
	}
	return loc, nil
}

func toNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}
