package db

type Chat struct {
	Id    int64
	Title string
}

type User struct {
	Id             int64
	ChatId         int64  `bson:"chat_id"`
	Username       string `bson:"username"`
	FirstName      string `bson:"first_name"`
	LastName       string `bson:"last_name"`
	CorrectAnswer  int8   `bson:"correct_answer"`
	CaptchaMessage int    `bson:"captcha_message"`
	IsBanned       bool   `bson:"is_banned"`
	DateJoined     int64  `bson:"date_joined"`
}
