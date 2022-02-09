package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	queryGetCustomerByID  string = "SELECT uuid, customer_name, country FROM customers where uuid = '%s';"
	queryGetAllCustomers  string = "SELECT uuid, customer_name, country FROM customers;"
	executeInsertCustomer string = "INSERT INTO customers (uuid, customer_name, country) VALUES (?, ?, ?);"
	executeDeleteCustomer string = "DELETE FROM customers WHERE uuid = ?;"
)

type Customer struct {
	UUID    string `json:"uuid"`
	Name    string `json:"customer_name"`
	Country string `json:"country"`
}

type Customers []Customer

func assertRowsAffected(result sql.Result, expectedNumRows int64) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected != expectedNumRows {
		return fmt.Errorf("Unexpected number of DB Records affected: %d", (rowsAffected))
	}

	return nil
}

func InsertCustomerEntry(db *sql.DB, customer *Customer) error {

	prepStmt, err := db.Prepare(executeInsertCustomer)
	if err != nil {
		return err
	}
	defer prepStmt.Close()
	result, err := prepStmt.Exec(customer.UUID, customer.Name,
		customer.Country)

	if err != nil {
		return err
	}

	err = assertRowsAffected(result, 1)
	if err != nil {
		return err
	}

	log.Printf("Successfully inserted Customer '%s' with uuid: %s", customer.Name, customer.UUID)

	return nil

}

func GetCustomerEntry(db *sql.DB, uuid string) (*Customer, error) {

	rows, err := db.Query(fmt.Sprintf(queryGetCustomerByID, uuid))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	exists := rows.Next()
	if !exists {
		err = rows.Err()
		if err != nil {
			return nil, err
		}
		// No rows matched
		return nil, errors.New("No rows matched")

	}

	var UUID, customerName, country string

	if err := rows.Scan(&UUID, &customerName, &country); err != nil {
		return nil, err
	}

	return &Customer{
		UUID:    UUID,
		Name:    customerName,
		Country: country,
	}, nil
}

func GetAllCustomerEntries(db *sql.DB) ([]*Customer, error) {

	rows, err := db.Query(queryGetAllCustomers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	exists := rows.Next()
	if !exists {
		err = rows.Err()
		if err != nil {
			return nil, err
		}
		return nil, errors.New("No rows matched")

	}

	var customers []*Customer

	for ok := true; ok; ok = rows.Next() {
		var UUID, customerName, country string

		if err := rows.Scan(&UUID, &customerName, &country); err != nil {
			return nil, err
		}

		customers = append(
			customers,
			&Customer{
				UUID:    UUID,
				Name:    customerName,
				Country: country,
			},
		)
	}

	return customers, nil

}

func DeleteCustomerEntry(db *sql.DB, uuid string) error {

	prepStmt, err := db.Prepare(executeDeleteCustomer)
	if err != nil {
		return err
	}
	defer prepStmt.Close()

	log.Printf("Preparing to delete Customer with uuid: %s", uuid)
	_, err = prepStmt.Exec(uuid)

	if err != nil {
		return err
	}

	return nil

}
