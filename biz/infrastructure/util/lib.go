package util

import (
	"github.com/xh-polaris/essay-show/biz/infrastructure/util/log"
	"strconv"
	"sync"

	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/hertz/pkg/common/json"
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
