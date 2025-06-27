#!/bin/bash

# GoModoro Uninstallation Script
# Run with: sudo ./uninstall.sh

set -e  # Exit on any error

BINARY_NAME="gomodoro"
INSTALL_DIR="/usr/local/bin"
DESKTOP_DIR="/usr/share/applications"
ICON_DIR="/usr/share/pixmaps"

echo "🗑️ Uninstalling GoModoro Pomodoro Timer..."

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "❌ This script must be run as root (use sudo)" 
   exit 1
fi

# Remove binary
if [[ -f "$INSTALL_DIR/$BINARY_NAME" ]]; then
    echo "🗑️ Removing binary from $INSTALL_DIR..."
    rm -f "$INSTALL_DIR/$BINARY_NAME"
else
    echo "⚠️ Binary not found in $INSTALL_DIR"
fi

# Remove desktop file
if [[ -f "$DESKTOP_DIR/GoModoro.desktop" ]]; then
    echo "🗑️ Removing desktop integration..."
    rm -f "$DESKTOP_DIR/GoModoro.desktop"
else
    echo "⚠️ Desktop file not found"
fi

# Remove icon
if [[ -f "$ICON_DIR/gomodoro.png" ]]; then
    echo "🗑️ Removing icon..."
    rm -f "$ICON_DIR/gomodoro.png"
else
    echo "⚠️ Icon not found"
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

echo "✅ GoModoro uninstalled successfully!"
echo ""
echo "💾 Note: User settings and state files are preserved in ~/.config/gomodoro/"
echo "   Remove manually if desired: rm -rf ~/.config/gomodoro/"
echo ""
echo "🏴‍☠️ Farewell, matey!"
