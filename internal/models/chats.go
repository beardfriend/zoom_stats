package models

import (
	"time"

	z "github.com/beardfriend/zoom_chatting_parser"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID        primitive.ObjectID   `json:"_id" bson:"_id"`
	ReactIds  []primitive.ObjectID `json:"react_ids" bson:"react_ids"`
	ReplyIds  []primitive.ObjectID `json:"reply_ids" bson:"reply_ids"`
	TextType  z.ChatContentType    `json:"text_type" bson:"text_type"`
	Sender    string               `json:"sender" bson:"sender"`
	Text      string               `json:"text" bson:"text"`
	Removed   bool                 `json:"removed" bson:"removed"`
	ChattedAt time.Time            `json:"chatted_at" bson:"chatted_at"`
}
