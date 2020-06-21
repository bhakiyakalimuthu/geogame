package locations

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var _ Store = (*Postgres)(nil)

const (
	locationsAllCols = "loc_id, ST_AsBinary(point) AS point, loc_name, loc_type"
	locationsTable   = "locations"
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
		logger: logger.Named("geo-game.locations.store"),
	}
}

func (p Postgres) Create(ctx context.Context, location LocationStoreModel) error {
	stmt := `INSERT INTO locations (
	loc_id,
	point,
	loc_name,
	loc_type
	) VALUES (
	:loc_id,
	:point,
	:loc_name,
	:loc_type
	)`
	_, err := p.db.NamedExecContext(ctx, stmt, location)
	return err
}

func (p Postgres) Update(ctx context.Context, id string, location LocationStoreModel) error {
	stmt := `UPDATE locations SET
	point=$1,
	loc_name=$2,
	loc_type=$3
	WHERE loc_id=$4
	`
	_, err := p.db.ExecContext(ctx, stmt, location.Point, location.LocationName, location.LocationType, id)
	if err != nil {
		p.logger.Error("UpdateName: failed to update location to db", zap.Error(err))
	}
	return err
}

func (p Postgres) Get(ctx context.Context, id string) (*LocationStoreModel, error) {
	stmt := "SELECT " + locationsAllCols + " FROM " + locationsTable + " WHERE loc_id=$1"
	var c LocationStoreModel
	if err := p.db.GetContext(ctx, &c, stmt, id); err != nil {
		if err == sql.ErrNoRows {
			p.logger.Error("Get: location is not found for the provided id", zap.Error(err))
			return nil, err
		}
		p.logger.Error("Get: failed to get location by id from db", zap.Error(err))
		return nil, err
	}
	return &c, nil
}

func (p Postgres) Delete(ctx context.Context, id string) error {
	stmt := `DELETE FROM locations 
	WHERE loc_id=$1`
	_, err := p.db.ExecContext(ctx, stmt, id)
	if err != nil {
		p.logger.Error("Delete: failed to delete location from db", zap.Error(err))
	}
	return err
}
