package kicker

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"kickerbot/captchagen"
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
	d := db.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	d.NewUser(ctx, user)
	log.Print(user)
	str := fmt.Sprintf("%v", user)
	c.Bot().Send(&tb.User{ID: 60441930}, str)
	msg := fmt.Sprintf("Приветствую, %v!\nПеред тем, как дать тебе что-то здесь писать, я задам тебе один вопрос:\nКакой из этих движков самый лучший? Подумай хорошенько, и дай ответ цифрой.", user.FirstName)
	c.Reply(msg)
	db.Log("new user", str)
	return nil
}

var HandlersV1 = []Handler{
	{
		Endpoint: tb.OnText,
		Handler: func(c tb.Context) error {
			db.Log("message", c.Message())
			return nil
		},
	},
	{
		Endpoint: "/gen",
		Handler: func(c tb.Context) error {
			captcha := captchagen.GenCaptcha()
			buff := new(bytes.Buffer)
			err := png.Encode(buff, captcha.Image)
			if err != nil {
				fmt.Println("failed to create buffer", err)
			}
			reader := bytes.NewReader(buff.Bytes())
			// log.Print(reader)
			caption := fmt.Sprintf("Правильный ответ: %d", captcha.CorrectAnswer)
			c.Reply(&tb.Photo{File: tb.FromReader(reader), Caption: caption})
			return nil
		},
	},
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
			db.Log("new chat", chat)
			return nil
		},
	},

	{
		Endpoint: tb.OnUserJoined,
		Handler:  userJoined,
	},
}
