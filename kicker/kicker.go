package kicker

import (
	"context"
	"kickerbot/db"
	"log"
	"strings"
	"time"

	tb "github.com/NicoNex/echotron/v3"
)

type bot struct {
	chatID       int64
	CaptchaTopic int64
	Me           *tb.User
	tb.API
}

func (b *bot) Update(update *tb.Update) {
	if update.Message != nil {
		if len(update.Message.NewChatMembers) != 0 {
			for _, user := range update.Message.NewChatMembers {
				if user.ID == b.Me.ID {
					botAdded(b, update)
				}
			}
			userJoined(b, update)
			return
		}
		if update.Message.LeftChatMember != nil {
			userLeft(b, update)
			return
		}
		if update.Message.Text != "" {
			if update.Message.Text == "/settopic" {
				setTopic(b, update)
				return
			}
			checkCaptcha(b, update)
		}
	}
}

// Базовая структура для бота
type Kicker struct {
	Token      string
	Dispatcher *tb.Dispatcher
}

func (b *Kicker) NewBot(chatID int64) tb.Bot {
	d := db.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	chat := db.Chat{
		Id:      chatID,
		Title:   "",
		TopicId: 0,
	}
	if !d.ChatExists(ctx, chat) {
		if err := d.NewChat(ctx, chat); err != nil {
			return &bot{}
		}
	}
	chat, _ = d.GetChat(ctx, chatID)
	CaptchaTopic := chat.TopicId
	result := &bot{
		chatID,
		CaptchaTopic,
		nil,
		tb.NewAPI(b.Token),
	}
	me, err := result.GetMe()
	if err != nil {
		log.Println(err)
	}
	result.Me = me.Result

	return result
}

// Initialize bot with token
func (b *Kicker) Init() error {
	dsp := tb.NewDispatcher(b.Token, b.NewBot)
	b.Dispatcher = dsp
	return nil
}

func (b *Kicker) Start() error {
	return b.Dispatcher.Poll()
}

func EscapeText(parseMode tb.ParseMode, text string) string {
	var replacer *strings.Replacer

	if parseMode == tb.HTML {
		replacer = strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;")
	} else if parseMode == tb.Markdown {
		replacer = strings.NewReplacer("_", "\\_", "*", "\\*", "`", "\\`", "[", "\\[")
	} else if parseMode == tb.MarkdownV2 {
		replacer = strings.NewReplacer(
			"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(",
			"\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">", "\\>",
			"#", "\\#", "+", "\\+", "-", "\\-", "=", "\\=", "|",
			"\\|", "{", "\\{", "}", "\\}", ".", "\\.", "!", "\\!",
		)
	} else {
		return ""
	}

	return replacer.Replace(text)
}
