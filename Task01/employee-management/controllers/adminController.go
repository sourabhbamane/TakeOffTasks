package controllers

import (
	"bufio"
	"time"

	"fmt"
	"os"
	"strconv"
	"strings"
)

// Add new Employee
func AddEmployee() {

	//Generating Unique Id
	//TO Find Maximum id in the slice
	max := 0
	for _, emp := range employees {
		if emp.ID > max {
			max = emp.ID
		}
	}
	//generating new ID
	id := max + 1

	reader := bufio.NewReader(os.Stdin)
	isName := true
	var fName string
	for isName {
		fmt.Printf("Enter First Name: ")
		i1, err := reader.ReadString('\n')
		CheckError(err)
		fName = strings.TrimSpace(i1)
		//Checking if the name provided is empty or not
		if fName == "" {
			fmt.Println("Plese Provide First Name")
		} else {
			isName = false
		}

	}

	isLname := true
	var lName string
	for isLname {
		fmt.Printf("Enter Last Name: ")
		i2, err := reader.ReadString('\n')
		CheckError(err)
		lName = strings.TrimSpace(i2)
		//Checking if the last name provided is empty or not
		if lName == "" {
			fmt.Println("Please Provide Last Name")
		} else {
			isLname = false
		}

	}

	var email string
	isEmail := true
	for isEmail {
		fmt.Printf("Enter Email: ")
		i3, err := reader.ReadString('\n')
		CheckError(err)
		email = strings.TrimSpace(i3)
		isPresent := isEmailExists(email)
		//Checking if the email provided is empty or not
		if email == "" {
			fmt.Println("Please Enter Email")
		} else if isPresent { //if non empty then again checking if it is already present or not
			fmt.Println("Email Id is Already Exist Please Provide Another email")
		} else {
			isEmail = false
		}

	}

	fmt.Printf("Enter password: ")
	i4, err := reader.ReadString('\n')
	pass := strings.TrimSpace(i4)
	CheckError(err)

	var phnNo string
	status2 := true
	//TO check if number is already in memory or not
	for status2 {
		fmt.Printf("Enter Phone Number: ")
		i5, err := reader.ReadString('\n')
		CheckError(err)
		phnNo = strings.TrimSpace(i5)
		isPresent := isPhoneAvail(phnNo)
		//Checking if the phone num provided is empty or not
		if phnNo == "" {
			fmt.Println("Please Provide Phone Number")
		} else if isPresent { //if non empty then again checking if it is already present or not
			fmt.Println("Phone no is Already Exist Please Provide Another Phone")
		} else {
			status2 = false
		}

	}

	role := ChooseRole()
	CheckError(err)

	fmt.Printf("Enter Salary: ")
	input, err := reader.ReadString('\n')
	CheckError(err)

	//convering input string to float
	salary, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
	CheckError(err)

	fmt.Printf("Enter DOB : \n")
	fmt.Printf("Enter Year of Birth: \n")
	i6, err := reader.ReadString('\n')
	CheckError(err)
	year, err := strconv.Atoi(strings.TrimSpace(i6))
	CheckError(err)

	fmt.Printf("Enter Month of Birth: \n")
	i7, err := reader.ReadString('\n')
	CheckError(err)
	month, err := strconv.Atoi(strings.TrimSpace(i7))
	CheckError(err)

	fmt.Printf("Enter day of Birth: \n")
	i8, err := reader.ReadString('\n')
	CheckError(err)
	day, err := strconv.Atoi(strings.TrimSpace(i8))
	CheckError(err)

	dob := Date{Year: year, Month: time.Month(month), Day: day}

	newEmp := Employee{ID: id, FirstName: fName, LastName: lName, Email: email, Password: pass, PhoneNo: phnNo, Role: role, Salary: salary, DOB: dob}

	employees = append(employees, newEmp)

	fmt.Println("Employee Added Successfully......................")

	//for logger file
	msg := "Added User Of Name " + fName + " " + lName
	LoggingFunc(msg)
}

// to view employee details by id
func GetEmployeeByID(empId int) {
	for _, emp := range employees {
		if emp.ID == empId {
			fmt.Printf("ID: %d\nName: %s %s\nEmail: %s\nRole: %s\nSalary: %.2f\nDOB: %d-%02d-%02d\n--------------------\n", emp.ID, emp.FirstName, emp.LastName, emp.Email, emp.Role, emp.Salary, emp.DOB.Year, emp.DOB.Month, emp.DOB.Day)
			msg := "Requested to get Employee of name " + emp.FirstName + " " + emp.LastName
			LoggingFunc(msg)
			return
		}
	}
	fmt.Println("Sorry Employee is not available with given id")

}

// fun To search Employee from his name
func SearchByName(empName string) {
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
			fmt.Printf("ID: %d\nName: %s %s\nEmail: %s\nRole: %s\nSalary: %.2f\nDOB: %d-%02d-%02d\n--------------------\n", emp.ID, emp.FirstName, emp.LastName, emp.Email, emp.Role, emp.Salary, emp.DOB.Year, emp.DOB.Month, emp.DOB.Day)
			msg := "Requested to get Employee of name " + emp.FirstName
			LoggingFunc(msg)
			return
		}
	} else {
		fmt.Println("No employees found of name " + empName)
	}
}

// TO dislpay All Employee Data
func Display() {

	fmt.Println("Present Employees: ")
	for _, emps := range employees {
		//fmt.Println(emps)

		fmt.Printf("ID: %d\nName: %s %s\nEmail: %s\nRole: %s\nSalary: %.2f\nDOB: %d-%02d-%02d\n--------------------\n", emps.ID, emps.FirstName, emps.LastName, emps.Email, emps.Role, emps.Salary, emps.DOB.Year, emps.DOB.Month, emps.DOB.Day)
	}

}

// TO delete the employee
func DeleteEmployee(empId int) {

	for index, emp := range employees {
		if emp.ID == empId {
			employees = append(employees[:index], employees[index+1:]...)
			fmt.Printf("Employee with id %d is deleted successfully\n", empId)
			msg := "Deleted Employee of name " + emp.FirstName + " " + emp.LastName
			LoggingFunc(msg)
			return
		}
	}
	fmt.Printf("Employee with id %d is not Available", empId)

}

// update employee by providing his id as input to search that emp and update
func UpdateEmployeeDetails(empId int) {

	//Check if Emp is Present or not
	//tofind the index of this employee
	var foundEmployee *Employee
	for index, emp := range employees {
		if emp.ID == empId {
			foundEmployee = &employees[index]
			break
		}
	}

	if foundEmployee == nil {
		fmt.Printf("Employee With Id %v is not Available", empId)
		return
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
	//To maintain log
	msg := "Updeted Employee details of " + foundEmployee.FirstName + " " + foundEmployee.LastName
	LoggingFunc(msg)

}

// function TO Print employees who have upcoming birthday in this month
func MonthBday() {

	//Get Current Month
	currentMonth := time.Now().Month()

	//iterating through employees and checking if there BirthDay is in this month

	for _, emp := range employees {
		//dob, err := time.Parse("2006-01-02", emp.DOB)
		//CheckError(err)

		//if employee have bday this month then we print his/her data
		if emp.DOB.Month == currentMonth {
			fmt.Printf("ID: %d\nName: %s %s\nEmail: %s\nRole: %s\nSalary: %.2f\nDOB: %d-%02d-%02d\n--------------------\n", emp.ID, emp.FirstName, emp.LastName, emp.Email, emp.Role, emp.Salary, emp.DOB.Year, emp.DOB.Month, emp.DOB.Day)

		}
	}

	msg := "Requested to get Employees Who are having BDay This Month"
	LoggingFunc(msg)

}

// func to sort the employees on the basis of their first Name
func sortEmployees(emps []Employee) []Employee {
	//if no employee awalable or only one available then return it
	if len(emps) <= 1 {
		return emps
	}

	index := len(emps) - 1
	pivot := emps[index]

	//Partitioning the slice into three : less,greater,equal
	var less, equal, greater []Employee

	for _, emp := range emps {
		switch {
		case emp.FirstName < pivot.FirstName:
			less = append(less, emp)
		case emp.FirstName == pivot.FirstName:
			equal = append(equal, emp)
		case emp.FirstName > pivot.FirstName:
			greater = append(greater, emp)
		}
	}

	//Recursive call
	//Calling method again to sort the greater and less slices
	less = sortEmployees(less)
	greater = sortEmployees(greater)

	//Combining the sorted slices
	sortedEmployees := append(less, equal...)
	sortedEmployees = append(sortedEmployees, greater...)
	return sortedEmployees

}
