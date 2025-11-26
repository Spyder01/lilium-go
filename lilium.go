package lilium

import (
	"context"

	"github.com/spyder01/lilium-go/pkg/config"
	"github.com/spyder01/lilium-go/pkg/core"
	"github.com/spyder01/lilium-go/pkg/module"
)

type LiliumModule module.Module

type LiliumModulePriority module.Module

type AppContext core.Context

type RequestContext core.RequestContext

type LiliuTask core.LiliumTask

func LoadConfig(path string) *config.LiliumConfig {
	return core.LoadLiliumConfig(path)
}

func New(config *config.LiliumConfig, ctx_ context.Context) *core.Lilium {
	return core.New(config, ctx_)
}

func NewRouter(app *core.Context) *core.Router {
	return core.NewRouter(app)
}
