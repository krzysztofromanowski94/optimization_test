package main

import (
	//"database/sql"
	"flag"
	"fmt"
	//_ "github.com/go-sql-driver/mysql"
	"bufio"
	"github.com/krzysztofromanowski94/optimization_test/client"
	"github.com/krzysztofromanowski94/optimization_test/server"
	"log"
	"os"
	//"runtime"
	//"time"
	"github.com/krzysztofromanowski94/optimization_test/reader"
)

var (
	address string
	dblogin string
	dbpass  string
	action  string
	scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
)

func init() {
	flag.StringVar(&address, "d", "localhost:2110", "Please provide correct ip and port (ex.: 127.0.0.1:1234")
	flag.StringVar(&dblogin, "u", "root", "Database username")
	flag.StringVar(&dbpass, "p", "pass", "Database password")
	flag.StringVar(&action, "a", "server", "server / client / reader")
	flag.Parse()
}

func main() {

	switch action {
	case "server":
		log.Println("ok, I'm server")
		server.Start(address, dblogin, dbpass)
	case "client":
		log.Println("ok, I'm client")
		client.Connect(address)
		client.Compute()
		//client.GetResults()
		client.CloseConnection()
		fmt.Println("Thank you, bye")
	case "reader":
		log.Println("ok, I'm reader")
		reader.Connect(address)
		reader.Read()
		reader.CloseConnection()

	default:
		log.Println("I don't know my purpose")
		os.Exit(42)
	}

	return

	//client.Compute()
	//
	//
	//fmt.Println("Hello")
	//database, err := sql.Open("mysql", "root:pass@/")
	//if err != nil {
	//	log.Fatal(err)
	//} else {
	//	log.Println("Connected")
	//}
	//defer func() {
	//	database.Close()
	//	log.Println("Database algorithms closed")
	//}()
	//defer func() {
	//	log.Println("www")
	//}()
	//
	//for _, str := range server.myQuery(database, "show databases") {
	//	fmt.Println(str)
	//}
	//
	//if _, err := database.Exec("USE mysql"); err != nil {
	//	log.Println(err)
	//}
	//
	//if query, err := database.Query("SHOW FROM user"); err != nil {
	//	log.Println(err)
	//} else {
	//	for query.Next() {
	//		var queryString string
	//		query.Scan(&queryString)
	//		fmt.Println(queryString)
	//	}
	//}

	log.Println("Finishing work")
}
