package kicker

import (
	"context"
	"fmt"
	"kickerbot/db"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v3"
)

func userJoined(c tb.Context) error {
	m := c.Message()
	user := db.User{
		Id:            m.Sender.ID,
		Username:      m.Sender.Username,
		FirstName:     m.Sender.FirstName,
		LastName:      m.Sender.LastName,
		IsBanned:      false,
		ChatId:        m.Chat.ID,
		CorrectAnswer: 0,
	}
	log.Print(user)
	str := fmt.Sprintf("%v", user)
	c.Bot().Send(&tb.User{ID: 60441930}, str)
	db.Log(str)
	return nil
}

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
			db.Log("message")
			return nil
		},
	},
	{
		Endpoint: tb.OnUserJoined,
		Handler:  userJoined,
	},
}
