package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"slices"
	"strings"
)

type Chain []func(http.Handler) http.Handler

func (c Chain) Then(h http.Handler) http.Handler {
	for _, fn := range slices.Backward(c) {
		h = fn(h)
	}

	return h
}

func (c Chain) ThenFunc(h http.HandlerFunc) http.Handler {
	return c.Then(h)
}

func Logger(next http.Handler) http.Handler {
	logger := log.New(log.Writer(), "", log.LstdFlags)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		logger.Printf("ip=%s - protocol=%s method=%s uri=%s", ip, proto, method, uri)
		next.ServeHTTP(w, r)
	})
}

func CheckCORS(next http.Handler) http.Handler {
	allowedMethods := []string{"GET", "DELETE", "PUT", "POST", "OPTIONS"}

	allowedOrigins := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
	}
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		allowedOrigins = strings.Split(origins, ",")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		if isPreflight(r) {
			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")

			method := r.Header.Get("Access-Control-Request-Method")
			if slices.Contains(allowedOrigins, origin) && slices.Contains(allowedMethods, method) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.WriteHeader(http.StatusOK)
				return
			}

			w.WriteHeader(http.StatusForbidden)
			return
		}

		if slices.Contains(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusForbidden)
	})
}

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-site")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		next.ServeHTTP(w, r)
	})
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 4096)
				n := runtime.Stack(buf, false)
				fmt.Printf("panic: %v\n%s", err, buf[:n])

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal Server Error"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func isPreflight(r *http.Request) bool {
	return r.Method == "OPTIONS" &&
		r.Header.Get("Origin") != "" &&
		r.Header.Get("Access-Control-Request-Method") != ""
}
