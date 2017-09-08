package gostuff

import (
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dchest/captcha"
	"golang.org/x/crypto/scrypt"
)

type UserInfo struct {
	Username        string
	Password        string
	Email           string
	IpAddress       string
	Token           string
	CaptchaId       string
	CaptchaSolution string
}

//processes the users input when signing up
func ProcessRegister(w http.ResponseWriter, r *http.Request) {

	userInfo, isValid := IsValidRegistration(w, r)
	if isValid {
		problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
		defer problems.Close()
		log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

		message, err := userInfo.Register(r.Host)
		if err == nil {

			message, err := userInfo.registerRating()
			if err == nil {
				//sends email to user
				go Sendmail(userInfo.Email, userInfo.Token, userInfo.Username, r.Host)

				message = "<script>$('#register').hide();</script><img src='img/ajax/available.png' /> Hello " +
					userInfo.Username + "! Please check email for instructions to verify your account."
				//if reached here just notify user to check his email and continue on with the account creation
				w.Write([]byte(message))
			} else {
				w.Write([]byte(message))
				log.Println(err)
			}
		} else {
			w.Write([]byte(message))
			log.Println(err)
		}
	}
}

func IsValidRegistration(w http.ResponseWriter, r *http.Request) (UserInfo, bool) {

	var userInfo UserInfo

	if r.Method != "POST" {
		w.WriteHeader(404)
		http.ServeFile(w, r, "404.html")
		return userInfo, false
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if !captcha.VerifyString(template.HTMLEscapeString(r.FormValue("captchaId")), template.HTMLEscapeString(r.FormValue("captchaSolution"))) {
		w.Write([]byte("<script>document.getElementById('captchaSolution').value = '';</script><img src='img/ajax/not-available.png' /> Wrong captcha solution"))

	} else {

		userInfo.Username = template.HTMLEscapeString(r.FormValue("username"))
		userInfo.Password = template.HTMLEscapeString(r.FormValue("pass"))
		confirm := template.HTMLEscapeString(r.FormValue("confirm"))
		userInfo.Email = template.HTMLEscapeString(r.FormValue("email"))
		userInfo.IpAddress, _, _ = net.SplitHostPort(r.RemoteAddr)

		if len(userInfo.Username) < 3 || len(userInfo.Username) > 12 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Please choose a username between 3 and 12 characters long."))

		} else if userInfo.Password != confirm {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Your password and confirm password did not match"))

		} else if len(userInfo.Password) < 5 || len(userInfo.Password) > 32 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Password must be at between 5 to 32 characters long"))

		} else if len(userInfo.Email) < 5 || len(userInfo.Email) > 30 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Please choose an email between 5 and 30 characters long"))
		} else {
			return userInfo, true
		}
	}
	return userInfo, false
}

// After all credentials are validated adds users info to database
// host param is used to handle corner case for localhost testing
// returns an error if there was a problem
func (userInfo *UserInfo) Register(host string) (string, error) {

	//check if database connection is open
	if db.Ping() != nil {
		return "<img src='img/ajax/not-available.png' /> Can't ping database.",
			errors.New("DATABASE DOWN in register.go Register 0")
	}

	//check if username already exists, if it does then break out and inform user
	//use javascript for as well as check this in backend in Golang
	var name string
	//checking if name exists
	_ = db.QueryRow("SELECT username FROM userinfo WHERE username=?", userInfo.Username).Scan(&name)

	if userInfo.Username == name {
		return "<img src='img/ajax/not-available.png' /> Username already exist. Please choose another username",
			fmt.Errorf("Prevented host %s from choosing duplicate username %s\n", userInfo.IpAddress, userInfo.Username)
	}

	key, err := hashPass(userInfo.Username, userInfo.Password)
	if err != nil {
		return "<img src='img/ajax/not-available.png' /> Can't hash password.", err
	}

	//inserting into database
	stmt, err := db.Prepare("INSERT userinfo SET username=?, password=?, email=?, date=?, time=?, host=?, verify=?, captcha=?")
	defer stmt.Close()

	if err != nil {
		return "<img src='img/ajax/not-available.png' /> Can't prepare database insertion.", err
	}

	date := time.Now()

	_, err = stmt.Exec(userInfo.Username, key, userInfo.Email, date, date, userInfo.IpAddress, "NO", 0)
	if err != nil {
		return "<img src='img/ajax/not-available.png' /> Username is already taken.", err
	}

	log.Printf("Account %s was created in userinfo table.\n", userInfo.Username)
	userInfo.Token = RandomString()

	//preparing token activation
	stmt, err = db.Prepare("INSERT activate SET username=?, token=?, email=?, expire=?")
	if err != nil {
		return "<img src='img/ajax/not-available.png' /> Can't prepare token activation", err
	}

	_, err = stmt.Exec(userInfo.Username, userInfo.Token, userInfo.Email, date)
	if err != nil {
		return "<img src='img/ajax/not-available.png' /> Can't execute token activation", err
	}

	if host == "localhost" { // handle corner case for localhost testing
		stmt, err = db.Prepare("UPDATE userinfo SET country=? WHERE username=?")
		if err != nil {
			return "<img src='img/ajax/not-available.png' /> Can't prepare to set country.", err
		}

		_, err = stmt.Exec("globe", userInfo.Username)
		if err != nil {
			return "<img src='img/ajax/not-available.png' /> Can't execute to set country.", err
		}
	} else {
		// updates players country in database when they register for the first time
		setCountry(userInfo.Username, userInfo.IpAddress)
	}
	return "", nil
}

func (userInfo *UserInfo) registerRating() (string, error) {

	//setting up player's rating
	stmt, err := db.Prepare("INSERT rating SET username=?, bullet=?, blitz=?, standard=?, correspondence=?, bulletRD=?, blitzRD=?, standardRD=?, correspondenceRD=?")
	if err != nil {
		return "<img src='img/ajax/not-available.png' /> Can't prepare to update player's rating.", err
	}

	_, err = stmt.Exec(userInfo.Username, "1500", "1500", "1500", "1500", "350.0", "350.0", "350.0", "350.0")
	if err != nil {
		return "<img src='img/ajax/not-available.png' /> Can't execute player's rating.", err
	}

	// add player to rating history table
	stmt, err = db.Prepare("INSERT ratinghistory SET username=?")
	if err != nil {
		return "<img src='img/ajax/not-available.png' /> Can't insert into rating history", err
	}

	_, err = stmt.Exec(userInfo.Username)
	if err != nil {
		return "<img src='img/ajax/not-available.png' /> Can't execute into rating history", err
	}

	return "", nil
}

// returns the password hash
func hashPass(username string, password string) (string, error) {

	dk, err := scrypt.Key([]byte(password), []byte(username), 16384, 8, 1, 32)

	if err != nil {
		log.Println(err)
		return "", err
	}
	return hex.EncodeToString(dk), nil
}
