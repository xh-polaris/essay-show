package provider

import (
	"github.com/google/wire"
)

var AllProvider = wire.NewSet(
	ControllerSet,
	ApplicationSet,
	InfrastructureSet,
)

var ControllerSet = wire.NewSet()

var ApplicationSet = wire.NewSet()

var InfrastructureSet = wire.NewSet()

var MapperSet = wire.NewSet()

var TransactionSet = wire.NewSet()
