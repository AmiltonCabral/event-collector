package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var eventsByType map[string][]map[string]interface{}
var quit = make(chan bool)

func main() {
	eventsByType = make(map[string][]map[string]interface{})

	http.HandleFunc("/events", handleEvents)

	go monitorUserInput()

	go func() {
		log.Println("Server started on port 3080")
		if err := http.ListenAndServe(":3080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	<-quit
	log.Println("Shutting down the application...")
}

func monitorUserInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		switch input {
		case "quit":
			quit <- true
			return
		case "save-csv":
			saveCSV()
		default:
			log.Println("Unknown command. Use 'save-csv' to save data or 'quit' to exit.")
		}
	}
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newEvents []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&newEvents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, event := range newEvents {
		eventType, ok := event["eventType"].(string)
		if !ok {
			log.Println("Event without valid type:", event)
			continue
		}
		eventsByType[eventType] = append(eventsByType[eventType], event)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Received %d events", len(newEvents))))
}

func saveCSV() {
	if len(eventsByType) == 0 {
		log.Println("No events to export")
		return
	}

	csvDir := "csv_exports"
	err := os.MkdirAll(csvDir, os.ModePerm)
	if err != nil {
		log.Println("Error creating directory for CSVs:", err)
		return
	}

	updatedFiles := []string{}

	for eventType, events := range eventsByType {
		if len(events) == 0 {
			continue
		}

		fileName := fmt.Sprintf("%s.csv", eventType)
		filePath := filepath.Join(csvDir, fileName)

		headers := getHeaders(events)
		existingEvents, existingHeaders := readExistingCSV(filePath)

		// Merge existing headers with new headers
		for _, header := range headers {
			if !contains(existingHeaders, header) {
				existingHeaders = append(existingHeaders, header)
			}
		}

		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Error creating/opening CSV file for %s: %v", eventType, err)
			continue
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write header
		writer.Write(existingHeaders)

		// Write existing events
		writeEvents(writer, existingEvents, existingHeaders)

		// Write new events
		writeEvents(writer, events, existingHeaders)

		updatedFiles = append(updatedFiles, fileName)

		// Clear events after writing them to CSV
		eventsByType[eventType] = nil
	}

	if len(updatedFiles) == 0 {
		log.Println("No CSV files were updated")
	} else {
		log.Printf("CSVs updated successfully: %v", updatedFiles)
	}
}

func getHeaders(events []map[string]interface{}) []string {
	headers := make(map[string]bool)
	for _, event := range events {
		for key := range event {
			headers[key] = true
		}
	}
	var headerSlice []string
	for header := range headers {
		headerSlice = append(headerSlice, header)
	}
	return headerSlice
}

func readExistingCSV(filePath string) ([]map[string]string, []string) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil || len(records) == 0 {
		return nil, nil
	}

	headers := records[0]
	var events []map[string]string

	for _, record := range records[1:] {
		event := make(map[string]string)
		for i, value := range record {
			if i < len(headers) {
				event[headers[i]] = value
			}
		}
		events = append(events, event)
	}

	return events, headers
}

func writeEvents(writer *csv.Writer, events interface{}, headers []string) {
	switch evts := events.(type) {
	case []map[string]interface{}:
		for _, event := range evts {
			writeEvent(writer, event, headers)
		}
	case []map[string]string:
		for _, event := range evts {
			writeEvent(writer, event, headers)
		}
	}
}

func writeEvent(writer *csv.Writer, event interface{}, headers []string) {
	row := make([]string, len(headers))
	switch e := event.(type) {
	case map[string]interface{}:
		for i, header := range headers {
			if value, ok := e[header]; ok {
				row[i] = fmt.Sprintf("%v", value)
			}
		}
	case map[string]string:
		for i, header := range headers {
			if value, ok := e[header]; ok {
				row[i] = value
			}
		}
	}
	writer.Write(row)
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
