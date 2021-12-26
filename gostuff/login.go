package gostuff

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// global SessionManager["username"] = sessionID
var SessionManager = make(map[string]string)
var DeviceManager = make(map[string]string)

// Keeps track of users that are locked out until the next tick cycle
var protectAccounts = struct {
	sync.RWMutex
	Users []UserInfo
}{}

type UserInfo struct {
	Username        string
	Password        string
	IpAddress       string
	CaptchaId       string
	CaptchaSolution string
	InvalidLogins   int // Number of times an account had incorrect logins, resets on lockup
	IsLocked        bool
}

func isUserInProtectAccounts(username string) bool {

	for _, value := range protectAccounts.Users {

		if value.Username == username {
			return true
		}
	}

	return false
}

func isAccountLocked(username string) bool {

	for _, value := range protectAccounts.Users {

		if value.Username == username && value.IsLocked {
			return true
		}
	}

	return false
}

//process user input when signing in
func ProcessLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Show404Page(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var userinfo UserInfo
	userinfo.Username = template.HTMLEscapeString(r.FormValue("user"))
	userinfo.Password = template.HTMLEscapeString(r.FormValue("password"))

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	userinfo.CaptchaId = template.HTMLEscapeString(r.FormValue("captchaId"))
	userinfo.CaptchaSolution = template.HTMLEscapeString(r.FormValue("captchaSolution"))
	userinfo.IpAddress, _, _ = net.SplitHostPort(r.RemoteAddr)

	result, err := userinfo.Login(r.Method, r.URL.Path, r.UserAgent(), r.Host)
	if err != nil {
		log.Println(err)
	}

	if result == "" {
		enterInside(w, userinfo.Username, userinfo.IpAddress)
		DeviceManager[userinfo.Username] = r.UserAgent()
	} else {
		w.Write([]byte(result))
	}
}

func (userinfo *UserInfo) Login(method string, url string, agent string, host string) (string, error) {

	if isAccountLocked(userinfo.Username) {
		return "<img src='img/ajax/not-available.png' /> Too many login attempts. Account is locked. Pleae try logging in later.", nil
	}

	//check if database connection is open
	if db.Ping() != nil {
		return "<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later.",
			errors.New("DATABASE DOWN!")
	}

	//hashing password
	key, err := hashPass(userinfo.Username, userinfo.Password)
	if err != nil {
		return "<img src='img/ajax/not-available.png' /> Something is wrong with the server.", err
	}

	var pass string

	//getting password, verify, captcha and email
	err = db.QueryRow("SELECT password FROM userinfo WHERE username=?", userinfo.Username).Scan(&pass)
	if err != nil {
		return "<img src='img/ajax/not-available.png' /> Wrong username/password combination.",
			errors.New("Incorrect login for " + userinfo.Username)
	}

	if pass != key {
		if isUserInProtectAccounts(userinfo.Username) == false {

			var user UserInfo
			user.InvalidLogins = 0
			user.Username = userinfo.Username
			user.IsLocked = false

			protectAccounts.Lock()
			protectAccounts.Users = append(protectAccounts.Users, user)
			protectAccounts.Unlock()

		} else {
			protectAccounts.Lock()
			for index, _ := range protectAccounts.Users {
				if protectAccounts.Users[index].Username == userinfo.Username {
					if protectAccounts.Users[index].InvalidLogins >= 3 {
						protectAccounts.Users[index].InvalidLogins = 0
						protectAccounts.Users[index].IsLocked = true
					} else {
						protectAccounts.Users[index].InvalidLogins += 1
					}
				}
			}
			protectAccounts.Unlock()
		}
		return "<img src='img/ajax/not-available.png' /> Wrong username/password combination.",
			fmt.Errorf("FAILED LOGIN IP: %s  Method: %s Location: %s Agent: %s\n",
				userinfo.IpAddress, method, url, agent)
	}

	return "", nil
}

func RandomString() string {
	size := 12 // change the length of the generated random string here 12 is actually 16

	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println("login.go RandomString 1 ", err)
	}

	return base64.URLEncoding.EncodeToString(rb)
}

func EnterGuest(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		Show404Page(w, r)
		return
	}

	username := template.HTMLEscapeString(r.FormValue("user"))
	sessionID := template.HTMLEscapeString(r.FormValue("pass"))

	if sessionID == "" || SessionManager[username] != sessionID {

		var index int = 0
		username = "guest0"
		for {
			if _, ok := SessionManager[username]; ok {
				index++
				username = "guest" + strconv.Itoa(index)
			} else {
				break
			}
		}

		expiration := time.Now().Add(5 * 24 * time.Hour)
		cookie := http.Cookie{Name: "username", Value: username, Secure: true, HttpOnly: true, Expires: expiration}
		http.SetCookie(w, &cookie)

		//generating random session ID to be stored in the backend
		sessionID = RandomString()
		cookie = http.Cookie{Name: "sessionID", Value: sessionID, Secure: true, HttpOnly: true, Expires: expiration}
		http.SetCookie(w, &cookie)

		userAgent := r.UserAgent()

		if userAgent == "okhttp/3.11.0" {
			DeviceManager[username] = "android"
		} else {
			DeviceManager[username] = userAgent
		}

		SessionManager[username] = sessionID
	}

	//w.Write([]byte("<script>window.location = '/server/lobby'</script>"))
	w.Write([]byte(username + "," + sessionID))
}

// after successfully identifying credentials, setup session and cookies
func enterInside(w http.ResponseWriter, username string, ipAddress string) {
	expiration := time.Now().Add(3 * 24 * time.Hour)
	cookie := http.Cookie{Name: "username", Value: username, Secure: true, HttpOnly: true, Expires: expiration}
	http.SetCookie(w, &cookie)

	//generating random session ID to be stored in the backend
	sessionID := RandomString()
	cookie = http.Cookie{Name: "sessionID", Value: sessionID, Secure: true, HttpOnly: true, Expires: expiration}
	http.SetCookie(w, &cookie)
	country := GetCountry(username)

	// country cookie can be modified by JS
	cookie = http.Cookie{Name: "country", Value: country, Secure: false, HttpOnly: false, Expires: expiration}
	http.SetCookie(w, &cookie)
	SessionManager[username] = sessionID

	w.Write([]byte("<script>window.location = '/server/lobby'</script>"))
}
