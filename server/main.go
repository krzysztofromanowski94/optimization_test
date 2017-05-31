package main

import (
        "fmt"
        "net"
        "log"
        "os"
        "github.com/krzysztofromanowski94/agentproto"
        "github.com/golang/protobuf/proto"
)

func main() {
        //c := make (chan *agentproto.AgentData)
        //go func (){
        //        for {
        //                message := <- c
        //                fmt.Println(message)
        //       }
        //}()
        listener, err := net.Listen("tcp", ":2110")
        if err != nil {
                log.Println(err)
        }
        defer listener.Close()
        if conn, err := listener.Accept(); err != nil {
                log.Println(err)
        } else {
                for {
                        stuff := make([]byte, 4096)
                        if _, err := conn.Read(stuff); err != nil {
                                log.Println(err)
                                os.Exit(-1)
                        }
                        protoobject := new(agentproto.AgentData)
                        err = proto.Unmarshal(stuff, protoobject)
			fmt.Println(protoobject)
                        //c <- protoobject
                }
        }
}

