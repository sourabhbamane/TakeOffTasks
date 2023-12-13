package addtask

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/gomail.v2"
)

func init() {
	functions.HTTP("AddTask", AddTask)
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
)

var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var uname string

// function to create task
func AddTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	// Get the value from the "token" cookie in the incoming request
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

	//declaring variable of type Tasks
	var task Task
	//now decoding the json body To object
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(CustomError{Message: "Error unmarshalling data"})
		//log.Write("Error While Marshelling Data")
		fmt.Println(err)
		return
	}

	//calling method of to save the task details
	result, err := CreateTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(CustomError{Message: "Error saving the Task"})
		//log.Write("Failed TO Save Employee")
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	stringId := strconv.Itoa(int(result.TaskId))
	json.NewEncoder(w).Encode(CustomError{Message: "task added succesfully of id: " + stringId})
	//log.Write("Employee Added Succesfully")

}

// func to create new task
func CreateTask(task *Task) (*Task, error) {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	//username := ctx.Value("Username").(string)
	username := uname
	task.TaskId = rand.Int63()

	_, _, err = client.Collection("tasks").Add(ctx, map[string]interface{}{
		"TaskId":      task.TaskId,
		"Title":       task.Title,
		"Description": task.Description,
		"Status":      task.Status,
		"CreatedAt":   time.Now(),
		"UpdatedAt":   time.Now(),
		"CreatedBy":   username,
		"AssignedTo":  task.AssignedTo,
	})
	if err != nil {
		return nil, err
	}
	SendEmail(username, task.AssignedTo, task.Description)
	return task, nil
}

type CustomError struct {
	Message string `json:"message"`
}

func SendEmail(From string, To string, task string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "managementtask557@gmail.com")
	m.SetHeader("To", To)
	m.SetHeader("Subject", "Task Assignment")
	m.SetBody("text", "Hello,\nNew Task Assignment,\nAssigned By: "+From+"\nTask Description: "+task)
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, "managementtask557@gmail.com", "mbsj nuwf valq mwiu")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

}
