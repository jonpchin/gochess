package gostuff

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
)

// For brand new installations this will setup credentials in the secret directory
func SetupSecretDir() {
	secretDir := "secret"
	if IsDirectory(secretDir) == false {
		CreateDirIfNotExist(secretDir)
	}

	setupSecretConfig()
}

// Prompts for config.txt if its not in secret folder
func setupSecretConfig() {

	// secretUserpath contains non root credentials of MySQL
	secretRootPath := "secret/root.txt"
	secretUserPath := "secret/config.txt"

	var databaseInfo DatabaseInfo
	databaseInfo.User = "root"

	// The following defaults below can be changed to connect to a remote database
	databaseInfo.Host = "localhost"
	databaseInfo.Port = "3306"
	databaseInfo.DbName = "gochess"
	scanner := bufio.NewScanner(os.Stdin)

	// Root MySQL
	if isFileExist(secretRootPath) == false {
		fmt.Println("Please enter your MySQL password for your root account:")

		for scanner.Scan() {
			databaseInfo.Password = scanner.Text()
			err := databaseInfo.writeSecretConfig(secretRootPath)
			if err != nil {
				fmt.Println("Could not create secret root file ", secretRootPath, err)
			} else {
				break
			}
		}
	}

	// Non root MySQL
	if isFileExist(secretUserPath) == false {
		// Setting up MySQL database should have a non root account for security reasons
		fmt.Println("Please enter your MySQL username for your non root account:")

		for scanner.Scan() {
			databaseInfo.User = scanner.Text()
			fmt.Println("Please enter your MySQL password for your non root account:")
			scanner.Scan()
			databaseInfo.Password = scanner.Text()
			err := databaseInfo.writeSecretConfig(secretUserPath)
			if err != nil {
				fmt.Println("Could not create secret user file ", secretRootPath, err)
			} else {
				break
			}
		}
	}
}

func (databaseInfo DatabaseInfo) writeSecretConfig(filePath string) error {

	info := databaseInfo.User + "\n" + getBase64Text(databaseInfo.Password) + "\n" +
		databaseInfo.Host + "\n" + databaseInfo.Port + "\n" + databaseInfo.DbName

	err := ioutil.WriteFile(filePath, []byte(info), 0666)
	if err != nil {
		return err
	}
	return nil
}

func getBase64Text(text string) string {

	encoded64 := base64.StdEncoding.EncodeToString([]byte(text))
	hexencode := hex.EncodeToString([]byte(encoded64))
	return hexencode
}
