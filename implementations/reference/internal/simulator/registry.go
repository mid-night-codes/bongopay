package simulator

import "fmt"

// DefaultScenarioName is used when a PaymentRequest targeting SIMULATOR carries no
// providerOptions.simulator.scenario, per specs/scenarios/scenario-format.md "Rules Specific to
// This Document": absence MUST default to a plain SUCCESS outcome.
const DefaultScenarioName = "success"

// Registry resolves a scenario name (as carried in providerOptions.simulator.scenario) to a
// Scenario. specs/scenarios/scenario-format.md leaves "scenario catalog vs. free-form name" as
// TODO(spec) — this is a fixed, in-memory registry for now, not the final answer to that
// question.
type Registry map[string]Scenario

// ErrUnknownScenario is returned when a scenario name isn't in the Registry.
type ErrUnknownScenario struct {
	Name string
}

func (e *ErrUnknownScenario) Error() string {
	return fmt.Sprintf("simulator: unknown scenario %q", e.Name)
}

// DefaultRegistry returns the scenarios this package actually implements. TIMEOUT,
// DUPLICATE_CALLBACK, OUT_OF_ORDER, and INVALID_SIGNATURE from
// specs/scenarios/scenario-format.md are intentionally absent — they need delay/callback-timing
// machinery this increment doesn't build yet (see implementations/reference/README.md).
func DefaultRegistry() Registry {
	return Registry{
		"success": {Name: "success", Outcome: OutcomeSuccess},
		"failure": {Name: "failure", Outcome: OutcomeFailure},
	}
}

// Resolve looks up name, or DefaultScenarioName if name is empty.
func (r Registry) Resolve(name string) (Scenario, error) {
	if name == "" {
		name = DefaultScenarioName
	}
	s, ok := r[name]
	if !ok {
		return Scenario{}, &ErrUnknownScenario{Name: name}
	}
	return s, nil
}
