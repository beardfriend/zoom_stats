package models

import (
	"time"

	z "github.com/beardfriend/zoom_chatting_parser"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UniversalTime time.Time

func NewUnivrersalTimeFromString(s string) UniversalTime {
	t, _ := time.Parse("2006-01-02T15:04:05", s)
	universalTime := UniversalTime(t.Add(9 * time.Hour))
	return universalTime
}

func NewUnivrersalTimeFromTime(t time.Time) UniversalTime {
	universalTime := UniversalTime(t.Add(-9 * time.Hour))
	return universalTime
}

func (d UniversalTime) ToKoreanTimeString() string {
	t := time.Time(d)
	t.Add(-9 * time.Hour)
	return t.Format("2006-01-02T15:04:05")
}

func (d *UniversalTime) UnMarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	uniTime := NewUnivrersalTimeFromString(string(data))
	*d = uniTime

	return nil
}

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
