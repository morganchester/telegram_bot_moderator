#!/bin/bash

# Installation script for Telegram Bot Chat Deleter
set -e

echo "Installing Telegram Bot Chat Deleter..."

# Create directory
INSTALL_DIR="/opt/scripts/telegramBotChatDeleter"
sudo mkdir -p "$INSTALL_DIR"

# Copy files
echo "Copying files to $INSTALL_DIR..."
sudo cp main.go go.mod "$INSTALL_DIR/"
sudo cp config.ini "$INSTALL_DIR/"

# Set ownership
sudo chown -R root:root "$INSTALL_DIR"

# Navigate to install directory
cd "$INSTALL_DIR"

# Install Go dependencies and build
echo "Building application..."
sudo go mod tidy
sudo go build -o telegram-bot-chat-deleter main.go

# Make executable
sudo chmod +x telegram-bot-chat-deleter

# Install systemd service
echo "Installing systemd service..."
sudo cp telegram.service /etc/systemd/system/telegram-bot-chat-deleter.service

# Reload systemd
sudo systemctl daemon-reload

echo "Installation complete!"
echo ""
echo "Next steps:"
echo "1. Edit the configuration file: sudo nano $INSTALL_DIR/config.ini"
echo "2. Add your bot token"
echo "   Add your forward group chat_id"
echo "3. For PRIVATE groups: Use chat_id (recommended)"
echo "   For PUBLIC channels: Use chat_name"
echo "4. Enable the service: sudo systemctl enable telegram-bot-chat-deleter"
echo "5. Start the service: sudo systemctl start telegram-bot-chat-deleter"
echo "6. Check status: sudo systemctl status telegram-bot-chat-deleter"
echo "7. View logs: sudo journalctl -u telegram-bot-chat-deleter -f"
echo ""
echo "To create a Telegram bot:"
echo "1. Message @BotFather on Telegram"
echo "2. Use /newbot command"
echo "3. Follow the instructions to get your bot token"
echo "4. Add your bot as an administrator to your channel/group"
echo "5. Give it permission to delete messages"
echo ""
echo "To get your private group's chat ID, before starting the bot:"
echo "1. Add your new bot to your group temporarily"
echo "2. Send any message in the group"
echo "3. Go to https://api.telegram.org/bot<TOKEN>/getUpdates"
echo "4. The message details will have the necessary chat.id"
echo "5. Use the negative chat ID in your config (e.g., -1001234567890)"