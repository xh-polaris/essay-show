package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/xh-polaris/essay-show/biz/adaptor"
	"github.com/xh-polaris/essay-show/biz/application/dto/essay/show"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/user"
	"github.com/xh-polaris/essay-show/biz/infrastructure/rpc/platform_sts"
	"github.com/xh-polaris/essay-show/biz/infrastructure/util"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/platform/sts"
	"net/http"
)

type IStsService interface {
	ApplySignedUrl(ctx context.Context, req *show.ApplySignedUrlReq) (*show.ApplySignedUrlResp, error)
	OCR(ctx context.Context, req *show.OCRReq) (*show.OCRResp, error)
	SendVerifyCode(ctx context.Context, req *show.SendVerifyCodeReq) (*show.Response, error)
}

type StsService struct {
	PlatformSts platform_sts.IPlatformSts
	UserMapper  *user.MongoMapper
}

var StsServiceSet = wire.NewSet(
	wire.Struct(new(StsService), "*"),
	wire.Bind(new(IStsService), new(*StsService)),
)

// ApplySignedUrl 向cos申请加签url
func (s *StsService) ApplySignedUrl(ctx context.Context, req *show.ApplySignedUrlReq) (*show.ApplySignedUrlResp, error) {
	// 获取用户信息
	aUser := adaptor.ExtractUserMeta(ctx)
	if aUser.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}
	// 构造响应
	resp := new(show.ApplySignedUrlResp)
	// 获取cos状态
	userId := aUser.GetUserId()
	data, err := s.PlatformSts.GenCosSts(ctx, &sts.GenCosStsReq{Path: "essays/" + userId + "/*"})
	if err != nil {
		return nil, err
	}

	// 生成加签url
	resp.SessionToken = data.SessionToken
	if req.Prefix != nil {
		*req.Prefix += "/"
	}
	data2, err := s.PlatformSts.GenSignedUrl(ctx, &sts.GenSignedUrlReq{
		SecretId:  data.SecretId,
		SecretKey: data.SecretKey,
		Method:    http.MethodPut,
		Path:      "essays/" + userId + "/" + req.GetPrefix() + uuid.New().String() + req.GetSuffix(),
	})
	if err != nil {
		return nil, err
	}
	// 返回响应
	resp.Url = data2.SignedUrl
	return resp, nil
}

func (s *StsService) OCR(ctx context.Context, req *show.OCRReq) (*show.OCRResp, error) {
	// 获取用户信息
	aUser := adaptor.ExtractUserMeta(ctx)
	if aUser.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}

	// 图片url与保留类型
	images := req.Ocr
	left := ""
	if req.LeftType != nil {
		left = *req.LeftType
	}

	// 调用ocr接口
	client := util.GetHttpClient()
	resp, err := client.BeeTitleUrlOCR(images, left)
	if err != nil {
		return nil, err
	}

	return &show.OCRResp{Title: resp["title"].(string), Text: resp["content"].(string)}, nil
}

// SendVerifyCode 发送验证码
func (s *StsService) SendVerifyCode(ctx context.Context, req *show.SendVerifyCodeReq) (*show.Response, error) {
	// 查找用户
	aUser, err := s.UserMapper.FindOneByPhone(ctx, req.AuthId)

	if req.Type == 1 { // 登录验证码
		// 查找数据库判断手机号是否注册过
		if errors.Is(err, consts.ErrNotFound) || aUser == nil { // 未找到，说明没有注册
			return nil, consts.ErrNotSignUp
		} else if err != nil {
			return nil, consts.ErrSend
		}
	} else { // 注册验证码
		if err == nil && aUser != nil {
			return nil, consts.ErrRepeatedSignUp
		} else if err != nil && !errors.Is(err, consts.ErrNotFound) {
			return nil, consts.ErrSignUp
		}
	}

	// 通过中台发送验证码
	httpClient := util.GetHttpClient()
	_, err = httpClient.SendVerifyCode(req.AuthType, req.AuthId)
	if err != nil {
		return nil, consts.ErrSend
	}

	return util.Succeed("发送验证码成功，请注意查收")
}
