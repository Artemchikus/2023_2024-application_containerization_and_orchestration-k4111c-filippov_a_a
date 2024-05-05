package storage

import (
	"context"
	"find-ship/config"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	db *pgxpool.Pool
}

var (
	pgInstance *PostgresStorage
	pgOnce     sync.Once
)

func MustPostgresStorage(ctx context.Context, conf *config.DBConfig) (Storage, error) {
	var err error

	pgOnce.Do(func() {
		var pgxConf *pgxpool.Config
		pgxConf, err = pgxpool.ParseConfig("postgres://" + conf.User + ":" + conf.Password + "@" + conf.Host + ":" + conf.Port + "/" + conf.Name + "?search_path=" + conf.SearchPath)
		if err != nil {
			return
		}

		var db *pgxpool.Pool
		db, err = pgxpool.NewWithConfig(ctx, pgxConf)
		if err != nil {
			return
		}

		err = db.Ping(ctx)
		if err != nil {
			return
		}

		pgInstance = &PostgresStorage{db: db}
	})

	return pgInstance, err
}

func (p PostgresStorage) Disconnect() {
	p.db.Close()
}

func (p *PostgresStorage) Reconnect(ctx context.Context) error {
	p.Disconnect()

	db, err := pgxpool.NewWithConfig(ctx, p.db.Config())
	if err != nil {
		return err
	}

	p.db = db

	return nil
}

func (p PostgresStorage) Stat() *pgxpool.Stat {
	return p.db.Stat()
}
