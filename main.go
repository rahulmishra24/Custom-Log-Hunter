// Import necessary packages

package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"

    "github.com/itchyny/gojq"
)

// LogEntry represents a log entry
type LogEntry struct {
    Timestamp string `json:"timestamp"`
    Message   string `json:"message"`
    // Add more fields as needed
}

// FilenameLogMap maps filenames to log entries
type FilenameLogMap map[string][]LogEntry

func main() {
    // Define command-line flags
    queryFlag := flag.String("query", "", "Gojq query to filter log entries")
    directoryFlag := flag.String("dir", ".", "Directory containing log files")
    flag.Parse()

    // Check if the query flag is provided
    if *queryFlag == "" {
        flag.PrintDefaults()
        os.Exit(1)
    }

    // Initialize map to hold logs grouped by filename
    filenameLogMap := make(FilenameLogMap)

    // Open the directory
    err := filepath.Walk(*directoryFlag, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            log.Printf("Error accessing path %s: %v", path, err)
            return nil
        }

        // Skip directories
        if info.IsDir() {
            return nil
        }

        // Process only text files
        if strings.HasSuffix(path, ".txt") {
            err := processLogFile(path, *queryFlag, filenameLogMap)
            if err != nil {
                log.Printf("Error processing log file %s: %v", path, err)
            }
        }

        return nil
    })

    if err != nil {
        log.Fatalf("Error walking directory: %v", err)
    }

    // Print logs grouped by filename
    printLogsByFilename(filenameLogMap)
}

// processLogFile processes the log file
func processLogFile(filePath string, query string, filenameLogMap FilenameLogMap) error {
    // Open the file
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    // Read the file content
    content, err := io.ReadAll(file)
    if err != nil {
        return fmt.Errorf("failed to read file: %w", err)
    }

    // Unmarshal the JSON data into a slice of LogEntry structs
    var entries []LogEntry
    if err := json.Unmarshal(content, &entries); err != nil {
        return fmt.Errorf("failed to parse JSON: %w", err)
    }

    // Group logs by filename
    filename := filepath.Base(filePath)
	for _,entry := range(entries){
		check:= matchQuery(entry,query)
		if err!=nil{
			fmt.Printf("Some Error Occurred")
		}
		// fmt.Printf(strconv.FormatBool(check))
		if check != false {
			filenameLogMap[filename] = append(filenameLogMap[filename],entry)
		}
	}

    return nil
}

// printLogsByFilename prints logs grouped by filename
func printLogsByFilename(filenameLogMap FilenameLogMap) {
    for filename, logs := range filenameLogMap {
        fmt.Printf("Logs for %s:\n", filename)
        for _, log := range logs {
            fmt.Printf("[%s] %s\n", log.Timestamp, log.Message)
        }
        fmt.Println()
    }
}

// matchQuery checks if a log entry matches the given query
func matchQuery(entry LogEntry, query string) bool {
    // Convert LogEntry struct to JSON
    jsonData, err := json.Marshal(entry)
    if err != nil {
        log.Printf("Error marshalling JSON: %v", err)
        return false
    }

    // Parse JSON data
    var data interface{}
    if err := json.Unmarshal(jsonData, &data); err != nil {
        log.Printf("Error unmarshalling JSON: %v", err)
        return false
    }

    // Parse JQ query
    queryParser, err := gojq.Parse(query)
    if err != nil {
        log.Printf("Error parsing query: %v", err)
        return false
    }

    // Apply JQ query
    iter := queryParser.Run(data)
    for {
        v, ok := iter.Next()
        if !ok {
            break
        }
        if v == true {
            return true
        }
    }
    return false
}

