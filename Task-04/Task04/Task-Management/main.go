package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/task-management/controller"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/login", controller.LoginHandler).Methods("POST")
	//r.Handle("/home", controller.AuthenticationMiddleware(controller.Home)).Methods("GET")
	r.HandleFunc("/logout", controller.Logout).Methods("POST")
	//r.HandleFunc("/home", controller.Home).Methods("GET")
	r.HandleFunc("/task", controller.AddTask).Methods("POST")
	r.HandleFunc("/task", controller.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/mytasks", controller.GetTasks).Methods("GET")
	r.HandleFunc("/task", controller.DeleteTask).Methods("DELETE")
	r.HandleFunc("/tasks", controller.GetAllTasks).Methods("GET")

	//To listen to port
	log.Fatal(http.ListenAndServe(":8080", r))

	//repository.SendEmail()
}
