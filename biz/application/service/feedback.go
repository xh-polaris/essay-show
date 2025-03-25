package service

import (
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
