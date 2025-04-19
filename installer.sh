#!/bin/bash

# ENV Dosyası Konumları
USER_ENV="$HOME/.docm.env"
DEFAULT_ENV_CONTENT="GEMINI_API_KEY=AIzaSyDKDg2dRq3-AJTZR6_bfIP4dxAkrrX31CI"

# Kullanıcının işletim sistemini tespit et
OS="$(uname -s)"
ARCH="amd64"
BINARY_URL=""

LATEST_VERSION=$(curl -s https://api.github.com/repos/Dert-Ops/Docme-Ag/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ "$OS" = "Linux" ]; then
    BINARY_URL="https://github.com/Dert-Ops/Docme-Ag/releases/download/$LATEST_VERSION/docm-linux-amd64"
elif [ "$OS" = "Darwin" ]; then
    BINARY_URL="https://github.com/Dert-Ops/Docme-Ag/releases/download/$LATEST_VERSION/docm-mac-amd64"
else
    echo "❌ Unsupported OS: $OS"
    exit 1
fi

echo "📥 Downloading docm $LATEST_VERSION for $OS..."
wget -O docm $BINARY_URL

# ~/.local/bin dizinini oluştur
BIN_DIR="$HOME/.local/bin"
mkdir -p "$BIN_DIR"

# Binary'yi kullanıcı dizinine taşı
echo "🚀 Installing to $BIN_DIR/..."
mv docm "$BIN_DIR/docm"
chmod +x "$BIN_DIR/docm"

# Kullanıcıya PATH'i güncellemesi gerektiğini hatırlat
if [ ":$PATH:" != *":$BIN_DIR:"* ]; then
    echo "🔧 Adding $BIN_DIR to your PATH. This change will be effective after restarting your terminal."
    # PATH'i .bashrc veya .zshrc'ye ekle
    if [ -n "$ZSH_VERSION" ]; then
        echo "export PATH=\$PATH:$BIN_DIR" >> "$HOME/.zshrc"
    else
        echo "export PATH=\$PATH:$BIN_DIR" >> "$HOME/.bashrc"
    fi
else
    echo "✅ $BIN_DIR is already in your PATH."
fi

# Kullanıcıya ENV dosyasının konumunu sormadan önce, var olup olmadığını kontrol et
if [ ! -f "$USER_ENV" ]; then
    echo "🔧 No existing .env file found. Creating a new one..."
    echo "$DEFAULT_ENV_CONTENT" > "$USER_ENV"
    chmod 600 "$USER_ENV"
    echo "✅ Please update your API key in $USER_ENV"
else
    echo "✅ Existing .env file found. No changes were made."
fi

echo "✅ Installation complete! Now you can use 'docm cm' or 'docm vs' after restarting your terminal."
