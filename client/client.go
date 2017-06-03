package client

import (
	"net"
	//"fmt"
	//"bufio"
	"log"
	//"github.com/krzysztofromanowski94/optimization_test/protomessage"
	//"github.com/golang/protobuf/proto"
	//"os"
	"time"
)

var (
	//agentProto = new(protomessage.AgentData)
	//data []byte
)

func Start(address string) {
	conn, err := net.Dial("tcp", address);
	if err != nil {
		log.Fatal(err)
	}
	//go func() {
	//	i := 0.0
	//	for {
	//		agentProto.X = proto.Float64(i)
	//		agentProto.Y = proto.Float64(i)
	//		agentProto.Average = proto.Int32(int32(i))
	//		agentProto.Fitness = proto.Float64(i / 2.)
	//
	//		data, err := proto.Marshal(agentProto) ; if err != nil {
	//			os.Stderr.Write([]byte(time.Now().String() + " " + err.Error() + "\n"))
	//			os.Exit(-1)
	//		}
	//		_, err = conn.Write(data) ; if err != nil {
	//			os.Stderr.Write([]byte(time.Now().String() + " " + err.Error() + "\n"))
	//			os.Exit(-1)
	//		}
	//		//time.Sleep(time.Second)
	//		i += 0.11235454
	//	}
	//}()
	dur, _ := time.ParseDuration("10s")
	time.Sleep(dur)
	conn.Close()

	//status, err := bufio.NewReader(conn).ReadString('\n')
}
