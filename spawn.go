package gostuff

import(
	"os/exec"
	"fmt"
	"runtime"
//	"os"
)

//spawns a donna background process for chess engine
func SpawnProcess(){
	
	if runtime.GOOS == "windows"{
		cmd := exec.Command("sh", "-c", "cd ../../michaeldv/donna/bin && ls")

		output, _ := cmd.Output()
	    fmt.Println(string(output))

	}else if runtime.GOOS == "linux"{ //then its a linux machine
		cmd, err := exec.Command("cd ../../michaeldv/donna/bin & ls").Output()
	    if err != nil {
	        panic(err)
	    }
	    fmt.Println("The output is ")
	    fmt.Println(string(cmd))
	}else{
		fmt.Println("Unknown operating system spawn.go SpawnProces 1")
	}
	
	
}