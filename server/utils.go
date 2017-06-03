package server

import (
	"database/sql"
	"log"
)

func MyQuery(database *sql.DB, queryStr string) []string {
	var resultSlice []string
	if query, err := database.Query(queryStr); err != nil { log.Println(err)
	} else {
		for query.Next() {
			var queryString string
			query.Scan(&queryString)
			resultSlice = append(resultSlice, queryString)
		}
	}
	return resultSlice
}

