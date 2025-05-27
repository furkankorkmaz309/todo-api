package handlers

import (
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/furkankorkmaz309/todo-api/internal/app"
)

func LogRequest(app *app.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ip, _, err := net.SplitHostPort(r.RemoteAddr)

			if err != nil {
				ip = r.RemoteAddr
			}

			next.ServeHTTP(w, r)

			duration := time.Since(start)
			ms := float64(duration.Microseconds()) / 1000

			app.InfoLog.Printf("%v - %v %.2fms from %v", r.Method, r.URL, ms, ip)
		})
	}
}

func RecoverPanic(app *app.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					w.Header().Set("Connection", "close")
					err := fmt.Errorf(" Panic : %v\nRequest : %v - %v from %v\nStack Trace : %v", err, r.Method, r.URL, r.RemoteAddr, debug.Stack())
					respondError(w, app.ErrorLog, http.StatusInternalServerError, err.Error(), err)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func LimitRequest(app *app.App) func(next http.Handler) http.Handler {
	const maxRequest = 60
	requestsLeft := maxRequest
	firstRequest := time.Now()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limitSecond := time.Duration(60) * time.Second

			duration := time.Since(firstRequest)
			if requestsLeft > 0 {
				requestsLeft--
			} else if duration < limitSecond && requestsLeft <= 0 {
				err := fmt.Errorf("much request in a minute")
				respondError(w, app.ErrorLog, http.StatusTooManyRequests, err.Error(), err)
				return
			} else if duration >= limitSecond {
				firstRequest = time.Now()
				requestsLeft = maxRequest
			}

			next.ServeHTTP(w, r)
		})
	}
}
