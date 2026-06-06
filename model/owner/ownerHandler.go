package owner

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	// Using pgconn for specific Postgres error codes for duplicates
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func InsertEntry(database *sql.DB, item string, code string, cost int) error {

	// Making an sql statement with either prepare or query (query returns rows, Exec with a statement does not)
	// Changed placeholders from ? in sqlite3 to $1, $2, $3 for Postgres
	statement, err := database.Prepare("INSERT INTO menu (item, code, cost) VALUES ($1, $2, $3)")

	if err != nil {
		return err
	}

	_, err = statement.Exec(item, code, cost)

	if err != nil {
		// Making a variable to hold the error that the postgres library will throw (23505 is the code for unique constraints)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("Duplicate! The code '%s' already exists", code)
			}
		}
		return err
	}
	return nil

}

func RemoveEntry(database *sql.DB, code string) {

	// TODO Make check if it was in the database or not!

	statement, _ := database.Prepare("DELETE FROM menu WHERE code = $1")

	statement.Exec(code)
}

func UpdateCost(database *sql.DB, code string, newPrice int) {

	// TODO Also need some error checking here as well

	statement, _ := database.Prepare("UPDATE menu SET cost = $1 Where code = $2")

	result, _ := statement.Exec(newPrice, code)

	// Show the change
	affectedRows, _ := result.RowsAffected()
	fmt.Printf("Updated %d\n", affectedRows)

}

func GetCost(database *sql.DB, code string) (int, error) {

	var cost int
	// Changed placeholder from ? to $1
	err := database.QueryRow("SELECT cost FROM menu WHERE code = $1", code).Scan(&cost)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("The item with code '%s' does not exist!", code)
		} else {
			return 0, err
		}
	}

	return cost, nil
}

func PrintAllRecords(database *sql.DB) {

	rows, _ := database.Query("SELECT id, item, code, cost FROM menu")

	var id int
	var item string
	var code string
	var cost int

	for rows.Next() {
		rows.Scan(&id, &item, &code, &cost)
		fmt.Println(strconv.Itoa(id) + ": " + item + " Item code: " + code + " Cost: " + strconv.Itoa(cost) + " yen")
	}
}
