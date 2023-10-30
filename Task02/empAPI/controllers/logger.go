package controllers

import (
	"log"
	"os"
)

func LoggingFunc(msg string) {
	//will open the file | 0644 wil give permission to read & write the file
	file, err := os.OpenFile("logFile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	CheckError(err)
	//closing the file
	defer file.Close()
	log.SetOutput(file)
	log.Println(msg)
}
