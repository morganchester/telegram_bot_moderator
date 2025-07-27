# Telegram Chat Moderator Bot

This is a simple Telegram bot written in Go that moderates messages in a specific topic thread of a supergroup. It automatically deletes messages that **do not contain a photo or a video of at least 30 seconds**. It also **forwards notifications** to another chat when a user leaves the group.

## Features

* Automatically deletes messages that:

    * do **not** include a photo;
    * and do **not** contain a video of 30 seconds or longer.
* Supports monitoring a **specific topic (thread)** in a supergroup.
* Supports automatic pin for last location sharing message in specific thread of supergroup
* Forwards a customizable message to another chat when a user leaves the group.
* Multilingual / customizable strings via a `strings.yaml` file.
* Configuration via `config.ini` (example provided)
* Service for systemd provided in `telegram.service`
* `install.sh` for fast compile
* `update.sh` for fast update with start/stop on remove server

---

## Configuration

### `config.ini`

```ini
[telegram]
token = YOUR_BOT_TOKEN
chat_id = -123456789          ; ID of the group to monitor
topic_id = 123                ; ID of the specific topic/thread
updates_chat_id = -987654321  ; ID of the chat where forward messages are sent
location_topic_id = 123  ; ID of the topic where you want to pin the latest location sharing message
```

### `strings.yaml`

All messages are configurable. Example:

```yaml
bot_started: "Bot started"
bot_stopped: "Bot stopped"
received_message_log: "Received message from @{username} (ID {msg_id}) in thread {thread_id}"
skipped_wrong_chat: "Skipped message from unknown chat ID: {chat_id}"
skipped_wrong_thread: "Skipped message from unrelated thread: {thread_id}"
has_photo: "Photo detected in message"
video_duration: "Video duration: {duration} seconds"
delete_reason: "Deleting message {msg_id} – no valid media"
delete_success: "Deleted message {msg_id} successfully"
delete_error: "Error deleting message {msg_id}: {err}"
media_ok: "Message {msg_id} contains valid media – skipping"
left_chat_message: "User {name} left the chat"
forward_sent_log: "Forwarded leave message to chat {chat_id}"
error_sending_forward: "Error forwarding message: {err}"
```

---

## Usage

1. Build the bot:

```bash
go build -o telegram-chat-moderator-bot
```

2. Run the bot:

```bash
./telegram-chat-moderator-bot
```

Make sure `config.ini` and `strings.yaml` are in the same directory as the executable.

---

## Dependencies

* [go-telegram-bot](https://github.com/go-telegram/bot)
* [gopkg.in/ini.v1](https://pkg.go.dev/gopkg.in/ini.v1)
* [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3)

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---