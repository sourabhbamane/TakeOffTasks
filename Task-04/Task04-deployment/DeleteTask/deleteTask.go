package deletetask

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/api/iterator"
)

func init() {
	functions.HTTP("DeleteTask", deleteTask)
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

// declaring Constants
const (
	projectId string = "task-management-405310"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// to delete the task
func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Content-Type", "application/json")

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

	taskIdStrg := r.URL.Query().Get("id")
	//if employee id is empty then will show the message
	if taskIdStrg == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CustomError{Message: "Missing Task id in the query string"})
		//log.Write("Employee Id is not provided")
		return
	}

	// Converting the employeeIDStr to int64
	taskId, err := strconv.ParseInt(taskIdStrg, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CustomError{Message: "Invalid Task id in the query string"})
		//log.Write("Invalid Emp Id")
		return
	}

	// Calling the function to delete the employee using Firestore
	err = DeleteTask(taskId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(CustomError{Message: err.Error()})
		//log.Write("Failed To Delete Emp")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CustomError{Message: "task deleted Succesfully"})
	//log.Write("Employee Deleted Succesfully")
}

// dunction to delete the task
func DeleteTask(id int64) error {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		//log.Write("Failed to create Firestore Client")
		err = errors.New("failed to create firestore client")
		return err
	}
	defer client.Close()

	taskCollection := client.Collection("tasks")
	iter := taskCollection.Where("TaskId", "==", id).Documents(ctx)

	doc, err := iter.Next()
	if err != nil {
		if err == iterator.Done {
			err = errors.New("task not found")
			return err
		}
		return err
	}

	doc.Ref.Delete(ctx)
	return nil
}

type CustomError struct {
	Message string `json:"message"`
}
