package server

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/redhat-appstudio/quality-studio/api/server/middleware"
	"github.com/redhat-appstudio/quality-studio/api/server/router"
	"github.com/redhat-appstudio/quality-studio/api/server/router/database"
	"github.com/redhat-appstudio/quality-studio/api/server/router/failure"
	"github.com/redhat-appstudio/quality-studio/api/server/router/jira"
	"github.com/redhat-appstudio/quality-studio/api/server/router/prow"
	"github.com/redhat-appstudio/quality-studio/api/server/router/repositories"
	"github.com/redhat-appstudio/quality-studio/api/server/router/suites"
	"github.com/redhat-appstudio/quality-studio/api/server/router/teams"
	"github.com/redhat-appstudio/quality-studio/api/server/router/version"
	_ "github.com/redhat-appstudio/quality-studio/docs/swagger"
	"github.com/redhat-appstudio/quality-studio/pkg/connectors/codecov"
	"github.com/redhat-appstudio/quality-studio/pkg/connectors/github"
	jiraAPI "github.com/redhat-appstudio/quality-studio/pkg/connectors/jira"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils/errdefs"
	"github.com/rs/cors"
	"github.com/slack-go/slack"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
	"go.uber.org/zap"

	// Register postgres driver.
	_ "github.com/lib/pq"
)

// @title Quality Studio API
// @version v1alpha
// @description Go microservice API for Quality Studio Server.

// @contact.name Source Code
// @contact.url https://github.com/redhat-appstudio/quality-studio
// @license.name MIT License
// @license.url https://github.com/redhat-appstudio/quality-studio/blob/main/LICENSE

// @host 127.0.0.1:9898
// @BasePath /api/quality/
// @schemes http https

// Config provides the configuration for the API server
type Config struct {
	Logger      *zap.Logger
	Storage     storage.Storage
	Jira        jiraAPI.Jira
	CorsHeaders string
	Version     string
	SocketGroup string
	TLSConfig   *tls.Config
	Github      *github.Github
	CodeCov     *codecov.API
	Db          *sql.DB
	Slack       *slack.Client
}

// HTTPServer contains an instance of http server and the listener.
// srv *http.Server, contains configuration to create an http server and a mux router with all api end points.
// l   net.Listener, is a TCP or Socket listener that dispatches incoming request to the router.
type HTTPServer struct {
	srv *http.Server
	l   net.Listener
}

// Server contains instance details for the server
type Server struct {
	cfg         *Config
	servers     []*HTTPServer
	middlewares []middleware.Middleware
	routers     []router.Router
}

// New returns a new instance of the server based on the specified configuration.
// It allocates resources which will be needed for ServeAPI(ports, unix-sockets).
func New(cfg *Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

// UseMiddleware appends a new middleware to the request chain.
// This needs to be called before the API routes are configured.
func (s *Server) UseMiddleware(m middleware.Middleware) {
	s.middlewares = append(s.middlewares, m)
}

// Accept sets a listener the server accepts connections into.
func (s *Server) Accept(addr string, listeners ...net.Listener) {
	for _, listener := range listeners {
		httpServer := &HTTPServer{
			srv: &http.Server{
				Addr: addr,
			},
			l: listener,
		}
		s.servers = append(s.servers, httpServer)
	}
}

// Close closes servers and thus stop receiving requests
func (s *Server) Close() {
	for _, srv := range s.servers {
		if err := srv.Close(); err != nil {
			s.cfg.Logger.Error("Cannot close stop servers to receive requests", zap.Error(err))
		}
	}
}

// serveAPI loops through all initialized servers and spawns goroutine
// with Serve method for each. It sets createMux() as Handler also.
func (s *Server) serveAPI() error {
	var chErrors = make(chan error, len(s.servers))
	c := cors.New(cors.Options{
		AllowedOrigins:   make([]string, 0),
		AllowedHeaders:   []string{"X-Registry-Auth", "Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "DELETE", "PUT", "OPTIONS"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})
	for _, srv := range s.servers {
		srv.srv.Handler = c.Handler(s.createMux())
		go func(srv *HTTPServer) {
			var err error
			s.cfg.Logger.Info("Server API initialize", zap.String("Port", srv.l.Addr().String()), zap.String("Server Version", s.cfg.Version), zap.String("API Maturity", version.APIMaturity))

			if err = srv.Serve(); err != nil && strings.Contains(err.Error(), "use of closed network connection") {
				err = nil
			}
			chErrors <- err
		}(srv)
	}

	for range s.servers {
		err := <-chErrors
		if err != nil {
			return err
		}
	}
	return nil
}

// Serve starts listening for inbound requests.
func (s *HTTPServer) Serve() error {
	return s.srv.Serve(s.l)
}

// Close closes the HTTPServer from listening for the inbound requests.
func (s *HTTPServer) Close() error {
	return s.l.Close()
}

func (s *Server) makeHTTPHandler(handler httputils.APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		r = r.WithContext(ctx)
		handlerFunc := s.handlerWithGlobalMiddlewares(handler)

		vars := mux.Vars(r)
		if vars == nil {
			vars = make(map[string]string)
		}

		if err := handlerFunc(ctx, w, r, vars); err != nil {
			statusCode := errdefs.GetHTTPErrorStatusCode(err)
			if statusCode >= 500 {
				s.cfg.Logger.Error("Handler for route failed", zap.Error(err), zap.String("Path", r.URL.Path), zap.String("Method", r.Method))
			}
			httputils.MakeErrorHandler(err)(w, r)
		}
	}
}

// InitRouter initializes the list of routers for the server.
// This method also enables the Go profiler.
func (s *Server) InitRouter() {
	s.routers = append(s.routers,
		version.NewRouter(),
		repositories.NewRouter(s.cfg.Storage),
		prow.NewRouter(s.cfg.Storage),
		teams.NewRouter(s.cfg.Storage),
		jira.NewRouter(s.cfg.Storage),
		database.NewRouter(s.cfg.Db),
		suites.NewRouter(s.cfg.Storage),
		failure.NewRouter(s.cfg.Storage))
}

type pageNotFoundError struct{}

func (pageNotFoundError) Error() string {
	return "page not found"
}

func (pageNotFoundError) NotFound() {}

// createMux initializes the main router the server uses.
func (s *Server) createMux() *mux.Router {
	m := mux.NewRouter()
	s.InitRouter()

	s.cfg.Logger.Info("Initializing server routes")
	for _, apiRouter := range s.routers {
		for _, r := range apiRouter.Routes() {
			f := s.makeHTTPHandler(r.Handler())
			s.cfg.Logger.Info("Registering route", zap.String("Path", "/api/quality"+r.Path()), zap.String("Method", r.Method()))
			m.Path(fmt.Sprintf("/api/quality%s", r.Path())).Methods(r.Method()).Handler(f)
		}
	}

	m.PathPrefix("/api/quality" + "/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/api/quality" + "/swagger/doc.json"),
	))

	m.HandleFunc("/api/quality"+"/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		doc, err := swag.ReadDoc()
		if err != nil {
			s.cfg.Logger.Error("/api/quality/"+"swagger error", zap.Error(err), zap.String("path", "/swagger.json"))
		}
		// nolint:all
		w.Write([]byte(doc))
	})

	notFoundHandler := httputils.MakeErrorHandler(pageNotFoundError{})
	m.NotFoundHandler = notFoundHandler
	m.MethodNotAllowedHandler = notFoundHandler

	str := staticRotationStrategy()
	s.startUpdateStorage(context.TODO(), str, time.Now)
	s.SendBugSLIAlerts()

	return m
}

// Wait blocks the server goroutine until it exits.
// It sends an error message if there is any error during
// the API execution.
func (s *Server) Wait(waitChan chan error) {
	if err := s.serveAPI(); err != nil {
		s.cfg.Logger.Error("Server API error", zap.Error(err))
		waitChan <- err
		return
	}
	waitChan <- nil
}
