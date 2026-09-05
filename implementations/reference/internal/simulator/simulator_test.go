package simulator

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"github.com/mid-night-codes/bongopay/implementations/reference/internal/payment"
)

func newTestSimulator() (*Simulator, *payment.Service, *payment.InMemoryStore) {
	store := payment.NewInMemoryStore()
	svc := payment.NewService(store)
	return New(svc), svc, store
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
	sim, _, _ := newTestSimulator()

	result, err := sim.Initiate(requestWithScenario("idem-1", ""))
	if err != nil {
		t.Fatalf("Initiate() error = %v, want nil", err)
	}
	if result.Payment.Status != payment.StatusSuccess {
		t.Errorf("Status = %s, want %s", result.Payment.Status, payment.StatusSuccess)
	}
}

func TestInitiate_ExplicitSuccessScenario(t *testing.T) {
	sim, _, _ := newTestSimulator()

	result, err := sim.Initiate(requestWithScenario("idem-1", "success"))
	if err != nil {
		t.Fatalf("Initiate() error = %v, want nil", err)
	}
	if result.Payment.Status != payment.StatusSuccess {
		t.Errorf("Status = %s, want %s", result.Payment.Status, payment.StatusSuccess)
	}
}

func TestInitiate_FailureScenario(t *testing.T) {
	sim, _, _ := newTestSimulator()

	result, err := sim.Initiate(requestWithScenario("idem-1", "failure"))
	if err != nil {
		t.Fatalf("Initiate() error = %v, want nil", err)
	}
	if result.Payment.Status != payment.StatusFailed {
		t.Errorf("Status = %s, want %s", result.Payment.Status, payment.StatusFailed)
	}
}

func TestInitiate_UnknownScenario_NoPaymentCreated(t *testing.T) {
	sim, _, store := newTestSimulator()

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
	sim, _, store := newTestSimulator()

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
	sim, _, _ := newTestSimulator()

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

// TestInitiate_ConcurrentSameKey_NoTransitionErrors guards against the race documented on an
// earlier version of Initiate: concurrent calls for the same brand-new IdempotencyKey used to
// be able to interleave between separate Service.Create/ApplyTransition calls and produce a
// *payment.TransitionError instead of the final outcome. Initiate now delegates the whole
// sequence to Service.CreateAndAdvance under one lock acquisition, closing that race.
func TestInitiate_ConcurrentSameKey_NoTransitionErrors(t *testing.T) {
	sim, _, _ := newTestSimulator()

	const attempts = 50
	var wg sync.WaitGroup
	results := make([]payment.PaymentResult, attempts)
	errs := make([]error, attempts)

	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i], errs[i] = sim.Initiate(requestWithScenario("idem-concurrent", "success"))
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("attempt %d: Initiate() error = %v, want nil", i, err)
		}
	}

	first := results[0].Payment
	for i, r := range results {
		if r.Payment.ID != first.ID {
			t.Errorf("attempt %d returned ID %q, want %q", i, r.Payment.ID, first.ID)
		}
		if r.Payment.Status != payment.StatusSuccess {
			t.Errorf("attempt %d Status = %s, want %s", i, r.Payment.Status, payment.StatusSuccess)
		}
	}
}
