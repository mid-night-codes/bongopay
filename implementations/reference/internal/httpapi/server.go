// Package httpapi implements the REST contract in contracts/openapi/bongopay.yaml against
// internal/payment and internal/simulator. It is the only place HTTP-specific concerns (status
// codes, request/response bodies) live — the packages it wraps know nothing about HTTP.
package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mid-night-codes/bongopay/implementations/reference/internal/payment"
	"github.com/mid-night-codes/bongopay/implementations/reference/internal/simulator"
)

// Server implements contracts/openapi/bongopay.yaml. The only provider it can actually reach
// today is the simulator — there are no adapters yet (see adapters/README.md) — so every
// PaymentRequest is routed to it, and Simulator.Initiate's own Provider.ID check is what
// rejects a request aimed at any other provider.
type Server struct {
	service   *payment.Service
	simulator *simulator.Simulator
}

// NewServer returns a Server backed by service and sim.
func NewServer(service *payment.Service, sim *simulator.Simulator) *Server {
	return &Server{service: service, simulator: sim}
}

// Handler returns the http.Handler implementing contracts/openapi/bongopay.yaml's paths.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /payments", s.handleInitiatePayment)
	mux.HandleFunc("GET /payments/{id}", s.handleGetPayment)
	return mux
}

func (s *Server) handleInitiatePayment(w http.ResponseWriter, r *http.Request) {
	var req payment.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	result, err := s.simulator.Initiate(req)
	if err != nil {
		writeErrFor(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleGetPayment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	p, err := s.service.Get(id)
	if err != nil {
		writeErrFor(w, err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// errorResponse matches the Error schema in contracts/openapi/bongopay.yaml — a provisional,
// non-canonical shape, not specs/errors/error-model.md (still unwritten).
type errorResponse struct {
	Error string `json:"error"`
}

// writeErrFor maps a package-local error from payment/simulator onto an HTTP status, per
// contracts/openapi/bongopay.yaml's documented response codes. Everything not explicitly
// recognized here is a 500 — a new package-local error type needs a case added deliberately,
// not falling through to some guessed default.
func writeErrFor(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, payment.ErrMissingIdempotencyKey),
		errors.Is(err, simulator.ErrWrongProvider):
		writeError(w, http.StatusBadRequest, err.Error())
		return
	case errors.Is(err, payment.ErrPaymentNotFound):
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	var unknownScenario *simulator.ErrUnknownScenario
	if errors.As(err, &unknownScenario) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var transitionErr *payment.TransitionError
	if errors.As(err, &transitionErr) {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	writeError(w, http.StatusInternalServerError, "internal error")
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
