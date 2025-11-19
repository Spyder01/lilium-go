# 📦 **Lilium — Elegant, Fast, Modular Web Framework for Go**

Lilium is a lightweight, flexible, high-performance web framework for Go, built with:

* ⚡ **Chi-powered routing**
* 📦 **Modular application container**
* 🧩 **Powerful dependency injection context**
* 🗂️ **Static file serving**
* 🧵 **Asynchronous structured logging (Zerolog)**
* 🔀 **EventBus for in-app pub/sub**
* 🛡️ **Composable middleware system**
* 🧹 **Graceful shutdown**
* 🧪 **First-class testability**

Lilium aims to provide a clean, intuitive API while staying close to the Go standard library.

---

# ✨ Features

### 🚦 **Router**

* Chi-based router wrapped with a friendly API
* Route groups, sub-routers, middleware
* Structured request logging
* Centralized error handling
* Automatic JSON, HTML, text helpers
* Strongly typed request context

### 🌐 **Static File Serving**

Declare static directories in config:

```yaml
server:
  static:
    - route: "/"
      directory: "./public"
    - route: "/assets"
      directory: "./assets"
```

### ⚙️ **Config (YAML)**

Load environment-specific config files with:

```go
cfg := config.LoadLiliumConfig("lilium.yaml")
```

### 🧵 **Async Logging (Zerolog)**

* Non-blocking logger
* File + stdout support
* Log rotation-ready
* Structured logging (`InfoEvent()`)

### 📡 **EventBus**

* Per-topic publish/subscribe
* Buffered non-blocking events
* Graceful close
* Used internally for app communication

### 🧩 **App Context (DI Container)**

Store any value globally:

```go
app.Context.Provide("db", db)
db := app.Context.MustGet("db").(*sql.DB)
```

Also includes local request context via `RequestContext`.

### 🧹 **Graceful Shutdown**

* Catches SIGINT / SIGTERM
* Drains in-flight requests
* Runs `OnStop` tasks
* Flushes logger safely

---

# 📦 Installation

```sh
go get github.com/spyder01/lilium-go@latest
```

---

# 🚀 Quick Start

### 1. Load config

```go
cfg := core.LoadLiliumConfig("lilium.yaml")
app := core.New(cfg)
```

### 2. Create router

```go
router := core.NewRouter(app.Context)
```

### 3. Add middleware

```go
router.Use(middleware.RequestLogger(app.Logger))
```

### 4. Define routes

```go
router.GET("/hello/{name}", func(c *core.RequestContext) error {
    return c.JSON(200, map[string]string{
        "message": "Hello " + c.Param("name"),
    })
})
```

### 5. Serve static files

Already handled automatically from config:

```yaml
server:
  static:
    - route: "/"
      directory: "./public"
```

### 6. Start server

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

**main.go:**

```go
package main

import (
    "github.com/spyder01/lilium-go/pkg/core"
)

func main() {
    cfg := core.LoadLiliumConfig("lilium.yaml")
    app := core.New(cfg)
    router := core.NewRouter(app.Context)

    router.GET("/", func(c *core.RequestContext) error {
        return c.Text(200, "Welcome to Lilium!")
    })

    app.Start(router)
}
```

---

# 🛠 Configuration (lilium.yaml)

Example:

```yaml
name: "MyApp"

server:
  port: 8080

  cors:
    enabled: true
    origins: ["*"]
    allowedMetods: ["GET", "POST"]
    allowedHeaders: ["Content-Type"]
    allowCredentials: true
    maxAge: 3600

  static:
    - route: "/"
      directory: "./public"

logger:
  toFile: true
  filePath: "./logs/app.log"
  toStdout: true
  debugEnabled: true

logRoutes: true
```

---

# 🧩 Middleware Usage

```go
router.Use(RequestLogger(app.Logger))
router.Use(AuthMiddleware)
router.Use(CORSMiddleware(cfg.Cors))
```

Or route-level:

```go
router.GET("/admin", AuthMiddleware, adminHandler)
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

Lilium is fully testable thanks to:

* Chi router compatibility
* RequestContext abstractions
* Async logging flush
* EventBus controlled channels

Example:

```go
req := httptest.NewRequest("GET", "/hello/world", nil)
rec := httptest.NewRecorder()
router.ServeHTTP(rec, req)

assert.Equal(t, 200, rec.Code)
```

---

# 🧱 Roadmap

* [ ] Authentication module (JWT, sessions)
* [ ] Built-in validators
* [ ] More middleware: rate limiting, CSRF, caching
* [ ] WebSocket handler
* [ ] CLI tool (`lilium new`, `lilium migrate`)
* [ ] Auto-generated OpenAPI docs
* [ ] Improved DI with generics

---

# ❤️ Contributing

Pull requests are welcome!
Please open an issue if you’d like to propose a new feature or fix.

---

# 📄 License

MIT License.
