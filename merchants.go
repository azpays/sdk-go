package azpays

import (
	"context"
	"fmt"
	"time"
)

// MerchantService handles merchant self-management operations.
type MerchantService struct {
	client *Client
}

// StrategyWalletGen defines blockchain wallet generation strategy.
type StrategyWalletGen struct {
	Version            int      `json:"version"`
	Chain              string   `json:"chain"`
	TokenSymbol        string   `json:"token_symbol,omitempty"`
	Mode               string   `json:"mode"`
	MerchantWalletID   *string  `json:"merchant_wallet_id,omitempty"`
	DestinationAddress *string  `json:"destination_address,omitempty"`
	Threshold          *float64 `json:"threshold,omitempty"`
	IsActive           bool     `json:"is_active"`
}

// StrategyNotify defines notification channels for payouts.
type StrategyNotify struct {
	IsEmail        bool    `json:"is_email"`
	IsNotification bool    `json:"is_notification"`
	IsWebhook      bool    `json:"is_webhook"`
	WebhookURL     *string `json:"webhook_url,omitempty"`
}

// StrategyPayout defines automated settlement payout rules.
type StrategyPayout struct {
	Version            int            `json:"version"`
	Chain              string         `json:"chain"`
	TokenSymbol        string         `json:"token_symbol"`
	Schedule           string         `json:"schedule"`
	MerchantWalletID   *string        `json:"merchant_wallet_id,omitempty"`
	DestinationAddress *string        `json:"destination_address,omitempty"`
	ThresholdAmount    *float64       `json:"threshold_amount,omitempty"`
	FeePayer           string         `json:"fee_payer,omitempty"`
	Notify             StrategyNotify `json:"notify"`
	IsActive           bool           `json:"is_active"`
}

// Merchant represents a business account entity.
type Merchant struct {
	ID                       string              `json:"id"`
	UserID                   string              `json:"user_id"`
	Name                     string              `json:"name"`
	Logo                     *string             `json:"logo"`
	Type                     int                 `json:"type"`
	Status                   int                 `json:"status"`
	StrategyWalletGeneration []StrategyWalletGen `json:"strategy_wallet_generation"`
	StrategyPayouts          []StrategyPayout    `json:"strategy_payouts"`
	Tell                     *string             `json:"tell"`
	Domain                   *string             `json:"domain"`
	IP                       any                 `json:"ip"`
	Webhook                  *string             `json:"webhook"`
	WebhookSecret            *string             `json:"webhook_secret,omitempty"`
	Callback                 *string             `json:"callback"`
	Description              *string             `json:"description"`
	Support                  any                 `json:"support"`
	Color                    *string             `json:"color"`
	APIKey                   *string             `json:"api_key,omitempty"`
	AcceptedChains           []string            `json:"accepted_chains"`
	AcceptedTokens           []string            `json:"accepted_tokens"`
	CreatedAt                time.Time           `json:"created_at"`
	UpdatedAt                time.Time           `json:"updated_at"`
}

// CreateMerchantRequest creates a new merchant.
type CreateMerchantRequest struct {
	Name                     string              `json:"name"`
	Logo                     *string             `json:"logo,omitempty"`
	Type                     int                 `json:"type"`
	Tell                     *string             `json:"tell,omitempty"`
	Domain                   *string             `json:"domain,omitempty"`
	Webhook                  *string             `json:"webhook,omitempty"`
	Callback                 *string             `json:"callback,omitempty"`
	Description              *string             `json:"description,omitempty"`
	Color                    *string             `json:"color,omitempty"`
	AcceptedChains           []string            `json:"accepted_chains,omitempty"`
	AcceptedTokens           []string            `json:"accepted_tokens,omitempty"`
	StrategyWalletGeneration []StrategyWalletGen `json:"strategy_wallet_generation,omitempty"`
	StrategyPayouts          []StrategyPayout    `json:"strategy_payouts,omitempty"`
}

// UpdateMerchantRequest modifies an existing merchant.
type UpdateMerchantRequest struct {
	Name                     *string              `json:"name,omitempty"`
	Logo                     *string              `json:"logo,omitempty"`
	Tell                     *string              `json:"tell,omitempty"`
	Domain                   *string              `json:"domain,omitempty"`
	Webhook                  *string              `json:"webhook,omitempty"`
	Callback                 *string              `json:"callback,omitempty"`
	Description              *string              `json:"description,omitempty"`
	Color                    *string              `json:"color,omitempty"`
	AcceptedChains           []string             `json:"accepted_chains,omitempty"`
	AcceptedTokens           []string             `json:"accepted_tokens,omitempty"`
	StrategyWalletGeneration *[]StrategyWalletGen `json:"strategy_wallet_generation,omitempty"`
	StrategyPayouts          *[]StrategyPayout    `json:"strategy_payouts,omitempty"`
}

// ListMerchantsParams filters for listing merchants.
type ListMerchantsParams struct {
	ListParams
}

// Create creates a new merchant.
func (s *MerchantService) Create(ctx context.Context, req *CreateMerchantRequest) (*Merchant, error) {
	var resp Response[Merchant]
	err := s.client.transport.post(ctx, "/v1/merchants", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Get retrieves a merchant by ID.
func (s *MerchantService) Get(ctx context.Context, id string) (*Merchant, error) {
	var resp Response[Merchant]
	err := s.client.transport.get(ctx, fmt.Sprintf("/v1/merchants/%s", id), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Me retrieves the current authenticated merchant associated with the client's API key.
func (s *MerchantService) Me(ctx context.Context) (*Merchant, error) {
	return s.Get(ctx, "me")
}

// List returns a paginated list of merchants.
func (s *MerchantService) List(ctx context.Context, params *ListMerchantsParams) ([]Merchant, *PaginationMeta, error) {
	if params == nil {
		params = &ListMerchantsParams{}
	}

	qp := listParamsToMap(params.ListParams)
	path := addQueryParams("/v1/merchants", qp)
	var resp PaginatedResponse[Merchant]
	err := s.client.transport.getPaginated(ctx, path, &resp)
	if err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// Update modifies an existing merchant.
func (s *MerchantService) Update(ctx context.Context, id string, req *UpdateMerchantRequest) (*Merchant, error) {
	var resp Response[Merchant]
	err := s.client.transport.put(ctx, fmt.Sprintf("/v1/merchants/%s", id), req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// RegenerateAPIKey regenerates the API key for a merchant.
func (s *MerchantService) RegenerateAPIKey(ctx context.Context, id string) (*Merchant, error) {
	var resp Response[Merchant]
	err := s.client.transport.post(ctx, fmt.Sprintf("/v1/merchants/%s/regenerate-key", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// AssetDefinition defines a supported cryptocurrency or token.
type AssetDefinition struct {
	ID              int    `json:"id"`
	Symbol          string `json:"symbol"`
	Name            string `json:"name"`
	Chain           string `json:"chain"`
	ContractAddress string `json:"contract_address,omitempty"`
	Decimals        int    `json:"decimals"`
	IsNative        bool   `json:"is_native"`
	Standard        string `json:"standard"`
	Logo            string `json:"logo"`
}

// NetworkDefinition groups blockchain networks with their supported tokens and capabilities.
type NetworkDefinition struct {
	ID                    string            `json:"id"`
	Label                 string            `json:"label"`
	FullName              string            `json:"full_name"`
	Badge                 string            `json:"badge"`
	NetworkType           string            `json:"network_type"`
	Placeholder           string            `json:"placeholder"`
	CoinType              uint32            `json:"coin_type"`
	DerivationPath        string            `json:"derivation_path"`
	IsActive              bool              `json:"is_active"`
	IsHDSupported         bool              `json:"is_hd_supported"`
	IsSingleSupported     bool              `json:"is_single_supported"`
	IsPrivateKeySupported bool              `json:"is_private_key_supported"`
	Logo                  string            `json:"logo"`
	Tokens                []AssetDefinition `json:"tokens"`
}

// Assets returns the dynamic catalog of supported blockchain networks, native coins, and tokens.
// Always query this live catalog rather than hardcoding chain or token parameters.
func (s *MerchantService) Assets(ctx context.Context) ([]NetworkDefinition, error) {
	var resp Response[[]NetworkDefinition]
	err := s.client.transport.get(ctx, "/v1/merchants/assets", &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
