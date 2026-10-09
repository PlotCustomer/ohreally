// Package server wires the Communication Service that the Epileptic
// application component realizes.
package server

import (
	"github.com/gin-gonic/gin"

	"epileptic.com/internal/content"
)

// Deps are the collaborators the server needs to expose the service.
type Deps struct {
	// Content is the Content Management handler.
	Content *content.Handler
	// Check reports whether the backing Postgres service is reachable.
	Check func() error
}

// NewRouter builds the HTTP router of the Communication Service.
func NewRouter(deps Deps) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	registerHealth(router, deps.Check)
	deps.Content.Register(router.Group("/"))

	return router
}
