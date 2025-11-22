# 🌸 **Lilium — A Fast, Elegant & Modular Web Framework for Go**

Lilium is a lightweight, flexible, developer-friendly web framework for Go, built around clarity, modularity, and performance.

It brings together:

* ⚡ **Chi-powered high-performance router**
* 🧩 **Modular Application Container + Dependency Injection**
* 📄 **Simple YAML Configuration**
* 🧵 **Asynchronous Zerolog-based Logging**
* 🔔 **Built-in EventBus (pub/sub)**
* 🗂️ **Static File Serving**
* 🛡️ **Composable Middleware System**
* 🧹 **Graceful Shutdown & Lifecycle Hooks**
* 🧪 **First-class Testability**

Lilium aims to stay as close as possible to idiomatic Go while providing a clean, modern developer experience.

---

# 🚀 Features

## 🚦 **Router**

A zero-friction wrapper around Chi:

* Route groups & nesting
* Middleware chaining
* Centralized error handling
* Typed `RequestContext`
* Built-in helpers (`JSON`, `HTML`, `Text`, `Param`…)
* Optional automatic route logging
* Fully testable

---

## 🌐 **Static File Serving**

Declare static directories directly in your config:

```yaml
server:
  static:
    - route: "/"
      directory: "./public"
    - route: "/assets"
      directory: "./assets"
```

Automatically mounted at startup.

---

## ⚙️ **Config System (YAML)**

Load a strongly typed config:

```go
cfg := config.Load("lilium.yaml")
```

Supports:

* Server settings
* CORS
* Static files
* Logging
* App metadata

Detailed structure below.

---

## 🧵 **Asynchronous Logging**

Based on **Zerolog**, non-blocking, and configurable via YAML:

* File + STDOUT support
* Debug mode
* Prefix + flags
* Buffered writes
* Flush on shutdown

---

## 📡 **EventBus**

A simple, fast, in-process pub/sub system:

* Per-topic buffered channels
* Non-blocking publish
* Safe graceful close
* Perfect for background tasks and modular packages

---

## 🧩 **Dependency Injection / App Context**

Lilium provides a global DI container:

```go
app.Context.Provide("db", db)
db := app.Context.MustGet("db").(*sql.DB)
```

Plus a **per-request context** giving:

* Path/URL params
* Query params
* Body helpers
* Logging hooks
* Shared app dependencies

---

## 🧹 **Graceful Shutdown**

* Handles `SIGINT`, `SIGTERM`
* Drains in-flight requests
* Flushes logs
* Runs `OnStart` and `OnStop` tasks
* Shuts down EventBus cleanly

---

# 📦 Installation

```sh
go get github.com/spyder01/lilium-go@latest
```

---

# 🚀 Quick Start

### 1. Load configuration

```go
cfg := config.Load("lilium.yaml")
```

### 2. Initialize the application

```go
app := core.New(cfg)
```

### 3. Create a router

```go
router := core.NewRouter(app.Context)
```

### 4. Define routes

```go
router.GET("/hello/{name}", func(c *core.RequestContext) error {
    return c.JSON(200, map[string]string{
        "message": "Hello " + c.Param("name"),
    })
})
```

### 5. Start server

```go
app.Start(router)
```

---

# 📁 Example Project Structure

```
.
├── lilium.yaml
├── main.go
└── public/
    └── index.html
```

**main.go**

```go
package main

import (
    "github.com/spyder01/lilium-go/pkg/core"
    "github.com/spyder01/lilium-go/pkg/config"
)

func main() {
    cfg := config.Load("lilium.yaml")
    app := core.New(cfg)

    router := core.NewRouter(app.Context)
    router.GET("/", func(c *core.RequestContext) error {
        return c.Text(200, "Welcome to Lilium!")
    })

    app.Start(router)
}
```

---

# ⚙️ Configuration Reference

Below is the full YAML structure Lilium supports.

### **`LiliumConfig`**

```yaml
name: "MyApp"

server:
  port: 8080

  cors:
    enabled: true
    origins: ["*"]
    allowedMetods: ["GET", "POST"]
    allowedHeaders: ["Content-Type"]
    exposedHeaders: []
    allowCredentials: true
    maxAge: 3600

  static:
    - route: "/"
      directory: "./public"

logger:
  toFile: true
  filePath: "./logs/app.log"
  toStdout: true
  prefix: ""
  flags: 0
  debugEnabled: true

logRoutes: true
```

---

# 🔧 Config Structures

### `ServerConfig`

```go
type ServerConfig struct {
    Port   uint           `yaml:"port"`
    Cors   *CorsConfig    `yaml:"cors"`
    Static []StaticConfig `yaml:"static"`
}
```

### `CorsConfig`

```go
type CorsConfig struct {
    Enabled          bool     `yaml:"enabled"`
    Origins          []string `yaml:"origins"`
    AllowedMetods    []string `yaml:"allowedMetods"`
    AllowedHeaders   []string `yaml:"allowedHeaders"`
    ExposedHeaders   []string `yaml:"exposedHeaders"`
    AllowCredentials bool     `yaml:"allowCredentials"`
    MaxAge           uint     `yaml:"maxAge"`
}
```

### `StaticConfig`

```go
type StaticConfig struct {
    Route     string `yaml:"route"`
    Directory string `yaml:"directory"`
}
```

### `LogConfig`

```go
type LogConfig struct {
    ToFile       bool   `yaml:"toFile"`
    FilePath     string `yaml:"filePath"`
    ToStdout     bool   `yaml:"toStdout"`
    Prefix       string `yaml:"prefix"`
    Flags        int    `yaml:"flags"`
    DebugEnabled bool   `yaml:"debugEnabled"`
}
```

---

# 🔔 EventBus Example

```go
id, ch, _ := app.Context.Bus.Subscribe("notifications", 10)

go func() {
    for msg := range ch {
        fmt.Println("received:", msg)
    }
}()

app.Context.Bus.Publish("notifications", "hello world")
```

---

# 🧪 Testing

`RequestContext`, router, and logger are fully testable.

```go
req := httptest.NewRequest("GET", "/hello/world", nil)
rec := httptest.NewRecorder()

router.ServeHTTP(rec, req)

assert.Equal(t, 200, rec.Code)
```

---

# 🛣️ Roadmap

* [ ] Authentication (JWT, sessions)
* [ ] Built-in validators
* [ ] More middleware (rate-limiting, CSRF, caching)
* [ ] WebSockets
* [ ] CLI tool (`lilium new`, `lilium migrate`)
* [ ] Auto OpenAPI generation
* [ ] Generic DI improvements

---

# ❤️ Contributing

PRs welcome!
Please open an issue for features, bugs, or proposals.

---

# 📄 License

MIT License.
