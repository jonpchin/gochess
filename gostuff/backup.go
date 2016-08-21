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
}