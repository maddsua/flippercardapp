package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/url"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	db_gen "github.com/maddsua/flippercardapp/db/generated"

	_ "github.com/mattn/go-sqlite3"
)

func Open(dburl string) (*sql.DB, error) {

	path, query, _ := strings.Cut(dburl, "?")

	opts, err := url.ParseQuery(query)
	if err != nil {
		return nil, fmt.Errorf("parse database url '%s': %v", dburl, err)
	}

	opts.Set("_fk", "true")
	opts.Set("_journal", "WAL")

	dburl = strings.Join([]string{path, opts.Encode()}, "?")

	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		return nil, fmt.Errorf("open database at '%s': %v", dburl, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("open database at '%s': %v", dburl, err)
	}

	return db, nil
}

//go:embed migrations/*
var migrationFS embed.FS

type MigrationState struct {
	Version uint
	Dirty   bool
	Updated bool
}

func Migrate(db *sql.DB) (*MigrationState, error) {

	source, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return nil, err
	}

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, err
	}

	migration, err := migrate.NewWithInstance("iofs", source, "sqlite3", instance)
	if err != nil {
		return nil, err
	}

	var upToDate bool

	if err := migration.Up(); err != nil {
		if err == migrate.ErrNoChange {
			upToDate = true
		} else {
			return nil, err
		}
	}

	version, dirty, err := migration.Version()
	if err != nil {
		return nil, err
	}

	return &MigrationState{
		Version: version,
		Dirty:   dirty,
		Updated: !upToDate,
	}, nil
}

func NewWrapper(db *sql.DB) *Wrapper {
	return &Wrapper{
		Queries: *db_gen.New(db),
		db:      db,
	}
}

type Wrapper struct {
	db_gen.Queries
	db *sql.DB
}

func (db *Wrapper) BeginTx(ctx context.Context) (*TxWrapper, error) {

	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &TxWrapper{
		Queries: *db_gen.New(tx),
		tx:      tx,
	}, nil
}

type TxWrapper struct {
	db_gen.Queries
	tx *sql.Tx
}

func (tx *TxWrapper) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *TxWrapper) Commit() error {
	return tx.tx.Commit()
}

func IsNull(err error) bool {
	return err != nil && err.Error() == "sql: no rows in result set"
}

func IsConflict(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "unique constraint failed")
}
