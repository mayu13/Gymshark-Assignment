package calculate

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mayu13/gymshark-assignment/internal/packs"
	"github.com/sirupsen/logrus"
)

const (
	defaultPort = 8080
)

type Server struct {
	log  *logrus.Entry
	port int
	srv  *http.Server
	pm   packs.PacksManager
}

type Option func(s *Server)

func WithLogger(l *logrus.Entry) Option {
	return func(s *Server) {
		s.log = l
	}
}

func WithPort(p int) Option {
	return func(s *Server) {
		s.port = p
	}
}

func NewServer(pm packs.PacksManager, opts ...Option) (*Server, error) {
	s := Server{
		log:  logrus.NewEntry(logrus.StandardLogger()),
		port: defaultPort,
		pm:   pm,
	}

	for _, opt := range opts {
		opt(&s)
	}

	if err := s.validate(); err != nil {
		return nil, fmt.Errorf("server validation error: %s", err)
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", s.port),
		Handler: handlers.CORS(
			handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type"}),
			handlers.AllowedOrigins([]string{"*"}),
		)(s.routes()),
	}

	s.srv = srv

	return &s, nil
}

func (s *Server) validate() error {
	if s.pm == nil {
		return errors.New("packs manager is empty")
	}

	return nil
}

func (s *Server) routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", HealthCheckHandler()).Methods("GET")
	r.HandleFunc("/health", HealthCheckHandler()).Methods("GET")
	r.HandleFunc("/calculate", CalculatePacksHandler(s.pm)).Methods("POST")
	r.HandleFunc("/packs", SetPackSizesHandler(s.pm)).Methods("POST")

	return r
}

func (s *Server) Start() error {
	s.log.WithFields(
		logrus.Fields{
			"addr": s.srv.Addr,
		},
	).Info("Starting server")

	err := s.srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logrus.WithError(err).Error("Failed to listen and serve")
		return err
	}

	return nil
}

// GracefulStop gracefully shuts down the server without interrupting any
// active connections.
func (s *Server) GracefulStop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func renderStatusCode(w http.ResponseWriter, statusCode int) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return
}
