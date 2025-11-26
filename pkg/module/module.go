package module

import "github.com/spyder01/lilium-go/pkg/core"

type Module interface {
	Init(app *core.Lilium) error
	Start(app *core.Lilium) error
	Shutdown(app *core.Lilium) error
}
