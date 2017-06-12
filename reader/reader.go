package reader

import (
	"github.com/krzysztofromanowski94/optimization_test/protomessage"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"log"
	"fmt"
	"strings"
	"time"
	"bufio"
	"os"
	"io"
	"strconv"
)

var(
	grpcconn               *grpc.ClientConn
	client 			protomessage.OptimizationTestClient
	scanner                *bufio.Scanner

)

func agentToStr(a *protomessage.AgentType) string {
	var str, xstr, ystr string

	if a.X >=0 {
		xstr = fmt.Sprintf("X: %e          ", a.X)
	} else {
		xstr = fmt.Sprintf("X: %e         ", a.X)
	}
	if a.Y >= 0{
		ystr = fmt.Sprintf("Y: %e          ", a.Y)
	} else {
		ystr = fmt.Sprintf("Y: %e         ", a.Y)
	}


	if a.Best {
		str = xstr + ystr + fmt.Sprintf("Fitness: %e\tBest", a.Fitness)
	} else {
		str = xstr + ystr + fmt.Sprintf("Fitness: %e", a.Fitness)
	}

	return str
}

func Read(){
	initReader()

	fmt.Println()
	resultAmount := getResults()
	if resultAmount < 1 {
		fmt.Println("No results in database")
	}
	fmt.Print("Select result id you want to view: ")
	for scanner.Scan(){
		input := scanner.Text()
		inputInt, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if inputInt > 0 && inputInt <= resultAmount {
			fmt.Println("ok")
			getHistory(inputInt)
			break
		} else {
			fmt.Print("Not in range. Try again: ")
			continue
		}

	}

}

func getResults() (count int) {
	stream, err := client.GetResults(context.Background(), &protomessage.AskDummy{})
	if err != nil{
		log.Println(err)
	}
	for{
		result, err := stream.Recv()
		if err == io.EOF{
			return
		}
		count++
		if err != nil{
			log.Fatal("runGetResults err: ", err)
		}

		fmt.Printf("Id: %d\tFunc Name: %s\tCode (if available): %s\n" +
			"Agent Amount: %d\tSteps: %d\tBest Fitness: %e\n" +
			"Borders: %s\n" +
			"Date: %s\n\n",
			result.ResultsId, result.TestFunctionsName, result.TestFunctionsCode,
			result.ResultsAgentAmount, result.Steps, result.BestFitness,
			result.Borders,
			result.ResultDate,
		)
	}
}

func getHistory(result_id int){
	var(
		onlyBest bool
	)
	stream, err := client.GetHistory(context.Background())
	if err != nil {
		log.Println("getHistory init stream err: ", err)
		stream.CloseSend()
	}
	i := uint64(0)
	exitHistory := make(chan bool)
	go func() {
		for {
			err = stream.Send(&protomessage.AskHistory{uint64(result_id), i})
			if err != nil {
				log.Println("Send AskHistory err: ", err)
			}
			historyPage, err := stream.Recv()
			if err != nil {
				log.Println("Get historyPage err: ", err)
			}

			for _, agent := range historyPage.Agent {
				switch onlyBest{
				case true:
					if agent.Best {
						fmt.Println(agentToStr(agent))
					}
				case false:
					fmt.Println(agentToStr(agent))
				}

			}

			fmt.Printf("Agents in step %d: %d\n", i, len(historyPage.Agent))
			fmt.Println("Prev: a\tNext: d\tStep forward: d{+/-}n\tToggle view type: t\tQuit: q")
			for scanner.Scan() {
				str := scanner.Text()
				switch  {
				case str == "":
					fallthrough
				case str == "d":
					i++
				case str == "a":
					i--
				case str == "t":
					onlyBest = !onlyBest
				case str == "q":
					exitHistory <- true
					return
				case strings.HasPrefix(str,"d+"):
					for {
						step, err := strconv.Atoi(str[2:])
						if err != nil{
							log.Println(err)
							continue
						}
						i += uint64(step)
						break
					}
				case strings.HasPrefix(str,"d-"):
					for {
						step, err := strconv.Atoi(str[2:])
						if err != nil{
							log.Println(err)
							continue
						}
						if int(i) - step < 0{
							i = 0
							break
						} else {
							i -= uint64(step)
						}
						break
					}
				default:
					continue

				}
				break
			}
		}
	}()
	<- exitHistory
	err = stream.CloseSend()
	if err != nil {
		log.Println("getHistory close stream err: ", err)
	}
}


func initReader(){
	scanner = bufio.NewScanner(os.Stdin)

}


func Connect(address string) {
	var err error
	grpcconn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client = protomessage.NewOptimizationTestClient(grpcconn)

}

func CloseConnection() {
	fmt.Println("closing")
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