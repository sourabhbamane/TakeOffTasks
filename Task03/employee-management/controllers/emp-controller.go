package controllers

import (
	"encoding/json"

	"net/http"
	"strconv"

	"employee.com/myapp/entity"
	"employee.com/myapp/errors"
	"employee.com/myapp/log"
	"employee.com/myapp/service"
	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

type controller struct{}

var (
	empService service.EmpService
)

type EmpController interface {
	GetAllEmployees(w http.ResponseWriter, r *http.Request)
	AddEmployee(w http.ResponseWriter, r *http.Request)
	DeleteEmployee(w http.ResponseWriter, r *http.Request)
	UpdateEmployeeDetails(w http.ResponseWriter, r *http.Request)
	SearchEmployees(w http.ResponseWriter, r *http.Request)
}

func NewEmpController(serv service.EmpService) EmpController {
	empService = serv
	return &controller{}
}

// Func to get All Employees
// gcloud functions deploy GetAllEmployees --runtime go116  --trigger-http --entry-point GetAllEmployees
func (*controller) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//Writing logs in logger file
	log.Write("Request To Get All Employee")
	//calling service method
	employees, err := empService.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Error getting the Employees"})
		log.Write("Failed To Get All Employee")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(employees)
	log.Write("Fetched All Employees")
}

// Func TO Add new Employee
func (*controller) AddEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Write("Request To Add Employee")
	//declaring var of type employee
	var employee entity.Employee
	//now decoding the json body To object
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Error unmarshalling data"})
		log.Write("Error While Marshelling Data")
		return
	}

	//Checking if Fields are empty or not in json body
	err1 := empService.Validate(&employee)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: err1.Error()})
		return
	}
	result, err2 := empService.Create(&employee)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Error saving the employee"})
		log.Write("Failed TO Save Employee")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	log.Write("Employee Added Succesfully")
}

// Func to Delete the Employee
func (*controller) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Write("Request To Delete Employee")
	// Parse the query string
	employeeIDStr := r.URL.Query().Get("id")
	//if employee id is empty then will show the message
	if employeeIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Missing employee id in the query string"})
		log.Write("Employee Id is not provided")
		return
	}

	// Converting the employeeIDStr to int64
	employeeID, err := strconv.ParseInt(employeeIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Invalid employee id in the query string"})
		log.Write("Invalid Emp Id")
		return
	}

	// Calling the function to delete the employee using Firestore
	msg, err := empService.DeleteEmp(employeeID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: msg})
		log.Write("Failed To Delete Emp")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
	log.Write("Employee Deleted Succesfully")
}

// TO Update The Employee
func (*controller) UpdateEmployeeDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Write("Request To Update Employee Details")
	//parsing the employee id from query string
	empIDStr := r.URL.Query().Get("id")
	//converting id to integer
	empID, err := strconv.ParseInt(empIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Invalid emp id"})
		return
	}

	// Parse the updated employee details from the request body
	var updatedEmployee *entity.Employee
	//decoding the json body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatedEmployee); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Failed To Decode Json"})
		log.Write("Failed TO Decode Body")
		return
	}

	// calling function to Update  employee details in Firestore
	if err := empService.Update(empID, updatedEmployee); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Failed to update employee details"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Employee Updated Succesfully")
	log.Write("Employee Updated Succesfully")
}

// Search Employee On the Basis OF email,firstName,LastName,role
func (*controller) SearchEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Write("Request To Search Employee")
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
			json.NewEncoder(w).Encode(errors.ServiceError{Message: "Invalid 'id' parameter"})
			return
		}
		employees, _ := empService.FindEmpByID(id)
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
	employees, err := empService.SearchEmp(filters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Internal Server Error"})
		return
	}
	// Serialize the result to JSON and send it as a response.
	if employees != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(employees)
		log.Write(" Got The Employee Details ")
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Employee not found")
		log.Write("Employee not found")
	}
}
