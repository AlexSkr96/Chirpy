package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/AlexSkr96/Chirpy/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
}

var cfg = apiConfig{}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cfg.fileserverHits.Store(0)
	cfg.DB = database.New(db)

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	okFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	serveMux.Handle("GET /api/healthz", okFunc)
	serveMux.Handle("POST /api/users", http.HandlerFunc(createUserHandler))
	serveMux.Handle("POST /api/validate_chirp", http.HandlerFunc(validateChirpHandler))

	metricsFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`
			<html>
				<body>
    				<h1>Welcome, Chirpy Admin</h1>
        			<p>Chirpy has been visited %d times!</p>
           		</body>
            </html>`,
			// "Hits: %v",
			cfg.fileserverHits.Load()),
		))
	})
	serveMux.Handle("GET /admin/metrics", metricsFunc)

	serveMux.Handle("POST /admin/reset", http.HandlerFunc(resetHandler))

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
