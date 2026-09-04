package azpays

import (
	"context"
	"fmt"
	"time"
)

// CheckoutService handles public checkout operations.
type CheckoutService struct {
	client *Client
}

// CheckoutSession represents a checkout session with payment details.
type CheckoutSession struct {
	Token              string   `json:"token"`
	Status             int      `json:"status"`
	Chain              *string  `json:"chain,omitempty"`
	TokenSymbol        *string  `json:"token_symbol,omitempty"`
	OriginalFiatAmount *float64 `json:"original_fiat_amount,omitempty"`
	DiscountCode       *string  `json:"discount_code,omitempty"`
	DiscountAmount     *float64 `json:"discount_amount,omitempty"`
	FiatAmount         float64  `json:"fiat_amount"`
	Amount             *float64 `json:"amount"`
	ExchangeRate       *float64 `json:"exchange_rate,omitempty"`
	Currency           *int     `json:"currency"`
	CurrencyName       *string  `json:"currency_name"`
	Description        *string  `json:"description"`
	WalletAddress      *string  `json:"wallet_address"`
	QR                 *string  `json:"qr"`
	ExpiresAt          *string  `json:"expires_at"`
	PaidAmount         *float64 `json:"paid_amount,omitempty"`
	UnderpaidAmount    *float64 `json:"underpaid_amount,omitempty"`
	OverpaidAmount     *float64 `json:"overpaid_amount,omitempty"`
	AmountDue          *float64 `json:"amount_due,omitempty"`
	Merchant           struct {
		Name   string  `json:"name"`
		Domain *string `json:"domain"`
	} `json:"merchant"`
}

// CheckoutStatus represents the current checkout payment status.
type CheckoutStatus struct {
	Token       string     `json:"token"`
	Status      int        `json:"status"`
	StatusName  string     `json:"status_name"`
	IsCompleted bool       `json:"is_completed"`
	VerifiedAt  *time.Time `json:"verified_at"`
}

// SelectCoinRequest specifies the blockchain and token for checkout.
type SelectCoinRequest struct {
	Chain    string `json:"chain"`
	Symbol   string `json:"symbol"`
	Token    string `json:"token"`
	Currency int    `json:"currency"`
}

// ApplyDiscountRequest applies a discount code to a checkout session.
type ApplyDiscountRequest struct {
	Code       string  `json:"code"`
	PayerEmail *string `json:"payer_email,omitempty"`
}

// GetSession retrieves checkout session details by payment token.
func (s *CheckoutService) GetSession(ctx context.Context, token string) (*CheckoutSession, error) {
	var resp Response[CheckoutSession]
	err := s.client.transport.get(ctx, fmt.Sprintf("/v1/checkout/%s", token), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListCoins returns available coins for a checkout session.
func (s *CheckoutService) ListCoins(ctx context.Context, token string) ([]any, error) {
	var resp Response[[]any]
	err := s.client.transport.get(ctx, fmt.Sprintf("/v1/checkout/%s/coins", token), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// SelectCoin selects a blockchain and token for payment, locks the rate, and assigns a deposit address.
func (s *CheckoutService) SelectCoin(ctx context.Context, token string, req *SelectCoinRequest) (*CheckoutSession, error) {
	var resp Response[CheckoutSession]
	err := s.client.transport.post(ctx, fmt.Sprintf("/v1/checkout/%s/select-coin", token), req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetStatus polls the current status of a checkout session.
func (s *CheckoutService) GetStatus(ctx context.Context, token string) (*CheckoutStatus, error) {
	var resp Response[CheckoutStatus]
	err := s.client.transport.get(ctx, fmt.Sprintf("/v1/checkout/%s/status", token), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ApplyDiscount applies a discount code to a checkout session.
func (s *CheckoutService) ApplyDiscount(ctx context.Context, token string, req *ApplyDiscountRequest) (*CheckoutSession, error) {
	var resp Response[CheckoutSession]
	err := s.client.transport.post(ctx, fmt.Sprintf("/v1/checkout/%s/discount", token), req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// RemoveDiscount removes a previously applied discount from a checkout session.
func (s *CheckoutService) RemoveDiscount(ctx context.Context, token string) error {
	return s.client.transport.del(ctx, fmt.Sprintf("/v1/checkout/%s/discount", token), nil)
}
