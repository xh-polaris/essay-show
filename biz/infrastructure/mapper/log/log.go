package log

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Log struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title string             `bson:"title" json:"title"`
	Text  [][]string         `bson:"content" json:"content"`
	// TODO 完善log定义
	Status     int       `bson:"status" json:"status"`
	CreateTime time.Time `bson:"create_time,omitempty" json:"createTime"`
	UpdateTime time.Time `bson:"update_time,omitempty" json:"updateTime"`
	DeleteTime time.Time `bson:"delete_time,omitempty" json:"deleteTime"`
}
