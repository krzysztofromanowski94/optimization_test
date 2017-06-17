package client

import (
	"github.com/krzysztofromanowski94/BHKulak/optimization_test/protokulak"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"bufio"
	"fmt"
	"github.com/krzysztofromanowski94/BHKulak/optimization_test/client/goblackholes"
	"os"
	"strconv"
	"log"
	"strings"
	"time"
	"io"
)

var (
	grpcconn               *grpc.ClientConn
	client			protokulak.BHServiceClient
	testFunctions          []string = []string{"Rosenbrock", "Easom", "McCormick", "Write your own function"}
	blackHolesAgentChannel chan *goblackholes.Agent
	oneofToServerChan      chan *protokulak.Oneof
	quitBlackholes         chan bool
	quitClient             chan bool
	acceptAgets bool = true
	initVariables          goblackholes.InitVariables
	scanner                *bufio.Scanner
)


func newResult(/*client protokulak.BHService_DoBHAServer*/){
	defer func(){
		if rec := recover(); rec != nil{
			fmt.Println("newResult recover: ", rec)
		}
	}()

	//newResult := <-oneofToServerChan
	stream, err := client.DoBHA(context.Background())
	if err != nil{
		fmt.Println("NewResult init stream: ", err)
	}
	err = stream.Send(&protokulak.Oneof{&protokulak.Oneof_Init{}})
	if err != nil {
		fmt.Println("NewResult stream: ", err)
	}

	initRcv, err := stream.Recv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(initRcv)


	switch union := initRcv.Union.(type){
	case *protokulak.Oneof_NewResult:
		fmt.Println("yeah, new result: ", union)

	default:
		log.Fatal("Lol, what is it? ", union)
	}

	err = stream.Send(&protokulak.Oneof{&protokulak.Oneof_Return{&protokulak.ReturnType{true, "o kurwa kurwa dziala"}}})
	if err != nil {
		log.Fatal("stream.Send oneof o kurwa dziala: ", err)
	}


	err = stream.CloseSend()
	if err != nil {
		log.Fatal("Close stream err: ", err)
	}


	return
	for {
		newResult, ok := <-oneofToServerChan
		if !ok {
			break
		}
		if acceptAgets {
			err := stream.Send(newResult)
			if err == io.EOF{
				fmt.Println("Server closed stream")
				break
			}
			if err != nil {
				fmt.Println("NewResult send agent: ", err)
			}
		} else {
			close(oneofToServerChan)
			break
		}
	}

	err = stream.CloseSend()
	if err != nil {
		log.Println("CloseSend err: ", err)
	}
	quitClient <- true
}

func Compute() {
	initCompute()

	go func() {
		newResult()
	}()

	fmt.Println("\nTo stop press any key...\n")

	go func() {
		goblackholes.Start(blackHolesAgentChannel, quitBlackholes, initVariables)
	}()

	go func(){
		defer func(){
			if rec := recover(); rec != nil{
			}
		}()
		newAgent := &protokulak.AgentType{}
		for {
			if newBHAgent, ok := <- blackHolesAgentChannel; ok{
				newAgent.Step = newBHAgent.Times
				newAgent.Fitness = newBHAgent.Fitness
				newAgent.Best = newBHAgent.Best
				newAgent.X = newBHAgent.X
				newAgent.Y = newBHAgent.Y
				oneofToServerChan <- &protokulak.Oneof{Union: &protokulak.Oneof_Agent{Agent: newAgent}}
			} else {
				return
			}
		}
	}()

	go func() {
		scanner.Scan()
		fmt.Println("Waiting for server to save results...")
		quitBlackholes <- true
		acceptAgets = false
	}()

	<-quitClient

	close(blackHolesAgentChannel)

	return
}

func initCompute() {
	scanner = bufio.NewScanner(os.Stdin)

	fmt.Println("Select type of function:")
	for i, val := range testFunctions {
		fmt.Println(i+1, val)
	}
	var typeOfFunctionStr string
	func() {
		for scanner.Scan() {
			typeOfFunctionStr = scanner.Text()
			switch typeOfFunctionStr {
			case "1":
				initVariables.TypeOfFucntion.Rosenbrock = true
				typeOfFunctionStr = "Rosenbrock"
				return
			case "2":
				initVariables.TypeOfFucntion.Easom = true
				typeOfFunctionStr = "Easom"
				return
			case "3":
				initVariables.TypeOfFucntion.McCormick = true
				typeOfFunctionStr = "McCormick"
				return
			case "4":
				fmt.Println("Write your function:")
				scanner.Scan()
				functionString := scanner.Text()
				initVariables.TypeOfFucntion.StringEvaluation = functionString
				typeOfFunctionStr = "String evaluation"
				return
			default:
				fmt.Println("Sorry, try again")
				fmt.Println("Select type of function:")
				for i, val := range testFunctions {
					fmt.Println(i+1, val)
				}
			}
		}
	}()

	fmt.Print("Amount of agents: ")
	agentAmount := 0
	func() {
		for scanner.Scan() {
			agentAmountStr := scanner.Text()
			agentAmountInt, err := strconv.Atoi(agentAmountStr)
			if err != nil {
				fmt.Println(err)
				continue
			}
			initVariables.AgentAmount = agentAmountInt
			agentAmount = agentAmountInt
			return
		}
	}()

	borders := goblackholes.Border_s{}
	fmt.Println("Setup borders:")
	func() {
		fmt.Print("x greater than: ")
		for scanner.Scan() {
			x1str := scanner.Text()
			x1, err := strconv.ParseFloat(x1str, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			borders.X1 = x1
			break
		}
		fmt.Print("x lesser than: ")
		for scanner.Scan() {
			x2str := scanner.Text()
			x2, err := strconv.ParseFloat(x2str, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			borders.X2 = x2
			break
		}
		fmt.Print("y greater than: ")
		for scanner.Scan() {
			y1str := scanner.Text()
			y1, err := strconv.ParseFloat(y1str, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			borders.Y1 = y1
			break
		}
		fmt.Print("y lesser than: ")
		for scanner.Scan() {
			y2str := scanner.Text()
			y2, err := strconv.ParseFloat(y2str, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			borders.Y2 = y2
			break
		}
	}()
	initVariables.Border = borders

	//oneofToServerChan = make(chan *protomessage.Oneof, initVariables.AgentAmount)
	blackHolesAgentChannel = make(chan *goblackholes.Agent, initVariables.AgentAmount)
	quitBlackholes = make(chan bool, 1)
	quitClient = make(chan bool, 1)

	newResult := &protokulak.ResultType{}
	newResult.AgentAmount = uint64(initVariables.AgentAmount)
	newResult.TypeOfFunction = typeOfFunctionStr
	newResult.Code = initVariables.TypeOfFucntion.StringEvaluation
	newResult.Borders = initVariables.Border.ToStr()

	//oneofToServerChan <- &protomessage.Oneof{&protomessage.Oneof_Result{newResult}}
}

func Connect(address string) {
	var err error
	grpcconn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	//client = protomessage.NewOptimizationTestClient(grpcconn)

}

func CloseConnection() {
	fmt.Println("Closing client...")
	err := grpcconn.Close()
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return
		}
		log.Println(err)
	}
	time.Sleep(time.Second)
	fmt.Println("Closed")
}
