package lichess

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jonpchin/gochess/gostuff"
)

// Gets player account info
func GetAccount() {

	client := gostuff.TimeOutHttp(5)
	url := "https://lichess.org/api/account"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+getLichessToken())

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Could not get lichess account info", err)
		return
	}

	defer response.Body.Close()
	htmlData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(htmlData))
}

func getLichessToken() string {
	readFile, err := os.Open("secret/lichess.txt")
	defer readFile.Close()
	if err != nil {
		fmt.Println("getLichessToken ", err)
	}

	scanner := bufio.NewScanner(readFile)
	scanner.Scan()
	return scanner.Text()
}
