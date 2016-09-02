package gostuff

import(
	"os"
	"log"
	"os/exec"
	"fmt"
)

//exports database to an .sql file as a hot backup
func ExportDatabase(){
	
	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	_, err := exec.Command("cmd.exe", "/C", "cd.. && bash backup.sh").Output()
	if err != nil{
		log.Println("backup.go ExportDatabase 1", err)
		fmt.Println("Error in exporting database, please check logs")
	}
	
	result := compress("./../backup/gochess.zip", []string{"./../backup/gochess.sql"})
	if result == true{
		fmt.Println("Database file succesfully compressed!")
	}
}

//imports the main gochess database, returns true if sucessful
func importDatabase() bool{
	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)
	
	result := unzip("./../backup/gochess.zip", "./../backup")
	
	if result == false{
		return false
	}
	
	_, err := exec.Command("cmd.exe", "/C", "cd.. && bash importGoChess.sh").Output()
	if err != nil{
		log.Println("backup.go importDatabase 1", err)
		fmt.Println("Error in importing gochess database, please check logs")
		return false
	}
	return true
}

//imports template database, returns true if sucessful
func importTemplateDatabase() bool{
	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)
	
	_, err := exec.Command("cmd.exe", "/C", "cd.. && bash importTemplate.sh").Output()
	if err != nil{
		log.Println("backup.go importTemplateDatabase 1", err)
		fmt.Println("Error in importing template database, please check logs")
		return false
	}
	
	return true
}