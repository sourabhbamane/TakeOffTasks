package controllers

import "time"

type Employee struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Password  string
	PhoneNo   string
	Role      string
	Salary    float64
	DOB       Date
}

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

//creating memory datastructure as database

var employees []Employee

func AddingEmp() {
	employees = append(employees, Employee{ID: 1, FirstName: "Sourabh", LastName: "Bamane", Email: "sourabh@gmail.com", Password: "srb123", PhoneNo: "7757976500", Role: "ADMIN", Salary: 30000, DOB: Date{Year: 2000, Month: 8, Day: 29}})
	employees = append(employees, Employee{ID: 2, FirstName: "Prasanna", LastName: "Shekhar", Email: "prasanna@gmail.com", Password: "psn123", PhoneNo: "8855994466", Role: "DEVELOPER", Salary: 50000, DOB: Date{Year: 2000, Month: 10, Day: 25}})
	employees = append(employees, Employee{ID: 3, FirstName: "Meghna", LastName: "Shinde", Email: "meghna@gmail.com", Password: "mgn123", PhoneNo: "8899102530", Role: "MANAGER", Salary: 30000, DOB: Date{Year: 1999, Month: 10, Day: 20}})
	employees = append(employees, Employee{ID: 4, FirstName: "Mahesh", LastName: "Patil", Email: "mahesh@gmail.com", Password: "mhs123", PhoneNo: "950659988", Role: "TESTER", Salary: 50000, DOB: Date{Year: 1998, Month: 2, Day: 29}})

}
