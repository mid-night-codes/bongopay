package payment

import "testing"

// allStatuses must be kept in sync with the const block in lifecycle.go — used to exhaustively
// check every (from, to) pair against specs/state-machines/payment-lifecycle.md.
var allStatuses = []PaymentStatus{
	StatusCreated, StatusPending, StatusSuccess, StatusFailed, StatusExpired, StatusCancelled,
	StatusReversalPending, StatusReversed, StatusRefundPending, StatusRefunded,
}

func TestCanTransition_MatchesSpecGraph(t *testing.T) {
	want := map[PaymentStatus]map[PaymentStatus]bool{
		StatusCreated:         {StatusPending: true, StatusCancelled: true},
		StatusPending:         {StatusSuccess: true, StatusFailed: true, StatusExpired: true, StatusCancelled: true},
		StatusSuccess:         {StatusReversalPending: true, StatusRefundPending: true},
		StatusReversalPending: {StatusReversed: true},
		StatusRefundPending:   {StatusRefunded: true},
	}

	for _, from := range allStatuses {
		for _, to := range allStatuses {
			got := CanTransition(from, to)
			expected := want[from][to]
			if got != expected {
				t.Errorf("CanTransition(%s, %s) = %v, want %v", from, to, got, expected)
			}
		}
	}
}

func TestCanTransition_NoSelfTransitions(t *testing.T) {
	for _, s := range allStatuses {
		if CanTransition(s, s) {
			t.Errorf("CanTransition(%s, %s) = true, want false: no state self-transitions in the spec", s, s)
		}
	}
}

func TestIsTerminal(t *testing.T) {
	terminal := map[PaymentStatus]bool{
		StatusCreated:         false,
		StatusPending:         false,
		StatusSuccess:         false, // spec marks this "No*" — has real outgoing transitions
		StatusFailed:          true,
		StatusExpired:         true,
		StatusCancelled:       true,
		StatusReversalPending: false,
		StatusReversed:        true,
		StatusRefundPending:   false,
		StatusRefunded:        true,
	}

	for _, s := range allStatuses {
		if got := IsTerminal(s); got != terminal[s] {
			t.Errorf("IsTerminal(%s) = %v, want %v", s, got, terminal[s])
		}
	}
}

func TestValidNextStates_ReturnsIndependentCopy(t *testing.T) {
	next := ValidNextStates(StatusCreated)
	next[0] = StatusRefunded // mutate the returned slice

	again := ValidNextStates(StatusCreated)
	if again[0] != StatusPending {
		t.Fatalf("mutating a returned slice affected internal state: got %s, want %s", again[0], StatusPending)
	}
}

func TestIsValidStatus(t *testing.T) {
	for _, s := range allStatuses {
		if !IsValidStatus(s) {
			t.Errorf("IsValidStatus(%s) = false, want true", s)
		}
	}
	if IsValidStatus("NOT_A_REAL_STATUS") {
		t.Error(`IsValidStatus("NOT_A_REAL_STATUS") = true, want false`)
	}
}

// TestRefundAndReversalOnlyFromSuccess guards the state-machine doc's explicit rule: reversal
// and refund sub-flows are only reachable from SUCCESS.
func TestRefundAndReversalOnlyFromSuccess(t *testing.T) {
	for _, from := range allStatuses {
		if from == StatusSuccess {
			continue
		}
		if CanTransition(from, StatusReversalPending) {
			t.Errorf("CanTransition(%s, REVERSAL_PENDING) = true, want false: only reachable from SUCCESS", from)
		}
		if CanTransition(from, StatusRefundPending) {
			t.Errorf("CanTransition(%s, REFUND_PENDING) = true, want false: only reachable from SUCCESS", from)
		}
	}
}
