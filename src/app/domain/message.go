package domain

import "time"

type Message struct {
	Id          string    `bson:"_id" json:"id,omitempty"`
	Text        string    `json:"text"`
	CreatedTime time.Time `json:"createdTime"`
	UserId      string    `json:"userId"`
}

type MessageRepo interface {
	CreateMessage(text string, userId string) (Message, error)

	GetFiveLastMessages() []Message
}
