package service

import (
	"errors"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
	"github.com/xh-polaris/essay-show/biz/adaptor"
	"github.com/xh-polaris/essay-show/biz/application/dto/essay/show"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/exercise"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/log"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/user"
	eu "github.com/xh-polaris/essay-show/biz/infrastructure/util/exercise"
	"golang.org/x/net/context"
)

type IExerciseService interface {
	CreateExercise(ctx context.Context, req *show.CreateExerciseReq) (resp *show.CreateExerciseResp, err error)
	ListSimpleExercises(ctx context.Context, req *show.ListSimpleExercisesReq) (resp *show.ListSimpleExercisesResp, err error)
	GetExercise(ctx context.Context, req *show.GetExerciseReq) (resp *show.GetExerciseResp, err error)
	DoExercise(ctx context.Context, req *show.DoExerciseReq) (resp *show.DoExerciseResp, err error)
}

type ExerciseService struct {
	ExerciseMapper *exercise.MongoMapper
	LogMapper      *log.MongoMapper
	UserMapper     *user.MongoMapper
}

var ExerciseServiceSet = wire.NewSet(
	wire.Struct(new(ExerciseService), "*"),
	wire.Bind(new(IExerciseService), new(*ExerciseService)),
)

func (s ExerciseService) CreateExercise(ctx context.Context, req *show.CreateExerciseReq) (resp *show.CreateExerciseResp, err error) {
	// 获取批改记录
	l, err := s.LogMapper.FindOne(ctx, req.LogId)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	userMeta := adaptor.ExtractUserMeta(ctx)
	if userMeta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}
	u, err := s.UserMapper.FindOne(ctx, userMeta.UserId)
	if err != nil {
		return nil, err
	}

	// 调用生成接口
	e, err := eu.GenerateExercise(u.Grade, l)
	if err != nil {
		return nil, err
	}

	e.LogId = req.LogId
	e.UserId = userMeta.UserId

	err = s.ExerciseMapper.Insert(ctx, e)
	if err != nil {
		return nil, err
	}

	dto := &show.Exercise{}
	err = copier.Copy(dto, e)
	if err != nil {
		return nil, err
	}
	dto.Id = e.ID.Hex()
	dto.CreateTime = e.CreateTime.Unix()
	dto.UpdateTime = e.CreateTime.Unix()

	resp = &show.CreateExerciseResp{
		Exercise: dto,
	}
	return
}

func (s ExerciseService) ListSimpleExercises(ctx context.Context, req *show.ListSimpleExercisesReq) (resp *show.ListSimpleExercisesResp, err error) {
	// 获取用户信息
	userMeta := adaptor.ExtractUserMeta(ctx)
	if userMeta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}

	// 查询
	data, total, err := s.ExerciseMapper.FindManyByLogId(ctx, req.LogId, req.PaginationOptions)
	if err != nil && !errors.Is(err, consts.ErrNotFound) {
		return nil, err
	}

	resp = &show.ListSimpleExercisesResp{
		Code: 0,
		Msg:  "success",
	}

	dtos := make([]*show.ListSimpleExercisesResp_SimpleExercise, 0)
	for _, v := range data {
		// 获取最后一次提交记录
		lastRecord := v.History.Records[len(v.History.Records)-1]
		records := make([]*show.ListSimpleExercisesResp_Record, 0)
		// 获取最后一次记录各题的得分
		for _, r := range lastRecord.Records {
			records = append(records, &show.ListSimpleExercisesResp_Record{
				Id:    r.Id,
				Score: r.Score,
			})
		}
		dto := &show.ListSimpleExercisesResp_SimpleExercise{
			Id:         v.ID.Hex(),
			TotalScore: lastRecord.Score,
			Records:    records,
			FinishTime: lastRecord.CreateTime.Unix(),
			Like:       v.Like,
		}
		dtos = append(dtos, dto)
	}
	resp.Exercises = dtos
	resp.Total = total
	return

}

func (s ExerciseService) GetExercise(ctx context.Context, req *show.GetExerciseReq) (resp *show.GetExerciseResp, err error) {
	return nil, nil
}

func (s ExerciseService) DoExercise(ctx context.Context, req *show.DoExerciseReq) (resp *show.DoExerciseResp, err error) {
	return nil, nil
}
