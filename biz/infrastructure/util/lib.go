package util

import (
	"errors"
	"github.com/xh-polaris/essay-show/biz/application/dto/essay/show"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	"github.com/xh-polaris/essay-show/biz/infrastructure/util/log"
	"google.golang.org/grpc/codes"
	"strconv"
	"sync"

	"encoding/json"
	"github.com/bytedance/gopkg/util/gopool"
)

func JSONF(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		log.Error("JSONF fail, v=%v, err=%v", v, err)
	}
	return string(data)
}

func ParseInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func ParallelRun(fns ...func()) {
	wg := sync.WaitGroup{}
	wg.Add(len(fns))
	for _, fn := range fns {
		fn := fn
		gopool.Go(func() {
			defer wg.Done()
			fn()
		})
	}
	wg.Wait()
}

func NonNullString(s *string, obj string) error {
	if s == nil || *s == "" {
		return consts.NewErrno(codes.Code(9001), errors.New(obj+"不能为空"))
	}
	return nil
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func Succeed(msg string) (*show.Response, error) {
	return &show.Response{
		Code: 0,
		Msg:  msg,
	}, nil
}

func Fail(code int64, msg string) *show.Response {
	return &show.Response{
		Code: code,
		Msg:  msg,
	}
}
