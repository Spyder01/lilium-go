package config

func applyDefaults(cfg *LiliumConfig) {
	if cfg.Name == "" {
		cfg.Name = "Lilium" // Default app name if none provided
	}

	// ---------- Server ----------
	if cfg.Server == nil {
		cfg.Server = &ServerConfig{}
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}

	// Static array optional → do not override if empty

	// ---------- CORS ----------
	if cfg.Server.Cors == nil {
		cfg.Server.Cors = &CorsConfig{}
	}

	// Only override if NOT explicitly defined
	// Enabled default = false → no change needed

	if cfg.Server.Cors.MaxAge == 0 {
		cfg.Server.Cors.MaxAge = 600 // seconds
	}

	// ---------- Logger ----------
	if cfg.Logger == nil {
		cfg.Logger = &LogConfig{}
	}

	if cfg.Logger.Prefix == "" {
		cfg.Logger.Prefix = "[Lilium] "
	}

	// If no output target specified → default stdout
	if !cfg.Logger.ToFile && !cfg.Logger.ToStdout {
		cfg.Logger.ToStdout = true
	}

	if cfg.Env == nil {
		cfg.Env = &EnvironmentConfig{
			EnableFile: false,
		}
	}

	if cfg.Env.EnableFile && cfg.Env.FilePath == "" {
		cfg.Env.FilePath = ".env"
	}

}
