package core

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spyder01/lilium-go/pkg/config"
	"github.com/spyder01/lilium-go/pkg/logger"
)

type Lilium struct {
	Config        *config.LiliumConfig
	onStartTasks  []LiliumTask
	onStopTasks   []LiliumTask
	Lock          *sync.Mutex
	Logger        *logger.Logger
	Context       *Context
	isRunning     bool
	moduleManager *ModuleManager
}

func New(cfg *config.LiliumConfig, ctx_ context.Context) *Lilium {
	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		panic("Unable to instantiate logger")
	}

	app := &Lilium{
		Config:       cfg,
		onStartTasks: []LiliumTask{},
		onStopTasks:  []LiliumTask{},
		Lock:         &sync.Mutex{},
		Logger:       log,
		isRunning:    false,
	}

	ctx := &Context{
		app:       app,
		store:     make(map[string]any),
		Bus:       NewEventBus(),
		isRunning: false,
		mu:        sync.RWMutex{},
		Logger:    log,
		Ctx:       ctx_,
	}

	app.Context = ctx
	app.moduleManager = NewModuleManager(ctx)

	return app
}

const DEFUALT_LILIUM_CONFIG = "lilium.yaml"

func LoadLiliumConfig(path string) *config.LiliumConfig {
	cfg, err := config.Load(path)
	if err != nil {
		panic(fmt.Sprintf("Error while reading the lilium config at %s: %v\n", path, err))
	}

	return cfg
}

func (app *Lilium) OnStart(task LiliumTask) {
	if app.isRunning {
		panic("OnStart function can't be called once the lilium app is running.")
	}

	app.Lock.Lock()
	defer app.Lock.Unlock()

	app.onStartTasks = append(app.onStartTasks, task)
}

func (app *Lilium) OnStop(task LiliumTask) {
	app.Lock.Lock()
	defer app.Lock.Unlock()

	app.onStopTasks = append(app.onStopTasks, task)
}

func (app *Lilium) Start(router *Router) {
	err := app.moduleManager.InitAll()
	if err != nil {
		app.Logger.Errorf("Error while Initializing modules: %v", err)
		panic(err)
	}

	app.processCors(router.mux)
	if app.Config.LogRoutes {
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Server.Port),
		Handler: router,
	}

	app.Logger.Info("Mounting static files")
	for _, s := range app.Config.Server.Static {
		router.Static(s.Route, s.Directory)
	}
	app.Logger.Info("Mounted static files")

	// Start Modules
	app.Logger.Info("Starting all the modules...")
	err = app.moduleManager.StartAll()
	if err != nil {
		app.Logger.Errorf("Error while starting attached modules: %v", err)
		panic(err)
	}
	app.Logger.Info("Started all the attached modules...")

	// Run onStart tasks
	app.Logger.Info("Running onStart tasks...")
	for _, task := range app.onStartTasks {
		if err := task(app.Context); err != nil {
			app.Logger.Errorf("Error while running startup task: %v", err)
			panic(err)
		}
	}
	app.Logger.Info("Startup tasks complete.")

	// Start server in goroutine
	go func() {
		app.Logger.Infof("Listening on :%d", app.Config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Errorf("Server error: %v", err)
			panic(err)
		}
	}()

	// Wait for shutdown signal (Ctrl+C etc.)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.Logger.Info("Shutting down server...")

	// Gracefully shut down
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		app.Logger.Errorf("Server forced to shutdown: %v", err)
	}

	// Run stop tasks
	app.Logger.Info("Running onStop tasks...")
	for _, task := range app.onStopTasks {
		if err := task(app.Context); err != nil {
			app.Logger.Errorf("Error while running stop task: %v", err)
		}
	}

	app.Logger.Info("Stopping all the attached modules...")
	app.moduleManager.ShutdownAll()
	app.Logger.Info("Stopped all the attached modules...")

	// Close logger
	_ = app.Logger.Close()

	app.Logger.Info("Lilium shutdown complete.")
}

func (app *Lilium) UseModule(m Module) {
	app.moduleManager.Register(m)
}
