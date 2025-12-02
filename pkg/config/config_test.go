package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, name string, content string) string {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return tmp
}

func TestLoadConfig_WithEnvExpansion(t *testing.T) {
	// Arrange
	os.Setenv("PUBLIC_DIR", "/srv/app")
	os.Setenv("PORT", "9090")
	defer os.Unsetenv("PUBLIC_DIR")
	defer os.Unsetenv("PORT")

	yamlContent := `
name: "TestApp"
server:
  port: ${PORT:8080}
  static:
    - route: "/"
      directory: "${PUBLIC_DIR:./public}"
logger:
  toStdout: true
`
	cfgFile := writeTempFile(t, "config.yaml", yamlContent)

	// Act
	cfg, err := Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Assert
	if cfg.Name != "TestApp" {
		t.Errorf("Expected name TestApp, got %s", cfg.Name)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Expected port=9090, got %d", cfg.Server.Port)
	}

	if len(cfg.Server.Static) == 0 || cfg.Server.Static[0].Directory != "/srv/app" {
		t.Errorf("Expected static.directory=/srv/app, got %+v", cfg.Server.Static)
	}

	if cfg.Logger == nil || !cfg.Logger.ToStdout {
		t.Error("Expected logger.toStdout=true")
	}
}

func TestLoadConfig_WithDefaults(t *testing.T) {
	// Arrange: missing port + static should trigger defaults
	yamlContent := `
name: "AppNoServer"
logger:
  toStdout: true
`
	cfgFile := writeTempFile(t, "config.yaml", yamlContent)

	// no PORT in env
	os.Unsetenv("PORT")

	// Act
	cfg, err := Load(cfgFile)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	// Assert defaults
	if cfg.Server == nil {
		t.Fatal("Expected Server struct initialized")
	}
	if cfg.Server.Port == 0 {
		t.Error("Expected default port applied, got 0")
	}
	if cfg.Logger == nil {
		t.Fatal("Logger should not be nil")
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := Load("does_not_exist.yaml")
	if err == nil {
		t.Fatal("Expected error for missing config file")
	}
}
