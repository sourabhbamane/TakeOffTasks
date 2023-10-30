package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// helper Func to check any error is there or not

/*
func CheckError(err error) (string, error) {
	if err != nil {
		//panic(err)
		return "", fmt.Errorf("something went wrong : %s", err)
	}
	return "", nil
}
*/

func CheckError(err error) {
	if err != nil {
		//e := errors.New(err.Error())
		fmt.Printf("Something Went Wrong: %s\n", err.Error())
	}
}

// helper Func to select the role to add in the slice
func ChooseRole() string {
	start := true
	for start {
		fmt.Printf("SELECT ROLE:\n1.Admin\n2.Manager\n3.Developer\n4.Tester\n")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		CheckError(err)

		choice, err := strconv.Atoi(strings.TrimSpace(input))
		CheckError(err)
		fmt.Println("choice :", choice)

		switch choice {
		case 1:
			return "ADMIN"
		case 2:
			return "MANAGER"
		case 3:
			return "DEVELOPER"
		case 4:
			return "TESTER"
		default:
			fmt.Println("Role Not Available")
		}
	}
	return ""
}

// helper method to check if email is already exists in the slice
func isEmailExists(email string) bool {

	for _, emp := range employees {
		if emp.Email == email {
			return true
		}
	}
	return false
}

//helper method to check if the phone No is Present or Not

func isPhoneAvail(phn string) bool {

	for _, emp := range employees {
		if emp.PhoneNo == phn {
			return true
		}
	}
	return false
}
