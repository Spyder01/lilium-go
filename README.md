# 🌸 **Lilium — A Fast, Elegant & Modular Web Framework for Go**

Lilium is a lightweight and high-performance web framework for Go — built around clarity, modularity, and developer experience.

It brings together:

* ⚡ **Chi-powered high-performance router**
* 🧩 **Modular App Container + Dependency Injection**
* 📄 **Smart YAML Configuration with ENV expansion**
* 🌱 **Automatic Defaults for common fields**
* 🧩 **Extensible config via Extras (for plugins/modules)**
* 🧵 **Asynchronous Zerolog logging**
* 🔔 **Built-in EventBus (pub/sub)**
* 🗂️ **Static File Serving**
* 🛡️ **Composable Middleware**
* 🧹 **Graceful Shutdown & Lifecycle Hooks**
* 🧪 **First-class Testability**

Designed to stay idiomatic to Go while providing a clean, modern developer experience.

---

# 🚀 Features

## 🚦 Router

A frictionless wrapper around Chi:

* Grouped & nested routes
* Typed `RequestContext`
* Middleware chaining
* Helpers for JSON / Text / HTML responses
* Centralized error handling
* Optional automatic route logging

Example:

```go
router.GET("/hello/{name}", func(c *core.RequestContext) error {
    return c.JSON(200, map[string]string{
        "message": "Hello " + c.Param("name"),
    })
})
```

---

## 🌐 Static File Serving

Declare static directories directly in config:

```yaml
server:
  static:
    - route: "/"
      directory: "./public"
```

Automatically mounted at startup.

---

# ⚙️ Smart Config System (YAML)

Load config with one line:

```go
cfg := config.Load("lilium.yaml")
```

### ✨ Includes:

| Feature                              | Status |
| ------------------------------------ | :----: |
| Environment variable expansion       |    ✔   |
| Fallback default values              |    ✔   |
| Unknown fields preserved for modules |    ✔   |
| Strongly typed configuration         |    ✔   |

---

## 🔄 Environment Variable Expansion

Supports `${VAR}` and `${VAR:default}`:

```yaml
server:
  port: ${PORT:8080}
```

| Scenario       | Result         |
| -------------- | -------------- |
| `PORT` exists  | use PORT value |
| `PORT` missing | use `8080`     |

---

## 🧠 Sane Defaults (Auto-Applied)

If omitted:

| Field                            | Default           |
| -------------------------------- | ----------------- |
| `name`                           | `"Lilium"`        |
| `server.port`                    | `8080`            |
| `server.cors.maxAge`             | `600` seconds     |
| Logger output                    | `toStdout = true` |
| Logger prefix                    | `"[Lilium] "`     |
| `env.enableFile`                 | `false`           |
| If `.env` enabled filePath empty | `.env`            |

This means you can start with only:

```yaml
server:
  port: 9000
```

→ Completely valid 🚀

---

## 🧩 Extensible Config (Extras)

Unknown YAML fields are stored in:

```go
cfg.Extras map[string]any
```

Used for modules:

```yaml
auth:
  provider: google
  tokenTTL: 3600
```

Module usage:

```go
type AuthConfig struct {
    Provider string `yaml:"provider"`
    TokenTTL int    `yaml:"tokenTTL"`
}

var auth AuthConfig
_ = cfg.GetExtra("auth", &auth)
```

This enables plugin systems and forward-compatible configuration.

---

# 🧵 Logging (Zerolog-based)

* Async writes
* File + STDOUT targets
* Debug mode
* Auto-flush on shutdown

```yaml
logger:
  toStdout: true
  prefix: "[MyApp] "
  debugEnabled: true
```

---

# 📡 EventBus

Simple in-process pub/sub:

```go
id, ch, _ := app.Context.Bus.Subscribe("notifications", 10)

go func() {
    for msg := range ch {
        fmt.Println("received:", msg)
    }
}()

app.Context.Bus.Publish("notifications", "hello world")
```

Perfect for background processing and modular integrations.

---

# 🧩 Dependency Injection / App Context

Lightweight DI for shared dependencies:

```go
app.Context.Provide("db", db)
db := app.Context.MustGet("db").(*sql.DB)
```

Per-request context includes logging + utilities.

---

# 🧹 Graceful Shutdown

On termination:

* Stop HTTP server cleanly
* Drain in-flight requests
* Flush logs
* Close EventBus
* Trigger module lifecycle hooks

---

# 📦 Installation

```sh
go get github.com/spyder01/lilium-go@latest
```

---

# 🏃 Quick Start

### Project Structure

```
.
├── lilium.yaml
├── main.go
└── public/
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

# 🧪 Testing

All router + RequestContext behavior is testable:

```go
req := httptest.NewRequest("GET", "/ping", nil)
rec := httptest.NewRecorder()

router.ServeHTTP(rec, req)

assert.Equal(t, 200, rec.Code)
```

---

# 🗺️ Roadmap

* [x] ENV var expansion in config
* [x] Unknown field `Extras` for modules
* [ ] Authentication (sessions + JWT)
* [ ] Built-in validators
* [ ] WebSockets
* [ ] Rate-limiting & caching middleware
* [ ] Auto OpenAPI generation
* [ ] CLI tooling (`lilium new`, scaffolding)
* [ ] Stronger DI capabilities

---

# ❤️ Contributing

PRs welcome!
Open an issue for ideas or bugs.

---

# 📄 License

MIT © 2025

