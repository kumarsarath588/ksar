package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"tabsquare/db"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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

func writeJsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&data)
}

func db_connect() (*sql.DB, error) {
	db_client, err := sql.Open("mysql", "user:BcGH2Gj41J5VF1@tcp(10.46.142.201:3306)/tabsquare")

	if err != nil {
		log.Print(err.Error())
	}
	return db_client, nil
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {

	db_client, err := db_connect()
	if err != nil {
		log.Print(err.Error())
	}

	defer db_client.Close()

	err = db_client.Ping()
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

func ReturnAllCustomers(w http.ResponseWriter, r *http.Request) {

	var payload Payload

	db_client, err := db_connect()
	if err != nil {
		log.Print(err.Error())
	}

	defer db_client.Close()

	customers, err := db.GetAllCustomerEntries(db_client)
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

func ReturnSingleCustomer(w http.ResponseWriter, r *http.Request) {

	var payload Payload

	vars := mux.Vars(r)
	customer_uuid := vars["uuid"]

	if !IsValidUUID(customer_uuid) {
		payload.Message = fmt.Sprintf("uuid %s is invalid", customer_uuid)
		writeJsonResponse(w, http.StatusBadRequest, &payload)
		return
	}

	db_client, err := db_connect()
	if err != nil {
		log.Print(err.Error())
	}

	customer, err := db.GetCustomerEntry(db_client, customer_uuid)
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

func CreateNewCustomer(w http.ResponseWriter, r *http.Request) {
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

	db_client, err := db_connect()
	if err != nil {
		log.Print(err.Error())
	}

	err = db.InsertCustomerEntry(db_client, customer)
	if err != nil {
		payload.Message = err.Error()
		writeJsonResponse(w, http.StatusBadRequest, &payload)
	} else {
		writeJsonResponse(w, http.StatusOK, &customer)
	}
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	var payload Payload

	vars := mux.Vars(r)
	customer_uuid := vars["uuid"]

	if !IsValidUUID(customer_uuid) {
		payload.Message = fmt.Sprintf("uuid %s is invalid", customer_uuid)
		writeJsonResponse(w, http.StatusBadRequest, &payload)
		return
	}

	db_client, err := db_connect()
	if err != nil {
		log.Print(err.Error())
	}

	err = db.DeleteCustomerEntry(db_client, customer_uuid)
	if err != nil {
		payload.Message = err.Error()
		writeJsonResponse(w, http.StatusNotFound, &payload)
	} else {
		payload.Message = fmt.Sprintf("Customer entry with uuid '%s' delete successfully", customer_uuid)
		writeJsonResponse(w, http.StatusOK, &payload)
	}
}
