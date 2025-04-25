package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Candle []float64 // [time, low, high, open, close, volume]

func FetchHistoricalCandles(start, end int64) ([]Candle, error) {
	url := fmt.Sprintf("https://api.exchange.coinbase.com/products/btc-usdt/candles?granularity=86400&start=%d&end=%d", start, end)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching candles failed: %w", err)
	}
	defer resp.Body.Close()

	var candles []Candle
	if err := json.NewDecoder(resp.Body).Decode(&candles); err != nil {
		return nil, fmt.Errorf("decoding candles failed: %w", err)
	}
	return candles, nil
}

func ExtractClosePrices(candles []Candle, count int) []float64 {
	result := make([]float64, count)
	for i := 0; i < count; i++ {
		result[i] = candles[i][4] // close
	}
	return result
}
