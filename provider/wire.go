//go:build wireinject
// +build wireinject

package provider

import (
	"github.com/google/wire"
)

func NewProvider() (*adaptor.UserServer, error) {
	wire.Build(
		wire.Struct(new(adaptor.UserServer), "*"),
		AllProvider,
	)
	return nil, nil
}
