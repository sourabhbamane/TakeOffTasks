package main

import (
	"fmt"
	"log"
	"net/http"

	"context"

	"employee.com/myapp/controllers"
	"employee.com/myapp/repository"
	"employee.com/myapp/service"
	"github.com/gorilla/mux"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

var (
	empRepository repository.EmpRepository  = repository.NewFireStoreRepo()
	empService    service.EmpService        = service.NewEmpService(empRepository)
	empController controllers.EmpController = controllers.NewEmpController(empService)
)

func main() {
	opt := option.WithCredentialsFile("./serviceAccountKey.json")
	_, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v", err)
	}
	r := mux.NewRouter()

	r.HandleFunc("/employees", empController.GetAllEmployees).Methods("GET")
	r.HandleFunc("/employee", empController.AddEmployee).Methods("POST")
	r.HandleFunc("/employee", empController.DeleteEmployee).Methods("DELETE")
	r.HandleFunc("/employee", empController.UpdateEmployeeDetails).Methods("PUT")
	r.HandleFunc("/employees/search", empController.SearchEmployees).Methods("GET")

	//To listen to port
	log.Fatal(http.ListenAndServe(":8080", r))
}
