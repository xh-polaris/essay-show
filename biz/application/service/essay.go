package service

import (
	"context"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
	"github.com/xh-polaris/essay-show/biz/adaptor"
	"github.com/xh-polaris/essay-show/biz/application/dto/essay/show"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/log"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/user"
	"github.com/xh-polaris/essay-show/biz/infrastructure/util"
	logx "github.com/xh-polaris/essay-show/biz/infrastructure/util/log"
	"time"
)

type IEssayService interface {
	EssayEvaluate(ctx context.Context, req *show.EssayEvaluateReq) (resp *show.EssayEvaluateResp, err error)
	GetEvaluateLogs(ctx context.Context, req *show.GetEssayEvaluateLogsReq) (resp *show.GetEssayEvaluateLogsResp, err error)
	LikeEvaluate(ctx context.Context, req *show.LikeEvaluateReq) (resp *show.Response, err error)
}

type EssayService struct {
	LogMapper  *log.MongoMapper
	UserMapper *user.MongoMapper
}

var EssayServiceSet = wire.NewSet(
	wire.Struct(new(EssayService), "*"),
	wire.Bind(new(IEssayService), new(*EssayService)),
)

func (s *EssayService) EssayEvaluate(ctx context.Context, req *show.EssayEvaluateReq) (*show.EssayEvaluateResp, error) {
	// TODO 应该实现一个用户同时只能调用一次批改

	// 先判断用户是否有充足次数
	userMeta := adaptor.ExtractUserMeta(ctx)
	if userMeta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}
	aUser, err := s.UserMapper.FindOne(ctx, userMeta.GetUserId())
	if err != nil {
		return nil, consts.ErrNotFound
	}
	// 剩余次数不足
	if aUser.Count <= 0 {
		return nil, consts.ErrInSufficientCount
	}

	// 走计费层调用call接口批改作文
	httpClient := util.NewHttpClient()
	callResponse, err := httpClient.BetaEvaluate(req.Title, req.Text)
	if err != nil { // 调用call失败
		return nil, consts.ErrCall
	}

	// 获取批改的结果
	code := int64(callResponse["code"].(float64))
	msg := callResponse["msg"].(string)
	result := callResponse["result"].(string)

	// 批改失败，记录对应的情况
	if code != 0 {
		l := &log.Log{
			Grade:      req.Grade,
			Ocr:        req.Ocr,
			Response:   result,
			Status:     int(code),
			CreateTime: time.Now(),
		}
		err2 := s.LogMapper.Insert(ctx, l)
		if err2 != nil {
			return nil, consts.ErrCall
		}
		return nil, consts.ErrCall
	}

	resp := &show.EssayEvaluateResp{
		Code:     code,
		Msg:      msg,
		Response: result,
	}

	// 扣除用户剩余次数
	err = s.UserMapper.UpdateCount(ctx, userMeta.GetUserId(), -1)
	if err != nil {
		return nil, err //  扣除失败用户不应该拿到结果
	}

	// 批改成功，添加记录
	l := &log.Log{
		UserId:     userMeta.GetUserId(),
		Grade:      req.Grade,
		Ocr:        req.Ocr,
		Response:   result,
		Status:     int(code),
		CreateTime: time.Now(),
	}
	err = s.LogMapper.Insert(ctx, l)
	if err != nil {
		// 记录插入失败应该也要获得结果，因为剩余次数已经成功扣除。TODO 但是这里的错误要怎么处理呢？
		return resp, nil
	}

	return &show.EssayEvaluateResp{
		Code:     code,
		Msg:      msg,
		Response: result,
		Id:       l.ID.Hex(),
	}, nil
}

func (s *EssayService) GetEvaluateLogs(ctx context.Context, req *show.GetEssayEvaluateLogsReq) (resp *show.GetEssayEvaluateLogsResp, err error) {
	userMeta := adaptor.ExtractUserMeta(ctx)
	if userMeta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}

	data, total, err := s.LogMapper.FindMany(ctx, userMeta.GetUserId(), req.PaginationOptions)
	if err != nil {
		return nil, err
	}
	var logs []*show.Log
	for _, val := range data {
		l := &show.Log{}
		err = copier.Copy(l, val)
		if err != nil {
			return nil, err
		}
		l.Id = val.ID.Hex()
		l.CreateTime = val.CreateTime.Unix()
		logs = append(logs, l)
	}

	return &show.GetEssayEvaluateLogsResp{
		Total: total,
		Logs:  logs,
	}, nil
}

func (s *EssayService) LikeEvaluate(ctx context.Context, req *show.LikeEvaluateReq) (resp *show.Response, err error) {
	l, err := s.LogMapper.FindOne(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	l.Like = req.Like
	err = s.LogMapper.Update(ctx, l)
	if err != nil {
		logx.Error(err.Error())
		return &show.Response{
			Code: 0,
			Msg:  "标记失败",
		}, nil
	}
	return &show.Response{
		Code: 0,
		Msg:  "标记成功",
	}, nil
}
