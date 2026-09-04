// Package simulator implements the SIMULATOR provider from
// specs/scenarios/scenario-format.md: a stand-in provider driven by declarative scenarios
// instead of a real network call, per ARCHITECTURE.md §6.
package simulator

// Outcome is one of the six simulated outcomes named in specs/scenarios/scenario-format.md.
// Only Success and Failure are implemented so far — see Registry.
type Outcome string

const (
	OutcomeSuccess           Outcome = "SUCCESS"
	OutcomeFailure           Outcome = "FAILURE"
	OutcomeTimeout           Outcome = "TIMEOUT"
	OutcomeDuplicateCallback Outcome = "DUPLICATE_CALLBACK"
	OutcomeOutOfOrder        Outcome = "OUT_OF_ORDER"
	OutcomeInvalidSignature  Outcome = "INVALID_SIGNATURE"
)

// Scenario is a named outcome selection, per specs/scenarios/scenario-format.md. Delay/timing
// controls from that document are not implemented yet — every scenario here is synchronous.
type Scenario struct {
	Name    string
	Outcome Outcome
}
