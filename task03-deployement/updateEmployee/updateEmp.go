package function

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
)

func init() {
	functions.HTTP("UpdateEmployee", updateEmployeeDetails)
}

type Employee struct {
	ID        int64   `json:"id"`
	FirstName string  `json:"firstname"`
	LastName  string  `json:"lastname"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	PhoneNo   string  `json:"phone"`
	Role      string  `json:"role"`
	Salary    float64 `json:"salary"`
}

// declaring Constants
const (
	projectId      string = "employee-management-403415"
	collectionName string = "employees"
)

// TO Update The Employee
func updateEmployeeDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logs("Request To Update Employee Details")
	//parsing the employee id from query string
	empIDStr := r.URL.Query().Get("id")
	//converting id to integer
	empID, err := strconv.ParseInt(empIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid emp id")
		return
	}

	// Parse the updated employee details from the request body
	var updatedEmployee *Employee
	//decoding the json body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatedEmployee); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Failed To Decode Json")
		logs("Failed TO Decode Body")
		return
	}

	// calling function to Update  employee details in Firestore
	msg, err := Update(empID, updatedEmployee)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		logs("Failed TO Update Employee Details")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Employee Updated Succesfully")
	logs("Employee Updated Succesfully")
}

// TO Update the Employee From its id
func Update(id int64, updatedEmp *Employee) (string, error) {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		logs("Failed to create Firestore Client")
		return "Failed to create Firestore Client", err
	}
	defer client.Close()
	// Update the employee details in Firestore

	employeesCollection := client.Collection(collectionName)
	// Query Firestore to find the employee with the given email.
	iter := employeesCollection.Where("id", "==", id).Documents(ctx)
	d, err := iter.Next()
	if err != nil {
		// Handle errors, such as when the employee is not found.
		if err == iterator.Done {
			logs("employee not found to update")
			return "employee not found to update", err // Employee not found
		}
		return "employee not found to update", err
	}
	//TO Convert password in encrypted form
	pass, err := EncryptedPass(updatedEmp.Password)
	if err != nil {
		logs("Error While Encrypting the password")
		return "Error While Encrypting the password", err
	}

	d.Ref.Update(ctx, []firestore.Update{
		{Path: "firstname", Value: updatedEmp.FirstName},
		{Path: "lastname", Value: updatedEmp.LastName},
		{Path: "email", Value: updatedEmp.Email},
		{Path: "password", Value: pass},
		{Path: "phone", Value: updatedEmp.PhoneNo},
		{Path: "salary", Value: updatedEmp.Salary},
	})

	return "", nil
}

func EncryptedPass(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func logs(msg string) {
	// //will open the file | 0644 wil give permission to read & write the file
	// file, err := os.OpenFile("C:/Users/HP/OneDrive/Desktop/task03-deployement/logs/logFile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// //closing the file
	// defer file.Close()
	// log.SetOutput(file)
	log.Println(msg)
}
