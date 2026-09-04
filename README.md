# AzPays Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/azpays/sdk-go.svg)](https://pkg.go.dev/github.com/azpays/sdk-go)

The official Go SDK for the [AzPays](https://azpays.net) crypto payment platform.

## Installation

```bash
go get github.com/azpays/sdk-go
```

**Requirements:** Go 1.21+ | **Dependencies:** None (stdlib only)

## API Reference & Interactive Documentation

If you need to inspect raw endpoints, request/response models, or test requests interactively in your browser, check out our OpenAPI specifications and interactive portals:

* **Interactive Scalar Documentation:** [https://api.azpays.net/docs/scalar](https://api.azpays.net/docs/scalar) — Modern API reference with live code generation and testing.
* **Swagger UI Explorer:** [https://api.azpays.net/docs](https://api.azpays.net/docs) — Classic OpenAPI interactive documentation.
* **Raw OpenAPI 3.0.3 Specification:** [https://api.azpays.net/openapi.json](https://api.azpays.net/openapi.json) — Complete machine-readable schema for code generators, Postman, or Insomnia.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    azpays "github.com/azpays/sdk-go"
)

func main() {
    client := azpays.NewClient("az_live_your_api_key")

    // Create a payment — merchant is automatically inferred from your API key!
    payment, err := client.Payments.Create(context.Background(), &azpays.CreatePaymentRequest{
        FiatAmount:  29.99,
        Description: azpays.String("Order #1234"),
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Payment created: %s (token: %s)\n", payment.ID, payment.Token)
}
```

## Configuration

```go
// Production (default)
client := azpays.NewClient("az_live_...")

// Sandbox / Local development
client := azpays.NewClient("az_test_...", azpays.WithBaseURL("http://localhost:8080"))

// With options
client := azpays.NewClient("az_live_...",
    azpays.WithBaseURL("https://api.azpays.net"),
    azpays.WithDebug(),                   // Log requests to stderr
    azpays.WithMaxRetries(5),             // Retry on 429/5xx
    azpays.WithHTTPClient(customClient),  // Custom http.Client
)
```

## Built-in Merchant Scoping

Every API key strictly identifies one merchant. In `sdk-go`, all operations (creating payments, issuing invoices, generating payment links, configuring webhooks) are automatically and securely scoped to your merchant account. You never need to supply a `MerchantID`.

To inspect your authenticated merchant profile:

```go
merchant, err := client.Merchants.Me(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Authenticated as: %s (ID: %s)\n", merchant.Name, merchant.ID)
```

## Services

| Service | Description |
|:---|:---|
| `client.Payments` | Create, retrieve, and list crypto payments |
| `client.Checkout` | Public checkout sessions, coin selection, rate locking |
| `client.Invoices` | Full invoice lifecycle (draft → open → paid/void) |
| `client.PaymentLinks` | Reusable hosted payment walls |
| `client.Webhooks` | Delivery logs + HMAC-SHA256 signature verification |
| `client.Wallets` | HD wallet generation, balance queries, transfers |
| `client.Prices` | Real-time crypto price quotes |
| `client.Payouts` | Merchant settlement disbursements |
| `client.Merchants` | Merchant self-management and settings |

## Usage Examples

### Payments

```go
// Create a payment
payment, err := client.Payments.Create(ctx, &azpays.CreatePaymentRequest{
    FiatAmount:     49.99,
    Description:    azpays.String("Premium Plan"),
    AcceptedChains: []string{"trx", "bnb", "ton"},
    AcceptedTokens: []string{"USDT"},
})

// Get a payment by ID or token
payment, err := client.Payments.Get(ctx, "payment-id-or-token")

// List payments with filters (automatically scoped to your merchant)
payments, pagination, err := client.Payments.List(ctx, &azpays.ListPaymentsParams{
    ListParams: azpays.ListParams{Page: 1, PerPage: 20},
    Statuses:   []int{azpays.PaymentStatusConfirmed, azpays.PaymentStatusPaidOut},
    DateFrom:   "2024-01-01",
})

// Get payment statistics
stats, err := client.Payments.Stats(ctx, nil)
fmt.Printf("Success rate: %.1f%%\n", stats.SuccessRate)
```

### Checkout

```go
// Get checkout session details
session, err := client.Checkout.GetSession(ctx, "payment-token")

// List available coins
coins, err := client.Checkout.ListCoins(ctx, "payment-token")

// Select coin and lock rate (30 min)
session, err := client.Checkout.SelectCoin(ctx, "payment-token", &azpays.SelectCoinRequest{
    Chain:    "trx",
    Symbol:   "USDT",
    Currency: 0,
})
fmt.Printf("Send %f to %s\n", *session.Amount, *session.WalletAddress)

// Poll status
status, err := client.Checkout.GetStatus(ctx, "payment-token")
if status.IsCompleted {
    fmt.Println("Payment confirmed!")
}
```

### Invoices

```go
// Create an invoice
invoice, err := client.Invoices.Create(ctx, &azpays.CreateInvoiceRequest{
    CustomerName:  "John Doe",
    CustomerEmail: "john@example.com",
    Items: []azpays.InvoiceItemRequest{
        {Description: "Web Development", Quantity: 40, UnitPrice: 150.00},
        {Description: "Design Review", Quantity: 5, UnitPrice: 200.00},
    },
    TaxPercent:   azpays.Float64(10.0),
    PaymentTerms: azpays.String("net_30"),
    AutoFinalize: true,
})

// Finalize a draft invoice
invoice, err = client.Invoices.Finalize(ctx, invoice.ID)

// Send to customer
invoice, err = client.Invoices.Send(ctx, invoice.ID, &azpays.SendInvoiceRequest{
    Message: azpays.String("Please find your invoice attached."),
})
```

### Payment Links

```go
// Create a reusable payment link
link, err := client.PaymentLinks.Create(ctx, &azpays.CreatePaymentLinkRequest{
    Title:       "Pro Membership",
    Amount:      azpays.Float64(99.99),
    CustomSlug:  azpays.String("pro-membership"),
    RedirectURL: azpays.String("https://example.com/thank-you"),
})
fmt.Printf("Share: https://pay.azpays.net/%s\n", link.Slug)
```

### Webhook Verification

```go
import "net/http"

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    signature := r.Header.Get("X-AzPays-Signature")

    // Verify signature (standalone — no client needed)
    if !azpays.VerifyWebhookSignature(body, signature, "your_webhook_secret") {
        http.Error(w, "Invalid signature", http.StatusForbidden)
        return
    }

    // Parse the event
    event, err := azpays.ParseWebhookEvent(body)
    if err != nil {
        http.Error(w, "Bad payload", http.StatusBadRequest)
        return
    }

    switch event.Event {
    case "payment.confirmed":
        // Handle confirmed payment
    case "payment.failed":
        // Handle failed payment
    case "payout.completed":
        // Handle completed payout
    }

    w.WriteHeader(http.StatusOK)
}
```

### Wallets

```go
// Generate a wallet via API key
wallet, err := client.Wallets.Generate(ctx, &azpays.GenerateWalletRequest{
    Chain:       "trx",
    TokenSymbol: "USDT",
})

// Get balance
summary, err := client.Wallets.GetSummary(ctx, &azpays.WalletSummaryRequest{
    Chain:   "trx",
    Address: "TAddress...",
})

// Estimate transfer fees
fee, err := client.Wallets.EstimateFee(ctx, &azpays.EstimateFeeRequest{
    Chain:       "trx",
    TokenSymbol: "USDT",
    Amount:      100.0,
})
fmt.Printf("Total fee: %f USDT ($%.2f)\n", fee.TotalFee, fee.TotalFeeUSD)
```

### Prices

```go
// Get current price
quote, err := client.Prices.GetQuote(ctx, "BTC")
fmt.Printf("BTC: $%.2f\n", quote.Price)
```

## Error Handling

```go
payment, err := client.Payments.Get(ctx, "nonexistent")
if err != nil {
    if azpays.IsNotFound(err) {
        fmt.Println("Payment not found")
    } else if azpays.IsUnauthorized(err) {
        fmt.Println("Invalid API key")
    } else if azpays.IsRateLimited(err) {
        fmt.Println("Rate limited — SDK auto-retries, so this is after max retries exhausted")
    } else {
        // Access full error details
        apiErr, ok := err.(*azpays.Error)
        if ok {
            fmt.Printf("Error %d: %s (request: %s)\n", apiErr.StatusCode, apiErr.Message, apiErr.RequestID)
        }
    }
}
```

## Pointer Helpers

Optional fields use pointers. Use the provided helpers:

```go
azpays.String("value")    // *string
azpays.Float64(3.14)      // *float64
azpays.Int(42)            // *int
azpays.Bool(true)         // *bool
```

## Supported Blockchains & Assets (Dynamic)

AzPays dynamically supports a multi-chain catalog across Bitcoin, Ethereum, BNB Smart Chain, Solana, TRON, TON, Polygon, and stablecoins (USDT, USDC).

**Do not hardcode chains, token contracts, or decimals.** Supported networks, contract addresses, derivation paths, and decimal precisions should always be queried dynamically at runtime via the Assets API:

```go
// Fetch the live dynamic catalog of supported blockchains and tokens
networks, err := client.Merchants.Assets(ctx)
if err != nil {
    log.Fatal(err)
}

for _, net := range networks {
    fmt.Printf("Blockchain: %s (%s) [Native: %s, HD Derivation: %s]\n", 
        net.FullName, net.ID, net.Badge, net.DerivationPath)
    for _, token := range net.Tokens {
        fmt.Printf("  • %s (%s) — Standard: %s, Decimals: %d, Contract: %s\n", 
            token.Name, token.Symbol, token.Standard, token.Decimals, token.ContractAddress)
    }
}
```

## Contributing

Contributions, bug reports, and suggestions are warmly welcomed!

1. Fork the repository on GitHub.
2. Create your feature branch (`git checkout -b feature/my-feature`).
3. Ensure all tests and lint checks pass:
   ```bash
   go test -v -race ./...
   go vet ./...
   ```
4. Commit your changes (`git commit -m 'feat: describe your change'`).
5. Push to the branch (`git push origin feature/my-feature`).
6. Open a Pull Request with a clear summary of your work.

## Support

Need integration assistance, found an issue, or have a question?

* **Interactive Documentation:** [https://api.azpays.net/docs/scalar](https://api.azpays.net/docs/scalar)
* **Website:** [https://azpays.net](https://azpays.net)
* **Email Support:** [support@azpays.net](mailto:support@azpays.net)
* **GitHub Issues:** [github.com/azpays/sdk-go/issues](https://github.com/azpays/sdk-go/issues)

## Security

If you believe you have discovered a security vulnerability in this SDK or the AzPays payment infrastructure, please **do not** open a public issue.

Please send your report privately to our security team:
* **Security Contact:** [security@azpays.net](mailto:security@azpays.net)

We acknowledge security reports promptly and collaborate with reporters on coordinated disclosure.

## License

MIT License — see [LICENSE](LICENSE) for details.
