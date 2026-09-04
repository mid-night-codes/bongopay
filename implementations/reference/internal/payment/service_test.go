package payment

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func newTestService(t *testing.T) (*Service, *InMemoryStore) {
	t.Helper()
	store := NewInMemoryStore()

	var counter int
	svc := NewService(store,
		WithIDGenerator(func() string {
			counter++
			return "test-id-" + string(rune('0'+counter))
		}),
		WithClock(func() time.Time {
			return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		}),
	)
	return svc, store
}

func testRequest(idempotencyKey string) PaymentRequest {
	return PaymentRequest{
		Provider:          Provider{ID: "SIMULATOR"},
		Amount:            Money{Value: 5000, Currency: Currency{Code: "TZS"}},
		CustomerReference: CustomerReference{},
		IdempotencyKey:    idempotencyKey,
	}
}

func TestCreate_NewPayment(t *testing.T) {
	svc, _ := newTestService(t)

	got, err := svc.Create(testRequest("idem-1"))
	if err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}
	if got.Status != StatusCreated {
		t.Errorf("Status = %s, want %s", got.Status, StatusCreated)
	}
	if got.IdempotencyKey != "idem-1" {
		t.Errorf("IdempotencyKey = %q, want %q", got.IdempotencyKey, "idem-1")
	}
	if got.ID == "" {
		t.Error("ID is empty, want a generated ID")
	}
	if got.CreatedAt.IsZero() || got.UpdatedAt.IsZero() {
		t.Error("CreatedAt/UpdatedAt not set")
	}
}

func TestCreate_MissingIdempotencyKey(t *testing.T) {
	svc, _ := newTestService(t)

	_, err := svc.Create(testRequest(""))
	if !errors.Is(err, ErrMissingIdempotencyKey) {
		t.Fatalf("Create() error = %v, want ErrMissingIdempotencyKey", err)
	}
}

func TestCreate_IdempotentReplay(t *testing.T) {
	svc, store := newTestService(t)

	first, err := svc.Create(testRequest("idem-1"))
	if err != nil {
		t.Fatalf("first Create() error = %v", err)
	}

	second, err := svc.Create(testRequest("idem-1"))
	if err != nil {
		t.Fatalf("second Create() error = %v", err)
	}

	if first.ID != second.ID {
		t.Errorf("second Create() returned a different Payment: got ID %q, want %q", second.ID, first.ID)
	}

	if _, ok := store.FindByID(first.ID); !ok {
		t.Fatal("payment not found in store after replay")
	}
}

func TestCreate_ConcurrentSameKey_CreatesExactlyOnePayment(t *testing.T) {
	svc, store := newTestService(t)

	const attempts = 50
	var wg sync.WaitGroup
	ids := make([]string, attempts)

	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p, err := svc.Create(testRequest("idem-race"))
			if err != nil {
				t.Errorf("Create() error = %v", err)
				return
			}
			ids[i] = p.ID
		}(i)
	}
	wg.Wait()

	first := ids[0]
	for i, id := range ids {
		if id != first {
			t.Errorf("attempt %d returned ID %q, want %q (all concurrent calls with the same idempotency key must return the same Payment)", i, id, first)
		}
	}

	p, ok := store.FindByIdempotencyKey("idem-race")
	if !ok {
		t.Fatal("no payment found for idem-race")
	}
	if p.ID != first {
		t.Errorf("store's payment ID = %q, want %q", p.ID, first)
	}
}

func TestCreateAndAdvance_ConcurrentSameKey_NoTransitionErrors(t *testing.T) {
	svc, _ := newTestService(t)
	path := []PaymentStatus{StatusPending, StatusSuccess}

	const attempts = 50
	var wg sync.WaitGroup
	results := make([]Payment, attempts)
	errs := make([]error, attempts)

	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i], errs[i] = svc.CreateAndAdvance(testRequest("idem-advance-race"), path)
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("attempt %d: CreateAndAdvance() error = %v, want nil", i, err)
		}
	}

	first := results[0]
	for i, p := range results {
		if p.ID != first.ID {
			t.Errorf("attempt %d returned ID %q, want %q", i, p.ID, first.ID)
		}
		if p.Status != StatusSuccess {
			t.Errorf("attempt %d Status = %s, want %s", i, p.Status, StatusSuccess)
		}
	}
}
