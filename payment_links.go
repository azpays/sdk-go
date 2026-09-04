package azpays

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// PaymentLinkService handles payment link operations.
type PaymentLinkService struct {
	client *Client
}

// PaymentLink represents a reusable, hosted payment link.
type PaymentLink struct {
	ID                string          `json:"id"`
	MerchantID        string          `json:"merchant_id"`
	UserID            string          `json:"user_id"`
	Slug              string          `json:"slug"`
	Title             string          `json:"title"`
	Description       *string         `json:"description,omitempty"`
	Amount            *float64        `json:"amount,omitempty"`
	Currency          string          `json:"currency"`
	AllowCustomAmount bool            `json:"allow_custom_amount"`
	MinAmount         *float64        `json:"min_amount,omitempty"`
	MaxAmount         *float64        `json:"max_amount,omitempty"`
	AcceptedChains    []string        `json:"accepted_chains"`
	AcceptedTokens    []string        `json:"accepted_tokens"`
	RedirectURL       *string         `json:"redirect_url,omitempty"`
	IsActive          bool            `json:"is_active"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// CreatePaymentLinkRequest creates a new payment link.
type CreatePaymentLinkRequest struct {
	Title             string   `json:"title"`
	Description       *string  `json:"description,omitempty"`
	Amount            *float64 `json:"amount,omitempty"`
	Currency          *string  `json:"currency,omitempty"`
	AllowCustomAmount *bool    `json:"allow_custom_amount,omitempty"`
	MinAmount         *float64 `json:"min_amount,omitempty"`
	MaxAmount         *float64 `json:"max_amount,omitempty"`
	AcceptedChains    []string `json:"accepted_chains,omitempty"`
	AcceptedTokens    []string `json:"accepted_tokens,omitempty"`
	RedirectURL       *string  `json:"redirect_url,omitempty"`
	CustomSlug        *string  `json:"custom_slug,omitempty"`
	Metadata          any      `json:"metadata,omitempty"`
}

// UpdatePaymentLinkRequest modifies an existing payment link.
type UpdatePaymentLinkRequest struct {
	Title             *string  `json:"title,omitempty"`
	Description       *string  `json:"description,omitempty"`
	Amount            *float64 `json:"amount,omitempty"`
	AllowCustomAmount *bool    `json:"allow_custom_amount,omitempty"`
	MinAmount         *float64 `json:"min_amount,omitempty"`
	MaxAmount         *float64 `json:"max_amount,omitempty"`
	AcceptedChains    []string `json:"accepted_chains,omitempty"`
	AcceptedTokens    []string `json:"accepted_tokens,omitempty"`
	RedirectURL       *string  `json:"redirect_url,omitempty"`
	IsActive          *bool    `json:"is_active,omitempty"`
	Metadata          any      `json:"metadata,omitempty"`
}

// ListPaymentLinksParams filters for listing payment links.
type ListPaymentLinksParams struct {
	ListParams
}

// Create creates a new payment link.
func (s *PaymentLinkService) Create(ctx context.Context, req *CreatePaymentLinkRequest) (*PaymentLink, error) {
	var resp Response[PaymentLink]
	err := s.client.transport.post(ctx, "/v1/payment-links", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Get retrieves a payment link by ID.
func (s *PaymentLinkService) Get(ctx context.Context, id string) (*PaymentLink, error) {
	var resp Response[PaymentLink]
	err := s.client.transport.get(ctx, fmt.Sprintf("/v1/payment-links/%s", id), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns a paginated list of payment links.
func (s *PaymentLinkService) List(ctx context.Context, params *ListPaymentLinksParams) ([]PaymentLink, *PaginationMeta, error) {
	if params == nil {
		params = &ListPaymentLinksParams{}
	}

	qp := listParamsToMap(params.ListParams)

	path := addQueryParams("/v1/payment-links", qp)
	var resp PaginatedResponse[PaymentLink]
	err := s.client.transport.getPaginated(ctx, path, &resp)
	if err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// Update modifies an existing payment link.
func (s *PaymentLinkService) Update(ctx context.Context, id string, req *UpdatePaymentLinkRequest) (*PaymentLink, error) {
	var resp Response[PaymentLink]
	err := s.client.transport.put(ctx, fmt.Sprintf("/v1/payment-links/%s", id), req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete removes a payment link.
func (s *PaymentLinkService) Delete(ctx context.Context, id string) error {
	return s.client.transport.del(ctx, fmt.Sprintf("/v1/payment-links/%s", id), nil)
}
