package util

import (
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/basic"
)

func ParsePageOpt(p *basic.PaginationOptions) (skip int64, limit int64) {
	// 设置分页参数
	skip = int64(0)
	limit = int64(10) // 默认限制为10条数据

	if p.Page != nil && p.Limit != nil {
		skip = (*p.Page - 1) * *p.Limit
		limit = *p.Limit
	}
	return skip, limit
}
