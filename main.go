package main

import (
	"database/sql"
	"fmt"
	"log"

	"tabsquare/db"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db_client, err := sql.Open("mysql", "user:BcGH2Gj41J5VF1@tcp(10.46.142.201:3306)/tabsquare")

	if err != nil {
		log.Print(err.Error())
	}
	defer db_client.Close()

	customer := &db.Customer{
		UUID:    "85743dba-daa5-4301-ab84-bb4a2611a619",
		Name:    "XYZ company",
		Country: "India",
	}

	err = db.InsertCustomerEntry(db_client, customer)
	if err != nil {
		log.Print(err.Error())
	}

	test, err := db.GetCustomerEntry(db_client, "85743dba-daa5-4301-ab84-bb4a2611a619")

	if err != nil {
		log.Print(err.Error())
	}

	log.Printf(test.Name)

	customers, err := db.GetAllCustomerEntries(db_client)

	for _, s := range customers {
		fmt.Println(s.UUID, s.Name)
	}

	err = db.DeleteCustomerEntry(db_client, "85743dba-daa5-4301-ab84-bb4a2611a619")
	if err != nil {
		log.Print(err.Error())
	}
}
