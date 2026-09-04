package azpays

import (
	"context"
	"fmt"
	"time"
)

// PayoutService handles payout operations.
type PayoutService struct {
	client *Client
}

// Payout status constants.
const (
	PayoutStatusDraft           = "draft"
	PayoutStatusPendingApproval = "pending_approval"
	PayoutStatusApproved        = "approved"
	PayoutStatusProcessing      = "processing"
	PayoutStatusCompleted       = "completed"
	PayoutStatusRejected        = "rejected"
	PayoutStatusCancelled       = "cancelled"
	PayoutStatusFailed          = "failed"
)

// Payout represents a payout/settlement disbursement.
type Payout struct {
	ID                 string     `json:"id"`
	MerchantID         string     `json:"merchant_id"`
	PayoutNumber       string     `json:"payout_number"`
	DestinationAddress string     `json:"destination_address"`
	Chain              string     `json:"chain"`
	Network            string     `json:"network"`
	TokenSymbol        string     `json:"token_symbol"`
	TokenAddress       *string    `json:"token_address,omitempty"`
	Amount             float64    `json:"amount"`
	FiatAmount         float64    `json:"fiat_amount"`
	FeeAmount          float64    `json:"fee_amount"`
	NetAmount          float64    `json:"net_amount"`
	Status             string     `json:"status"`
	TxHash             *string    `json:"tx_hash,omitempty"`
	RequestedBy        *string    `json:"requested_by,omitempty"`
	ApprovedBy         *string    `json:"approved_by,omitempty"`
	RejectionReason    *string    `json:"rejection_reason,omitempty"`
	ScheduleType       string     `json:"schedule_type"`
	Memo               *string    `json:"memo,omitempty"`
	IdempotencyKey     *string    `json:"idempotency_key,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	ProcessedAt        *time.Time `json:"processed_at,omitempty"`
}

// PayoutStats provides aggregate payout metrics.
type PayoutStats struct {
	TotalPayouts       int64   `json:"total_payouts"`
	TotalVolumeUSD     float64 `json:"total_volume_usd"`
	PendingCount       int64   `json:"pending_count"`
	PendingVolumeUSD   float64 `json:"pending_volume_usd"`
	CompletedCount     int64   `json:"completed_count"`
	CompletedVolumeUSD float64 `json:"completed_volume_usd"`
	RejectedCount      int64   `json:"rejected_count"`
	FailedCount        int64   `json:"failed_count"`
	AverageFeeUSD      float64 `json:"average_fee_usd"`
}

// CreatePayoutRequest creates a new payout.
type CreatePayoutRequest struct {
	DestinationAddress string   `json:"destination_address"`
	Chain              string   `json:"chain"`
	Network            *string  `json:"network,omitempty"`
	TokenSymbol        string   `json:"token_symbol"`
	TokenAddress       *string  `json:"token_address,omitempty"`
	Amount             float64  `json:"amount"`
	FiatAmount         *float64 `json:"fiat_amount,omitempty"`
	ScheduleType       *string  `json:"schedule_type,omitempty"`
	Memo               *string  `json:"memo,omitempty"`
	IdempotencyKey     *string  `json:"idempotency_key,omitempty"`
}

// ListPayoutsParams filters for listing payouts.
type ListPayoutsParams struct {
	ListParams
	Status string `json:"status,omitempty"`
	Chain  string `json:"chain,omitempty"`
}

// Create creates a new payout request.
func (s *PayoutService) Create(ctx context.Context, req *CreatePayoutRequest) (*Payout, error) {
	var resp Response[Payout]
	err := s.client.transport.post(ctx, "/v1/payouts", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Get retrieves a payout by ID.
func (s *PayoutService) Get(ctx context.Context, id string) (*Payout, error) {
	var resp Response[Payout]
	err := s.client.transport.get(ctx, fmt.Sprintf("/v1/payouts/%s", id), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns a paginated list of payouts.
func (s *PayoutService) List(ctx context.Context, params *ListPayoutsParams) ([]Payout, *PaginationMeta, error) {
	if params == nil {
		params = &ListPayoutsParams{}
	}

	qp := listParamsToMap(params.ListParams)
	if params.Status != "" {
		qp["status"] = params.Status
	}
	if params.Chain != "" {
		qp["chain"] = params.Chain
	}

	path := addQueryParams("/v1/payouts", qp)
	var resp PaginatedResponse[Payout]
	err := s.client.transport.getPaginated(ctx, path, &resp)
	if err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// Stats returns aggregate payout statistics for the authenticated merchant.
func (s *PayoutService) Stats(ctx context.Context) (*PayoutStats, error) {
	var resp Response[PayoutStats]
	err := s.client.transport.get(ctx, "/v1/payouts/stats", &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
