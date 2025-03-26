package feedback

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feedback struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId     string             `bson:"user_id" json:"userId"`         // 提交反馈的用户ID
	Type       int64              `bson:"type" json:"type"`              // 反馈类型 1系统功能，2功能建议，3界面建议，4批改信度，5题目内容，6素材内容
	Content    string             `bson:"content" json:"content"`        // 反馈内容
	Status     int                `bson:"status" json:"status"`          // 处理状态
	Images     []string           `bson:"images" json:"images"`          // 用户上传的图片URL列表
	CreateTime time.Time          `bson:"create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time          `bson:"update_time" json:"updateTime"` // 更新时间
}
