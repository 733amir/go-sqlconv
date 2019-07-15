package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strconv"
	"testing"

	"database/sql"
	"encoding/csv"

	_ "github.com/mattn/go-sqlite3"
)

const (
	pathToSqliteFile = "./test.db"
	populateDB       = `
	create table Users (
		ID      int  primary key,
		Name    text,
		Balance real
	);
	insert into Users (
		ID, Name, Balance
	) values(
		1, 'amir', 1.2
	), (
		2, 'سلام', 132
	), (
		3, 'comma,comma', -12.34
	);
	`
	getAllUsers = `select ID, Name, Balance from Users`
)

type User struct {
	ID      int
	Name    string
	Balance float64
}

var db *sql.DB

func TestMain(m *testing.M) {
	db = setup(pathToSqliteFile)
	defer db.Close()
	code := m.Run()
	os.Remove(pathToSqliteFile)
	os.Exit(code)
}

func setup(path string) *sql.DB {
	// Connecting to database.
	os.Remove(path)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(populateDB)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func TestExportCSV(t *testing.T) {
	r, err := db.Query(getAllUsers)
	if err != nil {
		t.Fatal(err)
	}

	csv, err := RowsToCSV(r)
	if err != nil {
		t.Fatal(err)
	}

	if expected := allUsersCSV(); expected != csv {
		t.Errorf("Expected:\n%v\nGot:\n%v\n", expected, csv)
	}
}

func allUsersCSV() string {
	rows, err := db.Query(getAllUsers)
	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	c := csv.NewWriter(w)
	c.Write([]string{"ID", "Name", "Balance"})

	var row User
	for rows.Next() {
		err := rows.Scan(&row.ID, &row.Name, &row.Balance)
		if err != nil {
			log.Fatal(err)
		}
		c.Write([]string{
			strconv.FormatInt(int64(row.ID), 10),
			row.Name,
			strconv.FormatFloat(row.Balance, 'f', -1, 64),
		})
	}

	c.Flush()
	return b.String()
}
