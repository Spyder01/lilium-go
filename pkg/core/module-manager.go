package core

import (
	"fmt"
	"sort"
	"sync"
)

type Module interface {
	Name() string
	Priority() uint // Lower number = earlier init, later shutdown
	Init(app *Context) error
	Start(app *Context) error
	Shutdown(app *Context) error
}

type ModuleManager struct {
	modules []Module
	mu      sync.Mutex
	started bool
	app     *Context
}

func NewModuleManager(app *Context) *ModuleManager {
	return &ModuleManager{
		app:     app,
		modules: make([]Module, 0),
	}
}

func (m *ModuleManager) Register(module Module) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.app.Logger.Infof("Registering module: %s", module.Name())
	m.modules = append(m.modules, module)
}

// Sort modules by priority
func (m *ModuleManager) sortByPriority() {
	sort.Slice(m.modules, func(i, j int) bool {
		return m.modules[i].Priority() < m.modules[j].Priority()
	})
}

func (m *ModuleManager) InitAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sortByPriority()

	m.app.Logger.Info("Initializing modules...")

	for _, module := range m.modules {
		m.app.Logger.Infof("→ Init %s", module.Name())
		if err := module.Init(m.app); err != nil {
			return fmt.Errorf("init failed for %s: %w", module.Name(), err)
		}
	}

	return nil
}

func (m *ModuleManager) StartAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.started {
		return fmt.Errorf("modules already started")
	}
	m.started = true

	m.sortByPriority()

	m.app.Logger.Info("Starting modules...")

	for _, module := range m.modules {
		m.app.Logger.Infof("→ Start %s", module.Name())
		if err := module.Start(m.app); err != nil {
			return fmt.Errorf("start failed for %s: %w", module.Name(), err)
		}
	}

	return nil
}

func (m *ModuleManager) ShutdownAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.started {
		m.app.Logger.Warn("Shutdown called before Start")
	}

	m.app.Logger.Info("Shutting down modules...")

	// Shutdown in reverse order of start
	for i := len(m.modules) - 1; i >= 0; i-- {
		mod := m.modules[i]
		m.app.Logger.Infof("→ Shutdown %s", mod.Name())
		if err := mod.Shutdown(m.app); err != nil {
			m.app.Logger.Errorf("shutdown error for %s: %v", mod.Name(), err)
		}
	}

	m.started = false
}
