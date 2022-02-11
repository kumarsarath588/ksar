package main

import (
	"log"
	"net/http"
	"os"
	"tabsquare/api"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func handleRequests() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/customers", api.ReturnAllCustomers).Methods("GET")
	r.HandleFunc("/customers", api.CreateNewCustomer).Methods("POST")
	r.HandleFunc("/customers/{uuid}", api.DeleteCustomer).Methods("DELETE")
	r.HandleFunc("/customers/{uuid}", api.ReturnSingleCustomer).Methods("GET")
	r.HandleFunc("/health", api.HealthCheck).Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}

func main() {
	handleRequests()
}
