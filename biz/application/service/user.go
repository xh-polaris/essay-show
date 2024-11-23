package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/essay-show/biz/adaptor"
	"github.com/xh-polaris/essay-show/biz/application/dto/essay/show"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/user"
	"github.com/xh-polaris/essay-show/biz/infrastructure/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IUserService interface {
	SignUp(ctx context.Context, req *show.SignUpReq) (*show.SignUpResp, error)
}
type UserService struct {
	UserMapper *user.MongoMapper
}

var UserServiceSet = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)),
)

func (u *UserService) SignUp(ctx context.Context, req *show.SignUpReq) (*show.SignUpResp, error) {
	// 在中台注册账户
	httpClient := util.NewHttpClient()
	signUpResponse, err := httpClient.SignUp(req.AuthType, req.AuthId, &req.VerifyCode)
	if err != nil {
		return nil, consts.ErrSigunUp
	}

	// 在中台设置密码
	authorization := signUpResponse["accessToken"].(string)
	_, err = httpClient.SetPassword(authorization, req.Password)
	if err != nil {
		return nil, consts.ErrSigunUp
	}

	// 初始化用户
	userId := signUpResponse["userId"].(string)
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	aUser := user.User{
		ID:         oid,
		Username:   req.Name,
		Count:      consts.DefaultCount,
		Status:     0,
		CreateTime: now,
		UpdateTime: now,
	}

	// 向数据库中插入数据
	err = u.UserMapper.Insert(ctx, &aUser)
	if err != nil {
		return nil, consts.ErrSigunUp
	}

	// 返回响应
	return &show.SignUpResp{
		Id:           userId,
		AccessToken:  authorization,
		AccessExpire: signUpResponse["accessExpire"].(int64),
	}, nil
}

func (u *UserService) SignIn(ctx context.Context, req *show.SignInReq) (*show.SignInResp, error) {
	// 通过中台登录
	httpClient := util.NewHttpClient()
	signInResponse, err := httpClient.SignIn(req.AuthType, req.AuthId, req.VerifyCode, req.Password)
	if err != nil {
		return nil, consts.ErrSigunUp
	}

	return &show.SignInResp{
		Id:           signInResponse["userId"].(string),
		AccessToken:  signInResponse["accessToken"].(string),
		AccessExpire: signInResponse["accessExpire"].(int64),
	}, nil
}

func (u *UserService) GetUserInfo(ctx context.Context, s *show.GetUserInfoReq) (*show.GetUserInfoResp, error) {
	userMeta := adaptor.ExtractUserMeta(ctx)
	if userMeta.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}
	aUser, err := u.UserMapper.FindOne(ctx, userMeta.GetUserId())
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
			Name:  aUser.Username,
			Count: aUser.Count,
		},
	}, nil
}
