#!/bin/bash

# Kullanıcının işletim sistemini tespit et
OS="$(uname -s)"
ARCH="amd64"
BINARY_URL=""

# En son sürüm numarasını al
LATEST_VERSION=$(curl -s https://api.github.com/repos/Dert-Ops/Docme-Ag/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ "$OS" == "Linux" ]; then
    BINARY_URL="https://github.com/Dert-Ops/Docme-Ag/releases/download/$LATEST_VERSION/docm-linux-amd64"
elif [ "$OS" == "Darwin" ]; then
    BINARY_URL="https://github.com/Dert-Ops/Docme-Ag/releases/download/$LATEST_VERSION/docm-mac-amd64"
else
    echo "❌ Unsupported OS: $OS"
    exit 1
fi

echo "📥 Downloading docm $LATEST_VERSION for $OS..."
wget -O docm "$BINARY_URL"

# Binary'yi sistem dizinine taşı
echo "🚀 Installing to /usr/local/bin/..."
sudo mv docm /usr/local/bin/docm

# İzinleri ayarla
sudo chmod +x /usr/local/bin/docm

echo "✅ Installation complete! Now you can use 'docm cm' or 'docm vs'."
