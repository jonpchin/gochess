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
	"time"

	"github.com/dchest/captcha"
	_ "github.com/go-sql-driver/mysql"
)

// global SessionManager["username"] = sessionID
var SessionManager = make(map[string]string)

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
	} else {
		w.Write([]byte(result))
	}
}

func (userinfo *UserInfo) Login(method string, url string, agent string, host string) (string, error) {
	if userinfo.CaptchaSolution == "" { //then assume user was not displayed captcha

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
		var verify string
		var captcha int
		var email string

		//getting password, verify, captcha and email
		err = db.QueryRow("SELECT password, verify, captcha, email FROM userinfo WHERE username=?", userinfo.Username).Scan(&pass, &verify, &captcha, &email)
		if err != nil {
			return "<img src='img/ajax/not-available.png' /> Wrong username/password combination.",
				errors.New("Incorrect login for " + userinfo.Username)
		}

		//if user entered password incorrect two times or more then they need to enter captcha to login
		if captcha >= 2 {
			return "<script>$('#hiddenCap').show();</script><img src='img/ajax/not-available.png' />" +
				"You entered password incorrectly too many times. Now you need to enter captcha.", nil
		}
		//checking if password entered by user matches encrypted key
		if pass != key {
			var result string
			if captcha == 1 {
				result = "<script>$('#hiddenCap').show();</script><img src='img/ajax/not-available.png' />" +
					"You entered password incorrectly too many times. Now you need to enter captcha."
			} else {
				result = "<img src='img/ajax/not-available.png' /> Wrong username/password combination."
			}
			addOneToCaptcha(userinfo.Username, captcha)
			return result, fmt.Errorf("FAILED LOGIN IP: %s  Method: %s Location: %s Agent: %s\n",
				userinfo.IpAddress, method, url, agent)
		}
		if verify != "YES" {
			result, err := needToActivate(host, userinfo.Username)
			return result, err
		}
		// update captcha to zero since login was a success
		stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
		defer stmt.Close()
		if err != nil {
			return "<img src='img/ajax/not-available.png' /> Error in captcha section 3", err
		}

		_, err = stmt.Exec(0, userinfo.Username)
		if err != nil {
			return "<img src='img/ajax/not-available.png' /> Error in captcha section 4", err
		}
		return "", nil

	} else if !captcha.VerifyString(userinfo.CaptchaId, userinfo.CaptchaSolution) {
		return "<script>document.getElementById('captchaSolution').value = '';</script><img src='img/ajax/not-available.png' /> Wrong captcha solution! Please try again.", nil
	} else {
		//check if database connection is open
		if db.Ping() != nil {

			log.Println()
			return "<img src='img/ajax/not-available.png' /> We are having trouble with our server.",
				errors.New("DATABASE DOWN!")
		}

		//hashing password
		key, err := hashPass(userinfo.Username, userinfo.Password)
		if err != nil {
			return "<img src='img/ajax/not-available.png' /> Something is wrong with the server.", err
		}

		var pass string
		var verify string
		var captcha int
		var email string
		var token string

		//getting password, verify, captcha, email from database
		err = db.QueryRow("SELECT password, verify, captcha, email FROM userinfo WHERE username=?", userinfo.Username).Scan(&pass, &verify, &captcha, &email)
		if err != nil {
			//check if there was more then one incorrect login attempt
			return "<img src='img/ajax/not-available.png' /> Wrong username/password combination.", err
		}

		if captcha == 5 {

			addOneToCaptcha(userinfo.Username, captcha)

			//create activation token in database and send user notifying them that their was five incorrect login attempts
			stmt, err := db.Prepare("INSERT activate SET username=?, token=?, email=?, expire=?")
			if err != nil {
				return "<img src='img/ajax/not-available.png' /> Can't prepare activation token", err
			}

			token = RandomString()
			_, err = stmt.Exec(userinfo.Username, token, email, time.Now())
			if err != nil {
				return "<img src='img/ajax/not-available.png' /> Can't execute activation token", err
			}
			//sends email to user with the token activation
			go SendAttempt(email, token, userinfo.Username, userinfo.IpAddress, host)

			return "<img src='img/ajax/not-available.png' /> This account has been deactivated becaue too many incorrect login attempts were made. An email has been sent to you on instructions on reactivating your account.",
				nil

		} else if captcha > 5 { //tell user on the front end that this account has too many login attempts, resends activation token

			err2 := db.QueryRow("SELECT token FROM activate WHERE username=?", userinfo.Username).Scan(&token)
			if err2 != nil {
				return "<img src='img/ajax/not-available.png' /> Can't query token from user", err2
			}

			//sends email again to user with the token activation
			go SendAttempt(email, token, userinfo.Username, userinfo.IpAddress, host)

			return "<img src='img/ajax/not-available.png' /> This account has been deactivated because too many incorrect login attempts were made. An email has been sent again regarding how to reactivate your account.", nil
		}
		if pass != key {
			addOneToCaptcha(userinfo.Username, captcha)
			return "<img src='img/ajax/not-available.png' /> Wrong username/password combination.", nil
		}
		if verify != "YES" {
			result, err := needToActivate(host, userinfo.Username)
			return result, err
		}

		// update captcha to zero since login was a success
		stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
		defer stmt.Close()

		if err != nil {
			return "<img src='img/ajax/not-available.png' /> Can't prepare captcha to zero", err
		}

		_, err = stmt.Exec(0, userinfo.Username)
		if err != nil {
			return "<img src='img/ajax/not-available.png' /> Can't execute captcha to zero", err
		}

		return "", nil
	}
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

	var index int = 0
	username := "guest0"
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
	sessionID := RandomString()
	cookie = http.Cookie{Name: "sessionID", Value: sessionID, Secure: true, HttpOnly: true, Expires: expiration}
	http.SetCookie(w, &cookie)

	SessionManager[username] = sessionID

	w.Write([]byte("<script>window.location = '/server/lobby'</script>"))
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

// sends an email again to reactivate an inactivated account
func needToActivate(host string, username string) (string, error) {
	var tokenInDB string
	var email string

	log.Printf("%s needs to activate his account before logging in.\n", username)

	//checking if token matches the one entered by user
	err2 := db.QueryRow("SELECT token, email FROM activate WHERE username=?", username).Scan(&tokenInDB, &email)
	if err2 != nil {
		log.Println(err2)
	} else {
		go Sendmail(email, tokenInDB, username, host)
	}
	return "<img src='img/ajax/not-available.png' /> You must activate your account by entering the activation token in your email at the activation page. " +
		"An email has been sent again containing your activation code.", nil
}

//add 1 to captcha if password was incorrect
func addOneToCaptcha(username string, captcha int) {

	stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return
	}
	captcha = captcha + 1

	_, err = stmt.Exec(captcha, username)
	if err != nil {
		log.Println(err)
	}
}
