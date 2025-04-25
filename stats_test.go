package main

import (
	"testing"
)

func TestFetchMarketStats(t *testing.T) {
	stats, err := FetchMarketStats()
	if err != nil {
		t.Fatalf("FetchMarketStats failed: %v", err)
	}

	last, err := stats.ParseLastPrice()
	if err != nil {
		t.Fatalf("ParseLastPrice failed: %v", err)
	}

	if last <= 0 {
		t.Errorf("Invalid last price: %f", last)
	}
}
