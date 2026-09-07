package simulator

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mid-night-codes/bongopay/implementations/reference/internal/payment"
)

// ProviderID is the Provider.id a PaymentRequest must use to reach this simulator, per
// specs/payments/payment-contract.md ("SIMULATOR is always a valid Provider").
const ProviderID = "SIMULATOR"

// ErrWrongProvider is returned when Initiate is called with a PaymentRequest not targeting
// ProviderID.
var ErrWrongProvider = fmt.Errorf("simulator: PaymentRequest.Provider.ID must be %q", ProviderID)

// ErrInvalidCallbackSignature is returned by HandleCallback when the signature doesn't verify.
// Per ARCHITECTURE.md §12, this MUST be checked, and MUST reject, before any canonical state
// transition is attempted — see HandleCallback.
var ErrInvalidCallbackSignature = errors.New("simulator: invalid callback signature")

// options is the shape of providerOptions.simulator, per
// specs/scenarios/scenario-format.md "Selecting a Scenario".
type options struct {
	Scenario string `json:"scenario"`
}

// Simulator implements the SIMULATOR provider against a payment.Service: it drives Create and
// ApplyTransition the same way a real adapter would, just without a network call. It is a peer
// of real provider adapters behind the same conceptual Provider Interface (ARCHITECTURE.md §6),
// not a special case.
type Simulator struct {
	service  *payment.Service
	registry Registry
	verifier *CallbackVerifier
	sleep    Sleeper
}

// Sleeper pauses for at least d, simulating a Scenario's Delay. Injectable so tests exercising
// a delayed outcome (e.g. "timeout") don't have to wait through a real delay.
type Sleeper func(d time.Duration)

// Option configures optional Simulator behavior — see WithSecret and WithSleeper.
type Option func(*Simulator)

// WithSecret sets the callback-signing secret (see CallbackVerifier) instead of a random one.
// Useful for a long-running process (e.g. cmd/server) that wants to print its secret once at
// startup so it can be reproduced externally (openssl, a test script) to sign a callback body
// for POST /simulator/callbacks.
func WithSecret(secret []byte) Option {
	return func(s *Simulator) { s.verifier = NewCallbackVerifier(secret) }
}

// WithSleeper overrides how Initiate waits out a delayed scenario's Delay. Default: time.Sleep.
func WithSleeper(sleep Sleeper) Option {
	return func(s *Simulator) { s.sleep = sleep }
}

// New returns a Simulator backed by service, using DefaultRegistry for scenario resolution,
// time.Sleep for any delayed outcome, and, unless overridden with WithSecret, a fresh, random
// callback-signing secret.
func New(service *payment.Service, opts ...Option) *Simulator {
	s := &Simulator{
		service:  service,
		registry: DefaultRegistry(),
		verifier: NewCallbackVerifier(randomSecret()),
		sleep:    time.Sleep,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// finalStatusFor maps an Outcome onto the PaymentStatus Initiate should drive a payment to.
func finalStatusFor(outcome Outcome) (payment.PaymentStatus, error) {
	switch outcome {
	case OutcomeSuccess:
		return payment.StatusSuccess, nil
	case OutcomeFailure:
		return payment.StatusFailed, nil
	case OutcomeTimeout:
		return payment.StatusExpired, nil
	default:
		// DefaultRegistry never returns DuplicateCallback/OutOfOrder/InvalidSignature
		// scenarios (those are only reachable via HandleCallback), so this is unreachable
		// today — guarded explicitly rather than silently falling through.
		return "", fmt.Errorf("simulator: outcome %q has no implemented behavior", outcome)
	}
}

// Initiate creates (or, for a repeated IdempotencyKey, looks up) a Payment and, for a freshly
// created one, drives it through CREATED -> PENDING -> the scenario's outcome.
//
// The scenario is resolved *before* anything is created, so an unknown scenario name or a
// PaymentRequest targeting the wrong provider never leaves behind an orphan CREATED Payment.
//
// A zero-delay scenario (success/failure) runs the create-and-drive sequence under
// payment.Service.CreateAndAdvance's single lock acquisition, so concurrent Initiate calls for
// the same brand-new IdempotencyKey cannot interleave mid-sequence.
//
// A delayed scenario (timeout) cannot hold that lock for the whole delay — that would
// serialize every other Initiate/ApplyTransition/HandleCallback call against it. Instead it
// submits to PENDING, waits out the delay, and then applies the outcome as a second step. A
// real callback (via HandleCallback) can legitimately resolve the payment differently during
// that window; that resolution wins — a *payment.TransitionError from the second step is
// treated as "already resolved by something else," not as Initiate's own failure.
func (s *Simulator) Initiate(req payment.PaymentRequest) (payment.PaymentResult, error) {
	if req.Provider.ID != ProviderID {
		return payment.PaymentResult{}, ErrWrongProvider
	}

	scenario, err := s.resolveScenario(req)
	if err != nil {
		return payment.PaymentResult{}, err
	}

	final, err := finalStatusFor(scenario.Outcome)
	if err != nil {
		return payment.PaymentResult{}, err
	}

	if scenario.Delay <= 0 {
		p, err := s.service.CreateAndAdvance(req, []payment.PaymentStatus{payment.StatusPending, final})
		if err != nil {
			return payment.PaymentResult{}, fmt.Errorf("simulator: initiating payment: %w", err)
		}
		return payment.PaymentResult{Payment: p}, nil
	}

	p, err := s.service.CreateAndAdvance(req, []payment.PaymentStatus{payment.StatusPending})
	if err != nil {
		return payment.PaymentResult{}, fmt.Errorf("simulator: initiating payment: %w", err)
	}
	if p.Status != payment.StatusPending {
		// Idempotent replay of a payment already resolved beyond PENDING — return it as-is
		// rather than re-driving the state machine or waiting out the delay again.
		return payment.PaymentResult{Payment: p}, nil
	}

	s.sleep(scenario.Delay)

	// A separate variable, not "p, err =": on error ApplyTransition returns a zero Payment{},
	// which would otherwise clobber p.ID right before the fallback Get(p.ID) below needs it.
	updated, err := s.service.ApplyTransition(p.ID, final)
	if err != nil {
		var transitionErr *payment.TransitionError
		if errors.As(err, &transitionErr) {
			current, getErr := s.service.Get(p.ID)
			if getErr != nil {
				return payment.PaymentResult{}, fmt.Errorf("simulator: fetching payment after superseded outcome: %w", getErr)
			}
			return payment.PaymentResult{Payment: current}, nil
		}
		return payment.PaymentResult{}, fmt.Errorf("simulator: applying %s: %w", final, err)
	}
	p = updated

	return payment.PaymentResult{Payment: p}, nil
}

func (s *Simulator) resolveScenario(req payment.PaymentRequest) (Scenario, error) {
	var opts options
	if raw, ok := req.ProviderOptions["simulator"]; ok {
		if err := json.Unmarshal(raw, &opts); err != nil {
			return Scenario{}, fmt.Errorf("simulator: parsing providerOptions.simulator: %w", err)
		}
	}
	return s.registry.Resolve(opts.Scenario)
}

// SignCallback returns a valid signature for body under this Simulator's own callback-signing
// secret, for constructing test callbacks without reaching into Simulator's internals.
func (s *Simulator) SignCallback(body []byte) string {
	return s.verifier.Sign(body)
}

// HandleCallback processes an asynchronous provider notification: verify first, per
// specs/providers/adapter-contract.md's verifyCallback capability and
// ARCHITECTURE.md §12 ("MUST be called, and MUST reject on failure, before any canonical state
// transition is attempted") — an invalid signature returns ErrInvalidCallbackSignature without
// ever calling Service.ApplyTransition. Only once verified is the body parsed and applied.
//
// Applying the parsed Callback reuses Service.ApplyTransition's existing idempotency semantics
// (specs/state-machines/payment-lifecycle.md "Idempotency and Retries"), which is what makes
// this the same code path for the DUPLICATE_CALLBACK and OUT_OF_ORDER scenarios in
// specs/scenarios/scenario-format.md: a repeat of the current status is a no-op, and a stale
// status claim that conflicts with a later one already applied returns a *payment.TransitionError
// rather than silently regressing the Payment.
func (s *Simulator) HandleCallback(body []byte, signatureHex string) (payment.Payment, error) {
	if !s.verifier.Verify(body, signatureHex) {
		return payment.Payment{}, ErrInvalidCallbackSignature
	}

	cb, err := ParseCallback(body)
	if err != nil {
		return payment.Payment{}, err
	}

	p, err := s.service.ApplyTransition(cb.PaymentID, cb.Status)
	if err != nil {
		return payment.Payment{}, fmt.Errorf("simulator: applying callback: %w", err)
	}
	return p, nil
}
