package api

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/flacatus/qe-dashboard-backend/pkg/api/apis/codecov"
	"github.com/flacatus/qe-dashboard-backend/pkg/api/apis/github"
	_ "github.com/flacatus/qe-dashboard-backend/pkg/api/docs"
	"github.com/flacatus/qe-dashboard-backend/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	// Register postgres driver.
	_ "github.com/lib/pq"

	jiraFactory "github.com/flacatus/qe-dashboard-backend/pkg/api/apis/jira"
)

// @title Quality Backend API
// @version 1.0.0
// @description Simple Api Rest server to monitor github repositories.

// @contact.name Source Code
// @contact.url https://github.com/redhat-appstudio/qe-dashboard-backend

// @license.name MIT License
// @license.url https://github.com/redhat-appstudio/qe-dashboard-backend/blob/master/LICENSE

// @host localhost:9898
// @BasePath /
// @schemes http https

type Config struct {
	HttpServerTimeout         time.Duration `mapstructure:"http-server-timeout"`
	HttpServerShutdownTimeout time.Duration `mapstructure:"http-server-shutdown-timeout"`
	Host                      string        `mapstructure:"host"`
	Port                      string        `mapstructure:"port"`
	H2C                       bool          `mapstructure:"h2c"`
	Storage                   storage.Storage
}

type Server struct {
	router     *mux.Router
	logger     *zap.Logger
	config     *Config
	githubAPI  *github.API
	JiraApi    jiraFactory.Jira
	codecovAPI *codecov.API
	handler    http.Handler
}

func NewServer(config *Config, logger *zap.Logger) (*Server, error) {
	gh := github.NewGitubClient()
	codecov := codecov.NewCodeCoverageClient()
	factory := jiraFactory.NewJiraConfig()

	srv := &Server{
		router:     mux.NewRouter(),
		logger:     logger,
		config:     config,
		githubAPI:  gh,
		codecovAPI: codecov,
		JiraApi:    factory,
	}

	return srv, nil
}

func (s *Server) registerHandlers() {
	s.router.HandleFunc("/api/version", s.versionHandler).Methods("GET")
	s.router.HandleFunc("/api/jira/e2e-known/get", s.getE2eKnownIssues).Methods("GET")
	s.router.HandleFunc("/api/quality/workflows/get", s.listWorkflowsHandler).Methods("GET")
	s.router.HandleFunc("/api/quality/repositories/list", s.listRepositoriesHandler).Methods("GET")
	s.router.HandleFunc("/api/quality/repositories/get/{git_org}/{repo_name}", s.getRepositoryHandler).Methods("GET")
	s.router.HandleFunc("/api/quality/repositories/create", s.repositoriesCreateHandler).Methods("POST")
	s.router.HandleFunc("/api/quality/repositories/delete", s.deleteRepositoryHandler).Methods("DELETE")
	s.router.PathPrefix("/api/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/api/swagger/doc.json"),
	))
	s.router.PathPrefix("/api/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/api/swagger/doc.json"),
	))
	s.router.HandleFunc("/api/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		doc, err := swag.ReadDoc()
		if err != nil {
			s.logger.Error("swagger error", zap.Error(err), zap.String("path", "/api/swagger.json"))
		}
		w.Write([]byte(doc))
	})
}

func (s *Server) registerMiddlewares() {
	httpLogger := NewLoggingMiddleware(s.logger)
	s.router.Use(httpLogger.Handler)
	s.router.Use(versionMiddleware)
}

func (s *Server) ListenAndServe(stopCh <-chan struct{}) {
	s.registerHandlers()
	s.registerMiddlewares()

	c := cors.New(cors.Options{
		AllowedOrigins:   make([]string, 0),
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "DELETE"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})
	if s.config.H2C {
		s.handler = h2c.NewHandler(s.router, &http2.Server{})
	} else {
		//cors.Default().Handler(s.router)
		s.handler = c.Handler(s.router)
	}

	str := staticRotationStrategy()
	s.startUpdateStorage(context.TODO(), str, time.Now)

	// create the http server
	srv := s.startServer()

	// wait for SIGTERM or SIGINT
	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), s.config.HttpServerShutdownTimeout)
	defer cancel()

	s.logger.Info("Shutting down HTTP/HTTPS server", zap.Duration("timeout", s.config.HttpServerShutdownTimeout))
	// wait for Kubernetes readiness probe to remove this instance from the load balancer
	// the readiness check interval must be lower than the timeout
	if viper.GetString("level") != "debug" {
		time.Sleep(3 * time.Second)
	}

	// determine if the http server was started
	if srv != nil {
		if err := srv.Shutdown(ctx); err != nil {
			s.logger.Warn("HTTP server graceful shutdown failed", zap.Error(err))
		}
	}
}

func (s *Server) startServer() *http.Server {
	// determine if the port is specified
	if s.config.Port == "0" {
		// move on immediately
		return nil
	}

	srv := &http.Server{
		Addr:         s.config.Host + ":" + s.config.Port,
		WriteTimeout: s.config.HttpServerTimeout,
		ReadTimeout:  s.config.HttpServerTimeout,
		IdleTimeout:  2 * s.config.HttpServerTimeout,
		Handler:      s.handler,
	}

	// start the server in the background
	go func() {
		s.logger.Info("Starting HTTP Server.", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatal("HTTP server crashed", zap.Error(err))
		}
	}()

	// return the server and routine
	return srv
}

type ArrayResponse []string
type MapResponse map[string]string
