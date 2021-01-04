package server

import (
	"database/sql"

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
	db, err := sql.Open("postgres", conf.dbConnString)
	if err != nil {
		return nil, err
	}

	store := postgres.New(db)
	return &Server{
		Conf:      conf,
		Store:     store,
		scheduler: newScheduler(conf.nWorkers, store),
	}, nil
}

// Start - start the server
func (s *Server) Start() error {
	s.scheduler.initPull()
	s.initRouter()
	s.router.Setup().Run(s.Conf.serverPort)
	return nil
}

func (s *Server) initRouter() {
	router := NewRouter(s.Conf.serverPort, s.Store, s.scheduler)
	s.router = router
}
