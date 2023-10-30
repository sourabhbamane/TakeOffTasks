package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Get Employees based on their first name
func GetEmployeesByFirstName(w http.ResponseWriter, r *http.Request) {
	// Acquire the file lock to prevent race conditions when accessing the CSV file
	LoggingFunc("Request To get Employee Details From First Name")
	csvLock.Lock()
	defer csvLock.Unlock()

	// Read the CSV file
	records, err := readCSVFile()
	if err != nil {
		http.Error(w, "Error reading CSV file", http.StatusInternalServerError)
		return
	}

	// Get the first name from the query parameter
	firstName := r.URL.Query().Get("firstname")

	// Create a slice to store the matching employee data
	matchingEmployees := make([]map[string]string, 0)

	for _, record := range records {
		// Check if the first name matches (case-insensitive)
		if strings.EqualFold(record[1], firstName) {
			// Create a map to store employee data (using column names as keys)
			employeeData := map[string]string{
				"empid":     record[0],
				"firstname": record[1],
				"lastname":  record[2],
				"email":     record[3],
				"phone":     record[4],
				"role":      record[5],
				"salary":    record[6],
			}

			matchingEmployees = append(matchingEmployees, employeeData)
		}
	}

	// Convert the matching employees slice to JSON and return it
	w.Header().Set("Content-Type", "application/json")
	if len(matchingEmployees) == 0 {
		//json.NewEncoder(w).Encode("Sorry No record found of given ID")
		//	w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "Sorry No record found of given ID : %v", firstName)
		return
	}
	if err := json.NewEncoder(w).Encode(matchingEmployees); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}
