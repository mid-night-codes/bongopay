package bongopay_test

import (
	"context"
	"errors"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	bongopay "github.com/mid-night-codes/bongopay/sdks/go"
	"github.com/mid-night-codes/bongopay/sdks/go/generated"
)

// buildReferenceServer builds implementations/reference/cmd/server into a temp binary, so
// these tests exercise the client against the real server over real HTTP — not a mock — the
// same discipline used to verify the Docker image and CLI.
func buildReferenceServer(t *testing.T) string {
	t.Helper()

	refDir, err := filepath.Abs("../../implementations/reference")
	if err != nil {
		t.Fatalf("resolving reference implementation path: %v", err)
	}
	if _, statErr := os.Stat(refDir); statErr != nil {
		t.Skipf("implementations/reference not found at %s: %v", refDir, statErr)
	}

	bin := filepath.Join(t.TempDir(), "bongopay-server")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/server")
	cmd.Dir = refDir
	if out, buildErr := cmd.CombinedOutput(); buildErr != nil {
		t.Fatalf("building reference server: %v\n%s", buildErr, out)
	}
	return bin
}

// startServer runs bin listening on addr and waits for it to actually accept connections
// before returning, killing it when the test ends.
func startServer(t *testing.T, bin, addr string) {
	t.Helper()

	cmd := exec.Command(bin, "-addr", addr, "-callback-secret", "sdk-test-secret")
	if err := cmd.Start(); err != nil {
		t.Fatalf("starting reference server: %v", err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("reference server did not start listening on %s in time", addr)
}

func TestClient_InitiateAndGetPayment(t *testing.T) {
	bin := buildReferenceServer(t)
	addr := "127.0.0.1:18099"
	startServer(t, bin, addr)

	client := bongopay.New("http://" + addr)
	ctx := context.Background()

	result, err := client.InitiatePayment(ctx, bongopay.PaymentRequest{
		Provider:          generated.Provider{Id: "SIMULATOR"},
		Amount:            generated.Money{Value: 5000, Currency: generated.Currency{Code: "TZS"}},
		CustomerReference: generated.CustomerReference{},
		IdempotencyKey:    "sdk-test-1",
	})
	if err != nil {
		t.Fatalf("InitiatePayment() error = %v", err)
	}
	if result.Payment.Status != generated.SUCCESS {
		t.Errorf("Status = %s, want %s", result.Payment.Status, generated.SUCCESS)
	}
	if result.Payment.Id == "" {
		t.Error("Payment.Id is empty")
	}

	got, err := client.GetPayment(ctx, result.Payment.Id)
	if err != nil {
		t.Fatalf("GetPayment() error = %v", err)
	}
	if got.Id != result.Payment.Id {
		t.Errorf("GetPayment() Id = %s, want %s", got.Id, result.Payment.Id)
	}
	if got.Status != generated.SUCCESS {
		t.Errorf("GetPayment() Status = %s, want %s", got.Status, generated.SUCCESS)
	}
}

func TestClient_GetPayment_NotFound(t *testing.T) {
	bin := buildReferenceServer(t)
	addr := "127.0.0.1:18100"
	startServer(t, bin, addr)

	client := bongopay.New("http://" + addr)

	_, err := client.GetPayment(context.Background(), "does-not-exist")

	var apiErr *bongopay.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("GetPayment() error = %v, want *bongopay.APIError", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
}
