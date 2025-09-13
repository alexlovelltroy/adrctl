#!/bin/bash
set -e

# adrctl installation script
# Usage: curl -sSL https://raw.githubusercontent.com/alexlovelltroy/adrctl/main/install.sh | bash

REPO="alexlovelltroy/adrctl"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $OS in
    darwin)
        OS="Darwin"
        ;;
    linux)
        OS="Linux"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

case $ARCH in
    x86_64)
        ARCH="x86_64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Get latest release version
echo "🔍 Finding latest release..."
VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)

if [ -z "$VERSION" ]; then
    echo "❌ Failed to get latest release version"
    exit 1
fi

echo "📦 Latest version: $VERSION"

# Construct download URL
FILENAME="adrctl_${VERSION#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

echo "⬇️  Downloading $FILENAME..."

# Download and install
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

if curl -fsSL "$URL" | tar -xz; then
    echo "✅ Downloaded and extracted successfully"
else
    echo "❌ Failed to download or extract $URL"
    exit 1
fi

# Install binary
if [ -w "$INSTALL_DIR" ]; then
    mv adrctl "$INSTALL_DIR/"
    echo "🎉 Installed adrctl to $INSTALL_DIR/adrctl"
else
    echo "🔐 Installing to $INSTALL_DIR (requires sudo)..."
    sudo mv adrctl "$INSTALL_DIR/"
    echo "🎉 Installed adrctl to $INSTALL_DIR/adrctl"
fi

# Cleanup
cd /
rm -rf "$TEMP_DIR"

# Verify installation
if command -v adrctl >/dev/null 2>&1; then
    echo "✨ Installation successful!"
    echo "📖 Version: $(adrctl --version)"
    echo "🚀 Try: adrctl --help"
else
    echo "⚠️  Installation completed but adrctl not found in PATH"
    echo "   You may need to add $INSTALL_DIR to your PATH"
fi