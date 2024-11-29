package platform_sts

import (
	"github.com/google/wire"
	"github.com/xh-polaris/essay-show/biz/infrastructure/config"
	"github.com/xh-polaris/gopkg/kitex/client"
	sts "github.com/xh-polaris/service-idl-gen-go/kitex_gen/platform/sts/stsservice"
)

type IPlatformSts interface {
	sts.Client
}

type PlatformSts struct {
	sts.Client
}

var PlatformStsSet = wire.NewSet(
	NewPlatformSts,
	wire.Struct(new(PlatformSts), "*"),
	wire.Bind(new(IPlatformSts), new(*PlatformSts)),
)

func NewPlatformSts(config *config.Config) sts.Client {
	return client.NewClient(config.Name, "platform.sts", sts.NewClient)
}
