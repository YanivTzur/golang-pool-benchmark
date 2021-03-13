package main

import (
	"fmt"
	"go-pool-perf/server"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const listeningPort = 8080

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/basic-handler", server.BasicHandler).Methods(http.MethodPost)
	r.HandleFunc("/object-pool-handler", server.ObjectPoolHandler).Methods(http.MethodPost)
	r.HandleFunc("/bounded-pool-handler", server.BoundedPoolHandler).Methods(http.MethodPost)

	log.Printf("Listening on port %d", listeningPort)
	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", listeningPort),
	}
	log.Fatal(srv.ListenAndServe())
}
