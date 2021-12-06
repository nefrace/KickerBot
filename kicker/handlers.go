package kicker

import (
	"context"
	"kickerbot/db"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v3"
)

var HandlersV1 = []Handler{
	// {
	// 	Endpoint: tb.OnText,
	// 	Handler: func(c tb.Context) error {
	// 		m := c.Message()
	// 		c.Bot().Send(m.Sender, m.Text)
	// 		return nil
	// 	},
	// },
	{
		Endpoint: tb.OnAddedToGroup,
		Handler: func(c tb.Context) error {
			m := c.Message()
			chat := db.Chat{
				Id:    m.Chat.ID,
				Title: m.Chat.Title,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			d := db.GetDatabase()
			err := d.NewChat(ctx, chat)
			if err != nil {
				log.Print(err)
			}
			return nil
		},
	},
	{
		Endpoint: tb.OnText,
		Handler: func(c tb.Context) error {
			c.Bot().Send(&tb.User{ID: 60441930}, c.Message().Text)
			return nil
		},
	},
}
