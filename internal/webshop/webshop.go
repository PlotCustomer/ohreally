// Package webshop implements the Webshop ApplicationService of DVBD.com. It
// exposes the catalog of the Product DataObject and lets a Klant run the
// Selecteer producten business process.
package webshop

import (
	"context"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"dvbd.com/webshop/internal/klant"
	"dvbd.com/webshop/internal/product"
	"dvbd.com/webshop/internal/selection"
)

// Config holds the dependencies of the Webshop application service.
type Config struct {
	ProductRepository product.Repository
	Address           string
}

// Service is the Webshop application service. It is realised by the DVBD.com
// application component and runs on the Go/Echo system software.
type Service struct {
	echo      *echo.Echo
	config    Config
	selection *selection.Store
}

// New builds the Webshop application service.
func New(cfg Config) *Service {
	if cfg.Address == "" {
		cfg.Address = "0.0.0.0:8080"
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	service := &Service{
		echo:      e,
		config:    cfg,
		selection: selection.NewStore(),
	}
	service.routes()
	return service
}

// Address returns the address the Webshop listens on.
func (s *Service) Address() string { return s.config.Address }

// Start starts the Webshop application service.
func (s *Service) Start() error { return s.echo.Start(s.config.Address) }

// Shutdown stops the Webshop application service.
func (s *Service) Shutdown(ctx context.Context) error { return s.echo.Shutdown(ctx) }

// Handler returns the HTTP handler of the Webshop, useful for tests.
func (s *Service) Handler() http.Handler { return s.echo }

func (s *Service) routes() {
	e := s.echo
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(klant.Middleware())

	// Webshop pages of the ApplicationService.
	e.GET("/", s.home)
	e.GET("/selectie", s.showSelection)
	e.POST("/selecteren", s.selectProduct)
	e.POST("/selectie/verwijderen", s.removeSelectedProduct)

	// Product DataObject access for the DVBD.com application component.
	e.GET("/api/health", s.health)
	e.GET("/api/products", s.listProducts)
	e.GET("/api/products/:id", s.getProduct)

	// Selecteer producten as a JSON API.
	e.GET("/api/selection", s.getSelection)
	e.POST("/api/selection", s.addSelection)
	e.DELETE("/api/selection/:id", s.deleteSelection)

	static, _ := fs.Sub(staticFS, "static")
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(http.FS(static)))))
}
