// Command dvbd starts the DVBD.com application component. DVBD.com realises
// the Webshop application service on the Go/Echo system software and accesses
// the Product data object.
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

	"dvbd.com/webshop/internal/product"
	"dvbd.com/webshop/internal/webshop"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	repo := newProductRepository()
	defer func() {
		if err := repo.Close(); err != nil {
			log.Printf("product: close repository: %v", err)
		}
	}()

	service := webshop.New(webshop.Config{
		ProductRepository: repo,
		Address:           "0.0.0.0:" + port,
	})

	errCh := make(chan error, 1)
	go func() {
		log.Printf("DVBD.com Webshop application service listening on %s", service.Address())
		errCh <- service.Start()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("webshop: %v", err)
		}
	case <-stop:
		log.Printf("DVBD.com Webshop application service shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := service.Shutdown(ctx); err != nil {
			log.Printf("webshop: shutdown: %v", err)
		}
	}
}

// newProductRepository returns the repository of the Product data object. It
// prefers a Postgres database and falls back to an in-memory catalog so the
// Webshop can always be started and previewed.
func newProductRepository() product.Repository {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("POSTGRES_DSN")
	}
	if dsn != "" {
		repo, err := product.NewPostgresRepository(dsn)
		if err != nil {
			log.Printf("product: postgres unavailable (%v); using in-memory catalog", err)
		} else {
			log.Printf("product: using postgres repository")
			return repo
		}
	}
	log.Printf("product: using in-memory catalog")
	return product.NewMemoryRepository(product.SeedCatalog())
}
