package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username   string             `bson:"username" json:"username"`
	Count      int64              `bson:"count" json:"count"` // 剩余可用批改次数
	Status     int                `bson:"status" json:"status"`
	CreateTime time.Time          `bson:"create_time,omitempty" json:"createTime"`
	UpdateTime time.Time          `bson:"update_time,omitempty" json:"updateTime"`
	DeleteTime time.Time          `bson:"delete_time,omitempty" json:"deleteTime"`
}
