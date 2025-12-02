package core

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func (app *Lilium) processCors(r chi.Router) {
	if app.Config.Server == nil || app.Config.Server.Cors == nil {
		return
	}

	corsCfg := app.Config.Server.Cors

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsCfg.Origins,
		AllowedMethods:   corsCfg.AllowedMetods,
		AllowedHeaders:   corsCfg.AllowedHeaders,
		MaxAge:           int(corsCfg.MaxAge),
		ExposedHeaders:   corsCfg.ExposedHeaders,
		AllowCredentials: corsCfg.AllowCredentials,
	}))
}

func (app *Lilium) processApp() {
	if !app.Config.LogRoutes {
		return
	}
}

func (app *Lilium) processEnv() {
	envCfg := app.Config.Env
	if envCfg == nil || !envCfg.EnableFile {
		return
	}

	err := godotenv.Load(envCfg.FilePath)
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}
}

/*
func (app *Lilium) processDB() {
	if app.Config.Db == nil {
		app.Logger.Warn("No DB config provided. Skipping database initialization.")
		return
	}

	db := app.Config.Db

	host := db.Host
	if host == "" {
		app.Logger.Warn("No host assigned for database, assigning localhost.")
		host = "localhost"
	}

	port := db.Port
	if port == "" {
		app.Logger.Warn("No port assigned for database, assigning 5432.")
		port = "5432"
	}

	user := db.Username
	if user == "" {
		app.Logger.Warn("No username provided for database.")
	}

	password := db.Password
	dbName := db.DbName
	if dbName == "" {
		app.Logger.Warn("No database name provided for database.")
	}

	if db.Type == "" {
		app.Logger.Error("DB type not provided. Supported: postgres.")
		return
	}

	if db.Type != "postgres" {
		app.Logger.Errorf("Unsupported DB type: %s. Only 'postgres' is supported for now.", db.Type)
		return
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbName,
	)

	app.OnStart(func(app *Lilium) error {
		app.Logger.Infof("Using PostgreSQL at %s:%s (db: %s)", host, port, dbName)

		pool, err := pgxpool.New(context.Background(), dsn)
		if err != nil {
			return fmt.Errorf("❌ Failed to connect to Postgres: %v", err)
		}

		app.OnStop(func(app *Lilium) error {
			app.Logger.Info("Closing database")
			pool.Close()
			app.Logger.Info("Closed database")
			return nil
		})

		queries := db.New(pool)
		ctx_.SetQueries(queries)
		app.Logger.Info("Connected to the Database.")

		return nil
	})

}
*/
