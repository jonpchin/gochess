package gostuff

func IsMod(username string) bool {
	return getRole(username) == "mod"
}

func IsAdmin(username string) bool {
	return getRole(username) == "admin"
}

func getRole(username string) string {

	// Always default to guest for safety reasons
	role := "guest"

	db.QueryRow("SELECT role from userinfo where username=?", username).Scan(&role)
	return role
}
