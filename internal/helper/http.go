package helper

import (
	"net/http"
	"runtime/debug"

	"keight/internal/logging"
)

// Recoverer is recover panic http middleware
func Recoverer(log logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					var (
						err error
						ok  bool
					)
					err, ok = rvr.(error)
					if !ok {
						err = errUnknown
					}
					log.WithError(err).Error(string(debug.Stack()))
					WrapError(w, err, http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
