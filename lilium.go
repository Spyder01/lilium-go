package lilium

import (
	"context"

	"github.com/spyder01/lilium-go/pkg/config"
	"github.com/spyder01/lilium-go/pkg/core"
)

type (
	LiliumModule   = core.Module
	AppContext     = core.Context
	RequestContext = core.RequestContext
	LiliumTask     = core.LiliumTask
)

func LoadConfig(path string) *config.LiliumConfig {
	return core.LoadLiliumConfig(path)
}

func New(config *config.LiliumConfig, ctx_ context.Context) *core.Lilium {
	return core.New(config, ctx_)
}

func NewRouter(app *core.Context) *core.Router {
	return core.NewRouter(app)
}
