package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func loadDataFromCSV() {
	file, err := os.Open("employees.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	//first row of the CSV file contains headers
	headers, _ := reader.Read()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
			//http.Error(w, "Error Reading file", http.StatusInternalServerError)
		}

		employee := Employee{}
		// Use a switch statement to map each header to the corresponding struct field.
		for i, value := range record {
			switch headers[i] {
			case "EmpID":
				// Parse and assign the ID to the Employee struct.
				employee.ID, _ = strconv.Atoi(value)
			case "FirstName":
				employee.FirstName = value
			case "LastName":
				employee.LastName = value
			case "Email":
				employee.Email = value
			case "PhoneNo":
				employee.PhoneNo = value
			case "Role":
				employee.Role = value
			case "Password":
				employee.Password = value
			case "Salary":
				employee.Salary, _ = strconv.ParseFloat(value, 64)
			}
		}

		// Append the populated Employee struct to the employees slice.
		employees = append(employees, employee)
	}
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	loadDataFromCSV()
	LoggingFunc("Request TO Update Employee")

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Find the employee to update
	var updated bool
	for i, emp := range employees {
		if emp.ID == id {
			// Parse the JSON body to get the updated data
			var updatedEmployee Employee
			if err := json.NewDecoder(r.Body).Decode(&updatedEmployee); err != nil {
				http.Error(w, "Failed to parse JSON body", http.StatusBadRequest)
				return
			}

			// Update the employee record
			employees[i] = updatedEmployee

			// Write the updated data back to the CSV file
			if err := writeDataToCSV(); err != nil {
				http.Error(w, "Failed to update employee", http.StatusInternalServerError)
				LoggingFunc("Failed To Update Employee")
				return
			}
			fmt.Fprint(w, "Employee updated successfully.")
			LoggingFunc("Employee Updated")
			updated = true
			break
		}
	}

	if !updated {
		http.Error(w, "Employee not found", http.StatusNotFound)
	}
}

func writeDataToCSV() error {
	file, err := os.Create("employees.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	// Write the headers
	headers := []string{"EmpID", "FirstName", "LastName", "Email", "PhoneNo", "Role", "Password", "Salary"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write the updated data
	for _, emp := range employees {
		record := []string{
			strconv.Itoa(emp.ID),
			emp.FirstName,
			emp.LastName,
			emp.Email,
			emp.PhoneNo,
			emp.Role,
			emp.Password,
			strconv.FormatFloat(emp.Salary, 'f', -1, 64),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}
