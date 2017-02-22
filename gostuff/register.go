package gostuff

import (
	"database/sql"
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
	Username  string
	Password  string
	Email     string
	IpAddress string
	Token     string
}

//processes the users input when signing up
func ProcessRegister(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(404)
		http.ServeFile(w, r, "404.html")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if !captcha.VerifyString(template.HTMLEscapeString(r.FormValue("captchaId")), template.HTMLEscapeString(r.FormValue("captchaSolution"))) {
		w.Write([]byte("<script>document.getElementById('captchaSolution').value = '';</script><img src='img/ajax/not-available.png' /> Wrong captcha solution"))

	} else {
		var userInfo UserInfo
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

			problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
			defer problems.Close()
			log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

			err := userInfo.Register(w, r)
			if err == nil {
				//sends email to user
				go Sendmail(userInfo.Email, userInfo.Token, userInfo.Username, r.Host)

				message := "<script>$('#register').hide();</script><img src='img/ajax/available.png' /> Hello " +
					userInfo.Username + "! Please check email for instructions to verify your account."
				//if reached here just notify user to check his email and continue on with the account creation
				w.Write([]byte(message))
			} else {
				log.Println(err)
			}
		}
	}
}

//after all credentials are validated adds users info to database
// returns an error if there was a problem
func (userInfo *UserInfo) Register(w http.ResponseWriter, r *http.Request) error {

	//check if database connection is open
	if db.Ping() != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 27"))
		return errors.New("DATABASE DOWN in register.go Register 0")
	}

	//check if username already exists, if it does then break out and inform user
	//use javascript for as well as check this in backend in Golang
	var name string
	//checking if name exists
	checkName := db.QueryRow("SELECT username FROM userinfo WHERE username=?", userInfo.Username).Scan(&name)

	if checkName != sql.ErrNoRows {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Username already exist. Please choose another username"))
		return fmt.Errorf("Prevented host %s from choosing duplicate username %s\n", userInfo.IpAddress, userInfo.Username)
	}

	key, err := hashPass(userInfo.Username, userInfo.Password)
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 28"))
		return err
	}

	//inserting into database
	stmt, err := db.Prepare("INSERT userinfo SET username=?, password=?, email=?, date=?, time=?, host=?, verify=?, captcha=?")
	defer stmt.Close()

	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 29"))
		return err
	}

	date := time.Now()

	_, err = stmt.Exec(userInfo.Username, key, userInfo.Email, date, date, userInfo.IpAddress, "NO", 0)
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 30"))
		return err
	}

	log.Printf("Account %s was created in userinfo table.\n", userInfo.Username)

	userInfo.Token = RandomString()

	//preparing token activation
	stmt, err = db.Prepare("INSERT activate SET username=?, token=?, email=?, expire=?")
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 32"))
		return err
	}

	_, err = stmt.Exec(userInfo.Username, userInfo.Token, userInfo.Email, date)
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 33"))
		return err
	}

	//setting up player's rating
	stmt, err = db.Prepare("INSERT rating SET username=?, bullet=?, blitz=?, standard=?, correspondence=?, bulletRD=?, blitzRD=?, standardRD=?, correspondenceRD=?")
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 34"))
		return err
	}

	_, err = stmt.Exec(userInfo.Username, "1500", "1500", "1500", "1500", "350.0", "350.0", "350.0", "350.0")
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 35"))
		return err
	}

	// add player to rating history table
	stmt, err = db.Prepare("INSERT ratinghistory SET username=?")
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 36"))
		return err
	}

	_, err = stmt.Exec(userInfo.Username)
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 37"))
		return err
	}

	if r.Host == "localhost" { // handle corner case for localhost testing
		stmt, err = db.Prepare("UPDATE userinfo SET country=? WHERE username=?")
		if err != nil {
			return err
		}

		_, err = stmt.Exec("globe", userInfo.Username)
		if err != nil {
			return err
		}
	} else {
		// updates players country in database when they register for the first time
		setCountry(userInfo.Username, userInfo.IpAddress)
	}
	return nil
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
