package azpays

import (
	"context"
	"fmt"
)

// PriceService handles price quotes and candlestick data.
type PriceService struct {
	client *Client
}

// PriceQuote represents a real-time price quote.
type PriceQuote struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

// Candlestick represents a 24-hour candlestick price data point.
type Candlestick struct {
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume,omitempty"`
	Timestamp int64   `json:"timestamp"`
}

// GetQuote retrieves the current price for a symbol (e.g. "BTC", "ETH", "USDT").
func (s *PriceService) GetQuote(ctx context.Context, symbol string) (*PriceQuote, error) {
	var resp Response[PriceQuote]
	err := s.client.transport.get(ctx, fmt.Sprintf("/v1/prices/%s", symbol), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetCandlesticks retrieves 24h candlestick data for a symbol.
func (s *PriceService) GetCandlesticks(ctx context.Context, symbol string) ([]Candlestick, error) {
	var resp Response[[]Candlestick]
	err := s.client.transport.get(ctx, fmt.Sprintf("/v1/prices/%s/candles", symbol), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
