package models

import (
	"log"
	"io/ioutil"
	"fmt"
)

var connectionStr string

func Init() {
	fileBytes, err := ioutil.ReadFile("connection.txt")
	if err != nil {
		fmt.Errorf("Error reading connection file\n")
		log.Fatalln(err)
	}
	connectionStr = string(fileBytes)
}

func myLog(lines... string) {
	for _, line := range(lines) {
		log.Println(log.Ldate, line)
	}
}

func myLogError(lines... string) {
	for _, line := range(lines) {
		log.Printf("***ERROR*** %s %s", log.Ldate, line)
	}
}
