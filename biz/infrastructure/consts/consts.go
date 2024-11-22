package consts

import (
	"errors"
)

const DefaultPageSize int64 = 10

var ErrNotAuthentication = errors.New("not authentication")
var ErrForbidden = errors.New("forbidden")
var PageSize int64 = 10

// 数据库相关
const (
	ID           = "_id"
	UserID       = "user_id"
	Status       = "status"
	CreateTime   = "create_time"
	DeleteStatus = 3
	EffectStatus = 0
)
