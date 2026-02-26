#!/bin/bash
set -euo pipefail

# freshbox installer
# Usage: curl -fsSL https://raw.githubusercontent.com/kittors/freshbox/main/install.sh | bash

REPO="kittors/freshbox"
TMP_DIR=$(mktemp -d)

# Use /opt/homebrew/bin on Apple Silicon, /usr/local/bin on Intel
if [ "$(uname -m)" = "arm64" ]; then
    INSTALL_DIR="/opt/homebrew/bin"
else
    INSTALL_DIR="/usr/local/bin"
fi

# Ensure install dir exists
sudo mkdir -p "$INSTALL_DIR" 2>/dev/null || true

cleanup() { rm -rf "$TMP_DIR"; }
trap cleanup EXIT

CYAN='\033[0;36m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

info()  { echo -e "${CYAN}[freshbox]${NC} $1"; }
ok()    { echo -e "${GREEN}[freshbox]${NC} $1"; }
fail()  { echo -e "${RED}[freshbox]${NC} $1"; exit 1; }

echo ""
echo -e "${CYAN}  🍃 freshbox installer${NC}"
echo ""

# Detect OS and arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" != "darwin" ]; then
    fail "freshbox only supports macOS. Detected: $OS"
fi

case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    arm64)   ARCH="arm64" ;;
    aarch64) ARCH="arm64" ;;
    *)       fail "Unsupported architecture: $ARCH" ;;
esac

info "Detected: macOS $ARCH"

# Check if Go is available for building from source
if command -v go &>/dev/null; then
    info "Go found, building from source..."
    cd "$TMP_DIR"
    git clone --depth 1 --quiet "https://github.com/$REPO.git" freshbox
    cd freshbox
    go build -o freshbox ./cmd/freshbox

    if [ -w "$INSTALL_DIR" ]; then
        mv freshbox "$INSTALL_DIR/freshbox"
    else
        info "Need sudo to install to $INSTALL_DIR"
        sudo mv freshbox "$INSTALL_DIR/freshbox"
    fi

    ok "Installed freshbox to $INSTALL_DIR/freshbox"
else
    # Try to download pre-built binary from GitHub releases
    info "Go not found, trying pre-built binary..."

    LATEST=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" 2>/dev/null | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/' || echo "")

    if [ -z "$LATEST" ]; then
        # No releases yet, build from source with temp Go
        info "No pre-built binary available. Installing Go temporarily..."

        if command -v brew &>/dev/null; then
            brew install go
        else
            fail "Neither Go nor Homebrew found. Please install Go first: https://go.dev/dl/"
        fi

        cd "$TMP_DIR"
        git clone --depth 1 --quiet "https://github.com/$REPO.git" freshbox
        cd freshbox
        go build -o freshbox ./cmd/freshbox

        if [ -w "$INSTALL_DIR" ]; then
            mv freshbox "$INSTALL_DIR/freshbox"
        else
            sudo mv freshbox "$INSTALL_DIR/freshbox"
        fi

        ok "Installed freshbox to $INSTALL_DIR/freshbox"
    else
        DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST/freshbox-darwin-$ARCH"
        info "Downloading $LATEST for darwin/$ARCH..."

        curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/freshbox" || fail "Download failed. Try building from source with Go."
        chmod +x "$TMP_DIR/freshbox"

        if [ -w "$INSTALL_DIR" ]; then
            mv "$TMP_DIR/freshbox" "$INSTALL_DIR/freshbox"
        else
            sudo mv "$TMP_DIR/freshbox" "$INSTALL_DIR/freshbox"
        fi

        ok "Installed freshbox $LATEST to $INSTALL_DIR/freshbox"
    fi
fi

echo ""
ok "Run 'freshbox' to start setting up your Mac!"
echo ""
