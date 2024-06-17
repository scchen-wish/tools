package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: compare_db [id_file] [table_name] [db1_connection] [db2_connection]")
		return
	}

	idFile := os.Args[1]
	tableName := os.Args[2]
	db1Conn := os.Args[3]
	db2Conn := os.Args[4]

	ids, err := readIDsFromFile(idFile)
	if err != nil {
		fmt.Println("Error reading IDs from file:", err)
		return
	}

	db1, err := sql.Open("mysql", db1Conn)
	if err != nil {
		fmt.Println("Error connecting to database 1:", err)
		return
	}
	defer db1.Close()

	db2, err := sql.Open("mysql", db2Conn)
	if err != nil {
		fmt.Println("Error connecting to database 2:", err)
		return
	}
	defer db2.Close()

	for _, id := range ids {
		rows1, err := db1.Query("SELECT * FROM "+tableName+" WHERE id = ?", id)
		if err != nil {
			fmt.Println("Error querying database 1:", err)
			continue
		}
		defer rows1.Close()

		rows2, err := db2.Query("SELECT * FROM "+tableName+" WHERE id = ?", id)
		if err != nil {
			fmt.Println("Error querying database 2:", err)
			continue
		}
		defer rows2.Close()

		if !compareRows(rows1, rows2) {
			fmt.Printf("Fields do not match for ID: %s\n", id)
		} else {
			fmt.Printf("Fields match for ID: %s\n", id)
		}
	}
}

func readIDsFromFile(idFile string) ([]string, error) {
	var ids []string

	file, err := os.Open(idFile)
	if err != nil {
		return ids, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ids = append(ids, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return ids, err
	}

	return ids, nil
}

func compareRows(rows1, rows2 *sql.Rows) bool {
	for rows1.Next() {
		var values1 []interface{}
		var values2 []interface{}

		if err := rows1.Scan(values1...); err != nil {
			fmt.Println("Error scanning row from database 1:", err)
			continue
		}

		if err := rows2.Scan(values2...); err != nil {
			fmt.Println("Error scanning row from database 2:", err)
			continue
		}

		if len(values1) != len(values2) {
			return false
		}

		for i := 0; i < len(values1); i++ {
			if values1[i] != values2[i] {
				return false
			}
		}
	}

	return true
}
