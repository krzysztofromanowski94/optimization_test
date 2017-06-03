package server

import (
       // "net"
        "log"
        //"github.com/golang/protobuf/proto"
        //"github.com/krzysztofromanowski94/optimization_test/protomessage"
        "database/sql"
        //"fmt"
        "time"
        //"os"
        //"reflect"
)

var(
        tbTestFunctions = []string{""}
        tbResults = []string{""}
        tbHistory = []string{""}
)

func Start(address, dblogin, dbpass string) {
        /// Connect to client
        //listener, err := net.Listen("tcp", address) ; if err != nil {
        //        log.Fatal(err)
        //}
        //defer listener.Close()
        //conn, err := listener.Accept() ; if err != nil {
        //        log.Fatal(err)
        //}

        database, err := sql.Open("mysql", "user:pass@/black_hole_test") ; if err != nil { log.Fatal(err)
        } else {
                log.Println("Database opened")
        }
        defer func(){
                database.Close()
                log.Println("Database closed")
        }()

        err = initdatabase(database, tbTestFunctions, tbResults, tbHistory)
        if err != nil {
                log.Fatal(err)
        }
        return

        go func(){
                //protomessage.
                //protoData := new(protomessage.Data)
                //agent := protomessage.Data_Agent{}
                //agent.Agent.X = *proto.Float64(1.2)
                //protoData = &protomessage.Data{agent}
                //pro
                //switch u := protoData.(type){
		//
                //}
                //incomingData := make([]byte, 100000)
                //for {
                //        countBytes, err := conn.Read(incomingData); if err != nil {
                //                os.Stderr.Write([]byte(time.Now().String() + " read  " + err.Error() + "\n"))
                //                //close(agentChannel)
                //                return
                //        }
                //        err = proto.Unmarshal(incomingData[:countBytes], agentProto) ; if err != nil {
                //                os.Stderr.Write([]byte(time.Now().String() + " unmarshall " + err.Error() + "\n"))
                //                return
                //        }
                //        //agentChannel <- agentProto
                //}
        }()

        //go func(){
        //        file, err := os.Create("stderr.txt") ; if err != nil {
        //                log.Fatal(err)
        //        }
        //        defer func (){
        //                file.Close()
        //        }()
        //        for {
        //                a, ok := <- agentChannel ; if !ok {
        //                        return
        //                }
        //                check := reflect.ValueOf(a)
        //                if check.IsValid() && check.CanInterface(){
        //                        fmt.Println(a)
        //                } else {
        //                        fmt.Println("Not valid...")
        //                        file.Write([]byte(time.Now().String() + "" + err.Error() + "\n"))
        //                        os.Exit(-1)
        //                }
        //        }
        //}()

        log.Println("Sleeping")
        dur, err := time.ParseDuration("15s");
        if err != nil {
                log.Fatal(err)
        }
        time.Sleep(dur)

        log.Println("Not sleeping anymore")




}

