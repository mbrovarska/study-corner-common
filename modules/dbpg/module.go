package dbpg

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/fx"

	"study-corner-common/pkg/db"
)

type Params struct {
	fx.In

	DBConfig  db.Config
	Lifecycle fx.Lifecycle
}

type sqlDB struct {
	db *sql.DB
}

func (s *sqlDB) DB() *sql.DB {
	return s.db
}

func New(p Params) (db.SQLDB, error) {
	d, err := sql.Open("postgres", p.DBConfig.DSN)
	if err != nil {
		return nil, err
	}

	d.SetMaxOpenConns(p.DBConfig.MaxOpenConns)
	d.SetMaxIdleConns(p.DBConfig.MaxIdleConns)
	d.SetConnMaxLifetime(time.Duration(p.DBConfig.ConnMaxLifetimeSeconds) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := d.PingContext(ctx); err != nil {
		_ = d.Close()
		return nil, err
	}

	p.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return d.Close()
		},
	})

	return &sqlDB{db: d}, nil
}

var Module = fx.Module(
	"dbpg",
	fx.Provide(New),
)