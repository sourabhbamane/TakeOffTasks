package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoginUser() {
	for {
		fmt.Println("------------Welcome-----------")
		fmt.Printf("0.Exit the Application\n1.login\n")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		CheckError(err)
		choice, err := strconv.Atoi(strings.TrimSpace(input)) //.Atoi is equivalent to parseint
		CheckError(err)
		switch choice {
		case 0:
			fmt.Println("Thank You.... Visit Again")
			os.Exit(0)
		case 1:
			emp, isOk := isPresent()
			if !isOk {
				fmt.Println("Please Provide Proper Credentials")
			} else {
				//for authorization
				if emp.Role == "ADMIN" {
					AdminFunctionality()
				} else {
					NonAdminUser(emp)
				}
			}
		default:
			fmt.Println("Enter Valid Choice-----")
		}

	}
}

// method to authenticate user
func isPresent() (Employee, bool) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please Enter Your Login Credentials")
	fmt.Printf("UserName/Email : ")
	input, err := reader.ReadString('\n')
	email := strings.TrimSpace(input)
	CheckError(err)

	fmt.Printf("Password :")
	input1, err := reader.ReadString('\n')
	pass := strings.TrimSpace(input1)
	CheckError(err)

	//declaring var referencing to struct
	var foundEmployee *Employee
	//iterating through for loop to get index
	//and check whether employee is available or not
	for index, emp := range employees {
		if emp.Email == email && emp.Password == pass {
			foundEmployee = &employees[index]
			return *foundEmployee, true
		}
	}
	return Employee{}, false
}
