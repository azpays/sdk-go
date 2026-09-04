package azpays

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("az_test_key123")

	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.Payments == nil {
		t.Fatal("expected non-nil Payments service")
	}
	if client.Checkout == nil {
		t.Fatal("expected non-nil Checkout service")
	}
	if client.Invoices == nil {
		t.Fatal("expected non-nil Invoices service")
	}
	if client.PaymentLinks == nil {
		t.Fatal("expected non-nil PaymentLinks service")
	}
	if client.Webhooks == nil {
		t.Fatal("expected non-nil Webhooks service")
	}
	if client.Wallets == nil {
		t.Fatal("expected non-nil Wallets service")
	}
	if client.Prices == nil {
		t.Fatal("expected non-nil Prices service")
	}
	if client.Payouts == nil {
		t.Fatal("expected non-nil Payouts service")
	}
	if client.Merchants == nil {
		t.Fatal("expected non-nil Merchants service")
	}
}

func TestNewClientDefaults(t *testing.T) {
	client := NewClient("az_test_key")

	cfg := client.transport.config
	if cfg.BaseURL != DefaultBaseURL {
		t.Errorf("expected base URL %q, got %q", DefaultBaseURL, cfg.BaseURL)
	}
	if cfg.UserAgent != DefaultUserAgent {
		t.Errorf("expected user agent %q, got %q", DefaultUserAgent, cfg.UserAgent)
	}
	if cfg.APIKey != "az_test_key" {
		t.Errorf("expected API key %q, got %q", "az_test_key", cfg.APIKey)
	}
	if cfg.MaxRetries != 3 {
		t.Errorf("expected max retries 3, got %d", cfg.MaxRetries)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	customClient := &http.Client{Timeout: 60 * time.Second}
	client := NewClient("az_live_key",
		WithBaseURL("http://localhost:8080"),
		WithHTTPClient(customClient),
		WithDebug(),
		WithUserAgent("custom-agent/1.0"),
		WithMaxRetries(5),
	)

	cfg := client.transport.config
	if cfg.BaseURL != "http://localhost:8080" {
		t.Errorf("expected base URL %q, got %q", "http://localhost:8080", cfg.BaseURL)
	}
	if cfg.HTTPClient != customClient {
		t.Error("expected custom HTTP client")
	}
	if !cfg.Debug {
		t.Error("expected debug to be true")
	}
	if cfg.UserAgent != "custom-agent/1.0" {
		t.Errorf("expected user agent %q, got %q", "custom-agent/1.0", cfg.UserAgent)
	}
	if cfg.MaxRetries != 5 {
		t.Errorf("expected max retries 5, got %d", cfg.MaxRetries)
	}
}

func TestHelpers(t *testing.T) {
	s := String("hello")
	if *s != "hello" {
		t.Errorf("expected %q, got %q", "hello", *s)
	}

	f := Float64(3.14)
	if *f != 3.14 {
		t.Errorf("expected %f, got %f", 3.14, *f)
	}

	i := Int(42)
	if *i != 42 {
		t.Errorf("expected %d, got %d", 42, *i)
	}

	b := Bool(true)
	if !*b {
		t.Error("expected true")
	}
}
