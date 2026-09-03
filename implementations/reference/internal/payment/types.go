// Package payment implements the canonical concepts defined in
// specs/payments/payment-contract.md. Field names and shapes here must trace back to that
// document — see AGENTS.md §7: implementation follows spec, never the reverse.
package payment

import (
	"encoding/json"
	"time"
)

// Money is an amount in a currency's minor unit (e.g. cents). Never a float — see
// specs/payments/README.md.
type Money struct {
	Value    int64    `json:"value"`
	Currency Currency `json:"currency"`
}

// Currency is an ISO 4217 alphabetic code, e.g. "TZS", "KES", "USD".
type Currency struct {
	Code string `json:"code"`
}

// Metadata is free-form, opaque key/value pairs supplied by the caller and echoed back
// unmodified. It MUST NOT be interpreted by core or by adapters.
type Metadata map[string]string

// ProviderOptions holds namespaced, provider-specific data (providerOptions.<provider> in the
// spec). Each namespace's contents are opaque to core and orchestration — see
// specs/providers/extensions.md — which is why values are kept as raw JSON rather than
// unmarshaled into a concrete type here.
type ProviderOptions map[string]json.RawMessage

// Provider identifies a payment provider, e.g. "MPESA", "AIRTEL_MONEY", or "SIMULATOR".
// SIMULATOR is always a valid Provider — see ARCHITECTURE.md §6.
type Provider struct {
	ID string `json:"id"`
}

// PaymentMethodType is intentionally open-ended: specs/payments/payment-contract.md marks the
// enum values as TODO(spec) pending finalization. These constants are illustrative, not
// exhaustive — do not assume this is the complete set.
type PaymentMethodType string

const (
	PaymentMethodMobileMoney  PaymentMethodType = "MOBILE_MONEY"
	PaymentMethodCard         PaymentMethodType = "CARD"
	PaymentMethodBankTransfer PaymentMethodType = "BANK_TRANSFER"
)

// PaymentMethod describes how a payment is made.
type PaymentMethod struct {
	Type PaymentMethodType `json:"type"`
}

// CustomerReference identifies the customer/payer in a provider-neutral way. Per
// specs/payments/payment-contract.md, whether exactly one of MSISDN/ExternalID is required is
// still TODO(spec) — both are represented as optional (pointers) here until that's resolved.
type CustomerReference struct {
	MSISDN     *string  `json:"msisdn,omitempty"`
	ExternalID *string  `json:"externalId,omitempty"`
	Metadata   Metadata `json:"metadata,omitempty"`
}

// ProviderReference identifies how a specific provider refers to a payment, once known.
// ProviderRawStatus is for observability/debugging only — it MUST NOT be used as a canonical
// PaymentStatus. See specs/state-machines/payment-lifecycle.md.
type ProviderReference struct {
	Provider              Provider `json:"provider"`
	ProviderTransactionID *string  `json:"providerTransactionId,omitempty"`
	ProviderRawStatus     *string  `json:"providerRawStatus,omitempty"`
}

// PaymentRequest is the input to initiate a payment.
type PaymentRequest struct {
	Provider          Provider          `json:"provider"`
	Amount            Money             `json:"amount"`
	CustomerReference CustomerReference `json:"customerReference"`
	PaymentMethod     *PaymentMethod    `json:"paymentMethod,omitempty"`
	IdempotencyKey    string            `json:"idempotencyKey"`
	Metadata          Metadata          `json:"metadata,omitempty"`
	ProviderOptions   ProviderOptions   `json:"providerOptions,omitempty"`
}

// Payment is the canonical, persistent representation of a payment, once created.
type Payment struct {
	ID                string             `json:"id"`
	Status            PaymentStatus      `json:"status"`
	Amount            Money              `json:"amount"`
	Provider          Provider           `json:"provider"`
	ProviderReference *ProviderReference `json:"providerReference,omitempty"`
	CustomerReference CustomerReference  `json:"customerReference"`
	IdempotencyKey    string             `json:"idempotencyKey"`
	CreatedAt         time.Time          `json:"createdAt"`
	UpdatedAt         time.Time          `json:"updatedAt"`
	Metadata          Metadata           `json:"metadata,omitempty"`
	ProviderOptions   ProviderOptions    `json:"providerOptions,omitempty"`
}

// PaymentResult is returned synchronously from an initiate/query/refund/reversal operation.
// Deliberately thin — most state changes are observed via events, not synchronous results.
type PaymentResult struct {
	Payment Payment `json:"payment"`
}
