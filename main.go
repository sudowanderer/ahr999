package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sudowanderer/notikit/notifier"
	"os"
	"time"
)

func handleRequest(ctx context.Context, event json.RawMessage) error {
	// ✅ 检查必备环境变量
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatID == "" {
		return fmt.Errorf("missing required environment variables: TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID")
	}

	// Optional timezone for Telegram messages
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tg := notifier.NewTelegramNotifierWithLocation(botToken, chatID, loc)

	// === 核心逻辑 ===
	const days = 200
	start := GetTimestampDaysAgo(days)
	end := time.Now().UTC().Unix()

	candles, err := FetchHistoricalCandles(start, end)
	if err != nil {
		return fmt.Errorf("failed to fetch candles: %w", err)
	}
	if len(candles) < days {
		return fmt.Errorf("not enough candle data")
	}

	closePrices := ExtractClosePrices(candles, days)
	geomean, err := GeometricMean(closePrices)
	if err != nil {
		return fmt.Errorf("failed to compute geomean: %w", err)
	}

	stats, err := FetchMarketStats()
	if err != nil {
		return fmt.Errorf("failed to fetch market stats: %w", err)
	}
	latestPrice, err := stats.ParseLastPrice()
	if err != nil {
		return fmt.Errorf("invalid latest price: %w", err)
	}

	latestTime := int64(candles[0][0])
	estimatedValue := ComputeEstimatedValue(latestTime)
	ahr999 := ComputeAHR999(latestPrice, geomean, estimatedValue)

	// ✅ 构造并发送格式化消息
	msg := fmt.Sprintf(`📈 AHR999 Indicator Report
AHR999: %.4f
Latest Price: %.2f
200-day Cost (Geomean): %.2f
建议区间:
✅ 抄底区: < 0.45
💰 定投区: < 1.20`, ahr999, latestPrice, geomean)

	if err := tg.Notify(msg); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}

func main() {
	//for local test
	//err := handleRequest(context.Background(), nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	lambda.Start(handleRequest)
}
