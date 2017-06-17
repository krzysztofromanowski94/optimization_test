package main

import (
	//"database/sql"
	"flag"
	//"fmt"
	//_ "github.com/go-sql-driver/mysql"
	//"bufio"
	//"log"
	//"os"
	//"runtime"
	//"time"
	"github.com/krzysztofromanowski94/BHKulak/optimization_test/client"
)

var (
	address string
	//dblogin string
	//dbpass  string
	//action  string
)

func init() {
	flag.StringVar(&address, "d", ":2110", "Please provide correct ip and port (ex.: 127.0.0.1:1234")
	//flag.StringVar(&dblogin, "u", "user", "Database username")
	//flag.StringVar(&dbpass, "p", "pass", "Database password")
	//flag.StringVar(&action, "a", "server", "server / client / reader")
	flag.Parse()
}

func main() {

	//switch action {
	//case "client":
	//	log.Println("ok, I'm client")
		client.Connect(address)
		client.Compute()
		client.CloseConnection()
		//fmt.Println("Thank you, bye")
	//default:
	//	log.Println("I don't know my purpose")
	//	os.Exit(42)
	//}

	return
}
