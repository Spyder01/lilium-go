package module

import "github.com/spyder01/lilium-go/pkg/core"

type ModulePriority uint

type Module interface {
	Name() string
	Priority() ModulePriority
	Init(app *core.Context) error
	Start(app *core.Context) error
	Shutdown(app *core.Context) error
}
