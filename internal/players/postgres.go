package players

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"geogame/internal/locations"
)

var _ Store = (*Postgres)(nil)

const (
	clientsAllCols = "id, name, email, password, loc_id, ST_AsBinary(point) AS point, loc_name, loc_type"
	clientsTable   = "clients"
)

// Postgres holds the Postgres repository.
type Postgres struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewPostgres instantiates a new PostgreSQL repository.
func NewPostgres(db *sqlx.DB, logger *zap.Logger) *Postgres {
	return &Postgres{
		db:     db,
		logger: logger.Named("geo-game.clients.store"),
	}
}

func (p Postgres) CreateClient(ctx context.Context, model *ClientStoreModel) error {
	stmt := `INSERT INTO clients (
	id,
	name,
	email,
	password,
	loc_id,
	point,
	loc_name,
	loc_type
	) VALUES (
	:id,
	:name,
	:email,
	:password,
	:loc_id,
	:point,
	:loc_name,
	:loc_type
	)`
	_, err := p.db.NamedExecContext(ctx, stmt, *model)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return errors.New("emailID already registered")
			}
		}
	}
	return err
}

func (p Postgres) UpdateName(ctx context.Context, clientID, name string) error {
	stmt := `UPDATE clients SET
	name=$1
	WHERE id=$2
	`
	_, err := p.db.ExecContext(ctx, stmt, name, clientID)
	if err != nil {
		p.logger.Error("UpdateName: failed to update name to db", zap.Error(err))
	}
	return err
}

func (p Postgres) UpdateLocation(ctx context.Context, clientID string, point locations.LocationStoreModel) error {
	stmt := `UPDATE clients SET
	loc_id=$1,
	point=$2,
	loc_name=$3,
	loc_type=$4
	WHERE id=$5
	`
	_, err := p.db.ExecContext(ctx, stmt, point.ID, point.Point, point.LocationName, point.LocationType, clientID)
	if err != nil {
		p.logger.Error("UpdateName: failed to update name to db", zap.Error(err))
	}
	return err
}

func (p Postgres) GetClientByEmail(ctx context.Context, emailID string) (*ClientStoreModel, error) {

	stmt := "SELECT " + clientsAllCols + " FROM " + clientsTable + " WHERE email=$1"
	var c ClientStoreModel
	if err := p.db.GetContext(ctx, &c, stmt, emailID); err != nil {
		if err == sql.ErrNoRows {
			p.logger.Error("GetClientByEmail: client is not found for the provided email", zap.Error(err))
			return nil, err
		}
		p.logger.Error("GetClientByEmail: failed to get client by email from db", zap.Error(err))
		return nil, err
	}
	return &c, nil
}

func (p Postgres) GetClientByID(ctx context.Context, id string) (*ClientStoreModel, error) {
	stmt := "SELECT " + clientsAllCols + " FROM " + clientsTable + " WHERE id=$1"
	var c ClientStoreModel
	if err := p.db.GetContext(ctx, &c, stmt, id); err != nil {
		if err == sql.ErrNoRows {
			p.logger.Error("GetClientByID: client is not found for the provided id", zap.Error(err))
			return nil, err
		}
		p.logger.Error("GetClientByID: failed to get client by id from db", zap.Error(err))
		return nil, err
	}
	return &c, nil
}
