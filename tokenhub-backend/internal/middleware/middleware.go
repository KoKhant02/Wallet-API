package middleware

import (
	"net/http"
	"reflect"
	"runtime"
	"time"

	"go.uber.org/zap"
)

func NewZapLoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			handlerName := getHandlerName(next)

			logger.Info("Route Accessed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("handler", handlerName),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}

func getHandlerName(h http.Handler) string {
	if handlerFunc, ok := h.(http.HandlerFunc); ok {
		return runtime.FuncForPC(reflect.ValueOf(handlerFunc).Pointer()).Name()
	}
	return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
}
