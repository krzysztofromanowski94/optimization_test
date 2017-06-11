package server

import (
	"log"
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"github.com/krzysztofromanowski94/optimization_test/protomessage"
	"io"
	"net"
	"strconv"
)

var (
	tbTestFunctions = []string{
		"test_functions",
		"id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT",
		"name VARCHAR(50) NOT NULL",
		"code VARCHAR(255)",
		"UNIQUE KEY unique_test_func_key (name, code)",
	}
	tbResults = []string{
		"results",
		"id INT UNSIGNED NOT NULL AUTO_INCREMENT",
		"test_func_id INT UNSIGNED NOT NULL",
		"PRIMARY KEY (id, test_func_id)",
		"CONSTRAINT result_test_func_id_fk FOREIGN KEY (test_func_id) REFERENCES " + tbTestFunctions[0] + "(id)",
		"agent_amount INT UNSIGNED NOT NULL",
		"best VARCHAR(500)",
		"steps INT UNSIGNED",
		"borders VARCHAR(500)",
		"date TIMESTAMP",
	}
	tbHistory = []string{
		"history",
		"result_id INT UNSIGNED NOT NULL",
		"CONSTRAINT history_result_id_fk FOREIGN KEY (result_id) REFERENCES " + tbResults[0] + "(id)",
		"x VARCHAR(500) NOT NULL",
		"y VARCHAR(500) NOT NULL",
		"fitness VARCHAR(500) NOT NULL",
		"step INT UNSIGNED NOT NULL",
		"best BOOL",
		"UNIQUE KEY unique_result (x, y, fitness, step)",

	}

	mysqlOperations = [][]string{
		{
			"CREATE VIEW view_results as " +
				"SELECT " +
				"results.id, " +
				"test_functions.name, " +
				"test_functions.code, " +
				"results.agent_amount, " +
				"results.best, " +
				"results.steps, " +
				"results.borders, " +
				"results.date " +
				"FROM results " +
				"INNER JOIN test_functions ON results.test_func_id=test_functions.id ",
		},
	}

	listener net.Listener
	database *sql.DB
)

type optimizationTestServer struct {
	tbTestFunctions_id     int
	tbTestFunctions_name   string
	tbResults_id 		uint64
	tbResults_code         string
	tbResults_agent_amount uint64
	tbResults_borders string

}

func (t *optimizationTestServer) NewResult(oneofstream protomessage.OptimizationTest_NewResultServer) error {

	thisResult, err := oneofstream.Recv()
	if err != nil {
		log.Println("NewResult: thisResult, err := oneofstream.Recv() ", err)
		return oneofstream.SendAndClose(&protomessage.ReturnType{false, "NewResult: " + err.Error()})
	}
	switch union := thisResult.Union.(type){
	case *protomessage.Oneof_Result:
		log.Print("Getting new results: ")
		if union.Result.Code == "" {
			fmt.Printf("Function: %s\tAmount of agents: %d\n", union.Result.TestFunc, union.Result.AgentAmount)
		} else {
			fmt.Printf("Function: %s\tAmount of agents: %d\tCode: %s\n", union.Result.TestFunc, union.Result.AgentAmount, union.Result.Code)

		}

		t.tbTestFunctions_name = union.Result.TestFunc
		t.tbResults_agent_amount = union.Result.AgentAmount
		t.tbResults_code = union.Result.Code
		t.tbResults_borders = union.Result.Borders

		// pkt
		_, err := database.Exec(
			"INSERT IGNORE INTO " + tbTestFunctions[0] + " VALUES (NULL, ?, ?)", t.tbTestFunctions_name, t.tbResults_code)
		if err != nil {
			log.Println("INSERT IGNORE INTO " + tbTestFunctions[0], err)
		}

		// pkt
		err = database.QueryRow("SELECT id FROM " + tbTestFunctions[0] +
			" WHERE name=? AND code=?", t.tbTestFunctions_name, t.tbResults_code).Scan(&t.tbTestFunctions_id)
		if err != nil {
			log.Println("SELECT id FROM " + tbTestFunctions[0], err)
		}

		// pkt
		_, err = database.Exec(
			"INSERT INTO " + tbResults[0] + " VALUES (NULL, ?, ?, NULL, NULL, ?, NOW())", t.tbTestFunctions_id, t.tbResults_agent_amount, t.tbResults_borders)
		if err != nil {
			log.Println("INSERT INTO " + tbResults[0], err)
		}

		// pkt
		query := database.QueryRow("SELECT id FROM "+ tbResults[0] + " ORDER BY date DESC")
		err = query.Scan(&t.tbResults_id)
		if err != nil {
			log.Println("query scan err: ", err)
		}

	default:
		fmt.Println("Wrong receive")
		fmt.Println(union)
		return oneofstream.SendAndClose(&protomessage.ReturnType{false, "NewResult: wrong receive type"})

	}

	agentChannel := make(chan *protomessage.AgentType)
	go func(){
		for {
			tempOneof, err := oneofstream.Recv()
			if err == io.EOF {
				log.Println("Finished receiving agents")
				close(agentChannel)
				return
			}
			if err != nil {
				log.Printf("NewResult server get agent err type: %T %s: ", err, err)
				return
			}
			switch agent := tempOneof.Union.(type){
			case *protomessage.Oneof_Agent:
				agentChannel <- agent.Agent
			default:
				log.Printf("Unknown behaviour for %T\n", agent)
			}
		}

	}()

	finished := make(chan bool)
	for i := 0 ; i < 5 ; i++{
		go func(){
			for {
				agent, ok := <- agentChannel
				if !ok {
					finished <- true
					break
				}
				// pkt
				_, err = database.Exec(
					"INSERT IGNORE INTO " + tbHistory[0] + " VALUES (?, ?, ?, ?, ?, ?)",
					t.tbResults_id,
					agent.X,
					agent.Y,
					agent.Fitness,
					agent.Step,
					agent.Best,
				)
				if err != nil {
					log.Println("INSERT INTO " + tbHistory[0], err)
				}
			}
		}()
	}

	<- finished
	// pkt: UPDATE / podzapytanie / subquery / group by / limit
	magicstr := "UPDATE results " +
		"SET results.best=" +
		"(" +
		"SELECT fitness FROM history " +
		"WHERE result_id=" + strconv.FormatUint(t.tbResults_id, 10) + " " +
		"AND best=1 " +
		"ORDER BY step DESC LIMIT 1" +
		") " +
		"WHERE id=" + strconv.FormatUint(t.tbResults_id, 10)
	_, err = database.Exec(magicstr)
	if err != nil {
		fmt.Println(magicstr)
		log.Println("Update best in result: ", err)
	}
	magicstr = "UPDATE results " +
		"SET results.steps=" +
		"(" +
		"SELECT COUNT(DISTINCT step) FROM history " +
		"WHERE result_id=" + strconv.FormatUint(t.tbResults_id, 10) +
		") " +
		"WHERE id=" + strconv.FormatUint(t.tbResults_id, 10)
	_, err = database.Exec(magicstr)
	if err != nil {
		fmt.Println(magicstr)
		log.Println("Update best in result: ", err)
	}
	return oneofstream.SendAndClose(&protomessage.ReturnType{true, "Received all results"})
}

func (t *optimizationTestServer) GetResults(askDummy *protomessage.AskDummy, stream protomessage.OptimizationTest_GetResultsServer) error {
	// pkt: view / widok
	query, err := database.Query("SELECT * FROM view_results")
	if err != nil {
		fmt.Println("GetResults query: ", err)
		return err
	}
	result := &protomessage.TBResults{}
	for query.Next(){
		var(
			results_id uint64
			test_functions_name string
			test_functions_code string
			results_agent_amount uint64
			result_best_fitess string
			result_steps uint64
			result_borders string
			results_date string
		)
		err := query.Scan(&results_id, &test_functions_name, &test_functions_code, &results_agent_amount, &result_best_fitess, &result_steps, &result_borders, &results_date)
		if err != nil {
			fmt.Println("GetResults results.Scan err: ", err)
		}
		result.ResultsId = results_id
		result.TestFunctionsName = test_functions_name
		result.TestFunctionsCode = test_functions_code
		result.ResultsAgentAmount = results_agent_amount
		result.BestFitness, err = strconv.ParseFloat(result_best_fitess, 64)
		if err != nil{
			log.Println("Conversion error: ", err)
		}
		result.Steps = result_steps
		result.Borders = result_borders
		result.ResultDate = results_date
		err = stream.Send(result)
		if err != nil {
			log.Println("GetResults error: ", err)
			return err
		}
	}
	return nil
}

func (t *optimizationTestServer) GetHistory(stream protomessage.OptimizationTest_GetHistoryServer) error {
	var (
		x       float64
		y       float64
		fitness float64
		best    string
	)

	for {
		get, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Println("GetHistory receive err: ", err)
			return err
		}
		listAgent := make([]*protomessage.AgentType, 0)
		// pkt
		query, err := database.Query("SELECT x, y, fitness, best FROM history WHERE history.result_id=? AND history.step=? ORDER BY fitness ASC",
			get.ResultId, get.Step)
		if err != nil {
			log.Println("GetHistory queue selec from history err:\n", err)
			return err
		}
		for query.Next() {
			query.Scan(&x, &y, &fitness, &best)
			singleAgent := &protomessage.AgentType{x, y, fitness, get.Step, false}
			if best == "1"{
				singleAgent.Best = true
			}
			listAgent = append(listAgent, singleAgent)
		}
		err = stream.Send(&protomessage.HistoryPage{listAgent, uint64(len(listAgent))})
		if err != nil {
			log.Println("GetHistory send err: ", err)
		}
	}
	return nil
}

func newServer() *optimizationTestServer {
	srv := new(optimizationTestServer)
	srv.tbTestFunctions_id = -1
	return srv
}

func Start(address, dblogin, dbpass string) {

	var err error
	database, err = sql.Open("mysql", "user:pass@/black_hole_test")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database opened")
	}
	defer func() {
		database.Close()
		log.Println("Database closed")
	}()

	err = initdatabase(database, tbTestFunctions, tbResults, tbHistory)
	if err != nil {
		log.Fatal(err)
	}


	//Connect to client
	listener, err = net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	protomessage.RegisterOptimizationTestServer(grpcServer, newServer())

	grpcServer.Serve(listener)

	return
}
