package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"epileptic/internal/config"
	"epileptic/internal/content"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type fakePinger struct {
	err error
}

func (f fakePinger) Ping(context.Context) error { return f.err }

type emptyRepository struct{}

func (emptyRepository) Create(context.Context, *content.Content) error  { return nil }
func (emptyRepository) List(context.Context) ([]content.Content, error) { return nil, nil }
func (emptyRepository) Get(context.Context, uuid.UUID) (*content.Content, error) {
	return nil, content.ErrNotFound
}
func (emptyRepository) Update(context.Context, *content.Content) error { return nil }
func (emptyRepository) Delete(context.Context, uuid.UUID) error        { return content.ErrNotFound }

func TestHealthz(t *testing.T) {
	router := NewRouter(config.Config{Env: "test"}, fakePinger{}, emptyRepository{})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestReadyzUnavailable(t *testing.T) {
	router := NewRouter(config.Config{Env: "test"}, fakePinger{err: errors.New("down")}, emptyRepository{})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}
