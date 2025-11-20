package controllers

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"slices"
	"strings"
)

type Middleware func(http.Handler) http.Handler

func createStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			x := middlewares[i]
			next = x(next)
		}
		return next
	}
}

func responseHeadersMiddleware(next http.Handler) http.Handler {
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Del("X-Powered-By")

		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-Workspace-Id",
		)

		origin := r.Header.Get("Origin")

		if slices.Contains(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func globalErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Caught Panic: %+v \n%s\n", err, string(debug.Stack()))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

var MiddlewareStack = createStack(responseHeadersMiddleware, globalErrorHandlerMiddleware)
