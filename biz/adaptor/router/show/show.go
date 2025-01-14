// Code generated by hertz generator. DO NOT EDIT.

package show

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	show "github.com/xh-polaris/essay-show/biz/adaptor/controller/show"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	{
		_essay := root.Group("/essay", _essayMw()...)
		_essay.POST("/evaluate", append(_essayevaluateMw(), show.EssayEvaluate)...)
		_essay.POST("/logs", append(_getevaluatelogsMw(), show.GetEvaluateLogs)...)
	}
	{
		_sts := root.Group("/sts", _stsMw()...)
		_sts.POST("/apply", append(_applysignedurlMw(), show.ApplySignedUrl)...)
		_sts.POST("/ocr", append(_ocrMw(), show.OCR)...)
		_sts.POST("/send_verify_code", append(_sendverifycodeMw(), show.SendVerifyCode)...)
	}
	{
		_user := root.Group("/user", _userMw()...)
		_user.GET("/info", append(_getuserinfoMw(), show.GetUserInfo)...)
		_user.POST("/sign_in", append(_signinMw(), show.SignIn)...)
		_user.POST("/sign_up", append(_signupMw(), show.SignUp)...)
		_user.POST("/update", append(_updateuserinfoMw(), show.UpdateUserInfo)...)
		_user.POST("/update_password", append(_updatepasswordMw(), show.UpdatePassword)...)
	}
}
