package server

import (
	"github.com/Vysogota99/unit-merchant-experience/internal/app/store"
	"github.com/gin-gonic/gin"
)

// Router ...
type Router struct {
	router     *gin.Engine
	serverPort string
	store      store.Store
}

// NewRouter - helper for initialization http
func NewRouter(serverPort string, store store.Store) *Router {
	return &Router{
		router:     gin.Default(),
		serverPort: serverPort,
		store:      store,
	}
}

// Setup - найстройка роутера
func (r *Router) Setup() *gin.Engine {
	r.router.POST("/file", r.fileHandler)
	return r.router
}
