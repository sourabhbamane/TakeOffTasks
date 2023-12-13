package deleteEmp

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"google.golang.org/api/iterator"
)

func init() {
	functions.HTTP("DeleteEmployee", deleteEmployee)
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

// ServiceError is used to return business error messages
type ServiceError struct {
	Message string `json:"message"`
}

// declaring Constants
const (
	projectId      string = "employee-management-403415"
	collectionName string = "employees"
)

// Func to Delete the Employee
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	logs("Request To Delete Employee")
	// Parse the query string

	// Get the "id" parameter from the URL
	employeeIDStr := r.URL.Query().Get("id")
	// vars := mux.Vars(r)
	// employeeIDStr := vars["id"]
	//if employee id is empty then will show the message
	if employeeIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ServiceError{Message: "Missing employee id in the query string"})
		//log.Write("Employee Id is not provided")
		return
	}

	// Converting the employeeIDStr to int64
	employeeID, err := strconv.ParseInt(employeeIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ServiceError{Message: "Invalid employee id in the query string"})
		logs("Invalid Emp Id")
		return
	}

	// Calling the function to delete the employee using Firestore
	err = DeleteEmp(employeeID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ServiceError{Message: err.Error()})
		logs("Failed To Delete Emp")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ServiceError{Message: "Employee Deleted Succesfully"})
	logs("Employee Deleted Succesfully")
}

// To delete Employee
func DeleteEmp(id int64) error {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		//log.Write("Failed to create Firestore Client")
		err = errors.New("failed to delete employee")
		return err
	}
	defer client.Close()

	employeesCollection := client.Collection(collectionName)

	// Query Firestore to find the employee with the given email.
	iter := employeesCollection.Where("id", "==", id).Documents(ctx)
	d, err := iter.Next()

	if err != nil {
		// Handle errors, such as when the employee is not found.
		if err == iterator.Done {
			err = errors.New("employee not found")
			return err // Employee not found
		}
		return err
	}
	d.Ref.Delete(ctx)
	return nil

}

func logs(msg string) {
	//will open the file | 0644 wil give permission to read & write the file
	// file, err := os.OpenFile("C:/Users/HP/OneDrive/Desktop/task03-deployement/logs/logFile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// //closing the file
	// defer file.Close()
	// log.SetOutput(file)
	// log.Println(msg)

	log.Println(msg)
}
