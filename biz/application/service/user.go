package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/essay-show/biz/application/dto/essay/show"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/user"
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
	// TODO 实现注册逻辑
	return nil, nil
}
