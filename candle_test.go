package main

import (
	"testing"
	"time"
)

func TestFetchHistoricalCandles(t *testing.T) {
	start := GetTimestampDaysAgo(10)
	end := time.Now().UTC().Unix()

	candles, err := FetchHistoricalCandles(start, end)
	if err != nil {
		t.Fatalf("FetchHistoricalCandles failed: %v", err)
	}

	if len(candles) == 0 {
		t.Fatal("No candle data returned")
	}

	for i, c := range candles {
		if len(c) != 6 {
			t.Errorf("Candle at index %d does not contain 6 elements: %v", i, c)
		}
	}
}
