package apiv1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"tabsquare/db"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS tabsquare.customers (
    uuid VARCHAR(255) NOT NULL,
    customer_name VARCHAR(2048) NOT NULL,
    country VARCHAR(4096) NOT NULL,
    PRIMARY KEY (uuid)
);`

type Payload struct {
	Entities []*db.Customer `json:"entities,omitempty"`
	Message  string         `json:"message,omitempty"`
}

type healthCheckResponse struct {
	Status string `json:"status"`
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, host, port, dbname string) {
	log.Printf("Initalizing database with user '%s', host '%s', port '%s', dbname '%s'", user, host, port, dbname)
	connectionString :=
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.createTable()

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) createTable() {
	log.Printf("Executing db migrate")
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func (a *App) Run(addr string) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, a.Router)
	log.Fatal(http.ListenAndServe(addr, loggedRouter))
}

func (a *App) initializeRoutes() {
	log.Printf("Initializing routes")
	a.Router.HandleFunc("/api/v1/customers", a.ReturnAllCustomers).Methods("GET")
	a.Router.HandleFunc("/api/v1/customers", a.CreateNewCustomer).Methods("POST")
	a.Router.HandleFunc("/api/v1/customers/{uuid}", a.DeleteCustomer).Methods("DELETE")
	a.Router.HandleFunc("/api/v1/customers/{uuid}", a.ReturnSingleCustomer).Methods("GET")
	a.Router.HandleFunc("/health", a.HealthCheck).Methods("GET")

}

func writeJsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&data)
}

func (a *App) HealthCheck(w http.ResponseWriter, r *http.Request) {

	err := a.DB.Ping()
	if err != nil {
		log.Printf("DB conneciton is not ok %s", err.Error())
	} else {
		log.Printf("DB conneciton is ok")
	}

	if err != nil {
		data := &healthCheckResponse{Status: "Database unaccessible"}
		writeJsonResponse(w, http.StatusServiceUnavailable, data)
	} else {
		data := &healthCheckResponse{Status: "OK"}
		writeJsonResponse(w, http.StatusOK, data)
	}
}

func (a *App) ReturnAllCustomers(w http.ResponseWriter, r *http.Request) {

	var payload Payload

	customers, err := db.GetAllCustomerEntries(a.DB)
	if err != nil {
		log.Print(err.Error())
	}

	if len(customers) == 0 {
		payload.Message = err.Error()
	} else {
		payload.Entities = customers
	}
	writeJsonResponse(w, http.StatusOK, &payload)
}

func (a *App) ReturnSingleCustomer(w http.ResponseWriter, r *http.Request) {

	var payload Payload

	vars := mux.Vars(r)
	customer_uuid := vars["uuid"]

	if !IsValidUUID(customer_uuid) {
		payload.Message = fmt.Sprintf("uuid %s is invalid", customer_uuid)
		writeJsonResponse(w, http.StatusBadRequest, &payload)
		return
	}

	customer, err := db.GetCustomerEntry(a.DB, customer_uuid)
	if err != nil {
		log.Print(err.Error())
	}

	if customer != nil {
		payload.Entities = append(payload.Entities, customer)
		writeJsonResponse(w, http.StatusOK, &payload)
	} else {
		payload.Message = err.Error()
		writeJsonResponse(w, http.StatusNotFound, &payload)
	}
}

func (a *App) CreateNewCustomer(w http.ResponseWriter, r *http.Request) {
	var payload Payload

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		payload.Message = "Invalid json data provided"
		writeJsonResponse(w, http.StatusBadRequest, &payload)
		return
	}
	var customer *db.Customer
	json.Unmarshal(reqBody, &customer)

	if customer == nil {
		payload.Message = "Invalid json data provided"
		writeJsonResponse(w, http.StatusBadRequest, &payload)
		return
	}

	if customer.UUID == "" {
		id := uuid.New()
		customer.UUID = id.String()
	}

	if !IsValidUUID(customer.UUID) {
		payload.Message = fmt.Sprintf("uuid %s is invalid", customer.UUID)
		writeJsonResponse(w, http.StatusBadRequest, &payload)
		return
	}

	err = db.InsertCustomerEntry(a.DB, customer)
	if err != nil {
		payload.Message = err.Error()
		writeJsonResponse(w, http.StatusBadRequest, &payload)
	} else {
		writeJsonResponse(w, http.StatusOK, &customer)
	}
}

func (a *App) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	var payload Payload

	vars := mux.Vars(r)
	customer_uuid := vars["uuid"]

	if !IsValidUUID(customer_uuid) {
		payload.Message = fmt.Sprintf("uuid %s is invalid", customer_uuid)
		writeJsonResponse(w, http.StatusBadRequest, &payload)
		return
	}
	err := db.DeleteCustomerEntry(a.DB, customer_uuid)
	if err != nil {
		payload.Message = err.Error()
		writeJsonResponse(w, http.StatusNotFound, &payload)
	} else {
		payload.Message = fmt.Sprintf("Customer entry with uuid '%s' delete successfully", customer_uuid)
		writeJsonResponse(w, http.StatusOK, &payload)
	}
}
