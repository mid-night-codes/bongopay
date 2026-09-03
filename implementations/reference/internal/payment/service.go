package payment

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrMissingIdempotencyKey is returned by Service.Create when PaymentRequest.IdempotencyKey is
// empty. This is a provisional, package-local error, not the canonical error taxonomy —
// specs/errors/error-model.md is still TODO(ADR); do not treat this as that model's shape.
var ErrMissingIdempotencyKey = errors.New("payment: idempotencyKey is required")

// IDGenerator produces a new canonical Payment.ID. Injectable so tests get deterministic IDs.
type IDGenerator func() string

// Clock returns the current time. Injectable so tests get deterministic CreatedAt/UpdatedAt.
type Clock func() time.Time

func defaultIDGenerator() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand failing is a fatal, unrecoverable environment problem, not a normal
		// error path any caller could sensibly handle.
		panic(fmt.Sprintf("payment: failed to generate random payment ID: %v", err))
	}
	return hex.EncodeToString(b[:])
}

// Service implements payment orchestration against a Store: creating payments idempotently and
// applying canonical state transitions. It knows nothing about providers or transport — see
// ARCHITECTURE.md §5 and §6 for where those attach.
type Service struct {
	mu    sync.Mutex
	store Store
	genID IDGenerator
	now   Clock
}

// ServiceOption configures optional Service behavior, primarily for tests.
type ServiceOption func(*Service)

// WithIDGenerator overrides how new Payment IDs are generated. Default: 16 random bytes,
// hex-encoded.
func WithIDGenerator(gen IDGenerator) ServiceOption {
	return func(s *Service) { s.genID = gen }
}

// WithClock overrides how the current time is obtained. Default: time.Now.
func WithClock(clock Clock) ServiceOption {
	return func(s *Service) { s.now = clock }
}

// NewService returns a Service backed by store.
func NewService(store Store, opts ...ServiceOption) *Service {
	s := &Service{
		store: store,
		genID: defaultIDGenerator,
		now:   time.Now,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Create initiates a new payment, per specs/payments/payment-contract.md "Idempotency": given
// an idempotency key already seen, it returns the existing logical Payment rather than creating
// a new one. The check-then-create sequence is serialized per Service instance so concurrent
// calls with the same IdempotencyKey cannot race into two different Payments.
func (s *Service) Create(req PaymentRequest) (Payment, error) {
	if req.IdempotencyKey == "" {
		return Payment{}, ErrMissingIdempotencyKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, ok := s.store.FindByIdempotencyKey(req.IdempotencyKey); ok {
		return existing, nil
	}

	now := s.now()
	p := Payment{
		ID:                s.genID(),
		Status:            StatusCreated,
		Amount:            req.Amount,
		Provider:          req.Provider,
		CustomerReference: req.CustomerReference,
		IdempotencyKey:    req.IdempotencyKey,
		CreatedAt:         now,
		UpdatedAt:         now,
		Metadata:          req.Metadata,
		ProviderOptions:   req.ProviderOptions,
	}

	if err := s.store.Save(p); err != nil {
		return Payment{}, fmt.Errorf("payment: saving new payment: %w", err)
	}

	return p, nil
}
