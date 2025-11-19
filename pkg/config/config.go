package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type CorsConfig struct {
	Enabled          bool     `yaml:"enabled"`
	Origins          []string `yaml:"origins"`
	AllowedMetods    []string `yaml:"allowedMetods"`
	AllowedHeaders   []string `yaml:"allowedHeaders"`
	ExposedHeaders   []string `yaml:"exposedHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
	MaxAge           uint     `yaml:"maxAge"`
}

type StaticConfig struct {
	Route     string `yaml:"route"`     // e.g. "/static" or "/"
	Directory string `yaml:"directory"` // e.g. "./public"
}

type ServerConfig struct {
	Port   uint           `yaml:"port"`
	Cors   *CorsConfig    `yaml:"cors"`
	Static []StaticConfig `yaml:"static"` // <-- Add this
}

type DBConfig struct {
	Type       string `yaml:"type"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	DbName     string `yaml:"dbName"`
	Migrations string `yaml:"migrations"`
	Queries    string `yaml:"queries"`
}

type LogConfig struct {
	ToFile       bool   `yaml:"toFile"`
	FilePath     string `yaml:"filePath"`
	ToStdout     bool   `yaml:"toStdout"`
	Prefix       string `yaml:"prefix"`
	Flags        int    `yaml:"flags"`
	DebugEnabled bool   `yaml:"debugEnabled"`
}

type LiliumConfig struct {
	Name      string        `yaml:"name"`
	Server    *ServerConfig `yaml:"server"`
	Db        *DBConfig     `yaml:"db"`
	Logger    *LogConfig    `yaml:"logger"`
	LogRoutes bool          `yaml:"logRoutes"`
}

func Load(path string) (*LiliumConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg LiliumConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &cfg, nil
}
