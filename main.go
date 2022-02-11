package main

import (
	"log"
	"net/http"
	"os"

	apiv1 "tabsquare/api"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func handleRequests() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/api/v1/customers", apiv1.ReturnAllCustomers).Methods("GET")
	r.HandleFunc("/api/v1/customers", apiv1.CreateNewCustomer).Methods("POST")
	r.HandleFunc("/api/v1/customers/{uuid}", apiv1.DeleteCustomer).Methods("DELETE")
	r.HandleFunc("/api/v1/customers/{uuid}", apiv1.ReturnSingleCustomer).Methods("GET")
	r.HandleFunc("/health", apiv1.HealthCheck).Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}

func main() {
	handleRequests()
}
