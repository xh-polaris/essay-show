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
	"time"
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
		records := make([]*show.ListSimpleExercisesResp_Record, 0)
		dto := &show.ListSimpleExercisesResp_SimpleExercise{
			Id:         v.ID.Hex(),
			TotalScore: -1,
			Records:    records,
			FinishTime: time.Time{}.Unix(),
			Like:       v.Like,
		}
		// 有过提交
		if len(v.History.Records) > 0 {
			// 获取最后一次提交记录
			lastRecord := v.History.Records[len(v.History.Records)-1]
			// 获取最后一次记录各题的得分
			for _, r := range lastRecord.Records {
				records = append(records, &show.ListSimpleExercisesResp_Record{
					Id:    r.Id,
					Score: r.Score,
				})
			}
			dto.TotalScore = lastRecord.Score
			dto.FinishTime = lastRecord.CreateTime.Unix()
		} else {
			for _, cq := range v.Question.ChoiceQuestions {
				records = append(records, &show.ListSimpleExercisesResp_Record{
					Id:    cq.Id,
					Score: -1,
				})
			}
		}
		dto.Records = records
		dtos = append(dtos, dto)
	}
	resp.Exercises = dtos
	resp.Total = total
	return

}

func (s ExerciseService) GetExercise(ctx context.Context, req *show.GetExerciseReq) (resp *show.GetExerciseResp, err error) {
	e, err := s.ExerciseMapper.FindOneById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	// 处理题目
	cqs := make([]*show.ChoiceQuestion, 0)
	for _, v := range e.Question.ChoiceQuestions {
		// 处理选项
		os := make([]*show.Option, 0)
		for _, o := range v.Options {
			os = append(os, &show.Option{
				Option:  o.Option,
				Content: o.Content,
				Score:   o.Score,
			})
		}
		cq := &show.ChoiceQuestion{
			Id:          v.Id,
			Question:    v.Question,
			Explanation: v.Explanation,
			Options:     os,
		}
		cqs = append(cqs, cq)
	}

	// 处理答题记录
	rds := make([]*show.Records, 0)
	for _, v := range e.History.Records {
		rs := make([]*show.Record, 0)
		for _, r := range v.Records {
			rs = append(rs, &show.Record{
				Id:     r.Id,
				Option: r.Option,
				Score:  r.Score,
			})
		}
		rds = append(rds, &show.Records{
			Records:    rs,
			Score:      v.Score,
			CreateTime: v.CreateTime.Unix(),
		})
	}
	dto := &show.Exercise{
		Id:         e.ID.Hex(),
		UserId:     e.UserId,
		LogId:      e.LogId,
		Question:   &show.Question{ChoiceQuestions: cqs},
		History:    &show.History{Records: rds},
		Like:       e.Like,
		CreateTime: e.CreateTime.Unix(),
		UpdateTime: e.UpdateTime.Unix(),
		Status:     e.Status,
	}

	resp = &show.GetExerciseResp{
		Code:     0,
		Msg:      "success",
		Exercise: dto,
	}

	return
}

func (s ExerciseService) DoExercise(ctx context.Context, req *show.DoExerciseReq) (resp *show.DoExerciseResp, err error) {
	e, err := s.ExerciseMapper.FindOneById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 初始化
	if e.History == nil {
		rds := make([]*exercise.Records, 0)
		e.History = &exercise.History{Records: rds}
	}

	// 用map存储题目id与题目
	cqs := e.Question.ChoiceQuestions
	qMap := make(map[string]*exercise.ChoiceQuestion)
	for _, v := range cqs {
		qMap[v.Id] = v
	}
	// 做题记录
	rs := make([]*exercise.Record, 0)
	var sum int64
	for _, v := range req.Records {
		// 根据id获取题目
		if q, ok := qMap[v.Id]; ok {
			var score int64
			for _, o := range q.Options {
				if o.Option == v.Option {
					score = o.Score
				}
			}
			r := &exercise.Record{
				Id:     q.Id,
				Option: v.Option,
				Score:  score,
			}
			sum += score
			rs = append(rs, r)
		}
	}
	rds := &exercise.Records{
		Records:    rs,
		Score:      sum,
		CreateTime: time.Now(),
	}

	// 更新记录
	e.History.Records = append(e.History.Records, rds)
	err = s.ExerciseMapper.Update(ctx, e)
	if err != nil {
		return nil, err
	}

	// 将最新的记录返回
	rsDto := make([]*show.Record, 0)
	for _, v := range e.History.Records[len(e.History.Records)-1].Records {
		rsDto = append(rsDto, &show.Record{
			Id:     v.Id,
			Option: v.Option,
			Score:  v.Score,
		})
	}
	dto := &show.Records{
		Records:    rsDto,
		Score:      rds.Score,
		CreateTime: rds.CreateTime.Unix(),
	}
	resp = &show.DoExerciseResp{
		Code:    0,
		Msg:     "success",
		Records: dto,
	}
	return
}
