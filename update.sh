#!/bin/bash

systemctl stop telegram

# Install Go dependencies and build
echo "Building application..."
sudo go mod tidy
sudo go build -o telegram-bot-chat-deleter main.go


# Make executable
sudo chmod +x telegram-bot-chat-deleter

# Reload systemd
sudo systemctl daemon-reload

systemctl start telegram