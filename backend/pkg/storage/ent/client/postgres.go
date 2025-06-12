package client

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"

	entSQL "entgo.io/ent/dialect/sql"
	"github.com/konflux-ci/quality-dashboard/pkg/storage"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
	util "github.com/konflux-ci/quality-dashboard/pkg/utils"
)

// nolint
const (
	// postgres SSL modes
	pgSSLDisable    = "disable"
	pgSSLRequire    = "require"
	pgSSLVerifyCA   = "verify-ca"
	pgSSLVerifyFull = "verify-full"
)

const (
	PostgresEntHostEnv     = "POSTGRES_ENT_HOST"
	PostgresEntPortEnv     = "POSTGRES_ENT_PORT"
	PostgresEntDatabaseEnv = "POSTGRES_ENT_DATABASE"
	PostgresEntUserEnv     = "POSTGRES_ENT_USER"
	PostgresEntPasswordEnv = "POSTGRES_ENT_PASSWORD"
	DefaultGithubTokenEnv  = "GITHUB_TOKEN"
)

// Postgres options for creating an SQL db.
type Postgres struct {
	NetworkDB
	SSL SSL `json:"ssl"`
}

// Open always returns a new in sqlite3 storage.
func (p *Postgres) Open() (storage.Storage, *sql.DB, error) {
	drv, err := p.driver()
	if err != nil {
		return nil, nil, err
	}

	databaseClient := NewDatabase(
		WithClient(db.NewClient(db.Driver(drv))),
		WithHasher(sha256.New),
		// The default behavior for Postgres transactions is consistent reads, not consistent writes.
		// For each transaction opened, ensure it has the correct isolation level.
		//
		// See: https://www.postgresql.org/docs/9.3/static/sql-set-transaction.html
		WithTxIsolationLevel(sql.LevelSerializable),
	)

	if err := databaseClient.Schema().Create(context.TODO()); err != nil {
		return nil, nil, err
	}

	return databaseClient, drv.DB(), nil
}

func (p *Postgres) driver() (*entSQL.Driver, error) {
	drv, err := entSQL.Open("postgres", p.dsn())
	if err != nil {
		return nil, err
	}

	// set database/sql tunables if configured
	if p.ConnMaxLifetime != 0 {
		drv.DB().SetConnMaxLifetime(time.Duration(p.ConnMaxLifetime) * time.Second)
	}

	if p.MaxIdleConns == 0 {
		drv.DB().SetMaxIdleConns(5)
	} else {
		drv.DB().SetMaxIdleConns(p.MaxIdleConns)
	}

	if p.MaxOpenConns == 0 {
		drv.DB().SetMaxOpenConns(5)
	} else {
		drv.DB().SetMaxOpenConns(p.MaxOpenConns)
	}

	return drv, nil
}

func (p *Postgres) dsn() string {
	// detect host:port for backwards-compatibility
	host, port, err := net.SplitHostPort(p.Host)
	if err != nil {
		// not host:port, probably unix socket or bare address
		host = p.Host
		if p.Port != 0 {
			port = strconv.Itoa(int(p.Port))
		}
	}

	var parameters []string
	addParam := func(key, val string) {
		parameters = append(parameters, fmt.Sprintf("%s=%s", key, val))
	}

	addParam("connect_timeout", strconv.Itoa(p.ConnectionTimeout))

	if host != "" {
		addParam("host", dataSourceStr(host))
	}

	if port != "" {
		addParam("port", port)
	}

	if p.User != "" {
		addParam("user", dataSourceStr(p.User))
	}

	if p.Password != "" {
		addParam("password", dataSourceStr(p.Password))
	}

	if p.Database != "" {
		addParam("dbname", dataSourceStr(p.Database))
	}

	if p.SSL.Mode == "" {
		// Assume the strictest mode if unspecified.
		addParam("sslmode", dataSourceStr(pgSSLVerifyFull))
	} else {
		addParam("sslmode", dataSourceStr(p.SSL.Mode))
	}

	if p.SSL.CAFile != "" {
		addParam("sslrootcert", dataSourceStr(p.SSL.CAFile))
	}

	if p.SSL.KeyFile != "" {
		addParam("sslkey", dataSourceStr(p.SSL.KeyFile))
	}

	return strings.Join(parameters, " ")
}

var strEsc = regexp.MustCompile(`([\\'])`)

func dataSourceStr(str string) string {
	return "'" + strEsc.ReplaceAllString(str, `\$1`) + "'"
}

// GetPostgresConnectionDetails returns postgres configurations from given environments
func GetPostgresConnectionDetails() Postgres {
	return Postgres{
		NetworkDB: NetworkDB{
			Database: util.GetEnv(PostgresEntDatabaseEnv, "postgres"),
			User:     util.GetEnv(PostgresEntUserEnv, "postgres"),
			Password: util.GetEnv(PostgresEntPasswordEnv, "postgres"),
			Host:     util.GetEnv(PostgresEntHostEnv, "localhost"),
			Port:     util.GetPortEnv(PostgresEntPortEnv, 5432),
		},
		SSL: SSL{
			Mode: "disable", // Postgres container doesn't support SSL.
		},
	}
}
