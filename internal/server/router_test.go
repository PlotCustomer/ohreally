package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"epileptic.com/internal/content"
)

type stubRepository struct{}

func (stubRepository) Create(context.Context, *content.Content) error { return nil }
func (stubRepository) Get(context.Context, string) (*content.Content, error) {
	return nil, content.ErrNotFound
}
func (stubRepository) List(context.Context) ([]content.Content, error) { return nil, nil }
func (stubRepository) Update(context.Context, *content.Content) error  { return nil }
func (stubRepository) Delete(context.Context, string) error            { return nil }

func newTestRouter(check func() error) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return NewRouter(Deps{
		Content: content.NewHandler(stubRepository{}),
		Check:   check,
	})
}

func TestHealthz(t *testing.T) {
	router := newTestRouter(nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("/healthz status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestReadyzHealthy(t *testing.T) {
	router := newTestRouter(func() error { return nil })

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("/readyz status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestReadyzUnavailable(t *testing.T) {
	router := newTestRouter(func() error { return errors.New("postgres down") })

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("/readyz status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}

func TestContentRouteRegistered(t *testing.T) {
	router := newTestRouter(nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/content/missing", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("/content/missing status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
