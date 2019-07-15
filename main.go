package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"database/sql"

	_ "github.com/lib/pq"
)

func main() {
	driver := flag.String("engine", "postgres", "Select one of supported dirvers. { postgres }")
	connection := flag.String("connection", "", "The appropriate connection string.")
	query := flag.String("query", "", "The query to get data from database.")
	format := flag.String("format", "csv", "The format for printing the result of query. { csv }")
	flag.Parse()

	db, err := sql.Open(*driver, *connection)
	if err != nil {
		log.Fatal("connecting to database: ", err)
	}

	rows, err := db.Query(*query)
	if err != nil {
		log.Fatal("running the query: ", err)
	}

	var result string
	switch strings.ToLower(*format) {
	case "csv":
		result, err = RowsToCSV(rows)
	default:
		log.Fatal("unknown format: ", *format)
	}
	if err != nil {
		log.Fatal("converting result to ", *format, ": ", err)
	}

	fmt.Print(result)
}
