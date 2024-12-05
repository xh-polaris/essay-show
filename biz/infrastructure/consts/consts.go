package consts

var PageSize int64 = 10

// 数据库相关
const (
	ID           = "_id"
	UserID       = "user_id"
	Status       = "status"
	CreateTime   = "create_time"
	DeleteStatus = 3
	EffectStatus = 0
	Phone        = "phone"
)

// http
const (
	Post                      = "POST"
	PlatformSignInUrl         = "https://api.xhpolaris.com/platform/auth/sign_in"
	PlatformSetPasswordUrl    = "https://api.xhpolaris.com/platform/auth/set_password"
	PlatformSendVerifyCodeUrl = "https://api.xhpolaris.com/platform/auth/send_verify_code"
	ContentTypeJson           = "application/json"
	CharSetUTF8               = "UTF-8"
	Beta                      = "beta"
	OpenApiCallUrl            = "https://api.xhpolaris.com/openapi/call/"
	BeeOCRUrl                 = "http://open.mifengjiaoyu.com/api/sc/image/ocr"
)

// 默认值
const (
	DefaultCount = 10
	AppId        = 14
)
