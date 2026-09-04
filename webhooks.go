package azpays

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// WebhookService handles webhook delivery logs and signature verification.
type WebhookService struct {
	client *Client
}

// Webhook delivery status constants.
const (
	WebhookStatusPending   = "pending"
	WebhookStatusDelivered = "delivered"
	WebhookStatusFailed    = "failed"
)

// WebhookDelivery represents an outbox delivery event.
type WebhookDelivery struct {
	ID             string          `json:"id"`
	MerchantID     string          `json:"merchant_id"`
	PaymentID      *string         `json:"payment_id,omitempty"`
	Event          string          `json:"event"`
	URL            string          `json:"url"`
	Payload        json.RawMessage `json:"payload"`
	Signature      string          `json:"signature"`
	Attempts       int             `json:"attempts"`
	MaxAttempts    int             `json:"max_attempts"`
	NextRetryAt    *time.Time      `json:"next_retry_at,omitempty"`
	LastStatusCode *int            `json:"last_status_code,omitempty"`
	LastResponse   *string         `json:"last_response,omitempty"`
	Status         string          `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
	DeliveredAt    *time.Time      `json:"delivered_at,omitempty"`
}

// WebhookEvent represents a parsed webhook event payload.
type WebhookEvent struct {
	Event     string          `json:"event"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// SendTestWebhookRequest sends a test webhook to an endpoint.
type SendTestWebhookRequest struct {
	URL    *string `json:"url,omitempty"`
	Secret *string `json:"secret,omitempty"`
	Event  *string `json:"event,omitempty"`
}

// ListWebhookDeliveriesParams filters for webhook delivery logs.
type ListWebhookDeliveriesParams struct {
	ListParams
	Event  string `json:"event,omitempty"`
	Status string `json:"status,omitempty"`
}

// ListDeliveries returns a paginated list of webhook delivery logs.
func (s *WebhookService) ListDeliveries(ctx context.Context, params *ListWebhookDeliveriesParams) ([]WebhookDelivery, *PaginationMeta, error) {
	if params == nil {
		params = &ListWebhookDeliveriesParams{}
	}

	qp := listParamsToMap(params.ListParams)
	if params.Event != "" {
		qp["event"] = params.Event
	}
	if params.Status != "" {
		qp["status"] = params.Status
	}

	path := addQueryParams("/v1/webhooks/deliveries", qp)
	var resp PaginatedResponse[WebhookDelivery]
	err := s.client.transport.getPaginated(ctx, path, &resp)
	if err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// GetDelivery retrieves a single webhook delivery by ID.
func (s *WebhookService) GetDelivery(ctx context.Context, id string) (*WebhookDelivery, error) {
	var resp Response[WebhookDelivery]
	err := s.client.transport.get(ctx, fmt.Sprintf("/v1/webhooks/deliveries/%s", id), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Test sends a test webhook event to a configured endpoint.
func (s *WebhookService) Test(ctx context.Context, req *SendTestWebhookRequest) error {
	var resp Response[any]
	return s.client.transport.post(ctx, "/v1/webhooks/test", req, &resp)
}

// VerifySignature verifies an AzPays webhook signature.
//
// The signature header format is: t={timestamp},v1={hex_digest}
// where the HMAC-SHA256 digest is computed over: {timestamp}.{raw_body}
//
// Usage:
//
//	valid := client.Webhooks.VerifySignature(
//	    body,
//	    r.Header.Get("X-AzPays-Signature"),
//	    "your_webhook_secret",
//	)
func (s *WebhookService) VerifySignature(payload []byte, signatureHeader string, secret string) bool {
	return VerifyWebhookSignature(payload, signatureHeader, secret)
}

// VerifyWebhookSignature verifies an AzPays webhook signature without needing a client instance.
// This is a standalone function for use in webhook handler middleware.
//
// The signature header format is: t={timestamp},v1={hex_digest}
func VerifyWebhookSignature(payload []byte, signatureHeader string, secret string) bool {
	if signatureHeader == "" || secret == "" {
		return false
	}

	// Parse the signature header: t={timestamp},v1={hex_digest}
	var timestamp string
	var signature string

	parts := strings.Split(signatureHeader, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "t=") {
			timestamp = strings.TrimPrefix(part, "t=")
		} else if strings.HasPrefix(part, "v1=") {
			signature = strings.TrimPrefix(part, "v1=")
		}
	}

	if timestamp == "" || signature == "" {
		return false
	}

	// Verify timestamp is not too old (5 minute tolerance)
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}
	if time.Since(time.Unix(ts, 0)) > 5*time.Minute {
		return false
	}

	// Compute expected signature: HMAC-SHA256("{timestamp}.{payload}")
	message := fmt.Sprintf("%s.%s", timestamp, string(payload))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

// ParseWebhookEvent parses a raw webhook request body into a WebhookEvent.
func ParseWebhookEvent(body []byte) (*WebhookEvent, error) {
	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("azpays: failed to parse webhook event: %w", err)
	}
	return &event, nil
}
