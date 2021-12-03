package bot

import (
	"errors"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v3"
)

type Handler struct {
	Endpoint interface{}
	Handler  tb.HandlerFunc
}

type Bot struct {
	Bot   *tb.Bot
	Token string
}

func (b *Bot) Init() error {
	bot, err := tb.NewBot(tb.Settings{
		Token:  b.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Print(err)
		return err
	}
	b.Bot = bot
	return nil
}

// Add handler methods to the bot
func (b *Bot) AddHandlers(handlers []Handler) error {
	if len(handlers) != 0 {
		for i := range handlers {
			b.Bot.Handle(handlers[i].Endpoint, handlers[i].Handler)
		}
		return nil
	}
	return errors.New("no handlers are declared")
}
