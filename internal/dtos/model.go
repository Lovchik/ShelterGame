package dtos

type User struct {
	Id       int64  `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	ChatId   int64  `json:"chat_id"`
}
