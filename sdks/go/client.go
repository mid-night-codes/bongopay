// Package bongopay is a thin, hand-written HTTP client for contracts/openapi/bongopay.yaml,
// built on the types in the generated subpackage. It is deliberately NOT itself generated:
// oapi-codegen's client generator pulls in github.com/oapi-codegen/runtime solely to encode
// path parameters using OpenAPI's style/explode rules — machinery this API doesn't need, since
// its only path parameter (GET /payments/{id}'s id) is a plain string. Hand-writing this thin a
// client keeps the SDK dependency-free (see docs/development/dependency-policy.md) while every
// request/response shape still comes from the generated types, not a hand-typed guess.
package bongopay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mid-night-codes/bongopay/sdks/go/generated"
)

// Re-exported so callers only need to import this package, not generated too.
type (
	PaymentRequest = generated.PaymentRequest
	PaymentResult  = generated.PaymentResult
	Payment        = generated.Payment
	Money          = generated.Money
	Currency       = generated.Currency
	Provider       = generated.Provider
)

// Client is a BongoPay API client.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// New returns a Client targeting baseURL (e.g. "http://localhost:8080"), using
// http.DefaultClient. Set HTTPClient directly for a custom timeout/transport.
func New(baseURL string) *Client {
	return &Client{BaseURL: baseURL, HTTPClient: http.DefaultClient}
}

// InitiatePayment calls POST /payments.
func (c *Client) InitiatePayment(ctx context.Context, req PaymentRequest) (*PaymentResult, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("bongopay: encoding request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/payments", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("bongopay: building request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	var result PaymentResult
	if err := c.do(httpReq, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPayment calls GET /payments/{id}.
func (c *Client) GetPayment(ctx context.Context, id string) (*Payment, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/payments/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, fmt.Errorf("bongopay: building request: %w", err)
	}

	var result Payment
	if err := c.do(httpReq, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// APIError is returned when the server responds with a non-2xx status. It only carries
// StatusCode and Message (from the Error schema) — bongopay.yaml's Error shape is itself
// provisional (specs/errors/error-model.md is unwritten), so StatusCode is the only thing a
// caller can reliably branch on today.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("bongopay: server returned %d: %s", e.StatusCode, e.Message)
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("bongopay: request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("bongopay: reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr generated.Error
		if jsonErr := json.Unmarshal(data, &apiErr); jsonErr == nil && apiErr.Error != "" {
			return &APIError{StatusCode: resp.StatusCode, Message: apiErr.Error}
		}
		return &APIError{StatusCode: resp.StatusCode, Message: string(data)}
	}

	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("bongopay: decoding response: %w", err)
	}
	return nil
}
