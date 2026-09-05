package payment

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// These are provisional, package-local errors, not the canonical error taxonomy —
// specs/errors/error-model.md is still TODO(ADR); do not treat this as that model's shape.
var (
	// ErrMissingIdempotencyKey is returned by Service.Create when
	// PaymentRequest.IdempotencyKey is empty.
	ErrMissingIdempotencyKey = errors.New("payment: idempotencyKey is required")

	// ErrPaymentNotFound is returned by Service.ApplyTransition when no Payment exists for
	// the given ID.
	ErrPaymentNotFound = errors.New("payment: payment not found")
)

// TransitionError reports a genuinely conflicting status update — a target status that is
// neither the payment's current status (that case is a no-op, not an error; see
// Service.ApplyTransition) nor a valid next state per
// specs/state-machines/payment-lifecycle.md.
type TransitionError struct {
	From PaymentStatus
	To   PaymentStatus
}

func (e *TransitionError) Error() string {
	return fmt.Sprintf("payment: invalid transition from %s to %s", e.From, e.To)
}

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
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.createLocked(req)
}

// createLocked is Create's body, callable by other Service methods that already hold s.mu —
// see CreateAndAdvance.
func (s *Service) createLocked(req PaymentRequest) (Payment, error) {
	if req.IdempotencyKey == "" {
		return Payment{}, ErrMissingIdempotencyKey
	}

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

// ApplyTransition moves the payment identified by id to status to, per
// specs/state-machines/payment-lifecycle.md "Idempotency and Retries":
//
//   - to equals the payment's current status: a duplicate delivery of an already-applied
//     event. No-op — the store is not touched and the current Payment is returned unchanged.
//   - to is a valid next state per CanTransition: applied normally, UpdatedAt is refreshed.
//   - anything else: a genuinely conflicting update. Returns *TransitionError and does not
//     mutate the stored Payment.
//
// Like Create, this is mutex-serialized per Service so a concurrent duplicate and a concurrent
// valid transition can't race against each other.
func (s *Service) ApplyTransition(id string, to PaymentStatus) (Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.applyTransitionLocked(id, to)
}

// applyTransitionLocked is ApplyTransition's body, callable by other Service methods that
// already hold s.mu — see CreateAndAdvance.
func (s *Service) applyTransitionLocked(id string, to PaymentStatus) (Payment, error) {
	p, ok := s.store.FindByID(id)
	if !ok {
		return Payment{}, ErrPaymentNotFound
	}

	if p.Status == to {
		return p, nil
	}

	if !CanTransition(p.Status, to) {
		return Payment{}, &TransitionError{From: p.Status, To: to}
	}

	p.Status = to
	p.UpdatedAt = s.now()

	if err := s.store.Save(p); err != nil {
		return Payment{}, fmt.Errorf("payment: saving transition: %w", err)
	}

	return p, nil
}

// CreateAndAdvance creates (or looks up, per Create's idempotency rule) a payment and, only if
// it was freshly created by this call, applies each status in path in order — all under a
// single lock acquisition.
//
// This exists because Create followed by separate ApplyTransition calls has a race: another
// goroutine's Create-replay can observe the payment between those calls, at a status that then
// makes a *later* step invalid (e.g. it already reached SUCCESS by the time this goroutine
// tries to apply PENDING). Doing the whole sequence under one lock makes it atomic with respect
// to every other Service call — whichever caller's sequence starts first for a given
// IdempotencyKey runs to completion before any other caller for that same key can observe or
// advance it, so a concurrent replay always sees either the pre-creation state (and becomes the
// one to drive it) or the fully-applied final state (and simply returns it) — never something
// in between.
//
// If the payment already existed (idempotent replay) or path is empty, this behaves exactly
// like Create/is a lookup — path is not applied.
func (s *Service) CreateAndAdvance(req PaymentRequest, path []PaymentStatus) (Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, err := s.createLocked(req)
	if err != nil {
		return Payment{}, err
	}

	if p.Status != StatusCreated {
		// Idempotent replay of a payment already progressed by another call — return it
		// as-is rather than re-driving transitions against it.
		return p, nil
	}

	for _, to := range path {
		p, err = s.applyTransitionLocked(p.ID, to)
		if err != nil {
			return Payment{}, err
		}
	}

	return p, nil
}

// Get returns the Payment with the given canonical ID, or ErrPaymentNotFound if none exists.
func (s *Service) Get(id string) (Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.store.FindByID(id)
	if !ok {
		return Payment{}, ErrPaymentNotFound
	}
	return p, nil
}
