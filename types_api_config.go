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
		// What should happen here for each request?
		apiCfg.fileServerHits++
		// How do we call the next handler?
		next.ServeHTTP(w, r)
	})
}

func (apiCfg *apiConfig) handlerFileServerHits(w http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", apiCfg.fileServerHits)))
}

func (apiCfg *apiConfig) middlewareResetMetrics() {
	apiCfg.fileServerHits = 0
}

func (apiCfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset"))
	apiCfg.middlewareResetMetrics()
}
