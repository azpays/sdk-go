package azpays

import (
	"context"
	"fmt"
	"time"
)

// WalletService handles wallet generation, balance queries, and transfers.
type WalletService struct {
	client *Client
}

// Wallet represents a blockchain wallet.
type Wallet struct {
	ID         string    `json:"id"`
	QR         string    `json:"qr,omitempty"`
	Currency   int       `json:"currency"`
	Chain      string    `json:"chain,omitempty"`
	Type       int       `json:"type"`
	Status     int       `json:"status"`
	PublicKey  string    `json:"public_key"`
	Color      *string   `json:"color,omitempty"`
	Balance    string    `json:"balance"`
	Freeze     string    `json:"freeze"`
	Locked     string    `json:"locked"`
	IsPublic   bool      `json:"is_public"`
	IsWatching bool      `json:"is_watching"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// WalletBalance represents a wallet balance response.
type WalletBalance struct {
	Address     string  `json:"address"`
	Chain       string  `json:"chain"`
	Balance     string  `json:"balance"`
	BalanceUSD  float64 `json:"balance_usd"`
	TokenSymbol string  `json:"token_symbol,omitempty"`
}

// GenerateWalletRequest generates a new wallet.
type GenerateWalletRequest struct {
	Chain       string `json:"chain"`
	TokenSymbol string `json:"token_symbol,omitempty"`
}

// WalletSummaryRequest queries wallet balance summaries.
type WalletSummaryRequest struct {
	Chain       string `json:"chain"`
	Address     string `json:"address"`
	TokenSymbol string `json:"token_symbol,omitempty"`
}

// TransferRequest initiates a blockchain transfer.
type TransferRequest struct {
	Chain       string  `json:"chain"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	Amount      string  `json:"amount"`
	TokenSymbol string  `json:"token_symbol,omitempty"`
	Memo        *string `json:"memo,omitempty"`
}

// TransferResult represents the result of a transfer.
type TransferResult struct {
	TxHash  string `json:"tx_hash"`
	Status  string `json:"status"`
	Chain   string `json:"chain"`
	Amount  string `json:"amount"`
	Fee     string `json:"fee,omitempty"`
	Message string `json:"message,omitempty"`
}

// EstimateFeeRequest calculates gas/network fees.
type EstimateFeeRequest struct {
	Chain       string  `json:"chain"`
	TokenSymbol string  `json:"token_symbol"`
	Amount      float64 `json:"amount"`
	From        string  `json:"from,omitempty"`
	To          string  `json:"to,omitempty"`
}

// FeeEstimate is the fee estimation response.
type FeeEstimate struct {
	Chain           string  `json:"chain"`
	TokenSymbol     string  `json:"token_symbol"`
	Amount          float64 `json:"amount"`
	PlatformFee     float64 `json:"platform_fee"`
	PlatformFeeUSD  float64 `json:"platform_fee_usd"`
	EstimatedGasFee float64 `json:"estimated_gas_fee"`
	EstimatedGasUSD float64 `json:"estimated_gas_usd"`
	TotalFee        float64 `json:"total_fee"`
	TotalFeeUSD     float64 `json:"total_fee_usd"`
	NetAmount       float64 `json:"net_amount"`
	NetAmountUSD    float64 `json:"net_amount_usd"`
	UnitPriceUSD    float64 `json:"unit_price_usd"`
}

// HDSeed represents a BIP-39 master seed.
type HDSeed struct {
	ID         string    `json:"id"`
	MerchantID string    `json:"merchant_id"`
	Label      *string   `json:"label,omitempty"`
	IsBackedUp bool      `json:"is_backed_up"`
	CreatedAt  time.Time `json:"created_at"`
}

// GenerateHDRequest generates a new HD wallet with BIP-39 mnemonic.
type GenerateHDRequest struct {
	Label     *string `json:"label,omitempty"`
	WordCount *int    `json:"word_count,omitempty"` // 12 or 24
}

// ImportHDRequest imports an existing BIP-39 mnemonic.
type ImportHDRequest struct {
	Mnemonic string  `json:"mnemonic"`
	Label    *string `json:"label,omitempty"`
}

// DeriveChildRequest derives a child wallet at a specific index.
type DeriveChildRequest struct {
	SeedID       string `json:"seed_id"`
	Chain        string `json:"chain"`
	AccountIndex *int   `json:"account_index,omitempty"`
	AddressIndex *int   `json:"address_index,omitempty"`
}

// Generate creates a new wallet address via the legacy API-key-authenticated endpoint.
func (s *WalletService) Generate(ctx context.Context, req *GenerateWalletRequest) (*Wallet, error) {
	var resp Response[Wallet]
	err := s.client.transport.post(ctx, "/wallets/generate", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetSummary retrieves balance summary for a wallet address.
func (s *WalletService) GetSummary(ctx context.Context, req *WalletSummaryRequest) (any, error) {
	var resp Response[any]
	err := s.client.transport.post(ctx, "/wallets/summary", req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Transfer initiates a blockchain transfer.
func (s *WalletService) Transfer(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	var resp Response[TransferResult]
	err := s.client.transport.post(ctx, "/wallets/transfer", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// EstimateFee calculates gas and platform fees for a transfer.
func (s *WalletService) EstimateFee(ctx context.Context, req *EstimateFeeRequest) (*FeeEstimate, error) {
	var resp Response[FeeEstimate]
	err := s.client.transport.post(ctx, "/wallets/estimate-fee", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GenerateHD creates a new BIP-39 HD wallet (mother wallet).
func (s *WalletService) GenerateHD(ctx context.Context, req *GenerateHDRequest) (any, error) {
	var resp Response[any]
	err := s.client.transport.post(ctx, "/v1/wallets/hd/generate", req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// ImportHD imports an existing BIP-39 mnemonic.
func (s *WalletService) ImportHD(ctx context.Context, req *ImportHDRequest) (any, error) {
	var resp Response[any]
	err := s.client.transport.post(ctx, "/v1/wallets/hd/import", req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// DeriveChild derives a child wallet at a specific index.
func (s *WalletService) DeriveChild(ctx context.Context, req *DeriveChildRequest) (any, error) {
	var resp Response[any]
	err := s.client.transport.post(ctx, "/v1/wallets/hd/derive", req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// List returns all wallets for the authenticated user.
func (s *WalletService) List(ctx context.Context, params *ListParams) ([]Wallet, *PaginationMeta, error) {
	if params == nil {
		params = &ListParams{}
	}

	qp := listParamsToMap(*params)
	path := addQueryParams("/v1/wallets", qp)

	var resp PaginatedResponse[Wallet]
	err := s.client.transport.getPaginated(ctx, path, &resp)
	if err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// GetBalance retrieves the balance of a specific wallet.
func (s *WalletService) GetBalance(ctx context.Context, id string) (*WalletBalance, error) {
	var resp Response[WalletBalance]
	err := s.client.transport.get(ctx, fmt.Sprintf("/v1/wallets/%s/balance", id), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
