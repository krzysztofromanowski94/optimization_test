package server

import (
	"database/sql"
	"log"
	//"fmt"
	"github.com/krzysztofromanowski94/optimization_test/protomessage"
	"fmt"
	"time"
)

type protoBuffData struct {
	bytes []byte
	count int
}

func addToDB(database *sql.DB, tbTypes []string) error {
	if len(tbTypes) == 1{
		//fmt.Println("addToDB len 1: ", tbTypes)
		_, err := database.Exec(tbTypes[0])
		if err != nil {
			return err
		}
		return nil
	}
	parsedQuery := "CREATE TABLE " + tbTypes[0] + "(\n"

	//fmt.Println("\tjust log:\n", tbTypes)

	for _, val := range tbTypes[1 : len(tbTypes)-1] {
		parsedQuery += val + ",\n"
	}
	parsedQuery += tbTypes[len(tbTypes)-1] + ");"

	_, err := database.Exec(parsedQuery)
	if err != nil {
		return err
	}

	return nil
}

func newResult(database *sql.DB, oneOfChan chan *protomessage.Oneof) {
	//var agentAmount uint64
	tbTestFunctions_id := -1
	//tbResults_id := 0
	tbResults_agent_amount := uint64(0)
	tbResults_code := ""
	testFunc := ""
	latestID := 0

	err := database.QueryRow("SELECT max(id) FROM " + tbTestFunctions[0]).Scan(&latestID)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Lastest id: ", latestID)

	for {
		oneOf, ok := <-oneOfChan
		if !ok {
			log.Fatal("newResult service not ok")
		}

		switch message := oneOf.Union.(type) {
		case *protomessage.Oneof_Agent:
			if message.Agent.Best {
				fmt.Println(message.Agent)
			}

		case *protomessage.Oneof_Result:
			tbResults_agent_amount = message.Result.AgentAmount
			tbResults_code = message.Result.Code
			testFunc = message.Result.TestFunc
			fmt.Println("new result:")
			fmt.Println(tbResults_agent_amount, tbResults_code, testFunc)

			_, err := database.Exec(
				"INSERT IGNORE INTO " + tbTestFunctions[0] + " VALUES (?, ?)", latestID, testFunc)
			if err != nil {
				log.Println(err)
			} else {
				latestID++
			}
			fmt.Println("after insert")


			err = database.QueryRow("SELECT id FROM " + tbTestFunctions[0] + " WHERE name=?", testFunc).Scan(&tbTestFunctions_id)
			if err != nil {
				log.Println(err)
			}
			fmt.Println("tbTestFunctions_id: ", tbTestFunctions_id)

			entries, err := database.Query("SELECT id, name FROM " + tbTestFunctions[0])
			if err != nil{
				fmt.Println(err)
			}
			for entries.Next(){
				v1 := 0
				v2 := ""
				err := entries.Scan(&v1, &v2)
				if err != nil{
					log.Println(err)
				}
				fmt.Println(v1, v2)
			}
			fmt.Println("after entries")



			time.Sleep(time.Minute)

		//default:
		//	fmt.Printf("Type %T\n", message)
		}
	}

}

func myQuery(database *sql.DB, queryStr string) []string {
	var resultSlice []string
	if query, err := database.Query(queryStr); err != nil {
		log.Println(err)
	} else {
		for query.Next() {
			var queryString string
			query.Scan(&queryString)
			resultSlice = append(resultSlice, queryString)
		}
	}
	return resultSlice
}
