package getalltasks

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/api/iterator"
)

func init() {
	functions.HTTP("GetAllTasks", GetAllTasks)
}

type Task struct {
	TaskId      int64     `json:"taskid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
	CreatedBy   string    `json:"createdby"`
	AssignedTo  string    `json:"assignedTo"`
}

// declaring Constants
const (
	projectId string = "task-management-405310"
	//collectionName string = "employees2"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

//Function to get all the tasks

func GetAllTasks(w http.ResponseWriter, r *http.Request) {

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
		json.NewEncoder(w).Encode(CustomError{Message: "invalid token"})
		return
	}

	//Calling Method to get All Tasks present in db
	tasks, err := FetchAllTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(CustomError{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func FetchAllTasks() ([]Task, error) {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		//log.Write("Failed to create Firestore Client")
		err = errors.New("failed to create firestore client")
		return nil, err
	}
	defer client.Close()

	//creating slice to store task details
	var tasks []Task

	//declaring iterator
	iter := client.Collection("tasks").Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		t := Task{
			TaskId:      doc.Data()["TaskId"].(int64),
			CreatedBy:   doc.Data()["CreatedBy"].(string),
			Title:       doc.Data()["Title"].(string),
			Description: doc.Data()["Description"].(string),
			Status:      doc.Data()["Status"].(string),
			CreatedAt:   doc.Data()["CreatedAt"].(time.Time),
			UpdatedAt:   doc.Data()["UpdatedAt"].(time.Time),
			AssignedTo:  doc.Data()["AssignedTo"].(string),
		}
		tasks = append(tasks, t)
	}
	return tasks, nil

}

type CustomError struct {
	Message string `json:"message"`
}
