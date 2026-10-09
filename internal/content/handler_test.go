package content

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type fakeRepository struct {
	items map[uuid.UUID]Content
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{items: make(map[uuid.UUID]Content)}
}

func (f *fakeRepository) Create(_ context.Context, c *Content) error {
	f.items[c.ID] = *c
	return nil
}

func (f *fakeRepository) List(_ context.Context) ([]Content, error) {
	items := make([]Content, 0, len(f.items))
	for _, c := range f.items {
		items = append(items, c)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID.String() < items[j].ID.String() })
	return items, nil
}

func (f *fakeRepository) Get(_ context.Context, id uuid.UUID) (*Content, error) {
	c, ok := f.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &c, nil
}

func (f *fakeRepository) Update(_ context.Context, c *Content) error {
	if _, ok := f.items[c.ID]; !ok {
		return ErrNotFound
	}
	f.items[c.ID] = *c
	return nil
}

func (f *fakeRepository) Delete(_ context.Context, id uuid.UUID) error {
	if _, ok := f.items[id]; !ok {
		return ErrNotFound
	}
	delete(f.items, id)
	return nil
}

func newTestRouter(repo Repository) *gin.Engine {
	r := gin.New()
	NewHandler(repo).Register(r.Group("/content"))
	return r
}

func TestCreateContent(t *testing.T) {
	repo := newFakeRepository()
	router := newTestRouter(repo)

	body := `{"title":"First post","body":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/content", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if len(repo.items) != 1 {
		t.Fatalf("repository holds %d items, want 1", len(repo.items))
	}
}

func TestCreateContentRequiresTitle(t *testing.T) {
	router := newTestRouter(newFakeRepository())

	req := httptest.NewRequest(http.MethodPost, "/content", strings.NewReader(`{"body":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestGetContentNotFound(t *testing.T) {
	router := newTestRouter(newFakeRepository())

	req := httptest.NewRequest(http.MethodGet, "/content/"+uuid.NewString(), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestGetContentInvalidID(t *testing.T) {
	router := newTestRouter(newFakeRepository())

	req := httptest.NewRequest(http.MethodGet, "/content/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestDeleteContent(t *testing.T) {
	repo := newFakeRepository()
	id := uuid.New()
	repo.items[id] = Content{ID: id, Title: "temp"}
	router := newTestRouter(repo)

	req := httptest.NewRequest(http.MethodDelete, "/content/"+id.String(), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if _, ok := repo.items[id]; ok {
		t.Fatal("content was not deleted")
	}
}
