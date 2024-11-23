package log

import (
	"context"
	"github.com/xh-polaris/essay-show/biz/application/dto/basic"
	"github.com/xh-polaris/essay-show/biz/infrastructure/config"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	util "github.com/xh-polaris/essay-show/biz/infrastructure/util/page"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	prefixKeyCacheKey = "cache:log"
	CollectionName    = "log"
)

type IMongoMapper interface {
	Insert(ctx context.Context, l *Log) error
	FindMany(ctx context.Context, userId string, p basic.PaginationOptions) (logs []*Log, total int64, err error)
}

type MongoMapper struct {
	conn *monc.Model
}

func NewMongoMapper(config *config.Config) *MongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.Cache)
	return &MongoMapper{conn: conn}
}

func (m *MongoMapper) Insert(ctx context.Context, l *Log) error {
	if l.ID.IsZero() {
		l.ID = primitive.NewObjectID()
		l.CreateTime = time.Now()
	}
	key := prefixKeyCacheKey + l.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, l)
	return err
}

func (m *MongoMapper) FindMany(ctx context.Context, userId string, p *basic.PaginationOptions) (logs []*Log, total int64, err error) {
	skip, limit := util.ParsePageOpt(p)
	logs = make([]*Log, 0, limit)
	err = m.conn.Find(ctx, &logs,
		bson.M{
			consts.UserID: userId,
		}, &options.FindOptions{
			Skip:  &skip,
			Limit: &limit,
			Sort:  bson.M{consts.CreateTime: -1},
		})
	if err != nil {
		return nil, 0, err
	}

	total, err = m.conn.CountDocuments(ctx, bson.M{
		consts.UserID: userId,
	})
	if err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
