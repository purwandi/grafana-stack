package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/purwandi/grafana-stack/src/pkg/logger"
	"go.uber.org/zap"
)

// Logger ...
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		start := time.Now()

		if err := next(c); err != nil {
			c.Error(err)
		}

		res := c.Response()
		stop := time.Now()

		p := req.URL.Path
		if p == "" {
			p = "/"
		}

		cl := req.Header.Get(echo.HeaderContentLength)
		if cl == "" {
			cl = "0"
		}
		bi, _ := strconv.Atoi(cl)

		logger.InfoWithContext(req.Context(), "response info",
			zap.String("time", time.Now().Format(time.RFC3339Nano)),
			zap.String("request_id", res.Header().Get(echo.HeaderXRequestID)),
			zap.String("remote_ip", c.RealIP()),
			zap.String("host", req.Host),
			zap.String("method", req.Method),
			zap.String("uri", req.RequestURI),
			zap.String("user_agent", req.UserAgent()),
			zap.Int("status", res.Status),
			zap.Int64("latency", int64(stop.Sub(start))),
			zap.String("latency_human", stop.Sub(start).String()),
			zap.Int64("bytes_in", int64(bi)),
			zap.Int64("bytes_out", res.Size),
		)

		return nil
	}
}
