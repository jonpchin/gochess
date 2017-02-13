package gostuff

import (
	"database/sql"
	"encoding/hex"
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
	username  string
	password  string
	email     string
	ipAddress string
	token     string
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
		userInfo.username = template.HTMLEscapeString(r.FormValue("username"))
		userInfo.password = template.HTMLEscapeString(r.FormValue("pass"))
		confirm := template.HTMLEscapeString(r.FormValue("confirm"))
		userInfo.email = template.HTMLEscapeString(r.FormValue("email"))
		userInfo.ipAddress, _, _ = net.SplitHostPort(r.RemoteAddr)

		if len(userInfo.username) < 3 || len(userInfo.username) > 12 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Please choose a username between 3 and 12 characters long."))

		} else if userInfo.password != confirm {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Your password and confirm password did not match"))

		} else if len(userInfo.password) < 5 || len(userInfo.password) > 32 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Password must be at between 5 to 32 characters long"))

		} else if len(userInfo.email) < 5 || len(userInfo.email) > 30 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Please choose an email between 5 and 30 characters long"))

		} else {
			success := userInfo.register(w, r)
			if success {
				//sends email to user
				go Sendmail(userInfo.email, userInfo.token, userInfo.username, r.Host)

				message := "<script>$('#register').hide();</script><img src='img/ajax/available.png' /> Hello " +
					userInfo.username + "! Please check email for instructions to verify your account."
				//if reached here just notify user to check his email and continue on with the account creation
				w.Write([]byte(message))
			}
		}
	}
}

//after all credentials are validated adds users info to database
// returns false if there was a problem
func (userInfo *UserInfo) register(w http.ResponseWriter, r *http.Request) bool {
	problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 27"))
		log.Println("DATABASE DOWN!")
		return false
	}

	//check if username already exists, if it does then break out and inform user
	//use javascript for as well as check this in backend in Golang
	var name string
	//checking if name exists
	checkName := db.QueryRow("SELECT username FROM userinfo WHERE username=?", userInfo.username).Scan(&name)

	if checkName != sql.ErrNoRows {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Username already exist. Please choose another username"))
		log.Printf("Prevented host %s from choosing duplicate username %s\n", userInfo.ipAddress, userInfo.username)
		return false
	}

	key, err := hashPass(userInfo.username, userInfo.password)
	if err != nil {
		log.Println(err)
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 28"))
		return false
	}

	//inserting into database
	stmt, err := db.Prepare("INSERT userinfo SET username=?, password=?, email=?, date=?, time=?, host=?, verify=?, captcha=?")
	defer stmt.Close()

	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 29"))
		log.Println(err)
		return false
	}

	date := time.Now()

	_, err = stmt.Exec(userInfo.username, key, userInfo.email, date, date, userInfo.ipAddress, "NO", 0)
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 30"))
		log.Println(err)
		return false
	}

	log.Printf("Account %s was created in userinfo table.\n", userInfo.username)

	userInfo.token = RandomString()

	//preparing token activation
	stmt, err = db.Prepare("INSERT activate SET username=?, token=?, email=?, expire=?")
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 32"))
		log.Println(err)
		return false
	}

	_, err = stmt.Exec(userInfo.username, userInfo.token, userInfo.email, date)
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 33"))
		log.Println(err)
		return false
	}

	//setting up player's rating
	stmt, err = db.Prepare("INSERT rating SET username=?, bullet=?, blitz=?, standard=?, correspondence=?, bulletRD=?, blitzRD=?, standardRD=?, correspondenceRD=?")
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 34"))
		log.Println(err)
		return false
	}

	_, err = stmt.Exec(userInfo.username, "1500", "1500", "1500", "1500", "350.0", "350.0", "350.0", "350.0")
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 35"))
		log.Println(err)
		return false
	}

	// add player to rating history table
	stmt, err = db.Prepare("INSERT ratinghistory SET username=?")
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 36"))
		log.Println(err)
		return false
	}

	_, err = stmt.Exec(userInfo.username)
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 37"))
		log.Println(err)
		return false
	}

	if r.Host == "localhost" { // handle corner case for localhost testing
		stmt, err = db.Prepare("UPDATE userinfo SET country=? WHERE username=?")
		if err != nil {
			log.Println(err)
		}

		_, err = stmt.Exec("globe", userInfo.username)
		if err != nil {
			log.Println(err)
		}
	} else {
		// updates players country in database when they register for the first time
		setCountry(userInfo.username, userInfo.ipAddress)
	}
	return true
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
