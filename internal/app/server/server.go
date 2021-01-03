package server

import (
	"github.com/Vysogota99/unit-merchant-experience/internal/app/store"
	"github.com/Vysogota99/unit-merchant-experience/internal/app/store/postgres"
)

// Server ...
type Server struct {
	Conf      *Config
	router    *Router
	scheduler *scheduler
	Store     store.Store
}

// NewServer - helper to init server
func NewServer(conf *Config) (*Server, error) {
	return &Server{
		Conf:      conf,
		Store:     postgres.New(conf.dbConnString),
		scheduler: newScheduler(conf.nWorkers),
	}, nil
}

// Start - start the server
func (s *Server) Start() error {
	s.initRouter()
	s.router.Setup().Run(s.Conf.serverPort)
	return nil
}

func (s *Server) initRouter() {
	router := NewRouter(s.Conf.serverPort, s.Store, s.scheduler)
	s.router = router
}
