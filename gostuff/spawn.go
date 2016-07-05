package gostuff

import(
	"os/exec"
	"fmt"
)

//spawns a donna background process for chess engine
func SpawnProcess(){
	cmd, err := exec.Command("cmd", "/C", "cd ../../michaeldv/donna/bin & ls").Output()
    if err != nil {
        panic(err)
    }
    fmt.Println("The output is ")
    fmt.Println(string(cmd))
}