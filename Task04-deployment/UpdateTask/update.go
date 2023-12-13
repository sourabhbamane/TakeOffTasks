package updatetask

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/api/iterator"
)

func init() {
	functions.HTTP("UpdateTask", UpdateTask)
}

type Task struct {
	TaskId      int64     `json:"taskid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
	CreatedBy   string    `json:"createdby"`
}

type CustomError struct {
	Message string `json:"message"`
}

// declaring Constants
const (
	projectId string = "task-management-405310"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var uname string

// function TO update the task details
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS")

	cookie, err := r.Cookie("token")
	if err != nil {
		// If the cookie is not present, return Unauthorized status
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// If there's an error other than the cookie missing, return Bad Request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Extract the token string from the cookie value
	tokenStr := cookie.Value

	// Initialize a Claims struct to hold the claims extracted from the token
	claims := &Claims{}
	//with the help of ParseWithClaims func passing 3 values
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Provide the key (jwtkey) for validating the token signature
		return jwtKey, nil
	})

	uname = claims.Username

	// Checking for errors during token parsing
	if err != nil {
		// If the error is due to an invalid signature, return Unauthorized status
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// If there's any other error during parsing, return Bad Request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the token is valid
	if !token.Valid {
		// If the token is not valid, return Unauthorized status
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//w.Header().Set("Content-Type", "application/json")
	/*
		vars := mux.Vars(r)
		TaskIdStr := vars["id"]
	*/
	w.Header().Set("Content-Type", "application/json")
	TaskIdStr := r.URL.Query().Get("id")
	//converting id to integer
	taskId, err := strconv.ParseInt(TaskIdStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CustomError{Message: "Invalid Task id"})
		return
	}

	// Parse the updated employee details from the request body
	var updatedEmployee Task
	//decoding the json body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatedEmployee); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CustomError{Message: "Failed To Decode Json"})
		//log.Write("Failed TO Decode Body")
		return
	}

	// calling function to Update  employee details in Firestore
	err = Update(taskId, updatedEmployee)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(CustomError{Message: "Failed to update Task details"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CustomError{Message: "Task Updated Successfully"})
	//log.Write("Employee Updated Succesfully")
}

// Function TO update the task
func Update(id int64, updatedTask Task) error {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		//log.Write("Failed to create Firestore Client")
		return err
	}
	defer client.Close()
	// Update the employee details in Firestore

	employeesCollection := client.Collection("tasks")
	// Query Firestore to find the employee with the given email.
	iter := employeesCollection.Where("TaskId", "==", id).Documents(ctx)
	d, err := iter.Next()
	if err != nil {
		// Handle errors, such as when the employee is not found.
		if err == iterator.Done {
			return nil // Employee not found
		}
		return err
	}

	//To Update The Fields in task Collection
	d.Ref.Update(ctx, []firestore.Update{
		{Path: "Status", Value: updatedTask.Status},
		{Path: "Description", Value: updatedTask.Description},
		{Path: "UpdatedAt", Value: time.Now()},
	})

	//To Save The Updted History
	_, _, err = client.Collection("updatehistory").Add(ctx, map[string]interface{}{
		"TaskId":    id,
		"UpdatedBy": uname,
		"Status":    updatedTask.Status,
		"UpdatedAt": time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}
