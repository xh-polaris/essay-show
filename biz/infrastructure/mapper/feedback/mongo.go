package feedback

import (
	"context"
	"time"

	"github.com/xh-polaris/essay-show/biz/application/dto/basic"
	"github.com/xh-polaris/essay-show/biz/infrastructure/config"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	util "github.com/xh-polaris/essay-show/biz/infrastructure/util/page"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	prefixKeyCacheKey = "cache:feedback"
	CollectionName    = "feedback"
)

type IMongoMapper interface {
	Insert(ctx context.Context, f *Feedback) error
	FindMany(ctx context.Context, userId string, feedbackType *int64, p *basic.PaginationOptions) (feedbacks []*Feedback, total int64, err error)
}

type MongoMapper struct {
	conn *monc.Model
}

func NewMongoMapper(config *config.Config) *MongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.Cache)
	return &MongoMapper{conn: conn}
}

func (m *MongoMapper) Insert(ctx context.Context, f *Feedback) error {
	if f.ID.IsZero() {
		f.ID = primitive.NewObjectID()
		f.CreateTime = time.Now()
		f.UpdateTime = time.Now()
	}
	_, err := m.conn.InsertOneNoCache(ctx, f)
	return err
}

func (m *MongoMapper) FindMany(ctx context.Context, userId string, feedbackType *int64, p *basic.PaginationOptions) (feedbacks []*Feedback, total int64, err error) {
	skip, limit := util.ParsePageOpt(p)
	feedbacks = make([]*Feedback, 0, limit) // 预分配内存，减少 append 过程中内存重新分配的次数，提高性能。
	// 查询条件
	filter := bson.M{}

	// 只有在提供了用户ID时才添加过滤条件
	if userId != "" {
		filter[consts.UserID] = userId // 确保用户只能看到自己的反馈，而不是所有反馈。
	}

	// 只有在提供了反馈类型时才添加过滤条件
	if feedbackType != nil {
		filter["type"] = *feedbackType
	}

	// 查询数据
	err = m.conn.Find(ctx, &feedbacks, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  bson.M{consts.CreateTime: -1}, // 按创建时间降序排序，最新的数据排在最前面
	})
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err = m.conn.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}
