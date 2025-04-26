# Makefile for AHR999 AWS Lambda deployment

# Constants
FUNCTION_NAME = ahr999
BUILD_OUTPUT = bootstrap
ZIP_FILE = myFunction.zip

# Default target
all: build zip deploy

# Build the Go binary for AWS Lambda
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -tags lambda.norpc -o $(BUILD_OUTPUT) .

# Create a zip archive for deployment
zip: build
	zip -j $(ZIP_FILE) $(BUILD_OUTPUT)

# Deploy to AWS Lambda
deploy: zip
	aws lambda update-function-code --function-name $(FUNCTION_NAME) --zip-file fileb://$(ZIP_FILE)

# Clean build artifacts
clean:
	rm -f $(BUILD_OUTPUT) $(ZIP_FILE)