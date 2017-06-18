package main

import (
	"flag"
	"github.com/krzysztofromanowski94/BHKulak/optimization_test/client"
)

var (
	address string
)

func init() {
	flag.StringVar(&address, "d", "localhost:2110", "Please provide correct ip and port (ex.: 127.0.0.1:1234")
	flag.Parse()
}

func main() {
	client.Connect(address)
	client.Compute()
	client.CloseConnection()
	return
}
