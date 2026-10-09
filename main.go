// Command epileptic is the application component of the Communications
// Platform MVP. It realises the Communication Service and is composed of the
// Content Management function. The component is built on go/gin.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// defaultPort is the port the preview of this assignment is served on. It
	// must stay in sync with the preview.port declared in .plot/package.json.
	defaultPort = "8080"

	// serviceName identifies the Communication Service this component realises.
	serviceName = "epileptic.com"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	store, err := newContentStore()
	if err != nil {
		log.Fatalf("%s: content management store: %v", serviceName, err)
	}
	defer store.Close()

	server := &http.Server{
		// Bind the interface of the Agent Machine, not loopback alone, so the
		// delivery is reachable through the preview address.
		Addr:              "0.0.0.0:" + port,
		Handler:           newRouter(store),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("%s: communication service listening on %s", serviceName, server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("%s: server stopped: %v", serviceName, err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("%s: graceful shutdown failed: %v", serviceName, err)
	}
}
