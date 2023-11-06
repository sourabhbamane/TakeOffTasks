package function

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"golang.org/x/crypto/bcrypt"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("AddEmployee", addEmployee)
}

// var c = cors.New(cors.Options{ //setting up a CORS middleware for HTTP server
// 	AllowedOrigins: []string{"*"}, //allowedorigins option specifies which websuite are allowed to make requests to your server
// 	//*allows requests from any origin
// 	AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
// 	AllowedHeaders: []string{"Content-Type"}, //content-type in the request to tell about content
// })

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

// // Func TO Add new Employee
func addEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logs("Request To Add Employee")
	//declaring var of type employee
	var employee Employee
	//now decoding the json body To object
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error unmarshalling data")
		logs("Error While Marshelling Data")
		return
	}

	//Checking if Fields are empty or not in json body
	err1 := Validate(&employee)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err1.Error())
		return
	}

	result, err2 := Save(&employee)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err2.Error())
		logs("Failed TO Save Employee")
		return
	}
	w.WriteHeader(http.StatusCreated)

	stringId := strconv.Itoa(int(result.ID))
	json.NewEncoder(w).Encode("Employee Added Succesfully with id: " + stringId)
	logs("Employee Added Succesfully")
}

func Save(emp *Employee) (*Employee, error) {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		logs("Failed to create Firestore Client")
		return nil, err
	}
	defer client.Close()

	emp.ID = rand.Int63()

	//TO Convert password in encrypted form
	pass, err := EncryptedPass(emp.Password)
	if err != nil {
		logs("Error While Encrypting the password")
		return nil, err
	}

	//To Check if email is aready present or not
	present := isEmailExists(client, emp.Email)
	if present {
		err := errors.New("email is already present")
		return nil, err
	}

	_, _, err = client.Collection(collectionName).Add(ctx, map[string]interface{}{
		"id":        emp.ID,
		"firstname": emp.FirstName,
		"lastname":  emp.LastName,
		"email":     emp.Email,
		"password":  pass,
		"phone":     emp.PhoneNo,
		"role":      emp.Role,
		"salary":    emp.Salary,
	})
	if err != nil {
		logs("Failed to Add Employee")
		return nil, err
	}
	return emp, nil
}

// TO generate encrypted passord
func EncryptedPass(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// TO Check If the certain Field Are empyty or not
func Validate(emp *Employee) error {
	if emp == nil {
		err := errors.New("the employee is empty")
		return err
	}
	if emp.FirstName == "" {
		err := errors.New("the employee name is empty")
		return err
	}
	if emp.LastName == "" {
		err := errors.New("the employee Lastname is empty")
		return err
	}
	if emp.Email == "" {
		err := errors.New("the employee email is empty")
		return err
	}
	if emp.PhoneNo == "" {
		err := errors.New("the employee phone number is empty")
		return err
	}
	return nil
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

func isEmailExists(client *firestore.Client, email string) bool {
	ctx := context.Background()

	// Create a reference to the "employees" collection in Firestore
	collection := client.Collection("employees")

	// Create a query to find documents with the matching email
	query := collection.Where("email", "==", email)

	// Get the documents that match the query
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false
	}

	return len(docs) > 0
}
