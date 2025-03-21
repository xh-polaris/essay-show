package service

import (
	"context"
	"errors"
	"github.com/google/wire"
	"github.com/xh-polaris/essay-show/biz/adaptor"
	"github.com/xh-polaris/essay-show/biz/application/dto/essay/show"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/attend"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/invitation"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/user"
	"github.com/xh-polaris/essay-show/biz/infrastructure/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IUserService interface {
	SignUp(ctx context.Context, req *show.SignUpReq) (*show.SignUpResp, error)
	SignIn(ctx context.Context, req *show.SignInReq) (*show.SignInResp, error)
	GetUserInfo(ctx context.Context, req *show.GetUserInfoReq) (*show.GetUserInfoResp, error)
	UpdateUserInfo(ctx context.Context, req *show.UpdateUserInfoReq) (*show.Response, error)
	UpdatePassword(ctx context.Context, req *show.UpdatePasswordReq) (*show.UpdatePasswordResp, error)
	DailyAttend(ctx context.Context, req *show.DailyAttendReq) (*show.Response, error)
	GetDailyAttend(ctx context.Context, req *show.GetDailyAttendReq) (*show.GetDailyAttendResp, error)
	FillInvitationCode(ctx context.Context, req *show.FillInvitationCodeReq) (*show.Response, error)
	GetInvitationCode(ctx context.Context, req *show.GetInvitationCodeReq) (*show.GetInvitationCodeResp, error)
}
type UserService struct {
	UserMapper   *user.MongoMapper
	AttendMapper *attend.MongoMapper
	CodeMapper   *invitation.CodeMongoMapper
	LogMapper    *invitation.LogMongoMapper
}

var UserServiceSet = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)),
)

func (s *UserService) SignUp(ctx context.Context, req *show.SignUpReq) (*show.SignUpResp, error) {
	var oldUser *user.User
	var err error
	if req.AuthType == "phone" {
		// 查找数据库判断手机号是否注册过
		oldUser, err = s.UserMapper.FindOneByPhone(ctx, req.AuthId)
		if err == nil && oldUser != nil {
			return nil, consts.ErrRepeatedSignUp
		} else if err != nil && !errors.Is(err, consts.ErrNotFound) {
			return nil, consts.ErrSignUp
		}
	}

	// 在中台注册账户
	httpClient := util.NewHttpClient()
	signUpResponse, err := httpClient.SignUp(req.AuthType, req.AuthId, &req.VerifyCode)
	if err != nil {
		return nil, consts.ErrSignUp
	}
	userId := signUpResponse["userId"].(string)

	// 在中台设置密码
	authorization := signUpResponse["accessToken"].(string)
	if req.Password != "" {
		_, err = httpClient.SetPassword(authorization, req.Password)
		if err != nil {
			return nil, consts.ErrSignUp
		}
	}

	// 初始化用户
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	u := user.User{
		ID:         oid,
		Username:   req.Name,
		Phone:      req.AuthId,
		Count:      consts.DefaultCount,
		Status:     0,
		CreateTime: now,
		UpdateTime: now,
	}
	if req.AuthType == "wechat-phone" {
		u.Phone = signUpResponse["option"].(string)
	}

	// 向数据库中插入数据
	err = s.UserMapper.Insert(ctx, &u)
	if err != nil {
		return nil, consts.ErrSignUp
	}

	// 返回响应
	return &show.SignUpResp{
		Id:           userId,
		AccessToken:  authorization,
		AccessExpire: int64(signUpResponse["accessExpire"].(float64)),
		Name:         u.Username,
	}, nil
}

func (s *UserService) SignIn(ctx context.Context, req *show.SignInReq) (*show.SignInResp, error) {
	var u *user.User
	var err error
	if req.AuthType == "phone" {
		// 查找数据库判断手机号是否注册过
		u, err = s.UserMapper.FindOneByPhone(ctx, req.AuthId)
		if errors.Is(err, consts.ErrNotFound) || u == nil { // 未找到，说明没有注册
			return nil, consts.ErrNotSignUp
		} else if err != nil {
			return nil, consts.ErrSignUp
		}
	}

	// 通过中台登录
	httpClient := util.NewHttpClient()
	signInResponse, err := httpClient.SignIn(req.AuthType, req.AuthId, req.VerifyCode, req.Password)
	if err != nil {
		return nil, consts.ErrSignIn
	}
	userId, ok := signInResponse["userId"].(string)
	if userId == "" || !ok {
		return nil, consts.ErrSignIn
	}

	// 托底逻辑，如果不注册直接登录也行
	u, err = s.UserMapper.FindOne(ctx, userId)
	if errors.Is(err, consts.ErrNotFound) || u == nil {
		// 初始化用户
		oid, err2 := primitive.ObjectIDFromHex(userId)
		if err2 != nil {
			return nil, err2
		}
		now := time.Now()
		u = &user.User{
			ID:         oid,
			Username:   "未设置用户名",
			Phone:      req.AuthId,
			Count:      consts.DefaultCount,
			Status:     0,
			CreateTime: now,
			UpdateTime: now,
		}
		if req.AuthType == "wechat-phone" {
			u.Phone = signInResponse["option"].(string)
		}

		err = s.UserMapper.Insert(ctx, u)
		if err != nil {
			return nil, consts.ErrSignUp
		}
	} else if err != nil {
		return nil, consts.ErrSignIn
	}

	resp := &show.SignInResp{
		Id:           userId,
		AccessToken:  signInResponse["accessToken"].(string),
		AccessExpire: int64(signInResponse["accessExpire"].(float64)),
		Name:         u.Username,
	}

	return resp, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *show.GetUserInfoReq) (*show.GetUserInfoResp, error) {
	userMeta := adaptor.ExtractUserMeta(ctx)
	if userMeta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}
	u, err := s.UserMapper.FindOne(ctx, userMeta.GetUserId())
	if err != nil {
		return &show.GetUserInfoResp{
			Code:    -1,
			Msg:     "查询用户信息失败，请先登录或重试",
			Payload: nil,
		}, nil
	}
	return &show.GetUserInfoResp{
		Code: 0,
		Msg:  "查询成功",
		Payload: &show.GetUserInfoResp_Payload{
			Name:  u.Username,
			Count: u.Count,
			Phone: u.Phone,
		},
	}, nil
}

func (s *UserService) UpdateUserInfo(ctx context.Context, req *show.UpdateUserInfoReq) (*show.Response, error) {
	// 获取用户id
	userMeta := adaptor.ExtractUserMeta(ctx)
	if userMeta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}

	// 根据用户id查询这个用户
	u, err := s.UserMapper.FindOne(ctx, userMeta.GetUserId())
	if err != nil {
		return nil, consts.ErrNotFound
	}

	// 更新用户信息
	u.Username = req.Name
	u.School = req.School
	u.Grade = req.Grade

	// 存入新的用户信息
	err = s.UserMapper.Update(ctx, u)
	if err != nil {
		return nil, consts.ErrUpdate
	}

	// 返回响应
	return &show.Response{
		Code: 0,
		Msg:  "更新成功",
	}, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, req *show.UpdatePasswordReq) (*show.UpdatePasswordResp, error) {
	// 获取用户id
	userMeta := adaptor.ExtractUserMeta(ctx)
	if userMeta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}

	// 根据用户id查询这个用户
	u, err := s.UserMapper.FindOne(ctx, userMeta.GetUserId())
	if err != nil {
		return nil, consts.ErrNotFound
	}

	// 在中台注册账户
	httpClient := util.NewHttpClient()
	signInResponse, err := httpClient.SignUp(consts.Phone, u.Phone, &req.VerifyCode)
	if err != nil {
		return nil, consts.ErrVerifyCode
	}

	// 在中台设置密码
	authorization := signInResponse["accessToken"].(string)
	_, err = httpClient.SetPassword(authorization, req.Password)
	if err != nil {
		return nil, consts.ErrSignUp
	}
	return &show.UpdatePasswordResp{
		Id:           u.ID.Hex(),
		AccessToken:  authorization,
		AccessExpire: int64(signInResponse["accessExpire"].(float64)),
		Name:         u.Username,
	}, nil
}

func (s *UserService) DailyAttend(ctx context.Context, req *show.DailyAttendReq) (*show.Response, error) {
	// 用户信息
	meta := adaptor.ExtractUserMeta(ctx)
	if meta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}

	// 查询最近的attend记录
	a, err := s.findAttend(ctx, meta.GetUserId())
	if err != nil && !errors.Is(err, consts.ErrNotFound) {
		return nil, consts.ErrDailyAttend
	}

	// 今日有签到记录且不是第一次签到
	if a != nil && time.Unix(a.Timestamp.Unix(), 0).Day() == time.Now().Day() && !a.Timestamp.IsZero() {
		return nil, consts.ErrRepeatDailyAttend
	}

	// 插入新的签到记录
	_a := &attend.Attend{
		ID:        primitive.NewObjectID(),
		UserId:    meta.GetUserId(),
		Timestamp: time.Now(),
	}
	err = s.AttendMapper.Insert(ctx, _a)
	if err != nil {
		return nil, consts.ErrDailyAttend
	}

	// 增加次数
	err = s.UserMapper.UpdateCount(ctx, meta.GetUserId(), consts.AttendReward)
	if err != nil {
		return nil, consts.ErrDailyAttend
	}

	return util.Succeed("签到成功")
}

func (s *UserService) GetDailyAttend(ctx context.Context, req *show.GetDailyAttendReq) (*show.GetDailyAttendResp, error) {
	resp := &show.GetDailyAttendResp{
		Code:   0,
		Msg:    "success",
		Attend: 0,
	}

	// 用户信息
	meta := adaptor.ExtractUserMeta(ctx)
	if meta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}

	// 获取最新的, 确定今天的更新状态
	a, err := s.findAttend(ctx, meta.GetUserId())
	if err != nil {
		return nil, err
	}
	if !a.Timestamp.IsZero() && time.Unix(a.Timestamp.Unix(), 0).Day() == time.Now().Day() {
		resp.Attend = 1
	}

	// 获取所有的指定年月的所有签到记录
	data, total, err := s.AttendMapper.FindByYearAndMonth(ctx, meta.GetUserId(), int(req.Year), int(req.Month))
	if err != nil {
		return nil, err
	}

	dtos := make([]int64, 0, len(data))
	for _, d := range data {
		dtos = append(dtos, int64(d.Timestamp.Day()))
	}
	resp.History = dtos
	resp.Total = total

	return resp, nil
}

func (s *UserService) FillInvitationCode(ctx context.Context, req *show.FillInvitationCodeReq) (*show.Response, error) {
	// 用户信息
	userMeta := adaptor.ExtractUserMeta(ctx)
	if userMeta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}

	// 获取邀请码对应邀请者
	c, err := s.CodeMapper.FindOneByCode(ctx, req.InvitationCode)
	if err != nil {
		return nil, consts.ErrNotFound
	}

	inviter := c.UserId
	invitee := userMeta.GetUserId()

	if invitee == inviter {
		return nil, consts.ErrInvitation
	}

	// 尝试获取邀请记录
	l, err := s.LogMapper.FindOneByInvitee(ctx, invitee)
	if err == nil && l != nil {
		// 已填过邀请码
		return nil, consts.ErrRepeatInvitation
	} else if !errors.Is(err, consts.ErrNotFound) {
		// 异常
		return nil, err
	}

	// 插入邀请记录
	err = s.LogMapper.Insert(ctx, inviter, invitee)
	if err != nil {
		return nil, consts.ErrInvitation
	}

	err = s.UserMapper.UpdateCount(ctx, inviter, consts.InvitationReward)
	if err != nil {
		return nil, err
	}
	return util.Succeed("success")
}

func (s *UserService) GetInvitationCode(ctx context.Context, req *show.GetInvitationCodeReq) (*show.GetInvitationCodeResp, error) {
	// 用户信息
	userMeta := adaptor.ExtractUserMeta(ctx)
	if userMeta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}

	c, err := s.CodeMapper.FindOneByUserId(ctx, userMeta.GetUserId())
	if errors.Is(err, consts.ErrNotFound) {
		c, err = s.CodeMapper.Insert(ctx, userMeta.GetUserId())
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &show.GetInvitationCodeResp{
		Code:           0,
		Msg:            "success",
		InvitationCode: c.Code,
	}, nil
}

func (s *UserService) findAttend(ctx context.Context, userId string) (*attend.Attend, error) {
	a, err := s.AttendMapper.FindLatestOneByUserId(ctx, userId)
	return a, err
}
