package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Get Perticular employee based on emailid
func GetEmployeesByEmail(w http.ResponseWriter, r *http.Request) {
	// Acquire the file lock to prevent race conditions when accessing the CSV file
	LoggingFunc("Request To get Employee Details From emailId")
	csvLock.Lock()
	defer csvLock.Unlock()

	// Read the CSV file
	records, err := readCSVFile()
	if err != nil {
		//panic(err)
		http.Error(w, "Error reading CSV file", http.StatusInternalServerError)
		return
	}

	// Get the first name from the query parameter
	email := r.URL.Query().Get("email")

	// Create a slice to store the matching employee data
	matchingEmployees := make([]map[string]string, 0)

	for _, record := range records {
		// Check if the first name matches (case-insensitive)
		//fmt.Sprintf("%d", emp.ID), emp.FirstName, emp.LastName, emp.Email, emp.PhoneNo, emp.Role, emp.Password, fmt.Sprintf("%f", emp.Salary)
		if strings.EqualFold(record[3], email) {
			// Create a map to store employee data (using column names as keys)
			employeeData := map[string]string{
				"empid":     record[0],
				"firstname": record[1],
				"lastname":  record[2],
				"email":     record[3],
				"phone":     record[4],
				"role":      record[5],
				"salary":    record[7],
			}

			matchingEmployees = append(matchingEmployees, employeeData)
		}
	}

	// Convert the matching employees slice to JSON and return it
	w.Header().Set("Content-Type", "application/json")
	if len(matchingEmployees) == 0 {
		//w.WriteHeader(http.StatusNoContent)
		//json.NewEncoder(w).Encode("Sorry No record found of email")
		fmt.Fprintf(w, "Sorry No record found of given ID : %v", email)
		LoggingFunc("No Records Found of id")
		return
	}
	if err := json.NewEncoder(w).Encode(matchingEmployees); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}
