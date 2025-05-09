// Code generated by hertz generator.

package show

import (
	"context"
	"github.com/xh-polaris/essay-show/biz/adaptor"
	"github.com/xh-polaris/essay-show/biz/application/dto/essay/show"
	"github.com/xh-polaris/essay-show/provider"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// SignUp .
// @router /user/sign_up [POST]
func SignUp(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.SignUpReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.UserService.SignUp(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// SignIn .
// @router /user/sign_in [POST]
func SignIn(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.SignInReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.UserService.SignIn(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetUserInfo .
// @router /user/info [GET]
func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.GetUserInfoReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.UserService.GetUserInfo(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// EssayEvaluate .
// @router /essay/evaluate [POST]
func EssayEvaluate(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.EssayEvaluateReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.EssayService.EssayEvaluate(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetEvaluateLogs .
// @router /essay/logs [POST]
func GetEvaluateLogs(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.GetEssayEvaluateLogsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.EssayService.GetEvaluateLogs(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// LikeEvaluate .
// @router /essay/like [POST]
func LikeEvaluate(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.LikeEvaluateReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.EssayService.LikeEvaluate(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// UpdateUserInfo .
// @router /user/update [POST]
func UpdateUserInfo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.UpdateUserInfoReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.UserService.UpdateUserInfo(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// UpdatePassword .
// @router /user/update_password [POST]
func UpdatePassword(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.UpdatePasswordReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.UserService.UpdatePassword(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetInvitationCode .
// @router /user/invitation/code [GET]
func GetInvitationCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.GetInvitationCodeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.UserService.GetInvitationCode(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// FillInvitationCode .
// @router /user/invitation/fill [POST]
func FillInvitationCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.FillInvitationCodeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.UserService.FillInvitationCode(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetDailyAttend .
// @router /user/daily_attend/get [GET]
func GetDailyAttend(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.GetDailyAttendReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.UserService.GetDailyAttend(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// DailyAttend .
// @router /user/daily_attend [GET]
func DailyAttend(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.DailyAttendReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.UserService.DailyAttend(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// OCR .
// @router /sts/ocr [POST]
func OCR(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.OCRReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.StsService.OCR(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// ApplySignedUrl .
// @router /sts/apply [POST]
func ApplySignedUrl(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.ApplySignedUrlReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.StsService.ApplySignedUrl(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// SendVerifyCode .
// @router /sts/send_verify_code [POST]
func SendVerifyCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.SendVerifyCodeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	p := provider.Get()
	resp, err := p.StsService.SendVerifyCode(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// SubmitFeedback .
// @router /feedback/submit [POST]
func SubmitFeedback(ctx context.Context, c *app.RequestContext) {
	var err error
	var req show.SubmitFeedbackReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.FeedBackService.Submit(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}
