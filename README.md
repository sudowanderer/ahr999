# 📈 AHR999 - Bitcoin On-chain Valuation Lambda

This AWS Lambda function calculates the AHR999 indicator, a Bitcoin on-chain valuation metric. It fetches historical close prices from Coinbase, calculates a 200-day geometric mean, compares it to the estimated value model, and notifies results via Telegram and/or Bark.

---

## ✨ Features

- Fetches BTC/USDT daily close prices (200 days)
- Computes:
    - 📊 200-day Geometric Mean
    - 🧠 Estimated Value Model
    - 🔢 AHR999 Indicator
- Sends notification via:
    - Telegram Bot
    - Bark (optional)

---

## 🚀 Prerequisites

- Go 1.20+ installed
- AWS CLI installed and configured
- Telegram Bot Token and Chat ID
- An IAM Role for Lambda execution (`AWSLambdaBasicExecutionRole`)

---

## ⚙️ Environment Variables

| Name           | Description                  | Required |
|----------------|------------------------------|----------|
| `TELEGRAM_BOT_TOKEN`    | Telegram Bot token           | ✅       |
| `TELEGRAM_CHAT_ID`      | Telegram user/group chat ID  | ✅       |
| `BARK_URL`     | Bark push URL (optional)     | ❌       |

---

## 🛠️ Build & Deploy

1. Build the binary for AWS Lambda

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
go build -tags lambda.norpc -ldflags="-s -w" -o bootstrap .
```

2. (Optional) Compress with upx to reduce size
```bash
upx --best bootstrap
```

3. Package into zip

```bash
zip myFunction.zip bootstrap
```

4. Create the Lambda function
```bash
aws lambda create-function --function-name ahr999 \
  --runtime provided.al2023 \
  --handler bootstrap \
  --architectures arm64 \
  --role arn:aws:iam::<your-account-id>:role/LambdaGeneralPurposeRole \
  --zip-file fileb://myFunction.zip \
  --timeout 10 \
  --environment Variables="{TELEGRAM_BOT_TOKEN=your_TELEGRAM_BOT_TOKEN,TELEGRAM_CHAT_ID=your_TELEGRAM_CHAT_ID,BARK_URL=optional_bark_url}"
```
Replace <your-account-id> and environment values accordingly.

## 📩 Example Notification
```shell
📈 AHR999 Indicator Report
AHR999: 0.8734
Latest Price: 26251.32
200-day Cost (Geomean): 24798.12
估值模型: 30500.91

建议区间:
✅ 抄底区: < 0.45
💰 定投区: < 1.20
```