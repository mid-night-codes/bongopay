package testcontainer_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	bongopay "github.com/mid-night-codes/bongopay/sdks/go"
	"github.com/mid-night-codes/bongopay/sdks/go/generated"
	"github.com/mid-night-codes/bongopay/sdks/go/testcontainer"
)

// repoRoot resolves the BongoPay repository root from this test file's location, or skips the
// test if it isn't found there (e.g. this module ever being tested outside a full checkout).
func repoRoot(t *testing.T) string {
	t.Helper()

	root, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatalf("resolving repo root: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(root, "deploy", "docker", "reference.Dockerfile")); statErr != nil {
		t.Skipf("repo root not found at %s: %v", root, statErr)
	}
	return root
}

// TestRunReference actually starts a Docker container — building the real image from
// deploy/docker/reference.Dockerfile — and exercises it via sdks/go's client over real HTTP.
// Skips (not fails) if Docker isn't available, since that's an environment fact, not a bug.
func TestRunReference(t *testing.T) {
	root := repoRoot(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	container, err := testcontainer.RunReference(ctx, root)
	if err != nil {
		t.Skipf("could not start reference container (is Docker running?): %v", err)
	}
	t.Cleanup(func() {
		if termErr := container.Terminate(context.Background()); termErr != nil {
			t.Logf("terminating container: %v", termErr)
		}
	})

	client := bongopay.New(container.BaseURL)

	result, err := client.InitiatePayment(ctx, bongopay.PaymentRequest{
		Provider:          generated.Provider{Id: "SIMULATOR"},
		Amount:            generated.Money{Value: 1000, Currency: generated.Currency{Code: "TZS"}},
		CustomerReference: generated.CustomerReference{},
		IdempotencyKey:    "testcontainer-test-1",
	})
	if err != nil {
		t.Fatalf("InitiatePayment() error = %v", err)
	}
	if result.Payment.Status != generated.SUCCESS {
		t.Errorf("Status = %s, want %s", result.Payment.Status, generated.SUCCESS)
	}

	got, err := client.GetPayment(ctx, result.Payment.Id)
	if err != nil {
		t.Fatalf("GetPayment() error = %v", err)
	}
	if got.Id != result.Payment.Id {
		t.Errorf("GetPayment() Id = %s, want %s", got.Id, result.Payment.Id)
	}
}
