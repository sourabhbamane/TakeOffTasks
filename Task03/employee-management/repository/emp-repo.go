package repository

import "employee.com/myapp/entity"

type EmpRepository interface {
	//To save Employee Details
	Save(emp *entity.Employee) (*entity.Employee, error)

	//To Get All Eployees
	FindAll() ([]entity.Employee, error)

	//To Delete Employee
	DeleteEmployee(empID int64) (string, error)

	//To Update The Employee Details From Id
	UpdateEmp(empId int64, emp *entity.Employee) error

	//To Search Employees
	FindEmployees(filters map[string]string) ([]entity.Employee, error)

	//To Find Employee By Its ID
	FindById(id int64) (*entity.Employee, error)
}
