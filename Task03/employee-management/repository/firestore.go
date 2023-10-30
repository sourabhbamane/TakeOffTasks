package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"employee.com/myapp/entity"
	"employee.com/myapp/log"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type repo struct{}

// creating new repository
func NewFireStoreRepo() EmpRepository {
	return &repo{}
}

// declaring Constants
const (
	projectId      string = "employee-management-bf0ee"
	collectionName string = "employees"
)

// TO generate encrypted passord
func EncryptedPass(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Buisness logic for adding the employee details
func (*repo) Save(emp *entity.Employee) (*entity.Employee, error) {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Write("Failed to create Firestore Client")
		return nil, err
	}
	defer client.Close()

	//TO Convert password in encrypted form
	pass, err := EncryptedPass(emp.Password)
	if err != nil {
		log.Write("Error While Encrypting the password")
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
		log.Write("Failed to Add Employee")
		return nil, err
	}
	return emp, nil
}

// Buisness Logic TO Get All EMployees Present in Db
func (*repo) FindAll() ([]entity.Employee, error) {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Write("Failed to create Firestore Client")

		return nil, err
	}
	defer client.Close()

	//ceating slice to store all employee details
	var employees []entity.Employee
	//declaring iterator
	itr := client.Collection(collectionName).Documents(ctx)

	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Write("Failed to create Firestore Client")
			return nil, err
		}

		emp := entity.Employee{
			ID:        doc.Data()["id"].(int64),
			FirstName: doc.Data()["firstname"].(string),
			LastName:  doc.Data()["lastname"].(string),
			Email:     doc.Data()["email"].(string),
			Password:  doc.Data()["password"].(string),
			Role:      doc.Data()["role"].(string),
			PhoneNo:   doc.Data()["phone"].(string),
			Salary:    doc.Data()["salary"].(int64),
		}
		employees = append(employees, emp)
	}
	return employees, nil
}

// To delete Employee
func (*repo) DeleteEmployee(id int64) (string, error) {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId, option.WithCredentialsFile("./serviceAccountKey.json"))
	if err != nil {
		log.Write("Failed to create Firestore Client")
		return "Failed To delete Employee", err
	}
	defer client.Close()

	employeesCollection := client.Collection(collectionName)

	// Query Firestore to find the employee with the given email.
	iter := employeesCollection.Where("id", "==", id).Documents(ctx)
	d, err := iter.Next()

	if err != nil {
		// Handle errors, such as when the employee is not found.
		if err == iterator.Done {
			return "employee not found", nil // Employee not found
		}
		return "", err
	}
	d.Ref.Delete(ctx)
	return "Employee Deleted Successfully", nil

}

// TO Update the Employee From its id
func (*repo) UpdateEmp(id int64, updatedEmp *entity.Employee) error {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId, option.WithCredentialsFile("./serviceAccountKey.json"))
	if err != nil {
		log.Write("Failed to create Firestore Client")
		return err
	}
	defer client.Close()
	// Update the employee details in Firestore

	employeesCollection := client.Collection(collectionName)
	// Query Firestore to find the employee with the given email.
	iter := employeesCollection.Where("id", "==", id).Documents(ctx)
	d, err := iter.Next()
	if err != nil {
		// Handle errors, such as when the employee is not found.
		if err == iterator.Done {
			return nil // Employee not found
		}
		return err
	}
	//TO Convert password in encrypted form
	pass, err := EncryptedPass(updatedEmp.Password)
	if err != nil {
		log.Write("Error While Encrypting the password")
		return err
	}

	d.Ref.Update(ctx, []firestore.Update{
		{Path: "firstname", Value: updatedEmp.FirstName},
		{Path: "lastname", Value: updatedEmp.LastName},
		{Path: "email", Value: updatedEmp.Email},
		{Path: "password", Value: pass},
		{Path: "phone", Value: updatedEmp.PhoneNo},
		{Path: "salary", Value: updatedEmp.Salary},
	})

	return nil
}

func (*repo) FindEmployees(filters map[string]string) ([]entity.Employee, error) {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId, option.WithCredentialsFile("./serviceAccountKey.json"))
	if err != nil {
		log.Write("Failed to create Firestore Client")
		return nil, err
	}
	defer client.Close()

	collection := client.Collection(collectionName)
	query := collection.Query

	for field, value := range filters {
		query = query.Where(field, "==", value)
	}

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}

	var employees []entity.Employee
	for _, doc := range docs {
		var e entity.Employee
		doc.DataTo(&e)
		employees = append(employees, e)
	}

	return employees, nil
}

func (*repo) FindById(id int64) (*entity.Employee, error) {

	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId, option.WithCredentialsFile("./serviceAccountKey.json"))
	if err != nil {
		log.Write("Failed to create Firestore Client")
		return nil, err
	}
	defer client.Close()
	// Reference to the Firestore collection where employees are stored.
	employeesCollection := client.Collection(collectionName)

	// Query Firestore to find the employee with the given email.
	iter := employeesCollection.Where("id", "==", id).Documents(ctx)
	doc, err := iter.Next()

	if err != nil {
		// Handle errors, such as when the employee is not found.
		if err == iterator.Done {
			return nil, nil // Employee not found
		}
		return nil, err
	}

	// Create an Employee struct and return it.
	employee := entity.Employee{
		ID:        doc.Data()["id"].(int64),
		FirstName: doc.Data()["firstname"].(string),
		LastName:  doc.Data()["lastname"].(string),
		Email:     doc.Data()["email"].(string),
		Password:  doc.Data()["password"].(string),
		Role:      doc.Data()["role"].(string),
		PhoneNo:   doc.Data()["phone"].(string),
		Salary:    doc.Data()["salary"].(int64),
	}

	return &employee, nil
}
