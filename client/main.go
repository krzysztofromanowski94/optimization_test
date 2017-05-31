package main

import (
	"net"
	//"fmt"
	//"bufio"
	"log"
	"time"
	"github.com/krzysztofromanowski94/protopackage"
	"github.com/golang/protobuf/proto"
	"os"
)

func main() {
	//proto
	asd := new(protopackage.ProtoMessage)
	var stuff []byte
	var err error
	//asd.Message
	conn, err := net.Dial("tcp", "172.16.100.11:2110")
	if err != nil {
		log.Println(err)
		time.Sleep(time.Second)
	} else {
		i := 0.0
		for {
			//fmt.Fprintf(conn, "%d\n", i)
			i += 0.1
			asd.MyNiceFloat = proto.Float64(i)
			asd.MyNiceMessage = proto.String("asdasd")

			if stuff, err = proto.Marshal(asd) ; err != nil {
				log.Println(err)
				os.Exit(-1)
			}
			conn.Write(stuff)
			time.Sleep(time.Second)
		}
	}
	//status, err := bufio.NewReader(conn).ReadString('\n')
}
