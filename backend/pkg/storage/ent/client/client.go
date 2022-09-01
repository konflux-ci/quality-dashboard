package client

import (
	"context"
	"database/sql"
	"hash"

	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/migrate"
)

var _ storage.Storage = (*Database)(nil)

type Database struct {
	client    *db.Client
	txOptions *sql.TxOptions

	hasher func() hash.Hash
}

// NewDatabase returns new database client with set options.
func NewDatabase(opts ...func(*Database)) *Database {
	database := &Database{}
	for _, f := range opts {
		f(database)
	}
	return database
}

// WithClient sets client option of a Database object.
func WithClient(c *db.Client) func(*Database) {
	return func(s *Database) {
		s.client = c
	}
}

// WithHasher sets client option of a Database object.
func WithHasher(h func() hash.Hash) func(*Database) {
	return func(s *Database) {
		s.hasher = h
	}
}

// WithTxIsolationLevel sets correct isolation level for database transactions.
func WithTxIsolationLevel(level sql.IsolationLevel) func(*Database) {
	return func(s *Database) {
		s.txOptions = &sql.TxOptions{Isolation: level}
	}
}

// Schema exposes migration schema to perform migrations.
func (d *Database) Schema() *migrate.Schema {
	return d.client.Schema
}

// Close calls the corresponding method of the ent database client.
func (d *Database) Close() error {
	return d.client.Close()
}

// BeginTx is a wrapper to begin transaction with defined options.
func (d *Database) BeginTx(ctx context.Context) (*db.Tx, error) {
	return d.client.BeginTx(ctx, d.txOptions)
}
