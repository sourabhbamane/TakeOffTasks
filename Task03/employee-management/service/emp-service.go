package service

import (
	"errors"
	"math/rand"

	"employee.com/myapp/entity"
	"employee.com/myapp/repository"
)

type EmpService interface {
	Validate(emp *entity.Employee) error
	Create(emp *entity.Employee) (*entity.Employee, error)
	GetAll() ([]entity.Employee, error)
	DeleteEmp(id int64) (string, error)
	FindEmpByID(id int64) (*entity.Employee, error)
	Update(id int64, emp *entity.Employee) error
	SearchEmp(filters map[string]string) ([]entity.Employee, error)
}

type service struct{}

var (
	repo repository.EmpRepository
)

func NewEmpService(r repository.EmpRepository) EmpService {
	repo = r
	return &service{}
}

// TO Check If the certain Field Are empyty or not
func (*service) Validate(emp *entity.Employee) error {
	if emp == nil {
		err := errors.New("the employee is empty")
		return err
	}
	if emp.FirstName == "" {
		err := errors.New("the employee name is empty")
		return err
	}
	return nil
}

// To Save the details
func (*service) Create(emp *entity.Employee) (*entity.Employee, error) {
	//creating random id
	emp.ID = rand.Int63()
	return repo.Save(emp)
}

// TO Get All Employees
func (*service) GetAll() ([]entity.Employee, error) {
	return repo.FindAll()
}

func (*service) DeleteEmp(id int64) (string, error) {
	return repo.DeleteEmployee(id)
}

func (*service) Update(id int64, emp *entity.Employee) error {
	return repo.UpdateEmp(id, emp)
}

func (*service) SearchEmp(f map[string]string) ([]entity.Employee, error) {
	return repo.FindEmployees(f)
}

func (*service) FindEmpByID(id int64) (*entity.Employee, error) {
	return repo.FindById(id)
}
