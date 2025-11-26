package lilium

import (
	"github.com/spyder01/lilium-go/pkg/config"
	"github.com/spyder01/lilium-go/pkg/core"
	"github.com/spyder01/lilium-go/pkg/module"
)

type LiliumModule module.Module

type LiliumModulePriority module.Module

type AppContext core.Context

type RequestContext core.RequestContext

func LoadConfig(path string) *config.LiliumConfig {
	return core.LoadLiliumConfig(path)
}

func New(config *config.LiliumConfig) *core.Lilium {
	return core.New(config)
}

func NewRouter(app *core.Context) *core.Router {
	return core.NewRouter(app)
}
