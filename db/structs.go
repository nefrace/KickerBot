package db

type Chat struct {
	Id    int64
	Title string
}

type User struct {
	Id            int64
	ChatId        int64  `bson:"chat_id"`
	Username      string `bson:"username"`
	FirstName     string `bson:"first_name"`
	LastName      string `bson:"last_name"`
	CorrectAnswer int8   `bson:"correct_answer"`
	IsBanned      bool   `bson:"is_banned"`
}

type Captcha struct {
	MessageId     int  `bson:"message_id"`
	CorrectAnswer int8 `bson:"correct_answer"`
}
