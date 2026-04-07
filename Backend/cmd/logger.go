package main

import (
	"log"
	"net/http"
)

func logger(next http.Handler) http.Handler {
	logger := log.New(log.Writer(), "", log.LstdFlags)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		logger.Printf("ip=%s - protocal=%s method=%s uri=%s", ip, proto, method, uri)

		next.ServeHTTP(w, r)
	})
}
