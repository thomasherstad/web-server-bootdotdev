package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	dbPath = "./database.json"
)

func main() {
	mux := http.NewServeMux()

	const filePathRoot = "."
	const port = "8080"

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	apiCfg := apiConfig{
		fileServerHits: 0,
	}

	//Fileserver
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))

	//---API---
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerResetMetrics)

	//Chirps
	mux.HandleFunc("POST /api/chirps", handlerPostChirps)
	mux.HandleFunc("GET /api/chirps", handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", handlerGetChirpById)

	mux.HandleFunc("POST /api/users", handlerPostUsers)

	//Admin
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerFileServerHits)

	fmt.Println("Server running...")
	log.Fatal(srv.ListenAndServe())

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJson(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
