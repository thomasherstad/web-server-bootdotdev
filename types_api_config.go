package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func (apiCfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func (apiCfg *apiConfig) handlerFileServerHits(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	html := fmt.Sprintf(`<html>
				<body>
					<h1>Welcome, Chirpy Admin</h1>
					<p>Chirpy has been visited %d times!</p>
				</body>
			</html>`, apiCfg.fileServerHits)
	w.Write([]byte(html))
}

func (apiCfg *apiConfig) middlewareResetMetrics() {
	apiCfg.fileServerHits = 0
}

func (apiCfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, req *http.Request) {
	apiCfg.middlewareResetMetrics()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset hits to 0"))
}
