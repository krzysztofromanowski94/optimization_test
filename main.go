package main

import (
	"fmt"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"log"
)

func myQuery(database *sql.DB, queryStr string) []string {
	var resultSlice []string
	if query, err := database.Query(queryStr); err != nil { log.Println(err)
	} else {
		for query.Next() {
			var queryString string
			query.Scan(&queryString)
			resultSlice = append(resultSlice, queryString)
		}
	}
	return resultSlice
}

func main(){
	fmt.Println("Hello")
	database, err := sql.Open("mysql", "root:pass@/")
	if err != nil {	log.Fatal(err)
	} else {
		log.Println("Connected")
	}
	defer func(){
		database.Close()
		log.Println("Database algorithms closed")
	}()
	defer func(){
		log.Println("www")
	}()

	for _, str := range myQuery(database, "show databases"){
		fmt.Println(str)
	}



	if _, err := database.Exec("USE mysql"); err != nil { log.Println(err) }

	if query, err := database.Query("SHOW FROM user"); err != nil { log.Println(err)
	} else {
		for query.Next() {
			var queryString string
			query.Scan(&queryString)
			fmt.Println(queryString)
		}
	}



	log.Println("Finishing work")
}