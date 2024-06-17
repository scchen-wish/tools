package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: go run main.go <file_prefix> <output_dir> <mod_value> <concurrency>")
		os.Exit(1)
	}

	filePrefix := os.Args[1]
	outputDir := os.Args[2]
	modValue := os.Args[3]
	concurrency := atoi(os.Args[4])

	// Validate file prefix
	if !strings.HasPrefix(filePrefix, "/") {
		fmt.Println("Error: File prefix must be an absolute path.")
		os.Exit(1)
	}

	// Remove and recreate output directory
	err := os.RemoveAll(outputDir)
	if err != nil {
		fmt.Printf("Error removing output directory: %v\n", err)
		os.Exit(1)
	}
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	files, err := findFiles(filePrefix)
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	// Use a semaphore to limit concurrency
	semaphore := make(chan struct{}, concurrency)

	for _, file := range files {
		semaphore <- struct{}{}
		go func(file string) {
			defer func() { <-semaphore }()
			defer wg.Done()
			processFile(file, outputDir, modValue)
		}(file)
	}

	wg.Wait()
	fmt.Println("All files processed successfully.")
}

func findFiles(prefix string) ([]string, error) {
	var files []string

	err := filepath.Walk(filepath.Dir(prefix), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), filepath.Base(prefix)) {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func processFile(file string, outputDir string, modValue string) {
	fmt.Printf("Processing file: %s\n", file)

	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	counter, validCounter := 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		counter++
		if len(line) < 18 {
			fmt.Printf("Skipping invalid line: %s\n", line)
			continue
		}
		validCounter++
		// Extract the first 16 characters (excluding '0x' prefix) and convert to decimal
		hexValue := line[2:18]
		// Convert hexValue to decimal
		decimal := hexToDecimal(hexValue)

		// Calculate mod value
		modResult := decimal % atoi(modValue)

		// Write result to the appropriate file
		outputFile := fmt.Sprintf("%s/result_%d.txt", outputDir, modResult)
		err := appendToFile(outputFile, line)
		if err != nil {
			fmt.Printf("Error writing to file %s: %v\n", outputFile, err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	fmt.Printf("Finished processing file: %s., total lines:%d, valid lines: %d\n", file, counter, validCounter)
}

func hexToDecimal(hex string) int64 {
	// Convert hexValue to decimal
	var result int64
	for _, c := range hex {
		result = result*16 + int64(c-'0')
	}
	return result
}

func atoi(s string) int64 {
	var result int64
	for _, c := range s {
		result = result*10 + int64(c-'0')
	}
	return result
}

func appendToFile(filename string, line string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(line + "\n")
	if err != nil {
		return err
	}

	return nil
}
