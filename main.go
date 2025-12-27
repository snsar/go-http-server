package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHit atomic.Int32
}

const metricsTemplate = `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHit.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileServerHit.Load()

	w.Header().Set("Content-Type", "text/html")

	html := fmt.Sprintf(metricsTemplate, hits)
	_, err := w.Write([]byte(html))
	if err != nil {
		return
	}
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHit.Store(0)

	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	apiCfg := &apiConfig{}
	fileServer := http.FileServer(http.Dir(filepathRoot))
	prefixStrip := http.StripPrefix("/app/", fileServer)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(prefixStrip))

	// Metrics
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	// Reset
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetrics)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("charset", "utf-8")

	})
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
