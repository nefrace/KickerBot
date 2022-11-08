package kicker

import (
	"context"
	"fmt"
	"kickerbot/captchagen"
	"kickerbot/db"
	"log"
	"strconv"
	"time"

	tb "github.com/NicoNex/echotron/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func userJoined(b *bot, update *tb.Update) error {
	captcha := captchagen.GenCaptcha()
	bytes, err := captcha.ToBytes()
	if err != nil {
		fmt.Printf("Error creating captcha bytes: %v", bytes)
		b.SendMessage("Не могу создать капчу, @nefrace, проверь логи.", update.Message.From.ID, &tb.MessageOptions{MessageThreadID: update.Message.ThreadID})
	}
	message := update.Message
	user := db.User{
		Id:            message.From.ID,
		Username:      message.From.Username,
		FirstName:     message.From.FirstName,
		LastName:      message.From.LastName,
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
	msg := fmt.Sprintf("Приветствую тебя, *[%s](tg://user?id=%d)*\\!\nДля подтверждения, что ты человек, выбери логотип движка, которому посвящен данный чат, и отправь его номер сюда\\.\n*_Я дам тебе две минуты на это\\._*", EscapeText(tb.MarkdownV2, user.FirstName), user.Id)
	options := tb.PhotoOptions{
		Caption:   msg,
		ParseMode: tb.MarkdownV2,
	}
	if message.Chat.IsForum {
		options.MessageThreadID = int(b.CaptchaTopic)
	}
	result, err := b.SendPhoto(tb.NewInputFileBytes("logos.png", *bytes), message.Chat.ID, &options)
	if err != nil {
		return err
	}
	user.CaptchaMessage = result.Result.ID

	d.NewUser(ctx, user)
	return nil
}

func userLeft(b *bot, update *tb.Update) error {
	message := update.Message
	sender := message.From
	d := db.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if user, err := d.GetUser(ctx, db.User{Id: sender.ID, ChatId: message.Chat.ID}); err == nil {
		d.RemoveUser(ctx, user)
		b.DeleteMessage(message.Chat.ID, message.ID)
		b.DeleteMessage(message.Chat.ID, user.CaptchaMessage)
	}
	return nil
}

func checkCaptcha(b *bot, update *tb.Update) error {
	message := update.Message
	sender := message.From
	d := db.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if user, err := d.GetUser(ctx, db.User{Id: sender.ID, ChatId: message.Chat.ID}); err == nil {
		if message.Chat.IsForum {
			chat, err := d.GetChat(ctx, message.Chat.ID)
			if err != nil {
				return err
			}
			if message.ThreadID != int(chat.TopicId) {
				b.DeleteMessage(message.Chat.ID, message.ID)
				return nil
			}
		}
		text_runes := []rune(message.Text)
		guess := string(text_runes[0])
		solved := false
		if num, err := strconv.Atoi(guess); err == nil {
			if num == int(user.CorrectAnswer) {
				_ = d.RemoveUser(ctx, user)
				solved = true
				b.DeleteMessage(message.Chat.ID, message.ID)
				b.DeleteMessage(message.Chat.ID, user.CaptchaMessage)
			}
		} else {
			log.Println(err)
			return err
		}
		if !solved {
			b.DeleteMessage(message.Chat.ID, message.ID)
			b.DeleteMessage(message.Chat.ID, user.CaptchaMessage)
			b.DeleteMessage(message.Chat.ID, user.JoinedMessage)
			b.BanChatMember(message.Chat.ID, sender.ID, nil)
			_ = d.RemoveUser(ctx, user)
		}
	}
	return nil
}

func botAdded(b *bot, update *tb.Update) error {
	m := update.Message
	chat := db.Chat{
		Id:      m.Chat.ID,
		Title:   m.Chat.Title,
		TopicId: 0,
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

func setTopic(b *bot, update *tb.Update) error {
	m := update.Message
	d := db.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	chat, err := d.GetChat(ctx, m.Chat.ID)
	if err != nil {
		return err
	}
	upd := bson.D{{Key: "$set", Value: bson.D{{Key: "topic_id", Value: m.ThreadID}}}}
	b.CaptchaTopic = int64(m.ThreadID)
	err = d.UpdateChat(ctx, chat, upd)
	if err != nil {
		return err
	}
	b.DeleteMessage(m.Chat.ID, m.ID)
	b.SendMessage("Данный топик выбран в качестве проверочного для пользователей", m.Chat.ID, &tb.MessageOptions{MessageThreadID: m.ThreadID})
	return nil
}
