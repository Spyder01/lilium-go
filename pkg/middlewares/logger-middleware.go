package middlewares

import (
	"net/http"
	"time"

	"github.com/spyder01/lilium-go/pkg/core"
	"github.com/spyder01/lilium-go/pkg/logger"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(p []byte) (int, error) {
	if rr.status == 0 {
		rr.status = http.StatusOK
	}
	n, err := rr.ResponseWriter.Write(p)
	rr.size += n
	return n, err
}

func RequestLoggingMiddleware(l *logger.Logger) core.Middleware {
	return func(next core.HandlerFunc) core.HandlerFunc {
		return func(c *core.RequestContext) error {
			start := time.Now()

			// Wrap the ResponseWriter to capture status + size
			rr := &responseRecorder{ResponseWriter: c.Res}
			c.Res = rr

			err := next(c)

			duration := time.Since(start)
			status := rr.status
			if status == 0 {
				status = http.StatusOK
			}

			// Logging fields
			event := l.InfoEvent().
				Str("method", c.Method()).
				Str("path", c.Path()).
				Int("status", status).
				Int("size", rr.size).
				Dur("duration", duration).
				Str("ip", c.ClientIP()).
				Str("user_agent", c.Req.UserAgent())

			if err != nil {
				event = event.Err(err)
			}

			event.Msg("request completed")

			return err
		}
	}
}
