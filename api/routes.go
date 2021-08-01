package api

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// GetRoutes returns a muxer with all the routes
func (a *API) GetRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/auth", loggerPanic(httpMethod(http.MethodPost, http.HandlerFunc(a.auth))))
	mux.Handle("/user", loggerPanic(httpMethod(http.MethodGet, http.HandlerFunc(a.getUser))))

	// special endpoint that should be in a different server
	// and will be used as `jku` header
	mux.Handle("/jwks", loggerPanic(httpMethod(http.MethodGet, http.HandlerFunc(a.getJWKS))))

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
