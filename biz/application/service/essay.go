package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/essay-show/biz/application/dto/essay/show"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/log"
)

type IEssayService interface {
	EssayEvaluate(ctx context.Context, req *show.EssayEvaluateReq) (resp *show.EssayEvaluateResp, err error)
	GetEvaluateLogs(ctx context.Context, req *show.GetEssayEvaluateLogsReq) (resp *show.GetEssayEvaluateLogsResp, err error)
}

type EssayService struct {
	LogMapper *log.MongoMapper
}

var EssayServiceSet = wire.NewSet(
	wire.Struct(new(EssayService), "*"),
	wire.Bind(new(IEssayService), new(*EssayService)),
)

func (s *EssayService) EssayEvaluate(ctx context.Context, req *show.EssayEvaluateReq) (*show.EssayEvaluateResp, error) {
	// TODO 实现作文批改逻辑
	return nil, nil
}

func (s *EssayService) GetEvaluateLogs(ctx context.Context, req *show.GetEssayEvaluateLogsReq) (resp *show.GetEssayEvaluateLogsResp, err error) {
	// TODO 实现分页获取批改记录的逻辑
	return nil, nil
}
