package client

import (
	"github.com/krzysztofromanowski94/BHKulak/optimization_test/client/goblackholes"
	"github.com/krzysztofromanowski94/BHKulak/optimization_test/protokulak"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	grpcconn               *grpc.ClientConn
	client                 protokulak.BHServiceClient
	blackHolesAgentChannel chan *goblackholes.Agent
	oneofToServerChan      chan *protokulak.Oneof
	quitBlackholes         chan bool
	acceptAgets            bool = true
	initVariables          goblackholes.InitVariables
)

func newResult() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("newResult recover: ", rec)
		}
	}()

	stream, err := client.DoBHA(context.Background())
	if err != nil {
		log.Println("err: NewResult init stream: ", err)
	}

	err = stream.Send(&protokulak.Oneof{&protokulak.Oneof_Init{&protokulak.InitConnection{}}, "init"})
	if err != nil {
		log.Fatal("err: NewResult stream: ", err)
	}

	initRcv, err := stream.Recv()
	if err != nil {
		log.Fatal("err: initRcv, err := stream.Recv(): ", err)
	}

	switch union := initRcv.Union.(type) {
	case *protokulak.Oneof_NewResult:
		log.Println("new result: ", union)

		switch union.NewResult.TypeOfFunction {
		case "Rosenbrock":
			initVariables.TypeOfFucntion.Rosenbrock = true
		case "McCormick":
			initVariables.TypeOfFucntion.McCormick = true
		case "Easom":
			initVariables.TypeOfFucntion.Easom = true
		case "Rastrigin":
			initVariables.TypeOfFucntion.Rastrigin = true
		case "StringEvaluation":
			initVariables.TypeOfFucntion.StringEvaluation = union.NewResult.Code
		}
		initVariables.AgentAmount = int(union.NewResult.AgentAmount)
		borderSlice := strings.Split(union.NewResult.Borders, ":")
		log.Println("Border slice: ", borderSlice)
		initVariables.Border.X1, err = strconv.ParseFloat(borderSlice[0], 64)
		initVariables.Border.X2, err = strconv.ParseFloat(borderSlice[1], 64)
		initVariables.Border.Y1, err = strconv.ParseFloat(borderSlice[2], 64)
		initVariables.Border.Y2, err = strconv.ParseFloat(borderSlice[3], 64)
		if err != nil {
			log.Fatal("err: initVariables.Border.X/Y, err = strconv.ParseFloat(borderSlice[], 64)", err)
		}

		log.Println(initVariables)
	default:
		log.Fatal("err: Expected new result, got: ", union)
	}

	oneofToServerChan = make(chan *protokulak.Oneof, initVariables.AgentAmount)
	blackHolesAgentChannel = make(chan *goblackholes.Agent, initVariables.AgentAmount)
	quitBlackholes = make(chan bool, 1)

	go func() {
		goblackholes.Start(blackHolesAgentChannel, quitBlackholes, initVariables)
	}()

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
			}
		}()
		newAgent := &protokulak.AgentType{}
		for {
			if newBHAgent, ok := <-blackHolesAgentChannel; ok {
				newAgent.Step = newBHAgent.Times
				newAgent.Fitness = newBHAgent.Fitness
				newAgent.Best = newBHAgent.Best
				newAgent.X = newBHAgent.X
				newAgent.Y = newBHAgent.Y
				oneofToServerChan <- &protokulak.Oneof{&protokulak.Oneof_Agent{Agent: newAgent}, "agent"}
			} else {
				return
			}
		}
	}()

	go func() {
		for {
			newResult, ok := <-oneofToServerChan
			if !ok {
				break
			}
			if acceptAgets {
				err := stream.Send(newResult)
				if err == io.EOF {
					log.Println("Server closed stream")
					break
				}
				if err != nil {
					log.Println("err: NewResult send agent: ", err)
				}
			} else {
				break
			}
		}
	}()

	func() {
		for {
			initRcv, err := stream.Recv()
			if err != nil {
				log.Fatal("initRcv, err := stream.Recv(): ", err)
			}
			log.Println("Got close signal")
			switch union := initRcv.Union.(type) {
			case *protokulak.Oneof_Ret:
				quitBlackholes <- true
				acceptAgets = false
				return
			default:
				log.Println("warning: Expected return type, got: ", union)
			}
		}
	}()

	close(blackHolesAgentChannel)
	close(oneofToServerChan)
	err = stream.CloseSend()
	if err != nil {
		log.Println("err: CloseSend: ", err)
	}
	return
}

func Compute() {
	newResult()
}

func Connect(address string) {
	var err error
	grpcconn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client = protokulak.NewBHServiceClient(grpcconn)

}

func CloseConnection() {
	log.Println("Closing client...")
	err := grpcconn.Close()
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return
		}
		log.Println(err)
	}
	time.Sleep(time.Second)
	log.Println("Closed")
}
