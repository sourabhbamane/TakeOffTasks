package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//if the logged in user is admin then it will get the following permissions

func AdminFunctionality() {
	status := true
	for status {
		fmt.Printf("Choose Option: \n0.TO LogOut\n1.Display All Employees\n2.Add Employee\n3.View Perticular Employee\n4.Update Employee\n")
		fmt.Printf("5.Delete Employee\n6.List Employee in Sorted order\n7.Employee Who Are Having BDay in This Month\n8.Search Employee By Name\n")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		CheckError(err)
		choice, err := strconv.Atoi(strings.TrimSpace(input)) //.Atoi is equivalent to parseint
		CheckError(err)
		fmt.Println("choice :", choice)

		switch choice {
		case 0:
			fmt.Println("Logged out Succesfully")
			msg := "Admin Logged Out"
			LoggingFunc(msg)
			status = false
		case 1:
			msg := "Admin Requested To View All Employees"
			LoggingFunc(msg)
			Display()
		case 2:
			AddEmployee()
		case 3:
			fmt.Printf("Enter Employee id to search : ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			CheckError(err)
			eId, err := strconv.Atoi(strings.TrimSpace(input))
			CheckError(err)
			GetEmployeeByID(eId)
		case 4:
			fmt.Printf("Enter Employee id to Update : ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			CheckError(err)
			eId, err := strconv.Atoi(strings.TrimSpace(input))
			CheckError(err)
			UpdateEmployeeDetails(eId)
		case 5:
			fmt.Printf("Enter Employee id to delete : ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			CheckError(err)
			eId, err := strconv.Atoi(strings.TrimSpace(input))
			CheckError(err)
			DeleteEmployee(eId)
		case 6:
			//DisplayInSorted(employees)

			emps := sortEmployees(employees)
			for _, emp := range emps {
				fmt.Printf("ID: %d\nName: %s %s\nEmail: %s\nRole: %s\nSalary: %.2f\nDOB: %d-%02d-%02d\n--------------------\n", emp.ID, emp.FirstName, emp.LastName, emp.Email, emp.Role, emp.Salary, emp.DOB.Year, emp.DOB.Month, emp.DOB.Day)
			}
			msg := "Admin Viewed all employees in sorted order"
			LoggingFunc(msg)
		case 7:
			MonthBday()
		case 8:
			fmt.Printf("Enter Employee Name To Search : ")
			input, err := reader.ReadString('\n')
			fName := strings.TrimSpace(input)
			CheckError(err)
			SearchByName(fName)
		default:
			fmt.Println("Enter valid Choice :)-----")
		}

		fmt.Println("......................................................................")
	}
}
