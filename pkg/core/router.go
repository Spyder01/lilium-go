package core

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	mux *chi.Mux
	app *Context
}

func NewRouter(app *Context) *Router {
	r := &Router{
		mux: chi.NewRouter(),
		app: app,
	}

	r.mux.Use(middleware.Recoverer)

	return r
}

func (r *Router) adapt(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := NewRequestContext(r.app, w, req)

		// extract chi params
		for _, key := range chi.RouteContext(req.Context()).URLParams.Keys {
			ctx.Params[key] = chi.URLParam(req, key)
		}

		if err := h(ctx); err != nil {
			// default error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (r *Router) GET(path string, h HandlerFunc) {
	r.mux.Get(path, r.adapt(h))
}

func (r *Router) POST(path string, h HandlerFunc) {
	r.mux.Post(path, r.adapt(h))
}

func (r *Router) PUT(path string, h HandlerFunc) {
	r.mux.Put(path, r.adapt(h))
}

func (r *Router) DELETE(path string, h HandlerFunc) {
	r.mux.Delete(path, r.adapt(h))
}

func (r *Router) PATCH(path string, h HandlerFunc) {
	r.mux.Patch(path, r.adapt(h))
}

func (r *Router) OPTIONS(path string, h HandlerFunc) {
	r.mux.Options(path, r.adapt(h))
}

func (r *Router) Group(fn func(g *Router)) {
	r.mux.Group(func(cr chi.Router) {
		gr := &Router{mux: cr.(*chi.Mux), app: r.app}
		fn(gr)
	})
}

func (r *Router) Use(mws ...Middleware) {
	for _, mw := range mws {
		r.mux.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				h := mw(func(ctx *RequestContext) error {
					next.ServeHTTP(w, req)
					return nil
				})

				ctx := NewRequestContext(r.app, w, req)
				_ = h(ctx)
			})
		})
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) SubRouter(prefix string) *Router {
	subMux := chi.NewRouter()

	// inherit global middlewares from parent chi mux
	// chi uses the same middleware stack for nested routers automatically.

	r.mux.Mount(prefix, subMux)

	return &Router{
		mux: subMux,
		app: r.app,
	}
}

func (r *Router) Static(route, dir string) {
	// normalize
	if route == "" {
		route = "/"
	}
	if !strings.HasPrefix(route, "/") {
		route = "/" + route
	}

	// file server handler
	fs := http.FileServer(http.Dir(dir))

	// special-case: root
	if route == "/" {
		// Serve everything and strip leading slash
		r.mux.Handle("/*", http.StripPrefix("/", fs))
		return
	}

	// Ensure trailing slash _redirect_ if user hits route exactly
	// e.g. "/static" -> "/static/"
	if route[len(route)-1] != '/' {
		r.mux.Get(route, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, route+"/", http.StatusMovedPermanently)
		})
	}

	// Serve static files under route/*
	pattern := route + "/*"
	fileHandler := http.StripPrefix(route+"/", fs)

	r.mux.Handle(pattern, fileHandler)
}
