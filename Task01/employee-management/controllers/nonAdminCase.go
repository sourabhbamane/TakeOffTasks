package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//functionalities od non admin users

func NonAdminUser(emp Employee) {
	//ViewMyDetails(2)
	id := emp.ID
	status := true
	for status {
		fmt.Printf("Choose Option :\n0.LogOut\n1.View Your Profile\n2.Update Your Profile\n3.Search Employee\n4.Search Employee From Id\n")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		CheckError(err)
		choice, err := strconv.Atoi(strings.TrimSpace(input)) //.Atoi is equivalent to parseint
		CheckError(err)
		fmt.Println("choice :", choice)
		switch choice {
		case 0:
			msg := emp.FirstName + " As " + emp.Role + " Logged out"
			LoggingFunc(msg)
			fmt.Println("Logged Out Succesfully")
			// os.Exit(0)
			status = false
		case 1:
			msg := emp.FirstName + " As " + emp.Role + " Viewed His/Her Details"
			LoggingFunc(msg)
			ViewMyDetails(id)
		case 2:
			UpdateMyProfile(id)
			msg := emp.FirstName + " As " + emp.Role + " Updated His/Her Details"
			LoggingFunc(msg)
		case 3:
			fmt.Printf("Enter Employee Name To Search : ")
			input, err := reader.ReadString('\n')
			fName := strings.TrimSpace(input)
			CheckError(err)
			SearchEmployee(fName)
			msg := emp.FirstName + " As " + emp.Role + " Searched Employee of name " + fName
			LoggingFunc(msg)
		case 4:
			fmt.Printf("Enter Your ID TO search: ")
			input, err := reader.ReadString('\n')
			CheckError(err)
			id, err := strconv.Atoi(strings.TrimSpace(input)) //.Atoi is equivalent to parseint
			CheckError(err)
			GetEmployeeByID(id)
		default:
			fmt.Println("Enter Correct Choice......")

		}
	}
}
