package module

import "github.com/spyder01/lilium-go/pkg/core"

type Module interface {
	Init(app *core.Lilium) error
	Start(app *core.Lilium) error
	Stop(app *core.Lilium) error
	Close(app *core.Lilium) error
}
