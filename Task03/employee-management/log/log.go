package log

import (
	"fmt"
	"log"
	"os"
)

func Write(msg string) {
	//will open the file | 0644 wil give permission to read & write the file
	file, err := os.OpenFile("logFile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	//closing the file
	defer file.Close()
	log.SetOutput(file)
	log.Println(msg)
}
