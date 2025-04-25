package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

// 历史 K线数据结构
type Candle []float64 // [time, low, high, open, close, volume]

// 最新市场价格结构
type MarketStats struct {
	Open     string `json:"open"`
	High     string `json:"high"`
	Low      string `json:"low"`
	Last     string `json:"last"`
	Volume   string `json:"volume"`
	Volume30 string `json:"volume_30day"`
}

func GetTimestampDaysAgo(days int) int64 {
	now := time.Now().UTC()
	daysAgo := now.AddDate(0, 0, -days)
	beginningOfDay := time.Date(daysAgo.Year(), daysAgo.Month(), daysAgo.Day(), 0, 0, 0, 0, time.UTC)
	return beginningOfDay.Unix()
}
func main() {
	start := GetTimestampDaysAgo(200)
	end := time.Now().UTC().Unix()
	// === Step 1: 获取历史数据（最近200天） ===
	historyURL := fmt.Sprintf("https://api.exchange.coinbase.com/products/btc-usdt/candles?granularity=86400&start=%d&end=%d", start, end)
	resp, err := http.Get(historyURL)
	if err != nil {
		log.Fatal("Error fetching historical data:", err)
	}
	defer resp.Body.Close()

	var candles []Candle
	if err := json.NewDecoder(resp.Body).Decode(&candles); err != nil {
		log.Fatal("Error decoding candles:", err)
	}

	if len(candles) < 200 {
		log.Fatal("Not enough data to compute 200-day geometric mean")
	}

	closePrices := make([]float64, 200)
	for i := 0; i < 200; i++ {
		closePrices[i] = candles[i][4] // ✅ 最近200天的close
	}
	geomean, _ := GeometricMean(closePrices)

	// === Step 2: 获取最新价格（使用low字段） ===
	statsURL := "https://api.exchange.coinbase.com/products/btc-usdt/stats"
	statsResp, err := http.Get(statsURL)
	if err != nil {
		log.Fatal("Error fetching market stats:", err)
	}
	defer statsResp.Body.Close()

	var stats MarketStats
	if err := json.NewDecoder(statsResp.Body).Decode(&stats); err != nil {
		log.Fatal("Error decoding stats:", err)
	}

	latestPrice, err := strconv.ParseFloat(stats.Last, 64)
	if err != nil {
		log.Fatal("Error parsing latest low price:", err)
	}

	// === Step 3: 计算指数估值 ===
	genesis := int64(1230940800)
	latestTime := int64(candles[0][0])
	days := float64(latestTime-genesis) / 86400

	if days <= 0 {
		log.Fatal("Invalid days calculation: days must be greater than 0")
	}

	estimatedValue := math.Pow(10, 5.84*math.Log10(days)-17.01)

	// === Step 4: 计算 AHR999 ===
	// 使用原Python逻辑: ahr999 = (latestPrice / geomean) * (latestPrice / estimatedValue)
	ahr999 := (latestPrice / geomean) * (latestPrice / estimatedValue)

	// 输出结果
	fmt.Printf("AHR999: %.4f\n", ahr999)
	fmt.Printf("Latest Price: %.2f\n", latestPrice)
	fmt.Printf("200-day cost: %.2f\n", geomean)
	println("Buy at the bottom: 0.45")
	println("Fixed Investment zone: 1.20")
}
