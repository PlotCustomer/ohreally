// Package server wires the HTTP surface of the Epileptic application
// component, which realizes the Communication Service.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"epileptic/internal/config"
	"epileptic/internal/content"
)

// Pinger reports whether the data store is reachable.
type Pinger interface {
	Ping(ctx context.Context) error
}

// NewRouter builds the HTTP router of the Epileptic application component.
func NewRouter(cfg config.Config, db Pinger, repo content.Repository) *gin.Engine {
	gin.SetMode(cfg.Env)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/readyz", readiness(db))

	content.NewHandler(repo).Register(router.Group("/content"))

	return router
}

func readiness(db Pinger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
