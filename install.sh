#!/bin/bash

# GoModoro Installation Script
# Run with: sudo ./install.sh

set -e  # Exit on any error

BINARY_NAME="gomodoro"
INSTALL_DIR="/usr/local/bin"
DESKTOP_DIR="/usr/share/applications"
ICON_DIR="/usr/share/pixmaps"

echo "🍅 Installing GoModoro Pomodoro Timer..."

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "❌ This script must be run as root (use sudo)" 
   exit 1
fi

# Check if binary exists
if [[ ! -f "./$BINARY_NAME" ]]; then
    echo "❌ GoModoro binary not found. Please build first with: go build ."
    exit 1
fi

# Install binary
echo "📦 Installing binary to $INSTALL_DIR..."
cp "./$BINARY_NAME" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Install desktop file
echo "🖥️ Installing desktop integration..."
if [[ -f "./GoModoro.desktop" ]]; then
    cp "./GoModoro.desktop" "$DESKTOP_DIR/"
    chmod 644 "$DESKTOP_DIR/GoModoro.desktop"
else
    echo "⚠️ Desktop file not found, creating basic one..."
    cat > "$DESKTOP_DIR/GoModoro.desktop" << EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=GoModoro
GenericName=Pomodoro Timer
Comment=A pirate's Pomodoro timer built with Go and Fyne
Exec=/usr/local/bin/gomodoro
Icon=gomodoro
Terminal=false
Categories=Utility;Office;ProjectManagement;
Keywords=pomodoro;timer;productivity;focus;work;break;
StartupNotify=true
EOF
    chmod 644 "$DESKTOP_DIR/GoModoro.desktop"
fi

# Install icon (if available)
if [[ -f "./icon.png" ]]; then
    echo "🎨 Installing icon..."
    cp "./icon.png" "$ICON_DIR/gomodoro.png"
    chmod 644 "$ICON_DIR/gomodoro.png"
else
    echo "⚠️ No icon.png found, skipping icon installation"
fi

# Update desktop database
if command -v update-desktop-database &> /dev/null; then
    echo "🔄 Updating desktop database..."
    update-desktop-database "$DESKTOP_DIR"
fi

# Update icon cache
if command -v gtk-update-icon-cache &> /dev/null; then
    echo "🔄 Updating icon cache..."
    gtk-update-icon-cache -f -t "$ICON_DIR" 2>/dev/null || true
fi

echo "✅ GoModoro installed successfully!"
echo ""
echo "🚀 You can now:"
echo "   • Launch from applications menu"
echo "   • Run 'gomodoro' from terminal"
echo "   • Right-click and select 'Open' on the binary"
echo ""
echo "🏴‍☠️ Happy productivity, matey!"
