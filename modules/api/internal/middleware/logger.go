package middleware

import (
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

// Logger returns a request logging middleware
// taken and adapted from https://github.com/chi-middleware/logrus-logger
func Logger(logger logrus.FieldLogger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					remoteIP = r.RemoteAddr
				}
				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}
				fields := logrus.Fields{
					"status_code":      ww.Status(),
					"bytes":            ww.BytesWritten(),
					"duration":         int64(time.Since(t1)),
					"duration_display": time.Since(t1).String(),
					"remote_ip":        remoteIP,
					"proto":            r.Proto,
					"method":           r.Method,
				}
				if len(reqID) > 0 {
					fields["request_id"] = reqID
				}
				logger.WithFields(fields).Logf(logrus.InfoLevel, "%s://%s%s", scheme, r.Host, r.RequestURI)
			}()

			h.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
