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
		log.Println("Database closed")
		log.Println("Yeah! Modified with hp")
	}()


	for _, str := range myQuery(database, "show databases"){
		fmt.Println(str)
	}

	// create database school
	dbName := "school"
	if _, err := database.Exec("CREATE DATABASE " + dbName); err != nil { log.Println(err)}
	defer func(){
		if _, err := database.Exec("DROP DATABASE IF EXISTS school"); err != nil { log.Println(err)}
	}()

	if _, err := database.Exec("USE " + dbName); err != nil { log.Println(err)}

	for _, str := range myQuery(database, "SELECT DATABASE()"){
		fmt.Println("selected database:\t " + str)
	}



	if _, err := database.Exec("CREATE TABLE student(" +
		"first_name VARCHAR(30) NOT NULL," +
		"last_name VARCHAR(30) NOT NULL," +
		"email VARCHAR(60) NULL," +
		"street VARCHAR(50) NOT NULL," +
		"city VARCHAR(40) NOT NULL," +
		"zip MEDIUMINT UNSIGNED NOT NULL," +
		"phone VARCHAR(20) NOT NULL," +
		"birth_date DATE NOT NULL," +
		"sex ENUM('M', 'F') NOT NULL," +
		"date_entered TIMESTAMP," +
		"lunch_cost FLOAT NULL," +
		"student_id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY)") ; err != nil { log.Println(err) }


	// create database algorithms
	log.Println("Creating database table for algorithms code")
	dbName = "algorithms"
	tbName := "code"
	if _, err := database.Exec("CREATE DATABASE " + dbName); err != nil { log.Println(err) } else {
		log.Println("Created db " + dbName)
	}
	defer func(){
		if _, err := database.Exec("DROP DATABASE algorithms"); err != nil { log.Println(err) }
	}()

	if _, err := database.Exec("USE " + dbName); err != nil { log.Println(err) } else {
		log.Println("Use db " + dbName)
	}
	if _, err := database.Exec("CREATE TABLE " + tbName + " (" +
		"id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY," +
		"name VARCHAR(50) NOT NULL," +
		"code MEDIUMTEXT NULL)") ; err != nil { log.Println(err) } else {
		log.Println("Created table " + tbName)
	}

       myFunc := func(aa string) string { return aa }
       log.Println(myFunc("asd"))


       // Prepare statement for inserting data
       insertNewCode, err := database.Prepare("INSERT INTO " + tbName + " VALUES( ?, ? )") // ? = placeholder
       if err != nil {
           log.Println(err) // proper error handling instead of panic in your app
       }
       defer insertNewCode.Close() // Close the statement when we leave main() / the program terminates


       codeOut, err := database.Prepare("SELECT code FROM " + tbName + " WHERE name = ?")
       if err != nil {
           log.Println(err) // proper error handling instead of panic in your app
       }
       defer stmtOut.Close()


       if _, err := insertNewCode.Exec("new code", "blablablala"); err != nil { log.Println(err) }


       var row string
       // Query the square-number of 13
       err = stmtOut.QueryRow("new code").Scan(&row) // WHERE number = 13
       if err != nil {
           panic(err.Error()) // proper error handling instead of panic in your app
       }
       fmt.Printf("Result: %s", row) 


	//if _, err := database.Exec("USE mysql"); err != nil { log.Println(err) }





	log.Println("Finishing work")
}
