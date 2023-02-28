package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/redhat-appstudio/quality-studio/api/apis/codecov"
	"github.com/redhat-appstudio/quality-studio/api/apis/github"
	"github.com/redhat-appstudio/quality-studio/api/server"
	version "github.com/redhat-appstudio/quality-studio/api/server/router/version"
	"github.com/redhat-appstudio/quality-studio/pkg/logger"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/client"
	util "github.com/redhat-appstudio/quality-studio/pkg/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	PostgresEntHostEnv     = "POSTGRES_ENT_HOST"
	PostgresEntPortEnv     = "POSTGRES_ENT_PORT"
	PostgresEntDatabaseEnv = "POSTGRES_ENT_DATABASE"
	PostgresEntUserEnv     = "POSTGRES_ENT_USER"
	PostgresEntPasswordEnv = "POSTGRES_ENT_PASSWORD"
	DefaultGithubTokenEnv  = "GITHUB_TOKEN"
)

const DEFAULT_SERVER_PORT = 9898

var (
	listener net.Listener
)

type keyCacher struct {
	storage.Storage

	now func() time.Time
}

func main() {
	// flags definition
	fs := pflag.NewFlagSet("default", pflag.ContinueOnError)
	fs.Int("port", DEFAULT_SERVER_PORT, "HTTP port to bind service to")

	versionFlag := fs.BoolP("version", "v", false, "get version number")

	// parse flags
	err := fs.Parse(os.Args[1:])
	switch {
	case err == pflag.ErrHelp:
		os.Exit(0)
	case err != nil:
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		fs.PrintDefaults()
		os.Exit(2)
	case *versionFlag:
		fmt.Println(version.ServerVersion)
		os.Exit(0)
	}

	// bind flags and environment variables
	viper.BindPFlags(fs)

	// validate port
	if _, err := strconv.Atoi(viper.GetString("port")); err != nil {
		port, _ := fs.GetInt("port")
		viper.Set("port", strconv.Itoa(port))
	}

	logger, _ := logger.InitZap("level")
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	listener, err = net.Listen("tcp", fmt.Sprintf(":%s", viper.GetString("port")))
	if err != nil {
		logger.Fatal("tcp listen fail")
	}
	defer listener.Close()
	cfg := client.GetPostgresConnectionDetails()

	storage, db, err := cfg.Open()

	if err != nil {
		logger.Fatal("Server fail to initialize database connection", zap.Error(err))
	}

	server := server.New(&server.Config{
		Logger:  logger,
		Version: version.ServerVersion,
		Storage: newKeyCacher(storage, time.Now),
		Github:  github.NewGithubClient(util.GetEnv(DefaultGithubTokenEnv, "")),
		CodeCov: codecov.NewCodeCoverageClient(),
		Db:      db,
	})
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
