// Package azpays provides the official Go SDK for the AzPays crypto payment platform.
//
// For interactive API documentation and complete OpenAPI 3.0 specifications, visit:
//   - Scalar Interactive Docs: https://api.azpays.net/docs/scalar
//   - Swagger UI: https://api.azpays.net/docs
//   - OpenAPI 3.0 JSON Schema: https://api.azpays.net/openapi.json
//
// Create a client with your API key and start accepting crypto payments:
//
//	client := azpays.NewClient("az_live_your_api_key")
//	payment, err := client.Payments.Create(ctx, &azpays.CreatePaymentRequest{
//	    FiatAmount: 29.99,
//	})
package azpays

import (
	"net/http"
	"time"
)

const (
	// Version is the current SDK version.
	Version = "1.0.0"

	// DefaultBaseURL is the production AzPays API endpoint.
	DefaultBaseURL = "https://api.azpays.net"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second

	// DefaultUserAgent is the default User-Agent header.
	DefaultUserAgent = "azpays-sdk-go/" + Version
)

// Client is the AzPays API client. Access resources through service fields.
type Client struct {
	// Services — merchant-focused resources
	Payments     *PaymentService
	Checkout     *CheckoutService
	Invoices     *InvoiceService
	PaymentLinks *PaymentLinkService
	Webhooks     *WebhookService
	Wallets      *WalletService
	Prices       *PriceService
	Payouts      *PayoutService
	Merchants    *MerchantService

	// transport handles all HTTP communication
	transport *transport
}

// NewClient creates a new AzPays API client with the given API key.
//
//	client := azpays.NewClient("az_live_your_api_key")
//	client := azpays.NewClient("az_test_your_api_key", azpays.WithBaseURL("http://localhost:8080"))
func NewClient(apiKey string, opts ...Option) *Client {
	cfg := &Config{
		APIKey:    apiKey,
		BaseURL:   DefaultBaseURL,
		UserAgent: DefaultUserAgent,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	t := newTransport(cfg)
	c := &Client{transport: t}

	// Wire up services
	c.Payments = &PaymentService{client: c}
	c.Checkout = &CheckoutService{client: c}
	c.Invoices = &InvoiceService{client: c}
	c.PaymentLinks = &PaymentLinkService{client: c}
	c.Webhooks = &WebhookService{client: c}
	c.Wallets = &WalletService{client: c}
	c.Prices = &PriceService{client: c}
	c.Payouts = &PayoutService{client: c}
	c.Merchants = &MerchantService{client: c}

	return c
}
