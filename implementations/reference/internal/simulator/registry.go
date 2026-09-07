package simulator

import (
	"fmt"
	"time"
)

// defaultTimeoutDelay is DefaultRegistry's simulated delay before a "timeout" scenario expires.
// It's short enough not to make a demo/test annoying to wait through — the point is exercising
// the CREATED -> PENDING -> EXPIRED path and the possibility of a real callback racing it, not
// modeling a realistic real-world provider timeout window.
const defaultTimeoutDelay = 2 * time.Second

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

// DefaultRegistry returns the scenarios reachable through Simulator.Initiate's scenario
// selection. DUPLICATE_CALLBACK, OUT_OF_ORDER, and INVALID_SIGNATURE from
// specs/scenarios/scenario-format.md are intentionally absent — those are exercised through
// Simulator.HandleCallback directly instead (see implementations/reference/README.md), since
// they're about how a *second*, later event is handled, not an outcome Initiate itself resolves
// to.
func DefaultRegistry() Registry {
	return Registry{
		"success": {Name: "success", Outcome: OutcomeSuccess},
		"failure": {Name: "failure", Outcome: OutcomeFailure},
		"timeout": {Name: "timeout", Outcome: OutcomeTimeout, Delay: defaultTimeoutDelay},
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
