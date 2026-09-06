// Package simulator implements the SIMULATOR provider from
// specs/scenarios/scenario-format.md: a stand-in provider driven by declarative scenarios
// instead of a real network call, per ARCHITECTURE.md §6.
package simulator

import "time"

// Outcome is one of the six simulated outcomes named in specs/scenarios/scenario-format.md.
// Success, Failure, and Timeout are reachable through Simulator.Initiate's scenario selection
// (see Registry); DuplicateCallback, OutOfOrder, and InvalidSignature are exercised through
// Simulator.HandleCallback directly instead — see callback.go.
type Outcome string

const (
	OutcomeSuccess           Outcome = "SUCCESS"
	OutcomeFailure           Outcome = "FAILURE"
	OutcomeTimeout           Outcome = "TIMEOUT"
	OutcomeDuplicateCallback Outcome = "DUPLICATE_CALLBACK"
	OutcomeOutOfOrder        Outcome = "OUT_OF_ORDER"
	OutcomeInvalidSignature  Outcome = "INVALID_SIGNATURE"
)

// Scenario is a named outcome selection, per specs/scenarios/scenario-format.md. Delay is that
// document's Duration: a simulated processing delay before Outcome applies. Zero means instant.
type Scenario struct {
	Name    string
	Outcome Outcome
	Delay   time.Duration
}
