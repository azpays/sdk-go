package azpays

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransportGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("X-API-Key") != "test_key" {
			t.Errorf("expected X-API-Key %q, got %q", "test_key", r.Header.Get("X-API-Key"))
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("expected Accept header application/json")
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/payments/123" {
			t.Errorf("expected path /v1/payments/123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"ok":  true,
			"msg": "Payment retrieved",
			"data": map[string]any{
				"id":          "123",
				"token":       "tok_abc",
				"fiat_amount": 29.99,
				"status":      0,
			},
		})
	}))
	defer server.Close()

	client := NewClient("test_key", WithBaseURL(server.URL))
	payment, err := client.Payments.Get(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.ID != "123" {
		t.Errorf("expected ID %q, got %q", "123", payment.ID)
	}
	if payment.Token != "tok_abc" {
		t.Errorf("expected token %q, got %q", "tok_abc", payment.Token)
	}
	if payment.FiatAmount != 29.99 {
		t.Errorf("expected fiat_amount %f, got %f", 29.99, payment.FiatAmount)
	}
}

func TestTransportPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json")
		}

		var req CreatePaymentRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.FiatAmount != 50.00 {
			t.Errorf("expected fiat_amount %f, got %f", 50.00, req.FiatAmount)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"ok":  true,
			"msg": "Payment created successfully",
			"data": map[string]any{
				"id":          "new_123",
				"merchant_id": "m_auto_from_api_key",
				"fiat_amount": 50.00,
				"status":      0,
				"token":       "tok_new",
			},
		})
	}))
	defer server.Close()

	client := NewClient("test_key", WithBaseURL(server.URL))
	payment, err := client.Payments.Create(context.Background(), &CreatePaymentRequest{
		FiatAmount: 50.00,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.ID != "new_123" {
		t.Errorf("expected ID %q, got %q", "new_123", payment.ID)
	}
	if payment.MerchantID != "m_auto_from_api_key" {
		t.Errorf("expected MerchantID %q, got %q", "m_auto_from_api_key", payment.MerchantID)
	}
}

func TestTransportErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-Id", "req_abc123")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"ok":  false,
			"msg": "Payment not found",
		})
	}))
	defer server.Close()

	client := NewClient("test_key", WithBaseURL(server.URL))
	_, err := client.Payments.Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}

	if !IsNotFound(err) {
		t.Errorf("expected not found error, got %v", err)
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status code 404, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "Payment not found" {
		t.Errorf("expected message %q, got %q", "Payment not found", apiErr.Message)
	}
	if apiErr.RequestID != "req_abc123" {
		t.Errorf("expected request ID %q, got %q", "req_abc123", apiErr.RequestID)
	}
}

func TestTransportUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"ok":  false,
			"msg": "Invalid API key",
		})
	}))
	defer server.Close()

	client := NewClient("bad_key", WithBaseURL(server.URL))
	_, _, err := client.Payments.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsUnauthorized(err) {
		t.Errorf("expected unauthorized error, got %v", err)
	}
}

func TestTransportRetryOn429(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{
				"ok":  false,
				"msg": "Rate limited",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"ok":  true,
			"msg": "OK",
			"data": map[string]any{
				"id": "123",
			},
		})
	}))
	defer server.Close()

	client := NewClient("test_key", WithBaseURL(server.URL), WithMaxRetries(3))
	payment, err := client.Payments.Get(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error after retries: %v", err)
	}
	if payment.ID != "123" {
		t.Errorf("expected ID %q, got %q", "123", payment.ID)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestTransportPaginated(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("expected page=2, got %q", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("per_page") != "10" {
			t.Errorf("expected per_page=10, got %q", r.URL.Query().Get("per_page"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"message": "List retrieved",
			"data": []map[string]any{
				{"id": "p1", "fiat_amount": 10.0, "status": 4},
				{"id": "p2", "fiat_amount": 20.0, "status": 0},
			},
			"pagination": map[string]any{
				"total":        50,
				"per_page":     10,
				"current_page": 2,
				"last_page":    5,
				"from":         11,
				"to":           20,
			},
		})
	}))
	defer server.Close()

	client := NewClient("test_key", WithBaseURL(server.URL))
	payments, pagination, err := client.Payments.List(context.Background(), &ListPaymentsParams{
		ListParams: ListParams{Page: 2, PerPage: 10},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 2 {
		t.Errorf("expected 2 payments, got %d", len(payments))
	}
	if pagination.Total != 50 {
		t.Errorf("expected total 50, got %d", pagination.Total)
	}
	if pagination.CurrentPage != 2 {
		t.Errorf("expected current_page 2, got %d", pagination.CurrentPage)
	}
	if pagination.LastPage != 5 {
		t.Errorf("expected last_page 5, got %d", pagination.LastPage)
	}
}

func TestAddQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		params   map[string]string
		contains []string
	}{
		{
			name:     "no params",
			base:     "/v1/payments",
			params:   map[string]string{},
			contains: []string{"/v1/payments"},
		},
		{
			name:     "with params",
			base:     "/v1/payments",
			params:   map[string]string{"page": "1", "per_page": "10"},
			contains: []string{"/v1/payments?", "page=1", "per_page=10"},
		},
		{
			name:     "empty values filtered",
			base:     "/v1/payments",
			params:   map[string]string{"page": "1", "search": ""},
			contains: []string{"page=1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addQueryParams(tt.base, tt.params)
			for _, c := range tt.contains {
				if !contains(result, c) {
					t.Errorf("expected %q to contain %q", result, c)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestMerchantAssets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/merchants/assets" {
			t.Errorf("expected /v1/merchants/assets, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"ok":  true,
			"msg": "Supported blockchain networks and tokens retrieved",
			"data": []map[string]any{
				{
					"id":        "trx",
					"label":     "TRON",
					"full_name": "TRON Network",
					"badge":     "TRX",
					"tokens": []map[string]any{
						{
							"symbol":   "USDT",
							"name":     "Tether (TRC20)",
							"chain":    "trx",
							"decimals": 6,
							"standard": "TRC-20",
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("test_key", WithBaseURL(server.URL))
	assets, err := client.Merchants.Assets(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(assets) != 1 {
		t.Fatalf("expected 1 network, got %d", len(assets))
	}
	if assets[0].ID != "trx" {
		t.Errorf("expected network ID trx, got %s", assets[0].ID)
	}
	if len(assets[0].Tokens) != 1 || assets[0].Tokens[0].Symbol != "USDT" {
		t.Errorf("expected token USDT, got %v", assets[0].Tokens)
	}
}
