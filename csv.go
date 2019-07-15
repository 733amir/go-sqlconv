package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"

	"database/sql"
	"encoding/csv"
)

// RowsToCSV convert rows and columns of a query result to a CSV format.
func RowsToCSV(rows *sql.Rows) (string, error) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	c := csv.NewWriter(w)

	columns, err := rows.Columns()
	if err != nil {
		return "", errors.New("csv: columns: " + err.Error())
	}
	c.Write(columns)

	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i := range values {
		pointers[i] = &values[i]
	}

	csvRow := make([]string, len(columns))
	for rows.Next() {
		err := rows.Scan(pointers...)
		if err != nil {
			return "", err
		}

		for i := range columns {
			// TODO: Customize for different types.
			switch value := values[i].(type) {
			case []byte:
				csvRow[i] = string(value)
			default:
				csvRow[i] = fmt.Sprint(value)
			}
		}

		c.Write(csvRow)
	}

	c.Flush()
	return b.String(), nil
}
