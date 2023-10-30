package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

var csvFileName = "employees.csv"

// Mutex for file locking
var csvLock sync.Mutex

// Function to delete employee by their id
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	LoggingFunc("Request TO Delete Employee")
	fail := "Failed TO delete Employee"
	//vars := mux.Vars(r)
	//employeeID, err := strconv.Atoi(vars["id"])
	employeeID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		LoggingFunc("Invalid Employee ID")
		return
	}

	// Acquire the file lock to prevent race conditions when accessing the CSV file
	csvLock.Lock()
	defer csvLock.Unlock()

	// Read the CSV file
	records, err := readCSVFile()
	if err != nil {
		http.Error(w, "Error reading CSV file", http.StatusInternalServerError)
		LoggingFunc(fail)
		return
	}

	// Remove the record with the specified employee ID
	updatedRecords, deleted := removeRecordByEmployeeID(records, employeeID)

	if !deleted {
		http.Error(w, "Employee record not found", http.StatusNotFound)
		LoggingFunc("Employee Not Found")
		return
	}

	// Write the updated records back to the CSV file
	if err := writeCSVFile(updatedRecords); err != nil {
		http.Error(w, "Error writing CSV file", http.StatusInternalServerError)
		LoggingFunc(fail)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Employee record deleted successfully")
	LoggingFunc("Employee Deleted Succesfully")
}

func removeRecordByEmployeeID(records [][]string, employeeID int) ([][]string, bool) {
	updatedRecords := make([][]string, 0)
	deleted := false

	for _, record := range records {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			continue // Skip invalid records
		}

		if id == employeeID {
			deleted = true
		} else {
			updatedRecords = append(updatedRecords, record)
		}
	}

	return updatedRecords, deleted
}
