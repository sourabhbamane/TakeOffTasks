package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// function to Search Employee but showing only limited fields
func SearchEmployee(empName string) {

	// Search for the employee by name
	foundEmployees := []Employee{}
	for _, emp := range employees {
		if strings.Contains(strings.ToLower(emp.FirstName), strings.ToLower(empName)) {
			foundEmployees = append(foundEmployees, emp)
		}
	}

	// Display the search results
	if len(foundEmployees) > 0 {
		fmt.Println("Found employees with matching names:")
		for _, emp := range foundEmployees {
			fmt.Printf("Name: %s %s\nEmail: %s\nRole: %s\n--------------------\n", emp.FirstName, emp.LastName, emp.Email, emp.Role)
		}
	} else {
		fmt.Println("No employees found with the specified name.")
	}
}

// view His Details
func ViewMyDetails(empId int) {
	for _, emp := range employees {
		if emp.ID == empId {
			fmt.Printf("ID: %d\nName: %s %s\nEmail: %s\nRole: %s\nSalary: %.2f\nDOB: %d-%02d-%02d\n--------------------\n", emp.ID, emp.FirstName, emp.LastName, emp.Email, emp.Role, emp.Salary, emp.DOB.Year, emp.DOB.Month, emp.DOB.Day)
			return
		}
	}
}

// Update his Details
func UpdateMyProfile(empId int) {

	// TO Find the index of employee in slice
	var foundEmployee *Employee
	for index, emp := range employees {
		if emp.ID == empId {
			foundEmployee = &employees[index]
			break
		}
	}

	fmt.Println("Present Info of Employee: ")
	GetEmployeeByID(empId)

	fmt.Println("From Above Which Fields Do Want TO Update")

	status := true
	for status {
		fmt.Printf("0.To Save\n1.First Name\n2.Last Name\n3.Email\n4.Password\n5.Role\n6.Salary\n7.DOB\n8.Phone Number\n------------------\n")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		CheckError(err)
		choice, err := strconv.Atoi(strings.TrimSpace(input)) //.Atoi is equivalent to parseint
		CheckError(err)
		fmt.Println("choice :", choice)
		switch choice {
		case 0:
			fmt.Println("User Updated Succesfully")
			fmt.Println("To View Updated Profile Press 1")
			status = false
		case 1:
			status := true
			for status {
				fmt.Printf("Change First Name: ")
				fName, err := reader.ReadString('\n')
				CheckError(err)
				fName = strings.TrimSpace(fName)
				if fName == "" {
					fmt.Println("Cannot be Empty")
				} else {
					foundEmployee.FirstName = fName
					status = false
				}

			}
			fmt.Println("First Name Updated Succesfully")
		case 2:
			fmt.Printf("Change Last Name: ")
			lName, err := reader.ReadString('\n')
			CheckError(err)
			foundEmployee.LastName = strings.TrimSpace(lName)
			fmt.Println("Last Name Updated Succesfully")
		case 3:
			status := true
			for status {
				fmt.Printf("Change Email: ")
				email, err := reader.ReadString('\n')
				CheckError(err)
				email = strings.TrimSpace(email)
				isPresent := isEmailExists(email)
				if email == "" {
					fmt.Println("Please Enter Email")
				} else if isPresent { //if non empty then again checking if it is already present or not
					fmt.Println("Email Id is Already Exist Please Provide Another email")
				} else {
					foundEmployee.Email = email
					status = false
				}

			}
			fmt.Println("Email Updated Succesfully")
		case 4:
			fmt.Printf("Change password: ")
			pass, err := reader.ReadString('\n')
			CheckError(err)
			foundEmployee.Password = strings.TrimSpace(pass)
			fmt.Println("Password Updated Succesfully")
		case 5:
			role := ChooseRole()
			CheckError(err)
			foundEmployee.Role = strings.TrimSpace(role)
			fmt.Println("Role Updated Succesfully")
		case 6:
			fmt.Printf("Edit Salary: ")
			input, err := reader.ReadString('\n')
			CheckError(err)

			//convering input string to float
			salary, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
			CheckError(err)
			foundEmployee.Salary = salary
			fmt.Println("Salary Updated Succesfully")
		case 7:
			fmt.Printf("Change Year of Birth: \n")
			i6, err := reader.ReadString('\n')
			CheckError(err)
			year, err := strconv.Atoi(strings.TrimSpace(i6))
			CheckError(err)

			fmt.Printf("Change Month of Birth: \n")
			i7, err := reader.ReadString('\n')
			CheckError(err)
			month, err := strconv.Atoi(strings.TrimSpace(i7))
			CheckError(err)

			fmt.Printf("Change Day of Birth: \n")
			i8, err := reader.ReadString('\n')
			CheckError(err)
			day, err := strconv.Atoi(strings.TrimSpace(i8))
			CheckError(err)

			dob := Date{Year: year, Month: time.Month(month), Day: day}
			foundEmployee.DOB = dob
			fmt.Println("DOB Updated Succesfully")
		case 8:
			status := true
			for status {
				fmt.Printf("Change Phone Number: ")
				phnNo, err := reader.ReadString('\n')
				phnNo = strings.TrimSpace(phnNo)
				CheckError(err)
				if phnNo == "" {
					fmt.Println("Cannot be empty")
				} else if isPhoneAvail(phnNo) {
					fmt.Println("Number Alredy Exists")
				} else {
					//Update the Employee
					foundEmployee.PhoneNo = phnNo
					status = false
				}

			}
			fmt.Println("Phone No. Updated Succesfully")
		default:
			fmt.Println("Enter Valid Choice...")
		}
	}
}
