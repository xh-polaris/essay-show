package attend

import (
	"errors"
	"github.com/xh-polaris/essay-show/biz/infrastructure/config"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"time"
)

const (
	prefixKeyCacheKey = "cache:attend"
	CollectionName    = "attend"
)

type IMongoMapper interface {
	Insert(ctx context.Context, userId string) (*Attend, error)
	FindOneByUserId(ctx context.Context, userId string) (a *Attend, err error)
	Update(ctx context.Context, a *Attend) error
}

type MongoMapper struct {
	conn *monc.Model
}

func NewMongoMapper(config *config.Config) *MongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.Cache)
	return &MongoMapper{conn: conn}
}

func (m *MongoMapper) Insert(ctx context.Context, userId string) (*Attend, error) {
	a := &Attend{
		ID:        primitive.NewObjectID(),
		UserId:    userId,
		Timestamp: time.Time{},
	}
	key := prefixKeyCacheKey + userId
	_, err := m.conn.InsertOne(ctx, key, a)
	return a, err
}

func (m *MongoMapper) FindOneByUserId(ctx context.Context, userId string) (a *Attend, err error) {
	a = &Attend{}
	key := prefixKeyCacheKey + userId
	err = m.conn.FindOne(ctx, key, a, bson.M{consts.UserID: userId})
	switch {
	case err == nil:
		return a, nil
	case errors.Is(err, mongo.ErrNoDocuments):
		return nil, consts.ErrNotFound
	default:
		return nil, err
	}
}

func (m *MongoMapper) Update(ctx context.Context, a *Attend) error {
	key := prefixKeyCacheKey + a.UserId
	_, err := m.conn.UpdateByID(ctx, key, a.ID, bson.M{"$set": a})
	return err
}
