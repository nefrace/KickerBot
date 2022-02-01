package kicker

import (
	"context"
	"fmt"
	"kickerbot/captchagen"
	"kickerbot/db"
	"log"
	"strconv"
	"time"

	tb "gopkg.in/tucnak/telebot.v3"
)

func userJoined(c tb.Context) error {
	bot := c.Bot()
	captcha := captchagen.GenCaptcha()
	reader := captcha.ToReader()
	message := c.Message()
	user := db.User{
		Id:            message.Sender.ID,
		Username:      message.Sender.Username,
		FirstName:     message.Sender.FirstName,
		LastName:      message.Sender.LastName,
		IsBanned:      false,
		ChatId:        message.Chat.ID,
		CorrectAnswer: int8(captcha.CorrectAnswer),
		DateJoined:    time.Now().Unix(),
	}
	user.CorrectAnswer = int8(captcha.CorrectAnswer)
	d := db.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	log.Print(user)
	str := fmt.Sprintf("%v", user)
	msg := fmt.Sprintf("Приветствую, %v!\nПеред тем, как дать тебе что-то здесь писать, я задам тебе один вопрос:\nКакой из этих движков самый лучший? Подумай хорошенько, и дай ответ цифрой.", user.FirstName)
	photo := tb.Photo{File: tb.FromReader(reader), Caption: msg}
	result, err := bot.Send(tb.ChatID(message.Chat.ID), &photo, &tb.SendOptions{ReplyTo: message})
	if err != nil {
		return err
	}
	user.CaptchaMessage = result.ID
	db.Log("new user", str)

	d.NewUser(ctx, user)
	return nil
}

var HandlersV1 = []Handler{
	{
		Endpoint: tb.OnText,
		Handler: func(c tb.Context) error {
			sender := c.Sender()
			message := c.Message()
			d := db.GetDatabase()
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			if user, err := d.GetUser(ctx, db.User{Id: sender.ID, ChatId: message.Chat.ID, IsBanned: false}); err == nil {
				text := message.Text
				if num, err := strconv.Atoi(text); err == nil {
					if num == int(user.CorrectAnswer) {
						_ = d.RemoveUser(ctx, user)
						c.Reply("Капча пройдена!")
						c.Bot().Delete(&tb.Message{Chat: message.Chat, ID: user.CaptchaMessage})
					} else {
						c.Reply("Ещё разочек")
					}
				} else {
					log.Print(err)
				}
			} else {
				c.Bot().Ban(message.Chat, &tb.ChatMember{User: sender})
			}
			return nil
		},
	},
	{
		Endpoint: "/gen",
		Handler: func(c tb.Context) error {
			captcha := captchagen.GenCaptcha()
			reader := captcha.ToReader()
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
