package content

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeRepository struct {
	items map[string]Content
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{items: map[string]Content{}}
}

func (r *fakeRepository) Create(_ context.Context, content *Content) error {
	r.items[content.ID] = *content
	return nil
}

func (r *fakeRepository) Get(_ context.Context, id string) (*Content, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &item, nil
}

func (r *fakeRepository) List(_ context.Context) ([]Content, error) {
	items := make([]Content, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, nil
}

func (r *fakeRepository) Update(_ context.Context, content *Content) error {
	if _, ok := r.items[content.ID]; !ok {
		return ErrNotFound
	}
	current := r.items[content.ID]
	current.Title = content.Title
	current.Body = content.Body
	r.items[content.ID] = current
	return nil
}

func (r *fakeRepository) Delete(_ context.Context, id string) error {
	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	return nil
}

func newTestRouter(repo Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHandler(repo)
	handler.newID = func() string { return "fixed-id" }
	handler.Register(router.Group("/"))
	return router
}

func TestCreateAndGetContent(t *testing.T) {
	router := newTestRouter(newFakeRepository())

	body := []byte(`{"title":"Welcome","body":"Hello"}`)
	req := httptest.NewRequest(http.MethodPost, "/content", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /content status = %d, want %d (%s)", rec.Code, http.StatusCreated, rec.Body.String())
	}

	getReq := httptest.NewRequest(http.MethodGet, "/content/fixed-id", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("GET /content/fixed-id status = %d, want %d", getRec.Code, http.StatusOK)
	}

	var got Content
	if err := json.Unmarshal(getRec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Title != "Welcome" || got.Body != "Hello" {
		t.Fatalf("unexpected content: %+v", got)
	}
}

func TestCreateContentValidation(t *testing.T) {
	router := newTestRouter(newFakeRepository())

	req := httptest.NewRequest(http.MethodPost, "/content", bytes.NewReader([]byte(`{"body":"no title"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("POST /content status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

func TestGetMissingContent(t *testing.T) {
	router := newTestRouter(newFakeRepository())

	req := httptest.NewRequest(http.MethodGet, "/content/does-not-exist", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET missing status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestUpdateAndDeleteContent(t *testing.T) {
	repo := newFakeRepository()
	repo.items["abc"] = Content{ID: "abc", Title: "Old"}
	router := newTestRouter(repo)

	req := httptest.NewRequest(http.MethodPut, "/content/abc", bytes.NewReader([]byte(`{"title":"New","body":"updated"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("PUT status = %d, want %d", rec.Code, http.StatusOK)
	}

	delReq := httptest.NewRequest(http.MethodDelete, "/content/abc", nil)
	delRec := httptest.NewRecorder()
	router.ServeHTTP(delRec, delReq)
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("DELETE status = %d, want %d", delRec.Code, http.StatusNoContent)
	}

	if _, err := repo.Get(context.Background(), "abc"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get after delete error = %v, want ErrNotFound", err)
	}
}
