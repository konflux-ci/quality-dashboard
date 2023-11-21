package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/redhat-appstudio/quality-studio/api/server"
	"github.com/redhat-appstudio/quality-studio/api/server/middleware"
	version "github.com/redhat-appstudio/quality-studio/api/server/router/version"
	"github.com/redhat-appstudio/quality-studio/pkg/connectors/codecov"
	"github.com/redhat-appstudio/quality-studio/pkg/connectors/gcs"
	"github.com/redhat-appstudio/quality-studio/pkg/connectors/github"
	jiraAPI "github.com/redhat-appstudio/quality-studio/pkg/connectors/jira"
	"github.com/redhat-appstudio/quality-studio/pkg/logger"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/client"
	util "github.com/redhat-appstudio/quality-studio/pkg/utils"
	"github.com/slack-go/slack"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	DefaultGithubTokenEnv = "GITHUB_TOKEN"
	DexIssuerUrl          = "DEX_ISSUER"
	DexApplicationId      = "DEX_APPLICATION_ID"
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
	// nolint:all
	viper.BindPFlags(fs)

	// validate port
	if _, err := strconv.Atoi(viper.GetString("port")); err != nil {
		port, _ := fs.GetInt("port")
		viper.Set("port", strconv.Itoa(port))
	}

	logger, _ := logger.InitZap("level")
	//nolint:all
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

	jiraAPI := jiraAPI.NewJiraConfig()
	githubClient := github.NewGithubClient(util.GetEnv(DefaultGithubTokenEnv, ""))
	g := gcs.BucketHandleClient()

	server := server.New(&server.Config{
		Logger:  logger,
		Version: version.ServerVersion,
		Storage: newKeyCacher(storage, time.Now),
		Github:  githubClient,
		CodeCov: codecov.NewCodeCoverageClient(),
		Jira:    jiraAPI,
		Db:      db,
		GCS:     g,
		Slack:   slack.New(util.GetEnv("SLACK_TOKEN", ""), slack.OptionDebug(true)),
	})

	dexIssuer := os.Getenv("DEX_ISSUER")
	if dexIssuer == "" {
		panic("Dex issuer url is not defined. Use DEX_ISSUER env to define the issuer")
	}

	dexApplication := os.Getenv("DEX_APPLICATION_ID")
	if dexApplication == "" {
		panic("Dex Application is not defined. Use DEX_APPLICATION_ID env to define the appid")
	}

	server.UseMiddleware(middleware.NewAuthenticationMiddleware(dexIssuer, dexApplication))
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
