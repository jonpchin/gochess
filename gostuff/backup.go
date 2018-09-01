package gostuff

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

//exports database(without Grandmaster games) to an .sql file as a hot backup
//@param isTemplate If true then export template database
func ExportDatabase(isTemplate bool) {

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	command := "mysqldump --databases gochess --ignore-table=gochess.grandmaster > ./../backup/gochessNoGrandmaster.sql"

	if isTemplate {
		command = "cd config && mysqldump --databases gochess --no-data > ./../backup/gochessTemplate.sql"
	}

	if runtime.GOOS == "windows" {
		_, err := exec.Command("cmd.exe", "/C", command).Output()
		if err != nil {
			log.Println(err)
			fmt.Println("Error in exporting database, please check logs")
		}
	} else {
		_, err := exec.Command("/bin/bash", "-c", command).Output()
		if err != nil {
			log.Println(err)
			fmt.Println("Error in exporting database, please check logs")
		}
	}
}

// zips up exported database
func CompressDatabase() {
	result := compress("./backup/gochess.zip", []string{"./backup/gochess.sql"})
	if result {
		fmt.Println("Exported database file succesfully compressed!")
	}
}

//imports the main gochess database, returns true if successful
func importDatabase() bool {
	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	result := unzip("./backup/gochess.zip", "./backup")

	if result == false {
		return false
	}
	if runtime.GOOS == "windows" {
		_, err := exec.Command("cmd.exe", "/C", "mysql < ./backup/gochess.sql").Output()
		if err != nil {
			log.Println(err)
			fmt.Println("Error in importing gochess database, please check logs")
			return false
		}
	} else {
		_, err := exec.Command("/bin/bash", "-c", "mysql < ./backup/gochess.sql").Output()
		if err != nil {
			log.Println(err)
			fmt.Println("Error in importing gochess database, please check logs")
			return false
		}
	}
	return true
}

//imports template database, returns true if sucessful
func importTemplateDatabase() bool {
	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	//determine which operating system to execute appropriate shell command
	if runtime.GOOS == "windows" {
		_, err := exec.Command("cmd.exe", "/C", "mysql < ./backup/gochessTemplate.sql").Output()
		if err != nil {
			log.Println(err)
			fmt.Println("Error in importing template database, please check logs")
			return false
		}
	} else {
		_, err := exec.Command("/bin/bash", "-c", "mysql < ./backup/gochessTemplate.sql").Output()
		if err != nil {
			log.Println(err)
			fmt.Println("Error in importing template database, please check logs")
			return false
		}
	}
	return true
}
