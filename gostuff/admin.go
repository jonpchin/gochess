package gostuff

import (
	"log"
	"os"
)

func IsMod(username string) bool {
	return getRole(username) == "mod"
}

func IsAdmin(username string) bool {
	return getRole(username) == "admin"
}

func getRole(username string) string {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	// Always default to user for safety reasons
	role := "user"

	err := db.QueryRow("SELECT role from userinfo where username=?", username).Scan(&role)
	if err != nil {
		log.Println(err)
	}
	return role
}
