#!/bin/bash
# =========================================================================
# Deploy Script for Hetzner VPS
# =========================================================================
# Usage: ./deploy/deploy.sh
#
# Prerequisites:
#   1. Set DEPLOY_HOST in .env (e.g., root@your-server-ip)
#   2. Set DEPLOY_PATH in .env (e.g., /opt/myapp)
#   3. Ensure SSH key is configured for passwordless access
# =========================================================================

set -e

# Load environment variables
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Configuration
DEPLOY_HOST="${DEPLOY_HOST:-}"
DEPLOY_PATH="${DEPLOY_PATH:-/opt/myapp}"
BINARY_NAME="server"
SERVICE_NAME="myapp"

# Validate configuration
if [ -z "$DEPLOY_HOST" ]; then
    echo "❌ Error: DEPLOY_HOST not set in .env"
    echo "   Add: DEPLOY_HOST=root@your-server-ip"
    exit 1
fi

echo "🚀 Starting deployment to $DEPLOY_HOST..."

# Step 1: Build production binary for Linux
echo "📦 Building production binary..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/${BINARY_NAME}-linux ./cmd/server

# Step 2: Create remote directory if needed
echo "📁 Ensuring remote directory exists..."
ssh "$DEPLOY_HOST" "mkdir -p $DEPLOY_PATH"

# Step 3: Upload binary
echo "📤 Uploading binary..."
scp ./bin/${BINARY_NAME}-linux "$DEPLOY_HOST:$DEPLOY_PATH/$BINARY_NAME"

# Step 4: Upload assets
echo "📤 Uploading assets..."
scp -r ./assets "$DEPLOY_HOST:$DEPLOY_PATH/"

# Step 5: Set permissions
echo "🔐 Setting permissions..."
ssh "$DEPLOY_HOST" "chmod +x $DEPLOY_PATH/$BINARY_NAME"

# Step 6: Restart service
echo "🔄 Restarting service..."
ssh "$DEPLOY_HOST" "sudo systemctl restart $SERVICE_NAME"

# Step 7: Verify deployment
echo "✅ Verifying deployment..."
sleep 2
ssh "$DEPLOY_HOST" "sudo systemctl status $SERVICE_NAME --no-pager"

echo ""
echo "🎉 Deployment complete!"
echo "   Server: $DEPLOY_HOST"
echo "   Path: $DEPLOY_PATH"
