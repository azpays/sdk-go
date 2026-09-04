package azpays

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// InvoiceService handles invoice operations.
type InvoiceService struct {
	client *Client
}

// Invoice status constants.
const (
	InvoiceStatusDraft         = "draft"
	InvoiceStatusOpen          = "open"
	InvoiceStatusPaid          = "paid"
	InvoiceStatusVoid          = "void"
	InvoiceStatusUncollectible = "uncollectible"
	InvoiceStatusExpired       = "expired"
)

// Invoice represents a business invoice entity.
type Invoice struct {
	ID              string           `json:"id"`
	MerchantID      string           `json:"merchant_id"`
	UserID          string           `json:"user_id"`
	InvoiceNumber   string           `json:"invoice_number"`
	Token           string           `json:"token"`
	CustomerID      *string          `json:"customer_id,omitempty"`
	CustomerName    string           `json:"customer_name"`
	CustomerEmail   string           `json:"customer_email"`
	CustomerPhone   *string          `json:"customer_phone,omitempty"`
	CustomerAddress *CustomerAddress `json:"customer_address,omitempty"`
	CustomerTaxID   *string          `json:"customer_tax_id,omitempty"`
	Status          string           `json:"status"`
	Currency        string           `json:"currency"`
	Subtotal        float64          `json:"subtotal"`
	DiscountID      *string          `json:"discount_id,omitempty"`
	DiscountCode    *string          `json:"discount_code,omitempty"`
	DiscountAmount  float64          `json:"discount_amount"`
	TaxPercent      float64          `json:"tax_percent"`
	TaxAmount       float64          `json:"tax_amount"`
	TotalAmount     float64          `json:"total_amount"`
	AmountPaid      float64          `json:"amount_paid"`
	AmountDue       float64          `json:"amount_due"`
	PaymentID       *string          `json:"payment_id,omitempty"`
	PaymentToken    *string          `json:"payment_token,omitempty"`
	SubscriptionID  *string          `json:"subscription_id,omitempty"`
	AcceptedChains  []string         `json:"accepted_chains"`
	AcceptedTokens  []string         `json:"accepted_tokens"`
	PaymentTerms    string           `json:"payment_terms"`
	Memo            *string          `json:"memo,omitempty"`
	Footer          *string          `json:"footer,omitempty"`
	Metadata        json.RawMessage  `json:"metadata,omitempty"`
	DueAt           *time.Time       `json:"due_at,omitempty"`
	FinalizedAt     *time.Time       `json:"finalized_at,omitempty"`
	PaidAt          *time.Time       `json:"paid_at,omitempty"`
	VoidedAt        *time.Time       `json:"voided_at,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	Items           []InvoiceItem    `json:"items,omitempty"`
}

// InvoiceItem represents a line item within an invoice.
type InvoiceItem struct {
	ID             string          `json:"id"`
	InvoiceID      string          `json:"invoice_id"`
	Description    string          `json:"description"`
	Quantity       float64         `json:"quantity"`
	UnitPrice      float64         `json:"unit_price"`
	Amount         float64         `json:"amount"`
	TaxRate        float64         `json:"tax_rate"`
	TaxAmount      float64         `json:"tax_amount"`
	DiscountAmount float64         `json:"discount_amount"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// CustomerAddress represents structured billing address.
type CustomerAddress struct {
	Line1      string `json:"line1,omitempty"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country,omitempty"`
}

// InvoiceStats provides aggregate metrics for the invoicing dashboard.
type InvoiceStats struct {
	TotalInvoices      int64   `json:"total_invoices"`
	TotalBilledVolume  float64 `json:"total_billed_volume"`
	TotalPaidVolume    float64 `json:"total_paid_volume"`
	TotalOpenVolume    float64 `json:"total_open_volume"`
	DraftCount         int64   `json:"draft_count"`
	OpenCount          int64   `json:"open_count"`
	PaidCount          int64   `json:"paid_count"`
	VoidCount          int64   `json:"void_count"`
	UncollectibleCount int64   `json:"uncollectible_count"`
	OverdueCount       int64   `json:"overdue_count"`
}

// InvoiceItemRequest represents a line item input.
type InvoiceItemRequest struct {
	Description    string   `json:"description"`
	Quantity       float64  `json:"quantity"`
	UnitPrice      float64  `json:"unit_price"`
	TaxRate        *float64 `json:"tax_rate,omitempty"`
	DiscountAmount *float64 `json:"discount_amount,omitempty"`
	Metadata       any      `json:"metadata,omitempty"`
}

// CreateInvoiceRequest creates a new invoice.
type CreateInvoiceRequest struct {
	MerchantID      string               `json:"merchant_id"`
	InvoiceNumber   *string              `json:"invoice_number,omitempty"`
	CustomerID      *string              `json:"customer_id,omitempty"`
	SubscriptionID  *string              `json:"subscription_id,omitempty"`
	CustomerName    string               `json:"customer_name"`
	CustomerEmail   string               `json:"customer_email"`
	CustomerPhone   *string              `json:"customer_phone,omitempty"`
	CustomerAddress *CustomerAddress     `json:"customer_address,omitempty"`
	CustomerTaxID   *string              `json:"customer_tax_id,omitempty"`
	Currency        *string              `json:"currency,omitempty"`
	Items           []InvoiceItemRequest `json:"items"`
	DiscountID      *string              `json:"discount_id,omitempty"`
	DiscountCode    *string              `json:"discount_code,omitempty"`
	DiscountAmount  *float64             `json:"discount_amount,omitempty"`
	TaxPercent      *float64             `json:"tax_percent,omitempty"`
	PaymentTerms    *string              `json:"payment_terms,omitempty"`
	DueAt           *time.Time           `json:"due_at,omitempty"`
	Memo            *string              `json:"memo,omitempty"`
	Footer          *string              `json:"footer,omitempty"`
	AcceptedChains  []string             `json:"accepted_chains,omitempty"`
	AcceptedTokens  []string             `json:"accepted_tokens,omitempty"`
	Metadata        any                  `json:"metadata,omitempty"`
	AutoFinalize    bool                 `json:"auto_finalize,omitempty"`
}

// UpdateInvoiceRequest modifies an existing draft invoice.
type UpdateInvoiceRequest struct {
	CustomerName    *string              `json:"customer_name,omitempty"`
	CustomerEmail   *string              `json:"customer_email,omitempty"`
	CustomerPhone   *string              `json:"customer_phone,omitempty"`
	CustomerAddress *CustomerAddress     `json:"customer_address,omitempty"`
	CustomerTaxID   *string              `json:"customer_tax_id,omitempty"`
	Currency        *string              `json:"currency,omitempty"`
	Items           []InvoiceItemRequest `json:"items,omitempty"`
	DiscountID      *string              `json:"discount_id,omitempty"`
	DiscountCode    *string              `json:"discount_code,omitempty"`
	DiscountAmount  *float64             `json:"discount_amount,omitempty"`
	TaxPercent      *float64             `json:"tax_percent,omitempty"`
	PaymentTerms    *string              `json:"payment_terms,omitempty"`
	DueAt           *time.Time           `json:"due_at,omitempty"`
	Memo            *string              `json:"memo,omitempty"`
	Footer          *string              `json:"footer,omitempty"`
	AcceptedChains  []string             `json:"accepted_chains,omitempty"`
	AcceptedTokens  []string             `json:"accepted_tokens,omitempty"`
	Metadata        any                  `json:"metadata,omitempty"`
}

// SendInvoiceRequest specifies how to send an invoice to the customer.
type SendInvoiceRequest struct {
	Email   *string `json:"email,omitempty"`
	Message *string `json:"message,omitempty"`
}

// ListInvoicesParams filters for listing invoices.
type ListInvoicesParams struct {
	ListParams
	Status        string `json:"status,omitempty"`
	CustomerEmail string `json:"customer_email,omitempty"`
}

// Create creates a new invoice.
func (s *InvoiceService) Create(ctx context.Context, req *CreateInvoiceRequest) (*Invoice, error) {
	var resp Response[Invoice]
	err := s.client.transport.post(ctx, "/v1/invoices", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Get retrieves an invoice by ID.
func (s *InvoiceService) Get(ctx context.Context, id string) (*Invoice, error) {
	var resp Response[Invoice]
	err := s.client.transport.get(ctx, fmt.Sprintf("/v1/invoices/%s", id), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns a paginated list of invoices.
func (s *InvoiceService) List(ctx context.Context, params *ListInvoicesParams) ([]Invoice, *PaginationMeta, error) {
	if params == nil {
		params = &ListInvoicesParams{}
	}

	qp := listParamsToMap(params.ListParams)
	if params.Status != "" {
		qp["status"] = params.Status
	}
	if params.CustomerEmail != "" {
		qp["customer_email"] = params.CustomerEmail
	}

	path := addQueryParams("/v1/invoices", qp)
	var resp PaginatedResponse[Invoice]
	err := s.client.transport.getPaginated(ctx, path, &resp)
	if err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// Update modifies a draft invoice.
func (s *InvoiceService) Update(ctx context.Context, id string, req *UpdateInvoiceRequest) (*Invoice, error) {
	var resp Response[Invoice]
	err := s.client.transport.put(ctx, fmt.Sprintf("/v1/invoices/%s", id), req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Finalize transitions an invoice from draft to open.
func (s *InvoiceService) Finalize(ctx context.Context, id string) (*Invoice, error) {
	var resp Response[Invoice]
	err := s.client.transport.post(ctx, fmt.Sprintf("/v1/invoices/%s/finalize", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Send triggers notification delivery of an invoice to the customer.
func (s *InvoiceService) Send(ctx context.Context, id string, req *SendInvoiceRequest) (*Invoice, error) {
	var resp Response[Invoice]
	err := s.client.transport.post(ctx, fmt.Sprintf("/v1/invoices/%s/send", id), req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Void voids an open invoice.
func (s *InvoiceService) Void(ctx context.Context, id string) (*Invoice, error) {
	var resp Response[Invoice]
	err := s.client.transport.post(ctx, fmt.Sprintf("/v1/invoices/%s/void", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// MarkUncollectible marks an invoice as uncollectible.
func (s *InvoiceService) MarkUncollectible(ctx context.Context, id string) (*Invoice, error) {
	var resp Response[Invoice]
	err := s.client.transport.post(ctx, fmt.Sprintf("/v1/invoices/%s/mark-uncollectible", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Stats returns aggregate invoice statistics.
func (s *InvoiceService) Stats(ctx context.Context, merchantID string) (*InvoiceStats, error) {
	qp := make(map[string]string)
	if merchantID != "" {
		qp["merchant_id"] = merchantID
	}
	path := addQueryParams("/v1/invoices/stats", qp)
	var resp Response[InvoiceStats]
	err := s.client.transport.get(ctx, path, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
