package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/flacatus/qe-dashboard-backend/pkg/api"
	"github.com/flacatus/qe-dashboard-backend/pkg/signals"
	"github.com/flacatus/qe-dashboard-backend/pkg/storage"
	"github.com/flacatus/qe-dashboard-backend/pkg/storage/ent"
	"github.com/flacatus/qe-dashboard-backend/pkg/utils"
	"github.com/flacatus/qe-dashboard-backend/pkg/version"
)

const (
	PostgresEntHostEnv     = "POSTGRES_ENT_HOST"
	PostgresEntPortEnv     = "POSTGRES_ENT_PORT"
	PostgresEntDatabaseEnv = "POSTGRES_ENT_DATABASE"
	PostgresEntUserEnv     = "POSTGRES_ENT_USER"
	PostgresEntPasswordEnv = "POSTGRES_ENT_PASSWORD"
	SSL                    = "MODE_SSL"
	MAX_CONN               = "MAX_CONNECTIONS"
)

type keyCacher struct {
	storage.Storage

	now func() time.Time
}

func main() {
	// flags definition
	fs := pflag.NewFlagSet("default", pflag.ContinueOnError)
	fs.String("host", "", "Host to bind service to")
	fs.Int("port", 9898, "HTTP port to bind service to")
	fs.String("level", "info", "log level debug, info, warn, error, flat or panic")
	fs.String("config-file", "config/properties.yaml", "quality backend server configurations")
	fs.Duration("http-client-timeout", 2*time.Minute, "client timeout duration")
	fs.Duration("http-server-timeout", 30*time.Second, "server read and write timeout duration")
	fs.Duration("http-server-shutdown-timeout", 5*time.Second, "server graceful shutdown timeout duration")
	fs.String("config", "properties.yaml", "config file name")
	fs.Bool("h2c", false, "allow upgrading to H2C")
	fs.Bool("random-error", false, "1/3 chances of a random response error")

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
		os.Exit(0)
	}

	// bind flags and environment variables
	viper.BindPFlags(fs)
	hostname, _ := os.Hostname()
	viper.Set("hostname", hostname)
	viper.Set("version", version.VERSION)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// configure logging
	logger, _ := initZap(viper.GetString("level"))
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	// validate port
	if _, err := strconv.Atoi(viper.GetString("port")); err != nil {
		port, _ := fs.GetInt("port")
		viper.Set("port", strconv.Itoa(port))
	}

	// load HTTP server config
	var srvCfg api.Config
	if err := viper.Unmarshal(&srvCfg); err != nil {
		logger.Panic("config unmarshal failed", zap.Error(err))
	}

	// log version and port
	logger.Info("Starting qe-dashboard-backend",
		zap.String("version", viper.GetString("version")),
		zap.String("port", srvCfg.Port),
	)

	cfg := GetPostgresConnectionDetails()

	storage, err := cfg.Open()
	if err != nil {
		logger.Fatal("Server fail to initialize database connection", zap.Error(err))
	}
	srvCfg.Storage = newKeyCacher(storage, time.Now)
	// start HTTP server
	srv, _ := api.NewServer(&srvCfg, logger)
	stopCh := signals.SetupSignalHandler()
	srv.ListenAndServe(stopCh)
}

func initZap(logLevel string) (*zap.Logger, error) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	}

	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapConfig := zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zapEncoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return zapConfig.Build()
}

// Return postgres configurations from given environments
func GetPostgresConnectionDetails() ent.Postgres {
	return ent.Postgres{
		NetworkDB: ent.NetworkDB{
			Database:     utils.GetEnv(PostgresEntDatabaseEnv, "postgres"),
			User:         utils.GetEnv(PostgresEntUserEnv, "postgres"),
			Password:     utils.GetEnv(PostgresEntPasswordEnv, "postgres"),
			Host:         utils.GetEnv(PostgresEntHostEnv, "localhost"),
			Port:         utils.GetPortEnv(PostgresEntPortEnv, 5432),
			MaxOpenConns: utils.GetIntEnv(MAX_CONN, 5),
		},
		SSL: ent.SSL{
			Mode: utils.GetEnv(SSL, "disable"), // Control This value by environment because RDS supports SSL
		},
	}
}

// newKeyCacher returns a storage which caches keys so long as the next
func newKeyCacher(s storage.Storage, now func() time.Time) storage.Storage {
	if now == nil {
		now = time.Now
	}
	return &keyCacher{Storage: s, now: now}
}
