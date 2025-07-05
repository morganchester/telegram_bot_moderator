package main

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gopkg.in/ini.v1"
)

type Config struct {
	Token         string
	ChatID        int64
	ForwardChatID int64
	TopicID       int
}

func loadConfig(path string) (*Config, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}
	s := cfg.Section("telegram")
	return &Config{
		Token:         s.Key("token").String(),
		ChatID:        s.Key("chat_id").MustInt64(),
		TopicID:       s.Key("topic_id").MustInt(),
		ForwardChatID: s.Key("updates_chat_id").MustInt64(),
	}, nil
}

type Strings map[string]string

func loadStrings(path string) (Strings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s Strings
	err = yaml.Unmarshal(data, &s)
	return s, err
}

func format(msg string, vars map[string]string) string {
	for k, v := range vars {
		msg = strings.ReplaceAll(msg, "{"+k+"}", v)
	}
	return msg
}

func main() {
	cfg, err := loadConfig("config.ini")
	if err != nil {
		log.Fatal(err)
	}

	stringsMap, err := loadStrings("strings.yaml")
	if err != nil {
		log.Fatal(err)
	}

	b, err := bot.New(cfg.Token,
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			m := update.Message
			if m == nil {
				return
			}

			// User left the group
			if m.LeftChatMember != nil {
				user := m.LeftChatMember
				name := user.FirstName
				if user.LastName != "" {
					name += " " + user.LastName
				}
				if user.Username != "" {
					name += " (@" + user.Username + ")"
				}

				msg := format(stringsMap["left_chat_message"], map[string]string{
					"name": name,
				})

				// Forward message to forward chat
				_, err = b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: cfg.ForwardChatID,
					Text:   msg,
				})
				if err != nil {
					log.Printf(format(stringsMap["error_sending_forward"], map[string]string{
						"err": err.Error(),
					}))
				} else {
					log.Printf(format(stringsMap["forward_sent_log"], map[string]string{
						"chat_id": fmt.Sprint(cfg.ForwardChatID),
					}))
				}

				return
			}

			log.Printf(format(stringsMap["received_message_log"], map[string]string{
				"username":  m.From.Username,
				"msg_id":    fmt.Sprint(m.ID),
				"thread_id": fmt.Sprint(m.MessageThreadID),
			}))

			if m.Chat.ID != cfg.ChatID {
				log.Printf(format(stringsMap["skipped_wrong_chat"], map[string]string{
					"chat_id": fmt.Sprint(m.Chat.ID),
				}))
				return
			}
			if m.MessageThreadID != cfg.TopicID {
				log.Printf(format(stringsMap["skipped_wrong_thread"], map[string]string{
					"thread_id": fmt.Sprint(m.MessageThreadID),
				}))
				return
			}

			// Check the photo or the long video
			hasPhoto := m.Photo != nil
			hasLongVideo := m.Video != nil && m.Video.Duration >= 30
			if hasPhoto {
				log.Printf(stringsMap["has_photo"])
			}
			if m.Video != nil {
				log.Printf(format(stringsMap["video_duration"], map[string]string{
					"duration": fmt.Sprint(m.Video.Duration),
				}))
			}

			if !hasPhoto && !hasLongVideo {
				log.Printf(format(stringsMap["delete_reason"], map[string]string{
					"msg_id": fmt.Sprint(m.ID),
				}))
				ctx2 := context.Background()
				go func(chatID int64, msgID int) {
					time.Sleep(1 * time.Second)
					_, err := b.DeleteMessage(ctx2, &bot.DeleteMessageParams{
						ChatID:    chatID,
						MessageID: msgID,
					})
					if err != nil {
						log.Printf(format(stringsMap["delete_error"], map[string]string{
							"msg_id": fmt.Sprint(msgID),
							"err":    err.Error(),
						}))
					} else {
						log.Printf(format(stringsMap["delete_success"], map[string]string{
							"msg_id": fmt.Sprint(msgID),
						}))
					}
				}(m.Chat.ID, m.ID)
			} else {
				log.Printf(format(stringsMap["media_ok"], map[string]string{
					"msg_id": fmt.Sprint(m.ID),
				}))
			}

		}))

	if err != nil {
		log.Fatal(err)
	}

	log.Println(stringsMap["bot_started"])
	b.Start(context.Background())
	log.Println(stringsMap["bot_stopped"])
}
