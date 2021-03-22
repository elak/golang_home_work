package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/logger"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := newLoggingResponseWriter(w)

		reqStart := time.Now()

		next.ServeHTTP(lrw, r)

		reqLatency := time.Since(reqStart)

		statusCode := lrw.statusCode
		logMsg := fmt.Sprintf("%s [%s] %s %s %s %d %d %s", r.RemoteAddr, reqStart.Format("02/Jan/2006 15:04:05 -0700"), r.Method, r.RequestURI, r.Proto, statusCode, reqLatency, r.UserAgent())

		logger.Info(logMsg)
	})
}
