package simulator

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/mid-night-codes/bongopay/implementations/reference/internal/payment"
)

func newTestSimulator() (*Simulator, *payment.InMemoryStore) {
	store := payment.NewInMemoryStore()
	svc := payment.NewService(store)
	return New(svc), store
}

func requestWithScenario(idempotencyKey, scenario string) payment.PaymentRequest {
	req := payment.PaymentRequest{
		Provider:       payment.Provider{ID: ProviderID},
		Amount:         payment.Money{Value: 1000, Currency: payment.Currency{Code: "TZS"}},
		IdempotencyKey: idempotencyKey,
	}
	if scenario != "" {
		raw, _ := json.Marshal(map[string]string{"scenario": scenario})
		req.ProviderOptions = payment.ProviderOptions{"simulator": raw}
	}
	return req
}

func TestInitiate_DefaultScenario_Success(t *testing.T) {
	sim, _ := newTestSimulator()

	result, err := sim.Initiate(requestWithScenario("idem-1", ""))
	if err != nil {
		t.Fatalf("Initiate() error = %v, want nil", err)
	}
	if result.Payment.Status != payment.StatusSuccess {
		t.Errorf("Status = %s, want %s", result.Payment.Status, payment.StatusSuccess)
	}
}

func TestInitiate_ExplicitSuccessScenario(t *testing.T) {
	sim, _ := newTestSimulator()

	result, err := sim.Initiate(requestWithScenario("idem-1", "success"))
	if err != nil {
		t.Fatalf("Initiate() error = %v, want nil", err)
	}
	if result.Payment.Status != payment.StatusSuccess {
		t.Errorf("Status = %s, want %s", result.Payment.Status, payment.StatusSuccess)
	}
}

func TestInitiate_FailureScenario(t *testing.T) {
	sim, _ := newTestSimulator()

	result, err := sim.Initiate(requestWithScenario("idem-1", "failure"))
	if err != nil {
		t.Fatalf("Initiate() error = %v, want nil", err)
	}
	if result.Payment.Status != payment.StatusFailed {
		t.Errorf("Status = %s, want %s", result.Payment.Status, payment.StatusFailed)
	}
}

func TestInitiate_UnknownScenario_NoPaymentCreated(t *testing.T) {
	sim, store := newTestSimulator()

	_, err := sim.Initiate(requestWithScenario("idem-1", "does-not-exist"))

	var unknownErr *ErrUnknownScenario
	if !errors.As(err, &unknownErr) {
		t.Fatalf("Initiate() error = %v, want *ErrUnknownScenario", err)
	}
	if unknownErr.Name != "does-not-exist" {
		t.Errorf("ErrUnknownScenario.Name = %q, want %q", unknownErr.Name, "does-not-exist")
	}

	if _, ok := store.FindByIdempotencyKey("idem-1"); ok {
		t.Error("a Payment was created for a request with an unknown scenario, want none")
	}
}

func TestInitiate_WrongProvider_NoPaymentCreated(t *testing.T) {
	sim, store := newTestSimulator()

	req := requestWithScenario("idem-1", "success")
	req.Provider = payment.Provider{ID: "MPESA"}

	_, err := sim.Initiate(req)
	if !errors.Is(err, ErrWrongProvider) {
		t.Fatalf("Initiate() error = %v, want ErrWrongProvider", err)
	}

	if _, ok := store.FindByIdempotencyKey("idem-1"); ok {
		t.Error("a Payment was created for a request targeting the wrong provider, want none")
	}
}

func TestInitiate_SequentialIdempotentReplay(t *testing.T) {
	sim, _ := newTestSimulator()

	first, err := sim.Initiate(requestWithScenario("idem-1", "success"))
	if err != nil {
		t.Fatalf("first Initiate() error = %v", err)
	}

	second, err := sim.Initiate(requestWithScenario("idem-1", "success"))
	if err != nil {
		t.Fatalf("second Initiate() error = %v", err)
	}

	if first.Payment.ID != second.Payment.ID {
		t.Errorf("second call returned a different Payment: got ID %q, want %q", second.Payment.ID, first.Payment.ID)
	}
	if second.Payment.Status != payment.StatusSuccess {
		t.Errorf("replayed Status = %s, want %s", second.Payment.Status, payment.StatusSuccess)
	}
	if !second.Payment.UpdatedAt.Equal(first.Payment.UpdatedAt) {
		t.Errorf("replay re-drove the state machine: UpdatedAt changed from %v to %v", first.Payment.UpdatedAt, second.Payment.UpdatedAt)
	}
}
