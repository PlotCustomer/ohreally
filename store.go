package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// ContentStore persists the Content Data Object on behalf of the Content
// Management function.
type ContentStore interface {
	List(ctx context.Context) ([]Content, error)
	Get(ctx context.Context, id string) (Content, error)
	Create(ctx context.Context, item Content) (Content, error)
	Update(ctx context.Context, item Content) (Content, error)
	Delete(ctx context.Context, id string) error
	Close()
}

// newContentStore selects the persistence backend. When a database is
// configured (DATABASE_URL) the default local postgres service is used,
// otherwise the component falls back to an in-memory store so the preview can
// always start.
func newContentStore() (ContentStore, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Printf("%s: DATABASE_URL not set, using in-memory content store", serviceName)
		return newMemoryStore(), nil
	}

	store, err := newPostgresStore(dsn)
	if err != nil {
		log.Printf("%s: postgres unavailable (%v), falling back to in-memory content store", serviceName, err)
		return newMemoryStore(), nil
	}
	log.Printf("%s: using postgres content store", serviceName)
	return store, nil
}

// newContentID returns a random identifier for a Content record.
func newContentID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("content-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

// memoryStore keeps Content in memory and is safe for concurrent use.
type memoryStore struct {
	mu    sync.RWMutex
	items map[string]Content
}

func newMemoryStore() *memoryStore {
	return &memoryStore{items: make(map[string]Content)}
}

func (s *memoryStore) List(_ context.Context) ([]Content, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]Content, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}
	return items, nil
}

func (s *memoryStore) Get(_ context.Context, id string) (Content, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return Content{}, ErrContentNotFound
	}
	return item, nil
}

func (s *memoryStore) Create(_ context.Context, item Content) (Content, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[item.ID] = item
	return item, nil
}

func (s *memoryStore) Update(_ context.Context, item Content) (Content, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[item.ID]; !ok {
		return Content{}, ErrContentNotFound
	}
	s.items[item.ID] = item
	return item, nil
}

func (s *memoryStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrContentNotFound
	}
	delete(s.items, id)
	return nil
}

func (s *memoryStore) Close() {}
