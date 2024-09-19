package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	dbPath = "./database.json"
)

func main() {
	//Delete the db on server startup with the --debug flag for easy debugging
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	// BUG: Exits if there is no database
	if *dbg {
		err := os.Remove(dbPath)
		if err != nil {
			log.Fatal("Couldn't delete database file on startup")
		}
		log.Println("Database deleted.")
	}

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

	//Users
	mux.HandleFunc("POST /api/users", handlerPostUsers)
	// mux.HandleFunc("POST /api/login")

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
