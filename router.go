package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// newRouter wires the Content Management function and the service endpoints of
// the Epileptic application component.
func newRouter(store ContentStore) http.Handler {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/", serviceInfo)
	router.GET("/health", health)

	api := router.Group("/api")
	{
		api.GET("/", serviceInfo)
		api.GET("/content", listContent(store))
		api.POST("/content", createContent(store))
		api.GET("/content/:id", getContent(store))
		api.PUT("/content/:id", updateContent(store))
		api.DELETE("/content/:id", deleteContent(store))
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	return router
}

// serviceInfo describes the Communication Service realised by this component.
func serviceInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service":   serviceName,
		"component": "Epileptic",
		"function":  "Content Management",
		"status":    "ok",
	})
}

// health reports liveness of the application component.
func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
