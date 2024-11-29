package log

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Log struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId     string             `bson:"user_id" json:"user_id"`
	Grade      int64              `bson:"grade" json:"grade"`
	Ocr        []string           `bson:"ocr" json:"ocr"`
	Response   string             `bson:"response" json:"response"`
	Status     int                `bson:"status" json:"status"`
	CreateTime time.Time          `bson:"create_time,omitempty" json:"createTime"`
}
