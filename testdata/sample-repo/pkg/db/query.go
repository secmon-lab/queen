package db

import (
	"database/sql"
	"fmt"
	"net/http"
)

var db *sql.DB

func activeStatus() string {
	return "active"
}

func GetUserByName(r *http.Request) (*sql.Row, error) {
	name := r.URL.Query().Get("name")
	query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)
	row := db.QueryRow(query)
	return row, nil
}

func GetActiveUsers() (*sql.Rows, error) {
	status := activeStatus()
	query := fmt.Sprintf("SELECT * FROM users WHERE status = '%s'", status)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
