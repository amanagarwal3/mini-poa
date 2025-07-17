package main

import (
	"fmt"
	"log"
	"net/http"

	"mini-poa/db"
	"mini-poa/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db.Connect()

	r := mux.NewRouter()
	r.HandleFunc("/api/provision", handlers.CreateProvisionRequest).Methods("POST")
	r.HandleFunc("/api/status/{id}", handlers.GetProvisionStatus).Methods("GET")
	r.HandleFunc("/api/requests", handlers.GetAllProvisionRequests).Methods("GET")

	fmt.Println("🚀 Server running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
