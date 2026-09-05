package simulator

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/mid-night-codes/bongopay/implementations/reference/internal/payment"
)

// seedPending creates a payment (via svc directly, bypassing Initiate) and drives it to PENDING,
// as if it had been submitted to a provider and were now awaiting a callback.
func seedPending(t *testing.T, svc *payment.Service, idempotencyKey string) payment.Payment {
	t.Helper()
	p, err := svc.Create(payment.PaymentRequest{
		Provider:       payment.Provider{ID: ProviderID},
		Amount:         payment.Money{Value: 1000, Currency: payment.Currency{Code: "TZS"}},
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		t.Fatalf("seedPending: Create() error = %v", err)
	}
	p, err = svc.ApplyTransition(p.ID, payment.StatusPending)
	if err != nil {
		t.Fatalf("seedPending: ApplyTransition(PENDING) error = %v", err)
	}
	return p
}

func callbackBody(t *testing.T, paymentID string, status payment.PaymentStatus) []byte {
	t.Helper()
	body, err := json.Marshal(Callback{PaymentID: paymentID, Status: status})
	if err != nil {
		t.Fatalf("marshaling callback: %v", err)
	}
	return body
}

func TestHandleCallback_ValidSignature_AppliesTransition(t *testing.T) {
	sim, svc, _ := newTestSimulator()
	p := seedPending(t, svc, "idem-1")

	body := callbackBody(t, p.ID, payment.StatusSuccess)
	sig := sim.SignCallback(body)

	got, err := sim.HandleCallback(body, sig)
	if err != nil {
		t.Fatalf("HandleCallback() error = %v, want nil", err)
	}
	if got.Status != payment.StatusSuccess {
		t.Errorf("Status = %s, want %s", got.Status, payment.StatusSuccess)
	}
}

func TestHandleCallback_InvalidSignature_Rejected(t *testing.T) {
	sim, svc, store := newTestSimulator()
	p := seedPending(t, svc, "idem-1")

	body := callbackBody(t, p.ID, payment.StatusSuccess)

	_, err := sim.HandleCallback(body, "not-a-real-signature")
	if !errors.Is(err, ErrInvalidCallbackSignature) {
		t.Fatalf("HandleCallback() error = %v, want ErrInvalidCallbackSignature", err)
	}

	stored, _ := store.FindByID(p.ID)
	if stored.Status != payment.StatusPending {
		t.Errorf("stored Status = %s, want unchanged %s (invalid signature must not mutate)", stored.Status, payment.StatusPending)
	}
}

func TestHandleCallback_TamperedBody_Rejected(t *testing.T) {
	sim, svc, store := newTestSimulator()
	p := seedPending(t, svc, "idem-1")

	original := callbackBody(t, p.ID, payment.StatusFailed)
	sig := sim.SignCallback(original)

	// An attacker (or a transport bug) changes the outcome after signing.
	tampered := callbackBody(t, p.ID, payment.StatusSuccess)

	_, err := sim.HandleCallback(tampered, sig)
	if !errors.Is(err, ErrInvalidCallbackSignature) {
		t.Fatalf("HandleCallback() error = %v, want ErrInvalidCallbackSignature", err)
	}

	stored, _ := store.FindByID(p.ID)
	if stored.Status != payment.StatusPending {
		t.Errorf("stored Status = %s, want unchanged %s (tampered body must not mutate)", stored.Status, payment.StatusPending)
	}
}

// TestHandleCallback_Duplicate_IsNoOp exercises specs/scenarios/scenario-format.md's
// DUPLICATE_CALLBACK outcome: a second, identical valid callback must be a no-op.
func TestHandleCallback_Duplicate_IsNoOp(t *testing.T) {
	sim, svc, _ := newTestSimulator()
	p := seedPending(t, svc, "idem-1")

	body := callbackBody(t, p.ID, payment.StatusSuccess)
	sig := sim.SignCallback(body)

	first, err := sim.HandleCallback(body, sig)
	if err != nil {
		t.Fatalf("first HandleCallback() error = %v", err)
	}

	second, err := sim.HandleCallback(body, sig)
	if err != nil {
		t.Fatalf("second (duplicate) HandleCallback() error = %v, want nil (no-op)", err)
	}

	if !second.UpdatedAt.Equal(first.UpdatedAt) {
		t.Errorf("duplicate callback re-drove the state machine: UpdatedAt changed from %v to %v", first.UpdatedAt, second.UpdatedAt)
	}
}

// TestHandleCallback_OutOfOrder_DoesNotRegressState exercises
// specs/scenarios/scenario-format.md's OUT_OF_ORDER outcome: a stale PENDING claim delivered
// after SUCCESS has already been applied must be rejected as a conflicting update, not silently
// applied backward.
func TestHandleCallback_OutOfOrder_DoesNotRegressState(t *testing.T) {
	sim, svc, store := newTestSimulator()
	p := seedPending(t, svc, "idem-1")

	successBody := callbackBody(t, p.ID, payment.StatusSuccess)
	if _, err := sim.HandleCallback(successBody, sim.SignCallback(successBody)); err != nil {
		t.Fatalf("HandleCallback(SUCCESS) error = %v", err)
	}

	staleBody := callbackBody(t, p.ID, payment.StatusPending)
	_, err := sim.HandleCallback(staleBody, sim.SignCallback(staleBody))

	var transitionErr *payment.TransitionError
	if !errors.As(err, &transitionErr) {
		t.Fatalf("HandleCallback(stale PENDING) error = %v, want *payment.TransitionError", err)
	}

	stored, _ := store.FindByID(p.ID)
	if stored.Status != payment.StatusSuccess {
		t.Errorf("stored Status = %s, want unchanged %s (must not regress)", stored.Status, payment.StatusSuccess)
	}
}
