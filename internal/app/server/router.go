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
	scheduler  *scheduler
}

// NewRouter - helper for initialization http
func NewRouter(serverPort string, store store.Store, scheduler *scheduler) *Router {
	return &Router{
		router:     gin.Default(),
		serverPort: serverPort,
		store:      store,
		scheduler:  scheduler,
	}
}

// Setup - найстройка роутера
func (r *Router) Setup() *gin.Engine {
	r.router.POST("/offer", r.offerHandler)
	r.router.GET("/offer", r.getOfferHandler)
	r.router.GET("/status/:id", r.fileStatusHandler)
	return r.router
}
