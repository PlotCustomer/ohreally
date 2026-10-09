// Command server starts the Epileptic application component, the realisation
// of the Communication Service of the Communications Platform MVP.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"epileptic.com/internal/config"
	"epileptic.com/internal/content"
	"epileptic.com/internal/database"
	"epileptic.com/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("epileptic: %v", err)
	}
}

func run() error {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := database.Migrate(ctx, pool); err != nil {
		return err
	}

	router := server.NewRouter(server.Deps{
		Content: content.NewHandler(content.NewPostgresRepository(pool)),
		Check: func() error {
			return pool.Ping(ctx)
		},
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Print("shutting down")
	return srv.Shutdown(shutdownCtx)
}
