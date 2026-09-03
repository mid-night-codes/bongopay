package payment

// PaymentStatus is the canonical payment lifecycle state, as defined in
// specs/state-machines/payment-lifecycle.md. Provider-specific statuses are never valid
// PaymentStatus values — see ProviderReference.ProviderRawStatus for those.
type PaymentStatus string

const (
	StatusCreated         PaymentStatus = "CREATED"
	StatusPending         PaymentStatus = "PENDING"
	StatusSuccess         PaymentStatus = "SUCCESS"
	StatusFailed          PaymentStatus = "FAILED"
	StatusExpired         PaymentStatus = "EXPIRED"
	StatusCancelled       PaymentStatus = "CANCELLED"
	StatusReversalPending PaymentStatus = "REVERSAL_PENDING"
	StatusReversed        PaymentStatus = "REVERSED"
	StatusRefundPending   PaymentStatus = "REFUND_PENDING"
	StatusRefunded        PaymentStatus = "REFUNDED"
)

// validTransitions is the canonical state graph from
// specs/state-machines/payment-lifecycle.md "Valid Transitions". Keep this map in sync with
// that document — it is the single source of truth this code implements, not the reverse.
var validTransitions = map[PaymentStatus][]PaymentStatus{
	StatusCreated:         {StatusPending, StatusCancelled},
	StatusPending:         {StatusSuccess, StatusFailed, StatusExpired, StatusCancelled},
	StatusSuccess:         {StatusReversalPending, StatusRefundPending},
	StatusFailed:          {},
	StatusExpired:         {},
	StatusCancelled:       {},
	StatusReversalPending: {StatusReversed},
	StatusReversed:        {},
	StatusRefundPending:   {StatusRefunded},
	StatusRefunded:        {},
}

// ValidNextStates returns the states from may validly transition into, per the canonical
// state graph. An empty (non-nil) slice means from has no further transitions.
func ValidNextStates(from PaymentStatus) []PaymentStatus {
	next, ok := validTransitions[from]
	if !ok {
		return nil
	}
	out := make([]PaymentStatus, len(next))
	copy(out, next)
	return out
}

// CanTransition reports whether from -> to is a valid transition in the canonical state
// machine. Per specs/state-machines/payment-lifecycle.md "Idempotency and Retries", an invalid
// transition (including from == to, which is never listed as a self-transition) must be
// treated as a no-op by the caller, not applied — CanTransition only reports validity, it does
// not mutate anything.
func CanTransition(from, to PaymentStatus) bool {
	for _, s := range validTransitions[from] {
		if s == to {
			return true
		}
	}
	return false
}

// IsTerminal reports whether a status has no further valid transitions, per the "Terminal?"
// column in specs/state-machines/payment-lifecycle.md. SUCCESS is deliberately NOT terminal
// here (the spec marks it "No*"): it has real outgoing transitions to REVERSAL_PENDING and
// REFUND_PENDING, even though most payments never leave it.
func IsTerminal(s PaymentStatus) bool {
	return len(validTransitions[s]) == 0
}

// IsValidStatus reports whether s is one of the canonical PaymentStatus values.
func IsValidStatus(s PaymentStatus) bool {
	_, ok := validTransitions[s]
	return ok
}
