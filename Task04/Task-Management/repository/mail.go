package repository

import "gopkg.in/gomail.v2"

func SendEmail(From string, To []string, task string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "managementtask557@gmail.com")
	m.SetHeader("To", To...)
	m.SetHeader("Subject", "Task Assignment")
	m.SetBody("text", "Hello,\nNew Task Assignment,\nAssigned By: "+From+"\nTask Description: "+task)
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, "managementtask557@gmail.com", "mbsj nuwf valq mwiu")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

}
