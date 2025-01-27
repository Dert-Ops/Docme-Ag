#!/bin/bash

# Kurulum başlıyor
echo "🚀 Setting up AI Dev Agent..."

# Gerekli paketleri yükle
echo "🔍 Checking for Go installation..."
if ! command -v go &> /dev/null
then
    echo "⚠ Go is not installed. Please install Go and rerun this script."
    exit 1
fi

# Projeyi klonla ve içeri gir
if [ ! -d "docme-ag" ]; then
    echo "📥 Cloning repository..."
    git clone https://github.com/dert-ops/docme-ag.git
fi
cd docme-ag || exit

# Go bağımlılıklarını yükle
echo "📦 Installing dependencies..."
go mod tidy

# Go binary derle
echo "🔨 Building project..."
go build -o docm main.go

# Alias'ı kullanıcının shell profil dosyasına ekle
if [[ "$SHELL" == *"zsh"* ]]; then
    PROFILE_FILE="$HOME/.zshrc"
elif [[ "$SHELL" == *"bash"* ]]; then
    PROFILE_FILE="$HOME/.bashrc"
else
    PROFILE_FILE="$HOME/.profile"
fi

echo "🔧 Adding alias to $PROFILE_FILE..."
echo 'alias docm="$(pwd)/docm"' >> "$PROFILE_FILE"

# Değişiklikleri yükle
echo "🔄 Reloading shell configuration..."
source "$PROFILE_FILE"

echo "✅ Installation complete! You can now use 'docm cm' and 'docm vs' commands."
