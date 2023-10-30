package controllers

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// helper func to handle error
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

//method to check if the data sent from end user is empty or not

func (e *Employee) IsEmpty() bool {
	return e.FirstName == "" && e.LastName == "" && e.Email == ""
}

// function to read csv file
func readCSVFile() ([][]string, error) {
	file, err := os.Open(csvFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// function to write csv file
func writeCSVFile(records [][]string) error {
	file, err := os.Create(csvFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	if err := csvWriter.WriteAll(records); err != nil {
		return err
	}

	csvWriter.Flush()
	return csvWriter.Error()
}

// TO find max employee id
func MaxID() int {
	// Open and read the CSV file
	file, err := os.Open("employees.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return -1
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Initialize variables to track the maximum employee ID
	maxEmployeeID := 0

	// Loop through the CSV records
	for {
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}

		//employee ID is in the first column
		id, err := strconv.Atoi(record[0])
		if err != nil {
			fmt.Println("Error parsing ID:", err)
			continue
		}

		// Check if this ID is greater than the current maximum
		if id > maxEmployeeID {
			maxEmployeeID = id
		}
	}
	return maxEmployeeID
}

//TO generate encrypted passord

func EncryptedPass(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
