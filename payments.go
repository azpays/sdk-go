package azpays

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// PaymentService handles payment operations.
type PaymentService struct {
	client *Client
}

// Payment status constants matching the API state machine.
const (
	PaymentStatusRegistered = 0
	PaymentStatusStarted    = 1
	PaymentStatusDetected   = 2
	PaymentStatusVerifying  = 3
	PaymentStatusConfirmed  = 4
	PaymentStatusCanceled   = 5
	PaymentStatusPartial    = 6
	PaymentStatusPaidOut    = 7
)

// Payment represents a payment request entity.
type Payment struct {
	ID              string     `json:"id"`
	Token           string     `json:"token"`
	MerchantID      string     `json:"merchant_id"`
	PayeeID         string     `json:"payee_id"`
	PayerID         *string    `json:"payer_id"`
	IdempotencyKey  *string    `json:"idempotency_key,omitempty"`
	Chain           *string    `json:"chain,omitempty"`
	Network         *string    `json:"network,omitempty"`
	TokenSymbol     *string    `json:"token_symbol,omitempty"`
	TokenAddress    *string    `json:"token_address,omitempty"`
	FiatAmount      float64    `json:"fiat_amount"`
	Amount          *float64   `json:"amount"`
	PaidAmount      float64    `json:"paid_amount"`
	PaidFiatAmount  float64    `json:"paid_fiat_amount"`
	UnderpaidAmount float64    `json:"underpaid_amount"`
	OverpaidAmount  float64    `json:"overpaid_amount"`
	ExchangeRate    *float64   `json:"exchange_rate,omitempty"`
	RateExpiresAt   *time.Time `json:"rate_expires_at,omitempty"`
	Currency        int        `json:"currency"`
	Factor          *string    `json:"factor"`
	Description     *string    `json:"description"`
	AcceptedChains  []string   `json:"accepted_chains"`
	AcceptedTokens  []string   `json:"accepted_tokens"`
	Merchant        *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"merchant,omitempty"`
	StartedAt  *time.Time `json:"started_at"`
	VerifiedAt *time.Time `json:"verified_at"`
	Status     int        `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// PaymentStats provides aggregate metrics for payments.
type PaymentStats struct {
	TotalPayments   int64   `json:"total_payments"`
	TotalVolume     float64 `json:"total_volume"`
	TotalPaidVolume float64 `json:"total_paid_volume"`
	ConfirmedCount  int64   `json:"confirmed_count"`
	PendingCount    int64   `json:"pending_count"`
	CanceledCount   int64   `json:"canceled_count"`
	PartialCount    int64   `json:"partial_count"`
	SuccessRate     float64 `json:"success_rate"`
}

// PaymentReports provides reporting metrics for payments.
type PaymentReports struct {
	TotalPayments     int64   `json:"total_payments"`
	ConfirmedPayments int64   `json:"confirmed_payments"`
	TotalVolume       float64 `json:"total_volume"`
}

// CreatePaymentRequest is the input for creating a new payment.
type CreatePaymentRequest struct {
	FiatAmount     float64  `json:"fiat_amount"`
	Description    *string  `json:"description,omitempty"`
	IdempotencyKey *string  `json:"idempotency_key,omitempty"`
	AcceptedChains []string `json:"accepted_chains,omitempty"`
	AcceptedTokens []string `json:"accepted_tokens,omitempty"`
}

// ListPaymentsParams filters for listing payments.
type ListPaymentsParams struct {
	ListParams
	Statuses []int  `json:"statuses,omitempty"`
	DateFrom string `json:"date_from,omitempty"` // RFC3339 or YYYY-MM-DD
	DateTo   string `json:"date_to,omitempty"`   // RFC3339 or YYYY-MM-DD
}

// PaymentStatsParams filters for payment statistics.
type PaymentStatsParams struct {
	DateFrom string `json:"date_from,omitempty"`
	DateTo   string `json:"date_to,omitempty"`
}

// Create creates a new payment request.
//
//	payment, err := client.Payments.Create(ctx, &azpays.CreatePaymentRequest{
//	    FiatAmount:  29.99,
//	    Description: azpays.String("Order #1234"),
//	})
func (s *PaymentService) Create(ctx context.Context, req *CreatePaymentRequest) (*Payment, error) {
	var resp Response[Payment]
	err := s.client.transport.post(ctx, "/v1/payments", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Get retrieves a payment by ID or token.
func (s *PaymentService) Get(ctx context.Context, idOrToken string) (*Payment, error) {
	var resp Response[Payment]
	err := s.client.transport.get(ctx, fmt.Sprintf("/v1/payments/%s", idOrToken), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns a paginated list of payments.
func (s *PaymentService) List(ctx context.Context, params *ListPaymentsParams) ([]Payment, *PaginationMeta, error) {
	if params == nil {
		params = &ListPaymentsParams{}
	}

	qp := listParamsToMap(params.ListParams)
	if len(params.Statuses) > 0 {
		strs := make([]string, len(params.Statuses))
		for i, s := range params.Statuses {
			strs[i] = strconv.Itoa(s)
		}
		qp["statuses"] = strings.Join(strs, ",")
	}
	if params.DateFrom != "" {
		qp["date_from"] = params.DateFrom
	}
	if params.DateTo != "" {
		qp["date_to"] = params.DateTo
	}

	path := addQueryParams("/v1/payments", qp)
	var resp PaginatedResponse[Payment]
	err := s.client.transport.getPaginated(ctx, path, &resp)
	if err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// Stats returns aggregate payment statistics.
func (s *PaymentService) Stats(ctx context.Context, params *PaymentStatsParams) (*PaymentStats, error) {
	qp := make(map[string]string)
	if params != nil {
		if params.DateFrom != "" {
			qp["date_from"] = params.DateFrom
		}
		if params.DateTo != "" {
			qp["date_to"] = params.DateTo
		}
	}

	path := addQueryParams("/v1/payments/stats", qp)
	var resp Response[PaymentStats]
	err := s.client.transport.get(ctx, path, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Reports returns payment reporting metrics.
func (s *PaymentService) Reports(ctx context.Context, params *PaymentStatsParams) (*PaymentReports, error) {
	qp := make(map[string]string)
	if params != nil {
		if params.DateFrom != "" {
			qp["date_from"] = params.DateFrom
		}
		if params.DateTo != "" {
			qp["date_to"] = params.DateTo
		}
	}

	path := addQueryParams("/v1/payments/reports", qp)
	var resp Response[PaymentReports]
	err := s.client.transport.get(ctx, path, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
