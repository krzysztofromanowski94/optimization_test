package main

import (
	//"database/sql"
	"flag"
	"fmt"
	//_ "github.com/go-sql-driver/mysql"
	//"bufio"
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
)

func init() {
	flag.StringVar(&address, "d", ":2110", "Please provide correct ip and port (ex.: 127.0.0.1:1234")
	flag.StringVar(&dblogin, "u", "user", "Database username")
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
}
