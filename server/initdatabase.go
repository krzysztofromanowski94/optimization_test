package server

import (
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	//"fmt"
	"log"
	"fmt"
	"errors"
)

func initdatabase(database *sql.DB, tbTestFunctions, tbResults, tbHistory []string) error {
	// database name to be initialized:
	dbName := "opt_test"


	// check if @dbName exists
	query, err := database.Query("show databases")
	if err != nil {
		log.Fatal(err)
	} else {
		for query.Next() {
			var queryString string
			err := query.Scan(&queryString)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(queryString)
			if queryString == dbName {
				log.Println("Database", dbName, "exists")
				return errors.New("Database " + dbName + " exists")
			}
		}
	}

	// create table algorithms
	{
		if _, err := database.Exec("CREATE TABLE algorithms (" +
			"id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY," +
			"name VARCHAR(50) NOT NULL," +
			"code MEDIUMTEXT NULL)"); err != nil {
			log.Println(err)
		} else {
			log.Println("Created table algorithms")
		}
	}
	// insert something to algorithms
	{
		_, err := database.Exec("INSERT  INTO algorithms VALUES (NULL, 'new code', 'blabla')")
		if err != nil {
			log.Println(err)
		}

	}

	//
	{
		query, err := database.Query("SELECT * FROM algorithms")
		if err != nil {
			log.Println(err)
		}
		for query.Next(){
			var id, name, code string
			query.Scan(&id, &name, &code)
			fmt.Println(id, name, code)
	}

	}

	return nil
}
