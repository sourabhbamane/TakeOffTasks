package controllers

type Employee struct {
	ID        int     `json:"id"`
	FirstName string  `json:"firstname"`
	LastName  string  `json:"lastname"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	PhoneNo   string  `json:"phoneno"`
	Role      string  `json:"role"`
	Salary    float64 `json:"salary"`
}

//creating memory datastructure as database

var employees []Employee

// func AddingEmp() {
// 	employees = append(employees, Employee{ID: 1, FirstName: "Sourabh", LastName: "Bamane", Email: "sourabh@gmail.com", Password: "srb123", PhoneNo: "7757976500", Role: "ADMIN", Salary: 30000})
// 	employees = append(employees, Employee{ID: 2, FirstName: "Prasanna", LastName: "Shekhar", Email: "prasanna@gmail.com", Password: "psn123", PhoneNo: "8855994466", Role: "DEVELOPER", Salary: 50000})
// 	employees = append(employees, Employee{ID: 3, FirstName: "Meghna", LastName: "Shinde", Email: "meghna@gmail.com", Password: "mgn123", PhoneNo: "8899102530", Role: "MANAGER", Salary: 30000})
// 	employees = append(employees, Employee{ID: 4, FirstName: "Mahesh", LastName: "Patil", Email: "mahesh@gmail.com", Password: "mhs123", PhoneNo: "950659988", Role: "TESTER", Salary: 50000})

// }
