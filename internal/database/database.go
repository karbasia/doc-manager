package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/karbasia/doc-manager/internal/config"
)

type DB struct {
	Write *sql.DB
	Read  *sql.DB
}

func Open(ctx context.Context, cfg config.DBConfig) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(cfg.Path), 0o755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	write, err := openPool(ctx, dsn(cfg, false), 1)
	if err != nil {
		return nil, fmt.Errorf("open write pool: %w", err)
	}

	read, err := openPool(ctx, dsn(cfg, true), cfg.MaxReadConns)
	if err != nil {
		write.Close()
		return nil, fmt.Errorf("open read pool: %w", err)
	}

	return &DB{Write: write, Read: read}, nil
}

func (db *DB) Close() error {
	return errors.Join(db.Read.Close(), db.Write.Close())
}

func (db *DB) Ping(ctx context.Context) error {
	return errors.Join(db.Write.PingContext(ctx), db.Read.PingContext(ctx))
}

func openPool(ctx context.Context, dsn string, maxConns int) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(maxConns)
	db.SetConnMaxLifetime(0)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func dsn(cfg config.DBConfig, readOnly bool) string {
	pragmas := []string{
		fmt.Sprintf("_pragma=busy_timeout(%d)", cfg.BusyTimeout.Milliseconds()),
		"_pragma=journal_mode(WAL)",
		"_pragma=foreign_keys(ON)",
		"_pragma=synchronous(NORMAL)",
	}

	if readOnly {
		pragmas = append(pragmas, "_pragma=query_only(ON)")
	} else {
		pragmas = append(pragmas, "_txlock=immediate")
	}

	return "file:" + cfg.Path + "?" + strings.Join(pragmas, "&")
}
