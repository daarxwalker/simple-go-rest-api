package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

var (
	ErrParseUriFailed   = errors.New("parse uri failed")
	ErrConnectionFailed = errors.New("connection failed")
	ErrPingFailed       = errors.New("ping failed")
)

func Connect(uri string) (*pgxpool.Pool, error) {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	config, parseUriErr := pgxpool.ParseConfig(uri)
	if parseUriErr != nil {
		return nil, errors.Join(ErrParseUriFailed, parseUriErr)
	}
	pool, createPoolErr := pgxpool.NewWithConfig(c, config)
	if createPoolErr != nil {
		return nil, errors.Join(ErrConnectionFailed, createPoolErr)
	}
	if pingErr := pool.Ping(c); pingErr != nil {
		return nil, errors.Join(ErrPingFailed, pingErr)
	}
	log.Println("🛢️ database(" + config.ConnConfig.Database + "): ✅")
	return pool, nil
}
