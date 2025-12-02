package config

import (
	"fmt"
	"os"

	"github.com/spyder01/lilium-go/pkg/utils/env"
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

type LogConfig struct {
	ToFile       bool   `yaml:"toFile"`
	FilePath     string `yaml:"filePath"`
	ToStdout     bool   `yaml:"toStdout"`
	Prefix       string `yaml:"prefix"`
	Flags        int    `yaml:"flags"`
	DebugEnabled bool   `yaml:"debugEnabled"`
}

type EnvironmentConfig struct {
	EnableFile bool   `yaml:"enableFile"`
	FilePath   string `yaml:"filePath"`
}

func LoadEnv(path string) (*EnvironmentConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &EnvironmentConfig{}
	if err := yaml.Unmarshal([]byte(data), cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return cfg, nil
}

type LiliumConfig struct {
	Name      string             `yaml:"name"`
	Server    *ServerConfig      `yaml:"server"`
	Logger    *LogConfig         `yaml:"logger"`
	LogRoutes bool               `yaml:"logRoutes"`
	Env       *EnvironmentConfig `yaml:"env"`

	Extras map[string]any `yaml:",inline"` // store unknown fields here
}

func Load(path string) (*LiliumConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	expanded := env.ExpandEnvWithDefault(string(data))

	// --- First decode into a yaml.Node ---
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(expanded), &root); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// --- Second decode into typed struct ---
	cfg := &LiliumConfig{}
	if err := yaml.Unmarshal([]byte(expanded), cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// --- Extract unknown fields ---
	cfg.Extras = extractUnknownFields(&root, cfg)

	applyDefaults(cfg)
	return cfg, nil
}

func extractUnknownFields(root *yaml.Node, cfg *LiliumConfig) map[string]any {
	known := map[string]struct{}{
		"name":      {},
		"server":    {},
		"logger":    {},
		"logRoutes": {},
		"env":       {},
	}

	extras := make(map[string]any)
	if root.Kind != yaml.MappingNode {
		return extras
	}

	for i := 0; i < len(root.Content); i += 2 {
		key := root.Content[i].Value
		val := root.Content[i+1]
		if _, ok := known[key]; !ok {
			var v any
			_ = val.Decode(&v)
			extras[key] = v
		}
	}

	return extras
}

func GetExtra[T any](c *LiliumConfig, key string, out *T) error {
	extra, ok := c.Extras[key]
	if !ok {
		return fmt.Errorf("extra config not found: %s", key)
	}
	b, err := yaml.Marshal(extra)
	if err != nil {
		return fmt.Errorf("failed to marshal extra config: %w", err)
	}
	if err := yaml.Unmarshal(b, out); err != nil {
		return fmt.Errorf("failed to unmarshal extra config: %w", err)
	}
	return nil
}
