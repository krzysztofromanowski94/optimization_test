package server

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func initdatabase(database *sql.DB, tbTestFunctions, tbResults, tbHistory []string) error {

	err := addToDB(database, tbTestFunctions)
	if err != nil {
		log.Println(err)
	}
	err = addToDB(database, tbResults)
	if err != nil {
		log.Println(err)
	}
	err = addToDB(database, tbHistory)
	if err != nil {
		log.Println(err)
	}

	for i, exec := range mysqlOperations{
		err := addToDB(database, exec)
		if err != nil{
			log.Printf("mysqlOperations[%d] %s", i, err)
		}
	}

	return nil
}
