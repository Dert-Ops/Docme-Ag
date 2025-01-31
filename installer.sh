#!/bin/bash

# ENV Dosyası Konumları
USER_ENV="$HOME/.docm.env"
SYSTEM_ENV="/etc/docm.env"
DEFAULT_ENV_CONTENT="GEMINI_API_KEY=your-api-key-here"

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
sudo chmod +x /usr/local/bin/docm

# Kullanıcıya ENV dosyasının konumunu sormadan önce, var olup olmadığını kontrol et
if [ ! -f "$USER_ENV" ] && [ ! -f "$SYSTEM_ENV" ]; then
    echo "🔧 No existing .env file found. Creating a new one..."
    echo "$DEFAULT_ENV_CONTENT" > "$USER_ENV"
    chmod 600 "$USER_ENV"
    echo "✅ Please update your API key in $USER_ENV"
else
    echo "✅ Existing .env file found. No changes were made."
fi

echo "✅ Installation complete! Now you can use 'docm cm' or 'docm vs'."
