package function

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"google.golang.org/api/iterator"
)

func init() {
	functions.HTTP("SearchEmployee", searchEmployees)
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

// Search Employee On the Basis OF email,firstName,LastName,role
func searchEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logs("Request To Search Employee")
	queryParams := r.URL.Query()
	filters := make(map[string]string)

	// Parse query parameters
	for key, values := range queryParams {
		if len(values) > 0 {
			filters[key] = values[0]
		}
	}

	// Convert "id" to int64
	idStr, exists := filters["id"]
	//if id exists in query string then we call find by id method
	if exists {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Invalid 'id' parameter")
			return
		}
		employees, _ := FindById(id)
		if employees != nil {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(employees)
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("Employee not found")
			return
		}

	}

	//calling function to search employees from firestore
	employees, err := FindEmployees(filters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal Server Error")
		return
	}
	// Serialize the result to JSON and send it as a response.
	if employees != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(employees)
		logs(" Got The Employee Details ")
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Employee not found")
		logs("Employee not found")
	}
}

func FindEmployees(filters map[string]string) ([]Employee, error) {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		logs("Failed to create Firestore Client")
		return nil, err
	}
	// Close the Firestore client when done
	defer client.Close()

	collection := client.Collection(collectionName)
	// Create a query for the collection
	query := collection.Query

	// Apply filters to the query
	for field, value := range filters {
		query = query.Where(field, "==", value)
	}
	// Execute the query and get the documents
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}

	var employees []Employee
	// Iterate over the documents, convert to Employee structs, and append to the result
	for _, doc := range docs {
		// var e Employee
		// doc.DataTo(&e)
		// employees = append(employees, e)

		emp := Employee{
			ID:        doc.Data()["id"].(int64),
			FirstName: doc.Data()["firstname"].(string),
			LastName:  doc.Data()["lastname"].(string),
			Email:     doc.Data()["email"].(string),
			Password:  doc.Data()["password"].(string),
			Role:      doc.Data()["role"].(string),
			PhoneNo:   doc.Data()["phone"].(string),
			Salary:    doc.Data()["salary"].(float64),
		}
		employees = append(employees, emp)
	}

	return employees, nil
}

func FindById(id int64) (*Employee, error) {

	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		logs("Failed to create Firestore Client")
		return nil, err
	}
	defer client.Close()
	// Reference to the Firestore collection where employees are stored.
	employeesCollection := client.Collection(collectionName)

	// Query Firestore to find the employee with the given email.
	iter := employeesCollection.Where("id", "==", id).Documents(ctx)
	doc, err := iter.Next()

	if err != nil {
		// Handle errors, such as when the employee is not found.
		if err == iterator.Done {
			return nil, nil // Employee not found
		}
		return nil, err
	}

	// Create an Employee struct and return it.
	employee := Employee{
		ID:        doc.Data()["id"].(int64),
		FirstName: doc.Data()["firstname"].(string),
		LastName:  doc.Data()["lastname"].(string),
		Email:     doc.Data()["email"].(string),
		Password:  doc.Data()["password"].(string),
		Role:      doc.Data()["role"].(string),
		PhoneNo:   doc.Data()["phone"].(string),
		Salary:    doc.Data()["salary"].(float64),
	}

	return &employee, nil
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
