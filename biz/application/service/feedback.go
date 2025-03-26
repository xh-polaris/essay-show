package service

import (
	"strconv"

	"github.com/google/wire"
	"github.com/xh-polaris/essay-show/biz/adaptor"
	"github.com/xh-polaris/essay-show/biz/application/dto/essay/show"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/feedback"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/user"
	"github.com/xh-polaris/essay-show/biz/infrastructure/util"
	"golang.org/x/net/context"
)

type IFeedbackService interface {
	Submit(ctx context.Context, req *show.SubmitFeedbackReq) (*show.Response, error)
	ListFeedback(ctx context.Context, req *show.ListFeedbackReq) (*show.ListFeedbackResp, error)
}

type FeedBackService struct {
	FeedbackMapper *feedback.MongoMapper
	UserMapper     *user.MongoMapper
}

var FeedbackServiceSet = wire.NewSet(
	wire.Struct(new(FeedBackService), "*"),
	wire.Bind(new(IFeedbackService), new(*FeedBackService)),
)

func (s *FeedBackService) Submit(ctx context.Context, req *show.SubmitFeedbackReq) (*show.Response, error) {
	meta := adaptor.ExtractUserMeta(ctx)
	if meta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}

	f := &feedback.Feedback{
		UserId:  meta.UserId,
		Type:    req.Type,
		Content: req.Content,
		Status:  0,
		Images:  req.Images,
	}

	err := s.FeedbackMapper.Insert(ctx, f)
	if err != nil {
		return util.Fail(999, "反馈失败"), nil
	}
	return util.Succeed("反馈成功")
}

func (s *FeedBackService) ListFeedback(ctx context.Context, req *show.ListFeedbackReq) (*show.ListFeedbackResp, error) {
	meta := adaptor.ExtractUserMeta(ctx)
	if meta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}

	// 查询当前用户的反馈
	userId := meta.GetUserId()

	// 查询反馈列表
	feedbacks, total, err := s.FeedbackMapper.FindMany(ctx, userId, req.Type, req.PaginationOptions)
	if err != nil {
		return nil, err
	}

	// 构建响应
	respFeedbacks := make([]*show.ListFeedbackResp_FeedbackItem, 0, len(feedbacks))
	// make 创建一个切片 len(feedbacks) 预分配容量为 feedbacks 的长度 减少 append 过程中内存重新分配的次数
	for _, f := range feedbacks { // 遍历数据库查询到的反馈数据
		item := &show.ListFeedbackResp_FeedbackItem{
			Id:         f.ID.Hex(),
			Type:       f.Type,
			Content:    f.Content,
			Images:     f.Images,
			CreateTime: f.CreateTime.Format("2006-01-02 15:04:05"),
			Status:     strconv.Itoa(f.Status),
		}
		respFeedbacks = append(respFeedbacks, item)
	}

	return &show.ListFeedbackResp{
		Code:      0,
		Msg:       "查询成功",
		Total:     total,
		Feedbacks: respFeedbacks,
	}, nil
}
