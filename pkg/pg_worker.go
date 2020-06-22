package pkg

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/voi-oss/svc"
	"go.uber.org/zap"
)

type Option func(*PGWorker)

type Config struct {
	Host string `env:"PG_HOST"                   envDefault:"pg_geogame"`
	Port int    `env:"PG_PORT"                   envDefault:"5432"`
	Name string `env:"PG_DB"                     envDefault:"geo_game_db"`
	User string `env:"PG_USER"                   envDefault:"geogameuser"`
	Pass string `env:"PG_PASS"                   envDefault:"password"`
}

func (c *Config) ConnStr() string {
	connStrParts := []string{
		fmt.Sprintf("host=%s", c.Host),
		fmt.Sprintf("port=%d", c.Port),
		fmt.Sprintf("user=%s", c.User),
		fmt.Sprintf("password=%s", c.Pass),
		fmt.Sprintf("dbname=%s", c.Name),
	}

	return strings.Join(connStrParts, " ")
}

var _ svc.Worker = (*PGWorker)(nil)

type PGWorker struct {
	svcName string
	config  *Config

	db *sqlx.DB
}

func New(svcName string, options ...Option) *PGWorker {
	w := &PGWorker{
		svcName: svcName,
		config:  &Config{},
	}

	for _, opt := range options {
		opt(w)
	}

	return w
}

func Connect(svcName string, options ...Option) (*PGWorker, error) {
	w := New(svcName, options...)
	if err := w.Connect(); err != nil {
		return nil, err
	}
	return w, nil
}

func WithConfig(c *Config) Option {
	return func(w *PGWorker) {
		w.config = c
	}
}

func WithDB(db *sqlx.DB) Option {
	return func(w *PGWorker) {
		w.db = db
	}
}

func (w *PGWorker) Init(logger *zap.Logger) error {
	return nil
}

func (w *PGWorker) Run() error {
	if w.db == nil {
		return fmt.Errorf("no database connection here")
	}
	return nil
}

func (w *PGWorker) Terminate() error {

	db := w.db
	w.db = nil
	return db.Close()
}

func (w *PGWorker) DB() *sqlx.DB {
	return w.db
}

func (w *PGWorker) Connect() error {
	if w.db != nil {
		return nil
	}
	var connector driver.Connector
	var err error
	connector, err = pq.NewConnector(w.config.ConnStr())
	if err != nil {
		return err
	}
	db := sql.OpenDB(connector)
	w.db = sqlx.NewDb(db, "postgres")

	return nil
}
