package core

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type RequestContext struct {
	App        *Context // reference to global App Context
	Req        *http.Request
	Res        http.ResponseWriter
	Params     map[string]string // path params
	store      map[string]any    // per-request KV store
	formParsed bool
}

type HandlerFunc func(*RequestContext) error
type Middleware func(HandlerFunc) HandlerFunc

func NewRequestContext(app *Context, w http.ResponseWriter, r *http.Request) *RequestContext {
	return &RequestContext{
		App:    app,
		Req:    r,
		Res:    w,
		Params: make(map[string]string),
		store:  make(map[string]any),
	}
}

func (c *RequestContext) JSON(status int, v any) error {
	c.Res.Header().Set("Content-Type", "application/json")
	c.Res.WriteHeader(status)
	return json.NewEncoder(c.Res).Encode(v)
}

func (c *RequestContext) Text(status int, s string) error {
	c.Res.Header().Set("Content-Type", "text/plain")
	c.Res.WriteHeader(status)
	_, err := c.Res.Write([]byte(s))
	return err
}

func (c *RequestContext) BindJSON(v any) error {
	return json.NewDecoder(c.Req.Body).Decode(v)
}

func (c *RequestContext) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *RequestContext) Param(key string) string {
	return c.Params[key]
}

func (c *RequestContext) Set(key string, val any) {
	c.store[key] = val
}

func (c *RequestContext) Get(key string) (any, bool) {
	val, ok := c.store[key]
	return val, ok
}

func (c *RequestContext) Status(code int) {
	c.Res.WriteHeader(code)
}

func (c *RequestContext) Header(key, value string) {
	c.Res.Header().Set(key, value)
}

func (c *RequestContext) Redirect(status int, url string) error {
	http.Redirect(c.Res, c.Req, url, status)
	return nil
}

func (c *RequestContext) BodyBytes() ([]byte, error) {
	return io.ReadAll(c.Req.Body)
}

func (c *RequestContext) Path() string {
	return c.Req.URL.Path
}

func (c *RequestContext) Method() string {
	return c.Req.Method
}

func (c *RequestContext) ClientIP() string {
	if ip := c.Req.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := c.Req.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return c.Req.RemoteAddr
}

func (c *RequestContext) QueryAll(key string) []string {
	return c.Req.URL.Query()[key]
}

func (c *RequestContext) JSONError(status int, message string) error {
	return c.JSON(status, map[string]string{
		"error": message,
	})
}

func (c *RequestContext) HTML(status int, html string) error {
	c.Res.Header().Set("Content-Type", "text/html")
	c.Res.WriteHeader(status)
	_, err := c.Res.Write([]byte(html))
	return err
}

func (c *RequestContext) File(path string) error {
	http.ServeFile(c.Res, c.Req, path)
	return nil
}

func (c *RequestContext) Headers(h http.Header) {
	for k, v := range h {
		for _, vv := range v {
			c.Res.Header().Add(k, vv)
		}
	}
}

func (c *RequestContext) Stream(fn func(w io.Writer) error) error {
	c.Header("Content-Type", "text/event-stream")
	c.Res.WriteHeader(http.StatusOK)
	return fn(c.Res)
}

func (c *RequestContext) Deadline() (time.Time, bool) {
	return c.Req.Context().Deadline()
}

func (c *RequestContext) Done() <-chan struct{} {
	return c.Req.Context().Done()
}

func (c *RequestContext) Err() error {
	return c.Req.Context().Err()
}

func (c *RequestContext) ensureFormParsed() {
	if !c.formParsed {
		c.Req.ParseForm() // parses both Form + PostForm + Multipart
		c.formParsed = true
	}
}

func (c *RequestContext) Form(key string) string {
	c.ensureFormParsed()
	return c.Req.Form.Get(key)
}

func (c *RequestContext) PostForm(key string) string {
	c.ensureFormParsed()
	return c.Req.PostForm.Get(key)
}
