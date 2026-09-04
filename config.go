package azpays

import (
	"net/http"
	"time"
)

// Config holds all configuration for the AzPays API client.
type Config struct {
	// APIKey is the merchant API key (required).
	// Format: az_live_... (production) or az_test_... (sandbox).
	APIKey string

	// BaseURL is the API base URL. Defaults to https://api.azpays.net.
	BaseURL string

	// HTTPClient is the underlying HTTP client. Defaults to a client with 30s timeout.
	HTTPClient *http.Client

	// UserAgent is the User-Agent header sent with every request.
	UserAgent string

	// Debug enables request/response logging to stderr.
	Debug bool

	// MaxRetries is the maximum number of retries on 429/5xx errors. Defaults to 3.
	MaxRetries int
}

// Option configures the AzPays client.
type Option func(*Config)

// WithBaseURL sets a custom API base URL (e.g. for local development).
//
//	client := azpays.NewClient(key, azpays.WithBaseURL("http://localhost:8080"))
func WithBaseURL(url string) Option {
	return func(c *Config) {
		c.BaseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

// WithDebug enables debug logging of requests and responses.
func WithDebug() Option {
	return func(c *Config) {
		c.Debug = true
	}
}

// WithUserAgent sets a custom User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Config) {
		c.UserAgent = ua
	}
}

// WithMaxRetries sets the maximum number of retries on transient errors.
func WithMaxRetries(n int) Option {
	return func(c *Config) {
		c.MaxRetries = n
	}
}

// WithTimeout sets the HTTP request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Config) {
		if c.HTTPClient == nil {
			c.HTTPClient = &http.Client{Timeout: d}
		} else {
			c.HTTPClient.Timeout = d
		}
	}
}

