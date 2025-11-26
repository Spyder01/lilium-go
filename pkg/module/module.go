package module

import "github.com/spyder01/lilium-go/pkg/core"

type ModulePriority uint

type Module interface {
	Name() string
	Priority() ModulePriority
	Init(app *core.Lilium) error
	Start(app *core.Lilium) error
	Shutdown(app *core.Lilium) error
}
