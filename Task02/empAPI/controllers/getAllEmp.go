package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Function to get all employees
func GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	LoggingFunc("Request To Get All Employees")
	fmt.Println("Get All Records")
	// Acquire the file lock to prevent race conditions when accessing the CSV file
	csvLock.Lock()
	defer csvLock.Unlock()

	// Read the CSV file
	records, err := readCSVFile()
	if err != nil {
		http.Error(w, "Error reading CSV file", http.StatusInternalServerError)
		LoggingFunc("Error reading CSV file")
		return
	}

	// Create a slice to store the employee data
	employees := make([]map[string]string, 0)

	for _, record := range records {
		//6,Priya,Jain,p@gmail.com,9855664488,DEVELOPER,p123,30000.000000

		// Create a map to store employee data (using column names as keys)
		employeeData := map[string]string{
			"empid":     record[0],
			"firstname": record[1],
			"lastname":  record[2],
			"email":     record[3],
			"phone":     record[4],
			"role":      record[5],
			"password":  record[6],
			"salary":    record[7],
		}

		employees = append(employees, employeeData)
	}

	// Convert the employees slice to JSON and return it
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		LoggingFunc("Failed To Get Employee")
		return
	}
}
