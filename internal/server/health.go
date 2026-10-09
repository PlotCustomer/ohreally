package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// registerHealth exposes the liveness and readiness probes of the Epileptic
// application component.
func registerHealth(router *gin.Engine, check func() error) {
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/readyz", func(c *gin.Context) {
		if check != nil {
			if err := check(); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status": "unavailable",
					"error":  err.Error(),
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})
}
