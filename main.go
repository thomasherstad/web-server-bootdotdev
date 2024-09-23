package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"web-server-bootdotdev/internal/database"

	"github.com/joho/godotenv"
)

const (
	dbPath = "./database.json"
)

//TODO:
//	- Add respond with error where necessary
//	- Handle capitalization in emails

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

	//load environment
	godotenv.Load()

	mux := http.NewServeMux()

	const filePathRoot = "."
	const port = "8080"

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileServerHits: 0,
		DB:             db,
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	//Fileserver
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))

	//---API---
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerResetMetrics)

	//Chirps
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpById)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirpById)

	//Users
	mux.HandleFunc("POST /api/users", apiCfg.handlerPostUsers)
	mux.HandleFunc("PUT /api/users", apiCfg.HandlerUserUpdate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandlerUserRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandlerUserRevoke)

	//Polka
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.HandlerPolkaWebhooks)

	//Admin
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerFileServerHits)

	fmt.Println("Server running...")
	log.Fatal(srv.ListenAndServe())

}
