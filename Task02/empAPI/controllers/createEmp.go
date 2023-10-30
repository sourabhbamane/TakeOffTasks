package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// if we send just / request
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Welcome TO Employee Management Application")
}

// TO generate id
var max = MaxID()

// TO generate unique id
func generateUniqueID() int {
	max++
	return max
}

// Track used email addresses
var usedEmails = make(map[string]struct{})

// TO make sure email is unique
func isEmailUnique(email string) bool {
	_, exists := usedEmails[email]
	return !exists
}

// Func To add Employee
// w http.ResponseWriter-- w is an interface provided by the net/http package that allows u to construct an HTTP response
// r *http.Request ---r is a pointer to an http.Request struct, which represents the incoming HTTP request.
// r - this struct contains information about the request, such as the method, headers, URL, and request body.
func AddEmpoyee(w http.ResponseWriter, r *http.Request) {
	//TO add in logger file
	LoggingFunc("Request TO Create Employee")
	okmsg := "Record Added Succesfully"
	fail := "Failed To add Record"

	fmt.Println("Adding Employee")

	//CheckError(err)

	//TO set the headers
	w.Header().Set("Content-Type", "application/json")

	//Checking if bodu is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Empty Body...Send Something")
		LoggingFunc(fail)
	}

	var emp Employee

	//Decoding the data from json
	decoder := json.NewDecoder(r.Body).Decode(&emp)
	//json.NewDecoder(r.Body)--creates a new JSON decoder that reads from the request's body.
	//.Decode(&emp): Decodes the JSON data from the request body into the newEmployee variable and & we passing pointer reference to emp
	if decoder != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed To Decode JSON data : %v", decoder)
		LoggingFunc(fail)
		return
	}
	//if email is repeated
	if !isEmailUnique(emp.Email) {
		fmt.Fprintf(w, "Email is alrady exists try another")
		LoggingFunc(fail)
		return
	}
	//TO Add Data in csv file if the file is not present then we will create it
	//if file is present then we simply append the data
	file, err := os.OpenFile("employees.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to open csv file: %v", err)
		LoggingFunc(fail)
		return
	}
	//closing the file
	defer file.Close()

	//rand.Seed(time.Now().UnixNano())

	//TO Add empoyee data
	emp.ID = generateUniqueID()
	//To Add Password
	pass, err := EncryptedPass(emp.Password)
	if err != nil {
		fmt.Fprint(w, "Error While Encrypting the password", http.StatusInternalServerError)
		return
	}
	//TO Add Data from slice to .csv file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	//To set the headers in csv file
	header := []string{"EmpID", "FirstName", "LastName", "Email", "PhoneNo", "Role", "Password", "Salary"}
	f, _ := file.Stat()
	if f.Size() == 0 {
		writer.Write(header)
	}
	//CheckError(err)
	data := []string{fmt.Sprintf("%d", emp.ID), emp.FirstName, emp.LastName, emp.Email, emp.PhoneNo, emp.Role, pass, fmt.Sprintf("%f", emp.Salary)}

	if err := writer.Write(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed To Wrte in csv file : %v", err)
		LoggingFunc(fail)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Employee Added Succesfully")
	LoggingFunc(okmsg)

	json.NewEncoder(w).Encode(emp)
	//return

	// Mark the email as used
	usedEmails[emp.Email] = struct{}{}
}
