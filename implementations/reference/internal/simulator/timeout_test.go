package simulator

import (
	"testing"
	"time"

	"github.com/mid-night-codes/bongopay/implementations/reference/internal/payment"
)

func TestInitiate_TimeoutScenario_ExpiresAfterDelay(t *testing.T) {
	store := payment.NewInMemoryStore()
	svc := payment.NewService(store)

	var slept time.Duration
	sim := New(svc, WithSleeper(func(d time.Duration) { slept = d }))

	result, err := sim.Initiate(requestWithScenario("idem-1", "timeout"))
	if err != nil {
		t.Fatalf("Initiate() error = %v, want nil", err)
	}
	if result.Payment.Status != payment.StatusExpired {
		t.Errorf("Status = %s, want %s", result.Payment.Status, payment.StatusExpired)
	}
	if slept != defaultTimeoutDelay {
		t.Errorf("slept %v, want %v (the registry's configured delay)", slept, defaultTimeoutDelay)
	}
}

// TestInitiate_TimeoutScenario_RealCallbackDuringDelayWins proves that a real callback
// resolving the payment *during* the simulated timeout delay wins over the timeout itself,
// rather than Initiate erroring or clobbering it back to EXPIRED.
func TestInitiate_TimeoutScenario_RealCallbackDuringDelayWins(t *testing.T) {
	store := payment.NewInMemoryStore()
	svc := payment.NewService(store, payment.WithIDGenerator(func() string { return "fixed-id" }))

	sim := New(svc, WithSleeper(func(time.Duration) {
		// Simulate a real callback arriving while Initiate is "waiting out" the timeout.
		if _, err := svc.ApplyTransition("fixed-id", payment.StatusSuccess); err != nil {
			t.Fatalf("simulating a concurrent callback: %v", err)
		}
	}))

	result, err := sim.Initiate(requestWithScenario("idem-1", "timeout"))
	if err != nil {
		t.Fatalf("Initiate() error = %v, want nil", err)
	}
	if result.Payment.Status != payment.StatusSuccess {
		t.Errorf("Status = %s, want %s (the real callback should win over the timeout)", result.Payment.Status, payment.StatusSuccess)
	}
}

func TestInitiate_TimeoutScenario_IdempotentReplayAfterResolution(t *testing.T) {
	store := payment.NewInMemoryStore()
	svc := payment.NewService(store)
	sim := New(svc, WithSleeper(func(time.Duration) {}))

	first, err := sim.Initiate(requestWithScenario("idem-1", "timeout"))
	if err != nil {
		t.Fatalf("first Initiate() error = %v", err)
	}

	second, err := sim.Initiate(requestWithScenario("idem-1", "timeout"))
	if err != nil {
		t.Fatalf("second Initiate() error = %v", err)
	}

	if first.Payment.ID != second.Payment.ID {
		t.Errorf("second call returned a different Payment: got ID %q, want %q", second.Payment.ID, first.Payment.ID)
	}
	if second.Payment.Status != payment.StatusExpired {
		t.Errorf("replayed Status = %s, want %s", second.Payment.Status, payment.StatusExpired)
	}
	if !second.Payment.UpdatedAt.Equal(first.Payment.UpdatedAt) {
		t.Errorf("replay re-drove the state machine: UpdatedAt changed from %v to %v", first.Payment.UpdatedAt, second.Payment.UpdatedAt)
	}
}

func TestDefaultRegistry_TimeoutScenario(t *testing.T) {
	scenario, err := DefaultRegistry().Resolve("timeout")
	if err != nil {
		t.Fatalf("Resolve(%q) error = %v", "timeout", err)
	}
	if scenario.Outcome != OutcomeTimeout {
		t.Errorf("Outcome = %s, want %s", scenario.Outcome, OutcomeTimeout)
	}
	if scenario.Delay <= 0 {
		t.Errorf("Delay = %v, want > 0", scenario.Delay)
	}
}
