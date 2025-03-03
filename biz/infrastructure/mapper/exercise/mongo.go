package exercise

import (
	"github.com/xh-polaris/essay-show/biz/application/dto/basic"
	"github.com/xh-polaris/essay-show/biz/infrastructure/config"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	util "github.com/xh-polaris/essay-show/biz/infrastructure/util/page"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"time"
)

const (
	prefixKeyCacheKey = "cache:exercise"
	CollectionName    = "exercise"
)

type IMongoMapper interface {
	Insert(ctx context.Context, e *Exercise) error
	FindManyByLogId(ctx context.Context, logId string, p *basic.PaginationOptions) (exercise []*Exercise, total int64, err error)
	FindOneById(ctx context.Context, id string) (*Exercise, error)
}

type MongoMapper struct {
	conn *monc.Model
}

func NewMongoMapper(config *config.Config) *MongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.Cache)
	return &MongoMapper{conn: conn}
}

func (m *MongoMapper) Insert(ctx context.Context, e *Exercise) error {
	if e.ID.IsZero() {
		e.ID = primitive.NewObjectID()
		e.CreateTime = time.Now()
		e.UpdateTime = time.Now()
	}
	key := prefixKeyCacheKey + e.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, e)
	return err
}

func (m *MongoMapper) FindManyByLogId(ctx context.Context, logId string, p *basic.PaginationOptions) (exercise []*Exercise, total int64, err error) {
	skip, limt := util.ParsePageOpt(p)

	filter := bson.M{
		consts.LogId:  logId,
		consts.Status: bson.M{consts.NotEqual: consts.DeleteStatus},
	}

	opt := &options.FindOptions{
		Limit: &limt,
		Skip:  &skip,
		Sort:  bson.M{consts.CreateTime: -1},
	}

	var data []*Exercise

	err = m.conn.Find(ctx, data, filter, opt)
	if err != nil {
		return nil, 0, err
	}

	total, err = m.conn.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (m *MongoMapper) FindOneById(ctx context.Context, id string) (*Exercise, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		consts.ID: oid,
	}

	key := prefixKeyCacheKey + id

	var e *Exercise
	err = m.conn.FindOne(ctx, key, e, filter)
	return e, err
}
