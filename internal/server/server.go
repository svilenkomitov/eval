package server

import (
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"github.com/svilenkomitov/eval/internal/eval"
	"github.com/svilenkomitov/eval/internal/storage"
	"net/http"
	"strconv"
)

type Server struct {
	server *http.Server
	db     *storage.Database
}

func New(c *Config, db *storage.Database) *Server {
	server := setUpServer(c, db)
	return &Server{
		server: server,
		db:     db,
	}
}

func initRoutes(router *chi.Mux, db *storage.Database) {
	evalHandler := eval.Handler{
		Service: eval.New(eval.NewRepository(db)),
	}
	evalHandler.Routes(router)
}

func setUpServer(c *Config, db *storage.Database) *http.Server {
	router := chi.NewRouter()
	initRoutes(router, db)

	server := &http.Server{
		Addr:    "0.0.0.0:" + strconv.Itoa(int(c.Port)),
		Handler: router,
	}
	return server
}

func (s *Server) Start() error {
	log.Info("starting the HTTP server at addr: ", s.server.Addr)
	if err := s.server.ListenAndServe(); nil != err && err != http.ErrServerClosed {
		log.Errorf("failed to start server: %v", err)
	}
	return nil
}
