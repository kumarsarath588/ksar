// main_test.go

package main_test

import (
	"fmt"
	apiv1 "ksar/api"
	"log"
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

var a apiv1.App

func TestMain(m *testing.M) {
	a = apiv1.App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_PORT"),
		os.Getenv("APP_DB_NAME"))

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM ksar.customers")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS ksar.customers 
(
    uuid VARCHAR(255) NOT NULL,
    customer_name VARCHAR(2048) NOT NULL,
    country VARCHAR(4096) NOT NULL,
    PRIMARY KEY (uuid)
);`

// tom: next functions added later, these require more modules: net/http net/http/httptest
func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/customers", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "no customer records found" {
		t.Errorf("Expected the 'message' key of the response to be set to 'no customer records found'. Got '%s'", m["message"])
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentCustomer(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/customers/9ffafd16-9ece-11ec-b909-0242ac120002", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "no customer entry found with uuid '9ffafd16-9ece-11ec-b909-0242ac120002'" {
		t.Errorf("Expected the 'message' key of the response to be set to 'no customer entry found'. Got '%s'", m["message"])
	}
}

// tom: rewritten function
func TestCreateCustomer(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"uuid":"c1b7c03b-6795-48ca-8894-f7742057b1c3", "customer_name":"ABC company", "country": "India"}`)
	req, _ := http.NewRequest("POST", "/api/v1/customers", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["uuid"] != "c1b7c03b-6795-48ca-8894-f7742057b1c3" {
		t.Errorf("Expected customer uuid to be 'c1b7c03b-6795-48ca-8894-f7742057b1c3'. Got '%v'", m["uuid"])
	}

	if m["customer_name"] != "ABC company" {
		t.Errorf("Expected customer name to be 'ABC company'. Got '%v'", m["customer_name"])
	}

	if m["country"] != "India" {
		t.Errorf("Expected customer location to be 'India'. Got '%v'", m["country"])
	}
}

func TestGetCustomer(t *testing.T) {
	clearTable()
	uuid := "9ffafd16-9ece-11ec-b909-0242ac120002"
	addCustomer(uuid)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/customers/%s", uuid), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addCustomer(uuid string) {
	prepare, _ := a.DB.Prepare("INSERT INTO customers (uuid, customer_name, country) VALUES (?, ?, ?);")

	prepare.Exec(uuid, "ABC company", "India")
}

func TestDeleteCustomer(t *testing.T) {
	clearTable()
	uuid := "c1b7c03b-6795-48ca-8894-f7742057b1c3"
	addCustomer(uuid)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/customers/%s", uuid), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/v1/customers/%s", uuid), nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/customers/%s", uuid), nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
