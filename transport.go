package azpays

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// transport handles all HTTP communication with the AzPays API.
type transport struct {
	config *Config
}

func newTransport(cfg *Config) *transport {
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3
	}
	return &transport{config: cfg}
}

// get performs a GET request and decodes the response into the standard envelope.
func (t *transport) get(ctx context.Context, path string, result any) error {
	return t.doRequest(ctx, http.MethodGet, path, nil, result)
}

// getPaginated performs a GET request for paginated endpoints.
func (t *transport) getPaginated(ctx context.Context, path string, result any) error {
	return t.doRequest(ctx, http.MethodGet, path, nil, result)
}

// post performs a POST request with a JSON body.
func (t *transport) post(ctx context.Context, path string, body any, result any) error {
	return t.doRequest(ctx, http.MethodPost, path, body, result)
}

// put performs a PUT request with a JSON body.
func (t *transport) put(ctx context.Context, path string, body any, result any) error {
	return t.doRequest(ctx, http.MethodPut, path, body, result)
}

// del performs a DELETE request.
func (t *transport) del(ctx context.Context, path string, result any) error {
	return t.doRequest(ctx, http.MethodDelete, path, nil, result)
}

// doRequest executes an HTTP request with retries, auth, and error handling.
func (t *transport) doRequest(ctx context.Context, method, path string, body any, result any) error {
	var bodyReader io.Reader
	var bodyBytes []byte

	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("azpays: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	fullURL := strings.TrimRight(t.config.BaseURL, "/") + "/" + strings.TrimLeft(path, "/")

	var lastErr error
	for attempt := 0; attempt <= t.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 500ms, 1s, 2s, 4s...
			backoff := time.Duration(1<<uint(attempt-1)) * 500 * time.Millisecond
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}

			// Reset body reader for retry
			if bodyBytes != nil {
				bodyReader = bytes.NewReader(bodyBytes)
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
		if err != nil {
			return fmt.Errorf("azpays: failed to create request: %w", err)
		}

		// Set headers
		req.Header.Set("X-API-Key", t.config.APIKey)
		req.Header.Set("User-Agent", t.config.UserAgent)
		req.Header.Set("Accept", "application/json")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		if t.config.Debug {
			log.Printf("[azpays] %s %s", method, fullURL)
			if bodyBytes != nil {
				log.Printf("[azpays] Request Body: %s", string(bodyBytes))
			}
		}

		resp, err := t.config.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("azpays: request failed: %w", err)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("azpays: failed to read response body: %w", err)
			continue
		}

		if t.config.Debug {
			log.Printf("[azpays] Response %d: %s", resp.StatusCode, string(respBody))
		}

		// Retry on 429 or 5xx
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			apiErr := &Error{
				StatusCode: resp.StatusCode,
				Message:    extractMessage(respBody),
				RequestID:  resp.Header.Get("X-Request-Id"),
			}
			lastErr = apiErr

			// On 429, use Retry-After if available
			if resp.StatusCode == 429 {
				if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
					if secs, err := strconv.Atoi(retryAfter); err == nil {
						select {
						case <-ctx.Done():
							return ctx.Err()
						case <-time.After(time.Duration(secs) * time.Second):
						}
					}
				}
			}
			continue
		}

		// Non-retryable error
		if resp.StatusCode >= 400 {
			return &Error{
				StatusCode: resp.StatusCode,
				Message:    extractMessage(respBody),
				RequestID:  resp.Header.Get("X-Request-Id"),
			}
		}

		// Success — decode the response
		if result != nil {
			if err := json.Unmarshal(respBody, result); err != nil {
				return fmt.Errorf("azpays: failed to decode response: %w", err)
			}
		}

		return nil
	}

	// All retries exhausted
	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("azpays: request failed after %d retries", t.config.MaxRetries)
}

// extractMessage tries to parse the error message from a JSON response body.
func extractMessage(body []byte) string {
	// Try standard envelope: {"ok": false, "msg": "..."}
	var envelope struct {
		Msg     string `json:"msg"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil {
		if envelope.Msg != "" {
			return envelope.Msg
		}
		if envelope.Message != "" {
			return envelope.Message
		}
	}
	// Fallback to raw body
	if len(body) > 200 {
		return string(body[:200])
	}
	return string(body)
}

// addQueryParams builds a URL path with query parameters from ListParams and additional params.
func addQueryParams(basePath string, params map[string]string) string {
	if len(params) == 0 {
		return basePath
	}

	values := url.Values{}
	for k, v := range params {
		if v != "" {
			values.Set(k, v)
		}
	}

	encoded := values.Encode()
	if encoded == "" {
		return basePath
	}

	if strings.Contains(basePath, "?") {
		return basePath + "&" + encoded
	}
	return basePath + "?" + encoded
}

// listParamsToMap converts ListParams to a map for query string building.
func listParamsToMap(p ListParams) map[string]string {
	m := make(map[string]string)
	if p.Page > 0 {
		m["page"] = strconv.Itoa(p.Page)
	}
	if p.PerPage > 0 {
		m["per_page"] = strconv.Itoa(p.PerPage)
	}
	if p.Search != "" {
		m["search"] = p.Search
	}
	return m
}
