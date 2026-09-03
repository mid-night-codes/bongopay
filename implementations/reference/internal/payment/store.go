package payment

import (
	"fmt"
	"sync"
)

// Store persists Payment records and looks them up by canonical ID or by IdempotencyKey.
// It has no notion of provider, transport, or events — just the canonical record.
type Store interface {
	Save(p Payment) error
	FindByID(id string) (Payment, bool)
	FindByIdempotencyKey(key string) (Payment, bool)
}

// InMemoryStore is a Store backed by maps, safe for concurrent use. It exists for the
// simulator and for tests — a durable Store is a Phase 1+ implementation detail, not a
// canonical concept (see ARCHITECTURE.md §14, "Persistence architecture" is a TODO(ADR)).
type InMemoryStore struct {
	mu              sync.Mutex
	byID            map[string]Payment
	idByIdempotency map[string]string
}

// NewInMemoryStore returns an empty, ready-to-use InMemoryStore.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		byID:            make(map[string]Payment),
		idByIdempotency: make(map[string]string),
	}
}

// Save inserts or replaces the Payment identified by p.ID and records its IdempotencyKey
// mapping. Save does not itself enforce idempotency-key uniqueness across different payment
// IDs — that policy lives in Service.Create, which checks before ever calling Save.
func (s *InMemoryStore) Save(p Payment) error {
	if p.ID == "" {
		return fmt.Errorf("payment: cannot save a Payment with an empty ID")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.byID[p.ID] = p
	if p.IdempotencyKey != "" {
		s.idByIdempotency[p.IdempotencyKey] = p.ID
	}
	return nil
}

// FindByID returns the Payment with the given canonical ID, if any.
func (s *InMemoryStore) FindByID(id string) (Payment, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.byID[id]
	return p, ok
}

// FindByIdempotencyKey returns the Payment created for the given IdempotencyKey, if any.
func (s *InMemoryStore) FindByIdempotencyKey(key string) (Payment, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, ok := s.idByIdempotency[key]
	if !ok {
		return Payment{}, false
	}
	p, ok := s.byID[id]
	return p, ok
}
