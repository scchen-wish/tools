package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Check command line arguments
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run script.go <username>:<password>@tcp(<host>:<port>)/<dbname> <table_name> <num_parts>")
		return
	}

	// Get command line arguments
	dsn := os.Args[1]         // Database connection information
	tableName := os.Args[2]   // Table name
	numPartsStr := os.Args[3] // Number of parts

	// Parse numParts argument
	numParts, err := strconv.Atoi(numPartsStr)
	if err != nil || numParts <= 0 {
		log.Fatalf("Invalid number of parts: %s, must be a positive integer", numPartsStr)
	}

	// Connect to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Query for the minimum and maximum id
	var minID, maxID int
	err = db.QueryRow("SELECT MIN(id), MAX(id) FROM "+tableName).Scan(&minID, &maxID)
	if err != nil {
		log.Fatalf("Error querying min and max id: %v", err)
	}

	// Print minID and maxID before starting the loop
	fmt.Printf("minID: %d, maxID: %d\n", minID, maxID)

	// Calculate the range and count rows for each part
	partSize := int(math.Ceil(float64(maxID-minID+1) / float64(numParts)))

	// Record start time
	startTime := time.Now()

	for i := 0; i < numParts; i++ {
		startID := minID + i*partSize
		endID := startID + partSize - 1

		// Print current range and IDs
		fmt.Printf("Processing range %d-%d\n", startID, endID)

		// Query the count of rows for each part
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM "+tableName+" WHERE id BETWEEN ? AND ?", startID, endID).Scan(&count)
		if err != nil {
			log.Printf("Error querying count for range %d-%d: %v", startID, endID, err)
			continue
		}

		fmt.Printf("Range %d-%d: %d rows\n", startID, endID, count)

		// Calculate and print elapsed time
		fmt.Printf("Elapsed time for range %d-%d: %v\n", startID, endID, time.Since(startTime))
	}

	// Print total elapsed time
	fmt.Printf("Total elapsed time: %v\n", time.Since(startTime))
}
