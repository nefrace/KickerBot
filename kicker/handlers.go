package kicker

import (
	tb "gopkg.in/tucnak/telebot.v3"
)

var HandlersV1 = []Handler{
	{
		Endpoint: tb.OnText,
		Handler: func(c tb.Context) error {
			m := c.Message()
			c.Bot().Send(m.Sender, m.Text)
			return nil
		},
	},
}
