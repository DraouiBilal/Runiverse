package routing

import (
	"log"
	"net/http"
	"time"
)

type middlewares struct {
	Logging Middleware
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

var Middlewares *middlewares = &middlewares{
	Logging: loggingMiddleware,
}
