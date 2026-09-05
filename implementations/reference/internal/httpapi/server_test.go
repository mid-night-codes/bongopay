package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mid-night-codes/bongopay/implementations/reference/internal/payment"
	"github.com/mid-night-codes/bongopay/implementations/reference/internal/simulator"
)

func newTestServer() *Server {
	store := payment.NewInMemoryStore()
	svc := payment.NewService(store)
	sim := simulator.New(svc)
	return NewServer(svc, sim)
}

func doRequest(t *testing.T, h http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var reader *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshaling request body: %v", err)
		}
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, reader)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decodeBody[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatalf("decoding response body %q: %v", rec.Body.String(), err)
	}
	return v
}

func validRequest(idempotencyKey string) payment.PaymentRequest {
	return payment.PaymentRequest{
		Provider:       payment.Provider{ID: simulator.ProviderID},
		Amount:         payment.Money{Value: 1000, Currency: payment.Currency{Code: "TZS"}},
		IdempotencyKey: idempotencyKey,
	}
}

func TestHandleInitiatePayment_Success(t *testing.T) {
	s := newTestServer()

	rec := doRequest(t, s.Handler(), "POST", "/payments", validRequest("idem-1"))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	result := decodeBody[payment.PaymentResult](t, rec)
	if result.Payment.Status != payment.StatusSuccess {
		t.Errorf("Status = %s, want %s", result.Payment.Status, payment.StatusSuccess)
	}
	if result.Payment.ID == "" {
		t.Error("Payment.ID is empty")
	}
}

func TestHandleInitiatePayment_MalformedBody(t *testing.T) {
	s := newTestServer()

	req := httptest.NewRequest("POST", "/payments", bytes.NewReader([]byte("{not json")))
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	errResp := decodeBody[errorResponse](t, rec)
	if errResp.Error == "" {
		t.Error("expected a non-empty error message")
	}
}

func TestHandleInitiatePayment_MissingIdempotencyKey(t *testing.T) {
	s := newTestServer()

	rec := doRequest(t, s.Handler(), "POST", "/payments", validRequest(""))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestHandleInitiatePayment_WrongProvider(t *testing.T) {
	s := newTestServer()

	req := validRequest("idem-1")
	req.Provider = payment.Provider{ID: "MPESA"}

	rec := doRequest(t, s.Handler(), "POST", "/payments", req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestHandleGetPayment_NotFound(t *testing.T) {
	s := newTestServer()

	rec := doRequest(t, s.Handler(), "GET", "/payments/does-not-exist", nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
}

func TestRoundTrip_InitiateThenGet(t *testing.T) {
	s := newTestServer()
	h := s.Handler()

	initiateRec := doRequest(t, h, "POST", "/payments", validRequest("idem-1"))
	if initiateRec.Code != http.StatusOK {
		t.Fatalf("initiate status = %d, want %d; body = %s", initiateRec.Code, http.StatusOK, initiateRec.Body.String())
	}
	created := decodeBody[payment.PaymentResult](t, initiateRec).Payment

	getRec := doRequest(t, h, "GET", "/payments/"+created.ID, nil)
	if getRec.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d; body = %s", getRec.Code, http.StatusOK, getRec.Body.String())
	}
	got := decodeBody[payment.Payment](t, getRec)

	if got.ID != created.ID {
		t.Errorf("GET returned ID %q, want %q", got.ID, created.ID)
	}
	if got.Status != created.Status {
		t.Errorf("GET returned Status %s, want %s", got.Status, created.Status)
	}
}
