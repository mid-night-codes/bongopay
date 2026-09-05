package simulator

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/mid-night-codes/bongopay/implementations/reference/internal/payment"
)

// Callback is the canonical shape of an asynchronous provider notification, per
// specs/providers/adapter-contract.md's parseCallback capability: it claims a PaymentStatus for
// a given payment. It says nothing about *how* it was delivered (HTTP body, queue message,
// etc.) — that's Simulator.HandleCallback's concern.
type Callback struct {
	PaymentID string                `json:"paymentId"`
	Status    payment.PaymentStatus `json:"status"`
}

// ParseCallback decodes a raw callback body into a Callback. Signature verification is
// deliberately separate (CallbackVerifier) and must happen first — see Simulator.HandleCallback.
func ParseCallback(body []byte) (Callback, error) {
	var cb Callback
	if err := json.Unmarshal(body, &cb); err != nil {
		return Callback{}, fmt.Errorf("simulator: parsing callback: %w", err)
	}
	return cb, nil
}

// CallbackVerifier implements specs/providers/adapter-contract.md's verifyCallback capability
// for the simulator: an HMAC-SHA256 over the raw callback body, using a secret known only to
// this Simulator instance. This is a simulator-specific scheme for exercising the verification
// *behavior* real adapters must have — it is not, and does not need to be, any real provider's
// actual signing scheme (each provider has its own; that's an adapter-specific detail, not a
// canonical one).
type CallbackVerifier struct {
	secret []byte
}

// NewCallbackVerifier returns a CallbackVerifier using secret as the HMAC key.
func NewCallbackVerifier(secret []byte) *CallbackVerifier {
	return &CallbackVerifier{secret: secret}
}

// randomSecret returns a fresh, random HMAC key for a Simulator that doesn't need a specific
// (e.g. test-deterministic) one.
func randomSecret() []byte {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("simulator: failed to generate callback verifier secret: %v", err))
	}
	return b
}

// Sign returns the hex-encoded HMAC-SHA256 of body.
func (v *CallbackVerifier) Sign(body []byte) string {
	mac := hmac.New(sha256.New, v.secret)
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify reports whether signatureHex is a valid HMAC-SHA256 of body under this verifier's
// secret. Comparison is constant-time (hmac.Equal) to avoid leaking timing information about
// how much of the signature matched.
func (v *CallbackVerifier) Verify(body []byte, signatureHex string) bool {
	sig, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, v.secret)
	mac.Write(body)
	return hmac.Equal(sig, mac.Sum(nil))
}
