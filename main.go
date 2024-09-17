package main

import (
	"fmt"
	"log"
	"net/http"
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

	//Other
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerFileServerHits)
	mux.HandleFunc("/reset", apiCfg.handlerResetMetrics)

	fmt.Println("Server running...")
	log.Fatal(srv.ListenAndServe())

}
