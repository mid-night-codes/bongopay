package simulator

import (
	"encoding/json"
	"fmt"

	"github.com/mid-night-codes/bongopay/implementations/reference/internal/payment"
)

// ProviderID is the Provider.id a PaymentRequest must use to reach this simulator, per
// specs/payments/payment-contract.md ("SIMULATOR is always a valid Provider").
const ProviderID = "SIMULATOR"

// ErrWrongProvider is returned when Initiate is called with a PaymentRequest not targeting
// ProviderID.
var ErrWrongProvider = fmt.Errorf("simulator: PaymentRequest.Provider.ID must be %q", ProviderID)

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
}

// New returns a Simulator backed by service, using DefaultRegistry for scenario resolution.
func New(service *payment.Service) *Simulator {
	return &Simulator{service: service, registry: DefaultRegistry()}
}

// Initiate creates (or, for a repeated IdempotencyKey, looks up) a Payment and, for a freshly
// created one, drives it through CREATED -> PENDING -> the scenario's outcome.
//
// The scenario is resolved *before* anything is created, so an unknown scenario name or a
// PaymentRequest targeting the wrong provider never leaves behind an orphan CREATED Payment.
// The create-and-drive sequence itself runs under payment.Service.CreateAndAdvance's single
// lock acquisition, so concurrent Initiate calls for the same brand-new IdempotencyKey cannot
// interleave mid-sequence: whichever starts first runs to completion (or a genuine replay sees
// that completed result) before any other caller for that key can observe or advance it.
func (s *Simulator) Initiate(req payment.PaymentRequest) (payment.PaymentResult, error) {
	if req.Provider.ID != ProviderID {
		return payment.PaymentResult{}, ErrWrongProvider
	}

	scenario, err := s.resolveScenario(req)
	if err != nil {
		return payment.PaymentResult{}, err
	}

	var final payment.PaymentStatus
	switch scenario.Outcome {
	case OutcomeSuccess:
		final = payment.StatusSuccess
	case OutcomeFailure:
		final = payment.StatusFailed
	default:
		// DefaultRegistry only ever returns Success/Failure scenarios today, so this is
		// unreachable — guarded explicitly rather than silently falling through.
		return payment.PaymentResult{}, fmt.Errorf("simulator: outcome %q has no implemented behavior", scenario.Outcome)
	}

	p, err := s.service.CreateAndAdvance(req, []payment.PaymentStatus{payment.StatusPending, final})
	if err != nil {
		return payment.PaymentResult{}, fmt.Errorf("simulator: initiating payment: %w", err)
	}

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
