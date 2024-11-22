package user

import (
	"context"
	"github.com/xh-polaris/essay-show/biz/infrastructure/config"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	prefixUserCacheKey = "cache:user"
	CollectionName     = "user"
)

type IMongoMapper interface {
	Insert(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	FindOne(ctx context.Context, id string) (*User, error)
	Delete(ctx context.Context, id string) error
}

type MongoMapper struct {
	conn *monc.Model
}

func NewMongoMapper(config *config.Config) *MongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.Cache)
	return &MongoMapper{
		conn: conn,
	}
}

func (m *MongoMapper) Insert(ctx context.Context, user *User) error {
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
		user.CreateTime = time.Now()
		user.UpdateTime = user.CreateTime
	}
	_, err := m.conn.InsertOneNoCache(ctx, user)
	return err
}

func (m *MongoMapper) Update(ctx context.Context, user *User) error {
	user.UpdateTime = time.Now()
	_, err := m.conn.UpdateByIDNoCache(ctx, user.ID, bson.M{"$set": user})
	return err
}
