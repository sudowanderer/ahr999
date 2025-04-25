package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type MarketStats struct {
	Open     string `json:"open"`
	High     string `json:"high"`
	Low      string `json:"low"`
	Last     string `json:"last"`
	Volume   string `json:"volume"`
	Volume30 string `json:"volume_30day"`
}

func FetchMarketStats() (MarketStats, error) {
	url := "https://api.exchange.coinbase.com/products/btc-usdt/stats"
	resp, err := http.Get(url)
	if err != nil {
		return MarketStats{}, fmt.Errorf("fetching stats failed: %w", err)
	}
	defer resp.Body.Close()

	var stats MarketStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return MarketStats{}, fmt.Errorf("decoding stats failed: %w", err)
	}
	return stats, nil
}

func (s MarketStats) ParseLastPrice() (float64, error) {
	return strconv.ParseFloat(s.Last, 64)
}
