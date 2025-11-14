package dbpg

import (
	"database/sql"
	"study-corner-common/pkg/db"
	"time"

	"go.uber.org/fx"
)

type Params struct{
	fx.In
	DBConfig db.Config
}

type sqlDB struct {
	db *sql.DB
}

func (s *sqlDB) DB() *sql.DB { return s.db } 

func Module() fx.Option {
	return fx.Provide(
		func(p Params) (db.SQLDB, error) {
			d, err := sql.Open("postgres", p.DBConfig.DSN)
			if err != nil {
				return nil, err
			}
			d.SetMaxOpenConns(p.DBConfig.MaxOpenConns)
			d.SetMaxIdleConns(p.DBConfig.MaxIdleConns)
			d.SetConnMaxLifetime(time.Duration(p.DBConfig.ConnMaxLifetimeSeconds) * time.Second)
			
			if err := d.Ping(); err != nil {
				return nil, err
			}
			return &sqlDB{db: d}, nil
		},
	)
}