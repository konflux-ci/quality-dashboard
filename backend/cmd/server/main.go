package main

import (
	"net"
	"time"

	"github.com/redhat-appstudio/quality-studio/api/apis/codecov"
	"github.com/redhat-appstudio/quality-studio/api/apis/github"
	"github.com/redhat-appstudio/quality-studio/api/server"
	version "github.com/redhat-appstudio/quality-studio/api/server/router/version"
	"github.com/redhat-appstudio/quality-studio/pkg/logger"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent"
	util "github.com/redhat-appstudio/quality-studio/pkg/utils"
	"go.uber.org/zap"
)

const (
	PostgresEntHostEnv     = "POSTGRES_ENT_HOST"
	PostgresEntPortEnv     = "POSTGRES_ENT_PORT"
	PostgresEntDatabaseEnv = "POSTGRES_ENT_DATABASE"
	PostgresEntUserEnv     = "POSTGRES_ENT_USER"
	PostgresEntPasswordEnv = "POSTGRES_ENT_PASSWORD"
)

var (
	listener net.Listener
	err      error
)

type keyCacher struct {
	storage.Storage

	now func() time.Time
}

func main() {
	logger, _ := logger.InitZap("level")
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	listener, err = net.Listen("tcp", ":9898")
	if err != nil {
		logger.Fatal("tcp listen fail")
	}
	defer listener.Close()
	cfg := GetPostgresConnectionDetails()

	storage, err := cfg.Open()
	if err != nil {
		logger.Fatal("Server fail to initialize database connection", zap.Error(err))
	}

	githubAPI := github.NewGithubClient("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")
	codecovApi := codecov.NewCodeCoverageClient()
	server := server.New(&server.Config{Logger: logger, Version: version.ServerVersion, Storage: newKeyCacher(storage, time.Now), Github: githubAPI, CodeCov: codecovApi})
	server.Accept("", listener)

	wait := make(chan error)
	go server.Wait(wait)

	err = <-wait
	if err != nil {
		logger.Error("Server fail", zap.Error(err))
	}
}

// newKeyCacher returns a storage which caches keys so long as the next
func newKeyCacher(s storage.Storage, now func() time.Time) storage.Storage {
	if now == nil {
		now = time.Now
	}
	return &keyCacher{Storage: s, now: now}
}

// Return postgres configurations from given environments
func GetPostgresConnectionDetails() ent.Postgres {
	return ent.Postgres{
		NetworkDB: ent.NetworkDB{
			Database: util.GetEnv(PostgresEntDatabaseEnv, "postgres"),
			User:     util.GetEnv(PostgresEntUserEnv, "postgres"),
			Password: util.GetEnv(PostgresEntPasswordEnv, "postgres"),
			Host:     util.GetEnv(PostgresEntHostEnv, "localhost"),
			Port:     util.GetPortEnv(PostgresEntPortEnv, 5433),
		},
		SSL: ent.SSL{
			Mode: "disable", // Postgres container doesn't support SSL.
		},
	}
}
