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
		JoinedMessage: message.ID,
		CorrectAnswer: int8(captcha.CorrectAnswer),
		DateJoined:    time.Now().Unix(),
	}
	user.CorrectAnswer = int8(captcha.CorrectAnswer)
	d := db.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	log.Print(user)
	msg := fmt.Sprintf("Приветствую тебя, %v!\nДля подтверждения, что ты человек, выбери логотип движка, которому посвящен данный чат, и отправь его номер сюда.\nЯ дам тебе две минуты на это.", user.FirstName)
	photo := tb.Photo{File: tb.FromReader(reader), Caption: msg}
	result, err := bot.Send(tb.ChatID(message.Chat.ID), &photo, &tb.SendOptions{ReplyTo: message})
	if err != nil {
		return err
	}
	user.CaptchaMessage = result.ID

	d.NewUser(ctx, user)
	return nil
}

func userLeft(c tb.Context) error {
	bot := c.Bot()
	message := c.Message()
	sender := c.Sender()
	d := db.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if user, err := d.GetUser(ctx, db.User{Id: sender.ID, ChatId: message.Chat.ID}); err == nil {
		d.RemoveUser(ctx, user)
		bot.Delete(&tb.Message{Chat: message.Chat, ID: user.CaptchaMessage})
		bot.Delete(&tb.Message{Chat: message.Chat, ID: user.JoinedMessage})
	}
	return nil
}

func checkCaptcha(c tb.Context) error {
	sender := c.Sender()
	message := c.Message()
	bot := c.Bot()
	d := db.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if user, err := d.GetUser(ctx, db.User{Id: sender.ID, ChatId: message.Chat.ID}); err == nil {
		text_runes := []rune(message.Text)
		guess := string(text_runes[0])
		solved := false
		if num, err := strconv.Atoi(guess); err == nil {
			if num == int(user.CorrectAnswer) {
				_ = d.RemoveUser(ctx, user)
				solved = true
				bot.Delete(message)
				bot.Delete(&tb.Message{Chat: message.Chat, ID: user.CaptchaMessage})
			}
		} else {
			log.Print(err)
		}
		if !solved {
			bot.Delete(message)
			bot.Delete(&tb.Message{Chat: message.Chat, ID: user.CaptchaMessage})
			bot.Delete(&tb.Message{Chat: message.Chat, ID: user.JoinedMessage})
			bot.Ban(message.Chat, &tb.ChatMember{User: sender})
			_ = d.RemoveUser(ctx, user)
		}
	}
	return nil
}

func botAdded(c tb.Context) error {
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
}

var HandlersV1 = []Handler{
	{
		Endpoint: tb.OnText,
		Handler:  checkCaptcha,
	},
	{
		Endpoint: tb.OnAddedToGroup,
		Handler:  botAdded,
	},

	{
		Endpoint: tb.OnUserJoined,
		Handler:  userJoined,
	},
	{
		Endpoint: tb.OnUserLeft,
		Handler:  userLeft,
	},
}
