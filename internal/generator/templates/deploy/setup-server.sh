#!/bin/bash
# =========================================================================
# Server Setup Script for Hetzner VPS
# =========================================================================
# Usage: ./deploy/setup-server.sh
#
# This script sets up a fresh Hetzner VPS with:
#   - Caddy web server
#   - Systemd service for the Go application
#   - Firewall rules
# =========================================================================

set -e

# Load environment variables
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Configuration
DEPLOY_HOST="${DEPLOY_HOST:-}"
DEPLOY_PATH="${DEPLOY_PATH:-/opt/myapp}"
SERVICE_NAME="myapp"

# Validate configuration
if [ -z "$DEPLOY_HOST" ]; then
    echo "❌ Error: DEPLOY_HOST not set in .env"
    echo "   Add: DEPLOY_HOST=root@your-server-ip"
    exit 1
fi

echo "🔧 Setting up Hetzner VPS at $DEPLOY_HOST..."

# Step 1: Update system
echo "📦 Updating system packages..."
ssh "$DEPLOY_HOST" "apt update && apt upgrade -y"

# Step 2: Install Caddy
echo "🌐 Installing Caddy..."
ssh "$DEPLOY_HOST" << 'EOF'
apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
apt update
apt install -y caddy
EOF

# Step 3: Create application directory
echo "📁 Creating application directory..."
ssh "$DEPLOY_HOST" "mkdir -p $DEPLOY_PATH && chown www-data:www-data $DEPLOY_PATH"

# Step 4: Upload systemd service
echo "⚙️ Setting up systemd service..."
scp ./deploy/app.service "$DEPLOY_HOST:/etc/systemd/system/${SERVICE_NAME}.service"
ssh "$DEPLOY_HOST" "systemctl daemon-reload && systemctl enable $SERVICE_NAME"

# Step 5: Upload Caddyfile
echo "🔒 Configuring Caddy..."
scp ./deploy/Caddyfile "$DEPLOY_HOST:/etc/caddy/Caddyfile"

# Step 6: Setup firewall
echo "🛡️ Configuring firewall..."
ssh "$DEPLOY_HOST" << 'EOF'
apt install -y ufw
ufw allow ssh
ufw allow http
ufw allow https
ufw --force enable
EOF

# Step 7: Create log directory
echo "📋 Setting up logging..."
ssh "$DEPLOY_HOST" "mkdir -p /var/log/caddy && chown caddy:caddy /var/log/caddy"

# Step 8: Start services
echo "🚀 Starting services..."
ssh "$DEPLOY_HOST" "systemctl restart caddy"

echo ""
echo "✅ Server setup complete!"
echo ""
echo "📝 Next steps:"
echo "   1. Edit deploy/Caddyfile and replace YOUR_DOMAIN with your domain"
echo "   2. Upload .env to $DEPLOY_PATH/.env on the server"
echo "   3. Point your domain's DNS to this server"
echo "   4. Run: make deploy"
