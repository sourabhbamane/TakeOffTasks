package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/task-management/entity"
	"github.com/task-management/errors"
	"github.com/task-management/repository"
)

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
	var task entity.Task
	//now decoding the json body To object
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.CustomError{Message: "Error unmarshalling data"})
		//log.Write("Error While Marshelling Data")
		fmt.Println(err)
		return
	}

	//calling method of to save the task details
	result, err := repository.NewFireStoreRepo().CreateTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.CustomError{Message: "Error saving the Task"})
		//log.Write("Failed TO Save Employee")
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	stringId := strconv.Itoa(int(result.TaskId))
	json.NewEncoder(w).Encode(errors.CustomError{Message: "task added succesfully of id: " + stringId})
	//log.Write("Employee Added Succesfully")

}

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
		json.NewEncoder(w).Encode(errors.CustomError{Message: "Invalid Task id"})
		return
	}

	// Parse the updated employee details from the request body
	var updatedEmployee *entity.Task
	//decoding the json body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatedEmployee); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.CustomError{Message: "Failed To Decode Json"})
		//log.Write("Failed TO Decode Body")
		return
	}

	// calling function to Update  employee details in Firestore
	err = repository.NewFireStoreRepo().Update(taskId, updatedEmployee)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.CustomError{Message: "Failed to update Task details"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(errors.CustomError{Message: "Task Updated Successfully"})
	//log.Write("Employee Updated Succesfully")
}

//Function TO get Tasks

func GetTasks(w http.ResponseWriter, r *http.Request) {
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

	//Calling Method to get All Tasks created By Logged in employee
	tasks, err := repository.NewFireStoreRepo().GetTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.CustomError{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)

}

// to delete the task
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Content-Type", "application/json")

	taskIdStrg := r.URL.Query().Get("id")
	//if employee id is empty then will show the message
	if taskIdStrg == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.CustomError{Message: "Missing Task id in the query string"})
		//log.Write("Employee Id is not provided")
		return
	}

	// Converting the employeeIDStr to int64
	taskId, err := strconv.ParseInt(taskIdStrg, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.CustomError{Message: "Invalid Task id in the query string"})
		//log.Write("Invalid Emp Id")
		return
	}

	// Calling the function to delete the employee using Firestore
	err = repository.NewFireStoreRepo().DeleteTask(taskId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errors.CustomError{Message: err.Error()})
		//log.Write("Failed To Delete Emp")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(errors.CustomError{Message: "task deleted Succesfully"})
	//log.Write("Employee Deleted Succesfully")
}

//Function to get all the tasks

func GetAllTasks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	//Calling Method to get All Tasks present in db
	tasks, err := repository.NewFireStoreRepo().FetchAllTasks2()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.CustomError{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}
