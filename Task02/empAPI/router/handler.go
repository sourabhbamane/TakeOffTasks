package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sourabhbamane/restapi/controllers"
)

func Handlers() {

	/*
		//controllers.LoadCSVFile("employees.csv")

		// sw := gin.New()
		// sw.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		//s := gin.Default()

		// Your API routes go here

		// Serve Swagger UI on a specific route
		// s.GET("/api-docs", func(c *gin.Context) {
		//     c.HTML(http.StatusOK, "C:/swagger/swagger-editor", nil)
		// })
	*/
	r := mux.NewRouter()

	// docs := middleware.NewSwaggerUI(middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"})
	// r.PathPrefix("/swagger").Handler(docs)

	//routing
	r.HandleFunc("/employee/", controllers.HomeHandler)
	//to add employee
	r.HandleFunc("/employee/add", controllers.AddEmpoyee).Methods("POST")
	//to get all employees
	r.HandleFunc("/employees/get", controllers.GetAllEmployees).Methods("GET")
	//TO Update the employee details
	r.HandleFunc("/employee/update", controllers.UpdateEmployee).Methods("PUT")
	//TO delete the employee
	r.HandleFunc("/employee/delete", controllers.DeleteEmployee).Methods("DELETE")
	//TO get employee by their First name
	r.HandleFunc("/employees/fname", controllers.GetEmployeesByFirstName).Methods("GET")
	//TO get employee by their Last name
	r.HandleFunc("/employees/lname", controllers.GetEmployeesByLastName).Methods("GET")
	//TO get employee by their Email
	r.HandleFunc("/employee/email", controllers.GetEmployeesByEmail).Methods("GET")
	//TO get Employee by their id
	r.HandleFunc("/employee/id", controllers.GetEmployeesByID).Methods("GET")
	//TO get Employee by their role
	r.HandleFunc("/employees/role", controllers.GetEmployeesByRole).Methods("GET")

	//To listen to port
	log.Fatal(http.ListenAndServe(":3000", r))

}
