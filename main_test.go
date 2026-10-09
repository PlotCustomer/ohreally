package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestRouter() http.Handler {
	return newRouter(newMemoryStore())
}

func TestServiceInfo(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	newTestRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("health status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestContentLifecycle(t *testing.T) {
	router := newTestRouter()

	// Create.
	body, _ := json.Marshal(contentInput{Title: "Welcome", Body: "Hello", Author: "Setup"})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/content", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d (%s)", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var created Content
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created content: %v", err)
	}
	if created.ID == "" {
		t.Fatal("created content has no id")
	}

	// Read.
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/content/"+created.ID, nil)
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", rec.Code, http.StatusOK)
	}

	// List.
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/content", nil)
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", rec.Code, http.StatusOK)
	}

	// Delete.
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/api/content/"+created.ID, nil)
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", rec.Code, http.StatusNoContent)
	}

	// Read after delete.
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/content/"+created.ID, nil)
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("get deleted status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestCreateContentRequiresTitle(t *testing.T) {
	body, _ := json.Marshal(contentInput{Body: "no title"})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/content", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
