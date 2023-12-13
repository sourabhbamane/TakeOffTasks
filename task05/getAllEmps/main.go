package getAllEmployees

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("GetAllEmployees", getAllEmployees)
}
func main() {

}

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

// ServiceError is used to return business error messages
type ServiceError struct {
	Message string `json:"message"`
}

// declaring Constants
const (
	projectId      string = "employee-management-403415"
	collectionName string = "employees"
)

func getAllEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	//Writing logs in logger file
	logs("Request To Get All Employees")
	//calling service method
	employees, err := FindAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ServiceError{Message: "Error getting the Employees"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(employees)
	logs("employee retrieved succesfully")
}

func FindAll() ([]Employee, error) {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	//ceating slice to store all employee details
	var employees []Employee
	//declaring iterator
	itr := client.Collection(collectionName).Documents(ctx)

	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		emp := Employee{
			ID:        doc.Data()["id"].(int64),
			FirstName: doc.Data()["firstname"].(string),
			LastName:  doc.Data()["lastname"].(string),
			Email:     doc.Data()["email"].(string),
			Password:  doc.Data()["password"].(string),
			Role:      doc.Data()["role"].(string),
			PhoneNo:   doc.Data()["phone"].(string),
			Salary:    doc.Data()["salary"].(float64),
		}
		employees = append(employees, emp)
	}
	return employees, nil
}

func logs(msg string) {
	//will open the file | 0644 wil give permission to read & write the file
	// file, err := os.OpenFile("C:/Users/HP/OneDrive/Desktop/task03-deployement/logs/logFile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// //closing the file
	// defer file.Close()
	// log.SetOutput(file)
	log.Println(msg)
}
