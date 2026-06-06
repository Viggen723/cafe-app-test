package main

import (
	"cafe-app-backend/database"
	"cafe-app-backend/model/owner"

	"fmt"
	"strconv"
)

func main() {

	db := database.InitializeDB()
	defer db.Close()

	err := owner.InsertEntry(db, "Hot Coffee", "C01", 300)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Println("Insertion successful!")
	}

	// Try repeating the entry

	err = owner.InsertEntry(db, "Bagel", "P01", 500)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Println("Insertion successful!")
	}

	owner.PrintAllRecords(db)
	owner.RemoveEntry(db, "P01")
	owner.PrintAllRecords(db)

	var cost int
	cost, err = owner.GetCost(db, "C01")
	fmt.Println("Cost = " + strconv.Itoa(cost))

	owner.UpdateCost(db, "C01", 350)

	owner.PrintAllRecords(db)

}
