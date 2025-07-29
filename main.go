package main

import (
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	okFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	serveMux.Handle("/healthz/", okFunc)

	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
