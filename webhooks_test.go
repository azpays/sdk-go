package azpays

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

func TestVerifyWebhookSignature(t *testing.T) {
	secret := "my_webhook_secret"
	payload := []byte(`{"event":"payment.confirmed","data":{"id":"pay_123"}}`)
	timestamp := time.Now().Unix()

	// Compute valid signature matching the API's ComputeHMACSignature format
	message := fmt.Sprintf("%d.%s", timestamp, string(payload))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	digest := hex.EncodeToString(mac.Sum(nil))
	signature := fmt.Sprintf("t=%d,v1=%s", timestamp, digest)

	t.Run("valid signature", func(t *testing.T) {
		valid := VerifyWebhookSignature(payload, signature, secret)
		if !valid {
			t.Error("expected valid signature")
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		valid := VerifyWebhookSignature(payload, signature, "wrong_secret")
		if valid {
			t.Error("expected invalid signature with wrong secret")
		}
	})

	t.Run("tampered payload", func(t *testing.T) {
		tampered := []byte(`{"event":"payment.confirmed","data":{"id":"pay_999"}}`)
		valid := VerifyWebhookSignature(tampered, signature, secret)
		if valid {
			t.Error("expected invalid signature with tampered payload")
		}
	})

	t.Run("empty signature", func(t *testing.T) {
		valid := VerifyWebhookSignature(payload, "", secret)
		if valid {
			t.Error("expected invalid with empty signature")
		}
	})

	t.Run("empty secret", func(t *testing.T) {
		valid := VerifyWebhookSignature(payload, signature, "")
		if valid {
			t.Error("expected invalid with empty secret")
		}
	})

	t.Run("expired timestamp", func(t *testing.T) {
		oldTimestamp := time.Now().Add(-10 * time.Minute).Unix()
		oldMessage := fmt.Sprintf("%d.%s", oldTimestamp, string(payload))
		oldMac := hmac.New(sha256.New, []byte(secret))
		oldMac.Write([]byte(oldMessage))
		oldDigest := hex.EncodeToString(oldMac.Sum(nil))
		oldSignature := fmt.Sprintf("t=%d,v1=%s", oldTimestamp, oldDigest)

		valid := VerifyWebhookSignature(payload, oldSignature, secret)
		if valid {
			t.Error("expected invalid with expired timestamp (>5 min)")
		}
	})

	t.Run("malformed signature header", func(t *testing.T) {
		valid := VerifyWebhookSignature(payload, "not_a_valid_sig", secret)
		if valid {
			t.Error("expected invalid with malformed header")
		}
	})

	t.Run("instance method works same as standalone", func(t *testing.T) {
		client := NewClient("test")
		v1 := VerifyWebhookSignature(payload, signature, secret)
		v2 := client.Webhooks.VerifySignature(payload, signature, secret)
		if v1 != v2 {
			t.Error("expected instance method to match standalone function")
		}
	})
}

func TestParseWebhookEvent(t *testing.T) {
	body := []byte(`{"event":"payment.confirmed","timestamp":1693800000,"data":{"id":"pay_123","status":4}}`)

	event, err := ParseWebhookEvent(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.Event != "payment.confirmed" {
		t.Errorf("expected event %q, got %q", "payment.confirmed", event.Event)
	}
	if event.Timestamp != 1693800000 {
		t.Errorf("expected timestamp 1693800000, got %d", event.Timestamp)
	}
	if event.Data == nil {
		t.Error("expected non-nil data")
	}
}

func TestParseWebhookEventInvalid(t *testing.T) {
	_, err := ParseWebhookEvent([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
