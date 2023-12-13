package repository

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/task-management/entity"
	"google.golang.org/api/iterator"
)

// TaskRepository handles task-related database operations
type TaskRepository interface {
	ValidateUser(Username, password string) (*entity.User, error)
	CreateTask(task *entity.Task) (*entity.Task, error)
	Update(id int64, updatedTask *entity.Task) error
	GetTasks() ([]entity.Task, error)
	DeleteTask(id int64) error
	FetchAllTasks() ([]entity.Task, error)
	FetchAllTasks2() ([]entity.Task, error)
}

type repo struct{}

// creating new repository
func NewFireStoreRepo() TaskRepository {
	return &repo{}
}

// declaring Constants
const (
	projectId      string = "task-management-405310"
	collectionName string = "tasks2"
)

var uname string

// function to validate user credentials for authentication
func (*repo) ValidateUser(username, password string) (*entity.User, error) {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		//log.Write("Failed to create Firestore Client")
		return nil, err
	}
	defer client.Close()
	// Reference to the Firestore collection where employees are stored.
	collection := client.Collection("users")
	// Query Firestore to get the user with the provided username and password
	iter := collection.Where("Username", "==", username).Where("Password", "==", password).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, errors.New("user not found")
	}
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return nil, err
	}

	uname = username
	var user entity.User
	err = doc.DataTo(&user)
	if err != nil {
		log.Printf("Error converting user data: %v", err)
		return nil, err
	}

	return &user, nil
}

/*
// To Generate Unique Id
func GenerateID() int64 {
	// Get the current timestamp
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// Generate a random number
	randomNumber := rand.Int63n(10)

	// Combine timestamp and random number to create a unique ID
	id := timestamp*1000 + randomNumber
	return id
}
*/

// func to create new task
func (*repo) CreateTask(task *entity.Task) (*entity.Task, error) {
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

	_, _, err = client.Collection(collectionName).Add(ctx, map[string]interface{}{
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
	if task.AssignedTo != nil {
		SendEmail(username, task.AssignedTo, task.Description)
	}

	return task, nil

}

// Function TO update the task
func (*repo) Update(id int64, updatedTask *entity.Task) error {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		//log.Write("Failed to create Firestore Client")
		return err
	}
	defer client.Close()
	// Update the employee details in Firestore

	employeesCollection := client.Collection("tasks")
	// Query Firestore to find the employee with the given email.
	iter := employeesCollection.Where("TaskId", "==", id).Documents(ctx)
	d, err := iter.Next()
	if err != nil {
		// Handle errors, such as when the employee is not found.
		if err == iterator.Done {
			return nil // Employee not found
		}
		return err
	}

	//To Update The Fields in task Collection
	d.Ref.Update(ctx, []firestore.Update{
		{Path: "Status", Value: updatedTask.Status},
		{Path: "Description", Value: updatedTask.Description},
		{Path: "UpdatedAt", Value: time.Now()},
	})

	//To Save The Updted History
	_, _, err = client.Collection("updatehistory").Add(ctx, map[string]interface{}{
		"TaskId":    id,
		"UpdatedBy": uname,
		"Status":    updatedTask.Status,
		"UpdatedAt": time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

// Function TO get Tasks
func (*repo) GetTasks() ([]entity.Task, error) {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		//log.Write("Failed to create Firestore Client")
		return nil, err
	}
	defer client.Close()
	// Update the employee details in Firestore
	username := uname

	employeesCollection := client.Collection(collectionName)
	// Query Firestore to find the employee with the given email.
	if username == "" {
		err1 := errors.New("please log in first")
		return nil, err1
	}
	itr := employeesCollection.Where("CreatedBy", "==", username).Documents(ctx)

	if err != nil {
		return nil, err
	}
	var tasks []entity.Task

	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		t := entity.Task{
			TaskId:      doc.Data()["TaskId"].(int64),
			CreatedBy:   doc.Data()["CreatedBy"].(string),
			Title:       doc.Data()["Title"].(string),
			Description: doc.Data()["Description"].(string),
			Status:      doc.Data()["Status"].(string),
			CreatedAt:   doc.Data()["CreatedAt"].(time.Time),
			UpdatedAt:   doc.Data()["UpdatedAt"].(time.Time),
			AssignedTo:  doc.Data()["AssignedTo"].([]string),
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// function to delete the task
func (*repo) DeleteTask(id int64) error {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		//log.Write("Failed to create Firestore Client")
		err = errors.New("failed to create firestore client")
		return err
	}
	defer client.Close()

	taskCollection := client.Collection(collectionName)
	iter := taskCollection.Where("TaskId", "==", id).Documents(ctx)

	doc, err := iter.Next()
	if err != nil {
		if err == iterator.Done {
			err = errors.New("task not found")
			return err
		}
		return err
	}

	doc.Ref.Delete(ctx)
	return nil
}

// TO Get All the tasks regardless logged in user
func (*repo) FetchAllTasks() ([]entity.Task, error) {
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
	var tasks []entity.Task

	//declaring iterator
	iter := client.Collection(collectionName).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		t := entity.Task{
			TaskId:      doc.Data()["TaskId"].(int64),
			CreatedBy:   doc.Data()["CreatedBy"].(string),
			Title:       doc.Data()["Title"].(string),
			Description: doc.Data()["Description"].(string),
			Status:      doc.Data()["Status"].(string),
			CreatedAt:   doc.Data()["CreatedAt"].(time.Time),
			UpdatedAt:   doc.Data()["UpdatedAt"].(time.Time),
			AssignedTo:  doc.Data()["AssignedTo"].([]string),
		}
		tasks = append(tasks, t)
	}
	return tasks, nil

}

func (*repo) FetchAllTasks2() ([]entity.Task, error) {
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

	//declaring iterator
	docs, err := client.Collection(collectionName).Documents(ctx).GetAll()

	if err != nil {
		//fmt.Println(err)
		return nil, err
	}

	var tasks []entity.Task
	for _, doc := range docs {
		var e entity.Task
		doc.DataTo(&e)
		tasks = append(tasks, e)
	}
	return tasks, nil

}
