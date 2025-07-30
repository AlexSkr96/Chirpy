package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiConfig := apiConfig{}
	apiConfig.fileserverHits.Store(0)

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	okFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	serveMux.Handle("GET /healthz", okFunc)

	metricsFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hits: %v", apiConfig.fileserverHits.Load())))
	})
	serveMux.Handle("GET /metrics", metricsFunc)

	resetFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiConfig.fileserverHits.Store(0)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	serveMux.Handle("POST /reset", resetFunc)

	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
// 	cfg.fileserverHits.Add(1)
// 	return next
// }
