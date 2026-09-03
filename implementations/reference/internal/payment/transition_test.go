package payment

import (
	"errors"
	"testing"
	"time"
)

// seedTime is deliberately distinct from newTestService's fixed clock (2026-01-01) so tests can
// tell "the stored UpdatedAt is unchanged" apart from "it happens to match by coincidence."
var seedTime = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

func seedPayment(t *testing.T, store *InMemoryStore, id string, status PaymentStatus) Payment {
	t.Helper()
	p := Payment{
		ID:                id,
		Status:            status,
		Amount:            Money{Value: 1000, Currency: Currency{Code: "TZS"}},
		Provider:          Provider{ID: "SIMULATOR"},
		CustomerReference: CustomerReference{},
		IdempotencyKey:    "seed-" + id,
		CreatedAt:         seedTime,
		UpdatedAt:         seedTime,
	}
	if err := store.Save(p); err != nil {
		t.Fatalf("seedPayment: Save() error = %v", err)
	}
	return p
}

func TestApplyTransition_NotFound(t *testing.T) {
	svc, _ := newTestService(t)

	_, err := svc.ApplyTransition("does-not-exist", StatusPending)
	if !errors.Is(err, ErrPaymentNotFound) {
		t.Fatalf("ApplyTransition() error = %v, want ErrPaymentNotFound", err)
	}
}

func TestApplyTransition_DuplicateSameStatus_IsNoOp(t *testing.T) {
	svc, store := newTestService(t)
	seedPayment(t, store, "p1", StatusPending)

	got, err := svc.ApplyTransition("p1", StatusPending)
	if err != nil {
		t.Fatalf("ApplyTransition() error = %v, want nil", err)
	}
	if got.Status != StatusPending {
		t.Errorf("Status = %s, want %s", got.Status, StatusPending)
	}
	if !got.UpdatedAt.Equal(seedTime) {
		t.Errorf("UpdatedAt = %v, want unchanged seedTime %v (no-op must not mutate)", got.UpdatedAt, seedTime)
	}

	stored, _ := store.FindByID("p1")
	if !stored.UpdatedAt.Equal(seedTime) {
		t.Errorf("stored UpdatedAt = %v, want unchanged seedTime %v", stored.UpdatedAt, seedTime)
	}
}

func TestApplyTransition_ValidTransition_Applies(t *testing.T) {
	svc, store := newTestService(t)
	seedPayment(t, store, "p1", StatusCreated)

	got, err := svc.ApplyTransition("p1", StatusPending)
	if err != nil {
		t.Fatalf("ApplyTransition() error = %v, want nil", err)
	}
	if got.Status != StatusPending {
		t.Errorf("Status = %s, want %s", got.Status, StatusPending)
	}
	if got.UpdatedAt.Equal(seedTime) {
		t.Error("UpdatedAt was not refreshed on a real transition")
	}

	stored, ok := store.FindByID("p1")
	if !ok {
		t.Fatal("payment missing from store after transition")
	}
	if stored.Status != StatusPending {
		t.Errorf("stored Status = %s, want %s", stored.Status, StatusPending)
	}
}

func TestApplyTransition_ConflictingUpdate_ReturnsTransitionError(t *testing.T) {
	svc, store := newTestService(t)
	seedPayment(t, store, "p1", StatusSuccess)

	// SUCCESS's only valid next states are REVERSAL_PENDING and REFUND_PENDING — FAILED is a
	// genuinely conflicting update (e.g. a callback claiming failure for an already-successful
	// payment), not a duplicate of the current status.
	_, err := svc.ApplyTransition("p1", StatusFailed)

	var transitionErr *TransitionError
	if !errors.As(err, &transitionErr) {
		t.Fatalf("ApplyTransition() error = %v, want *TransitionError", err)
	}
	if transitionErr.From != StatusSuccess || transitionErr.To != StatusFailed {
		t.Errorf("TransitionError = {From: %s, To: %s}, want {From: %s, To: %s}",
			transitionErr.From, transitionErr.To, StatusSuccess, StatusFailed)
	}

	stored, _ := store.FindByID("p1")
	if stored.Status != StatusSuccess {
		t.Errorf("stored Status = %s, want unchanged %s (conflicting update must not mutate)", stored.Status, StatusSuccess)
	}
	if !stored.UpdatedAt.Equal(seedTime) {
		t.Errorf("stored UpdatedAt = %v, want unchanged seedTime %v", stored.UpdatedAt, seedTime)
	}
}

func TestApplyTransition_InvalidTransition_SkippingStates(t *testing.T) {
	svc, store := newTestService(t)
	seedPayment(t, store, "p1", StatusCreated)

	// CREATED can only reach PENDING or CANCELLED directly — REFUNDED is unreachable without
	// first passing through PENDING, SUCCESS, and REFUND_PENDING.
	_, err := svc.ApplyTransition("p1", StatusRefunded)

	var transitionErr *TransitionError
	if !errors.As(err, &transitionErr) {
		t.Fatalf("ApplyTransition() error = %v, want *TransitionError", err)
	}
}
