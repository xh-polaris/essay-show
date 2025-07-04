package consts

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Errno struct {
	err  error
	code codes.Code
}

// GRPCStatus 实现 GRPCStatus 方法
func (en *Errno) GRPCStatus() *status.Status {
	return status.New(en.code, en.err.Error())
}

// 实现 Error 方法
func (en *Errno) Error() string {
	return en.err.Error()
}

// NewErrno 创建自定义错误
func NewErrno(code codes.Code, err error) *Errno {
	return &Errno{
		err:  err,
		code: code,
	}
}

// 定义常量错误
var (
	ErrForbidden         = NewErrno(codes.PermissionDenied, errors.New("forbidden"))
	ErrNotAuthentication = NewErrno(codes.Code(1000), errors.New("not authentication"))
	ErrSignUp            = NewErrno(codes.Code(1001), errors.New("注册失败，请重试"))
	ErrSignIn            = NewErrno(codes.Code(1002), errors.New("登录失败，请先注册或重试"))
	ErrInSufficientCount = NewErrno(codes.Code(1003), errors.New("剩余调用次数不足，请充值或联系管理员添加"))
	ErrRepeatedSignUp    = NewErrno(codes.Code(1004), errors.New("该手机号已注册"))
	ErrOCR               = NewErrno(codes.Code(1005), errors.New("OCR识别失败，请重试"))
	ErrNotSignUp         = NewErrno(codes.Code(1006), errors.New("请确认手机号已注册"))
	ErrSend              = NewErrno(codes.Code(1007), errors.New("发送验证码失败，请重试"))
	ErrVerifyCode        = NewErrno(codes.Code(1008), errors.New("验证码错误"))
	ErrDailyAttend       = NewErrno(codes.Code(1009), errors.New("签到失败"))
	ErrRepeatDailyAttend = NewErrno(codes.Code(1010), errors.New("一天只能签到一次"))
	ErrRepeatInvitation  = NewErrno(codes.Code(1011), errors.New("已填写过邀请码"))
	ErrInvitation        = NewErrno(codes.Code(1012), errors.New("填写邀请码失败，请重试"))
	ErrGetInvitation     = NewErrno(codes.Code(1013), errors.New("获取邀请码失败，请重试"))
	ErrExerciseTimeout   = NewErrno(codes.Code(1014), errors.New("生成练习超时"))
	ErrExercise          = NewErrno(codes.Code(1015), errors.New("生成练习失败"))
)

// ErrInvalidParams 调用时错误
var (
	ErrInvalidParams = NewErrno(codes.InvalidArgument, errors.New("参数错误"))
	ErrCall          = NewErrno(codes.Unknown, errors.New("调用接口失败，请重试"))
	ErrOneCall       = NewErrno(codes.Code(3001), errors.New("同一时刻仅可以批改一篇作文, 请等待上一篇作文批改结束"))
)

// 数据库相关错误
var (
	ErrNotFound        = NewErrno(codes.NotFound, errors.New("not found"))
	ErrInvalidObjectId = NewErrno(codes.InvalidArgument, errors.New("无效的id "))
	ErrUpdate          = NewErrno(codes.Code(2001), errors.New("更新失败"))
)
