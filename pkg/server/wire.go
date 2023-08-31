//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"github.com/galgotech/fermions-workflow/pkg/bus"
	"github.com/galgotech/fermions-workflow/pkg/setting"
)

var wireBasicSet = wire.NewSet(
	bus.Provide,
	New,
)

func Initialize(setting.Setting) (*Server, error) {
	wire.Build(wireBasicSet)
	return &Server{}, nil
}
