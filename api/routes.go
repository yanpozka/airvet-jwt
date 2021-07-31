package api

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func (a *API) GetRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/auth", loggerPanic(httpMethod(http.MethodPost, http.HandlerFunc(a.auth))))

	return mux
}

func loggerPanic(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			val := recover()
			if val == nil {
				return
			}
			log.Printf("Recovering from panic with value: '%v'", val)
			debug.PrintStack()
		}()

		start := time.Now()

		next.ServeHTTP(w, r)

		delta := time.Since(start)
		log.Printf("[%s]: %s | Time consumed: %v", r.Method, r.RequestURI, delta)
	})
}

func httpMethod(method string, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		//
		next.ServeHTTP(w, r)
	})
}
