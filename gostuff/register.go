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

		username := template.HTMLEscapeString(r.FormValue("username"))
		password := template.HTMLEscapeString(r.FormValue("pass"))
		confirm := template.HTMLEscapeString(r.FormValue("confirm"))
		email := template.HTMLEscapeString(r.FormValue("email"))
		ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)

		if len(username) < 3 || len(username) > 12 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Please choose a username between 3 and 12 characters long."))

		} else if password != confirm {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Your password and confirm password did not match"))

		} else if len(password) < 5 || len(password) > 32 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Password must be at between 5 to 32 characters long"))

		} else if len(email) < 5 || len(email) > 30 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Please choose an email between 5 and 30 characters long"))

		} else {

			problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
			defer problems.Close()
			log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

			//check if database connection is open
			if db.Ping() != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 27"))
				log.Println("DATABASE DOWN!")
				return
			}

			//check if username already exists, if it does then break out and inform user
			//use javascript for as well as check this in backend in Golang
			var name string
			//checking if name exists
			checkName := db.QueryRow("SELECT username FROM userinfo WHERE username=?", username).Scan(&name)

			if checkName != sql.ErrNoRows {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> Username already exist. Please choose another username"))
				log.Printf("Prevented host %s from choosing duplicate username %s\n", ipAddress, username)
				return

			}
			message := "<script>$('#register').hide();</script><img src='img/ajax/available.png' /> Hello " + username + "! Please check email for instructions to verify your account."
			//if reached here just notify user to check his email and continue on with the account creation
			w.Write([]byte(message))

			//hashing password
			dk, err1 := scrypt.Key([]byte(password), []byte(username), 16384, 8, 1, 32)
			key := hex.EncodeToString(dk)
			if err1 != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 28"))
				log.Println(err)
				return
			}

			//inserting into database
			stmt, err := db.Prepare("INSERT userinfo SET username=?, password=?, email=?, date=?, time=?, host=?, verify=?, captcha=?")
			defer stmt.Close()

			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 29"))
				log.Println(err)
				return
			}

			date := time.Now()

			_, err = stmt.Exec(username, key, email, date, date, ipAddress, "NO", 0)
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 30"))
				log.Println(err)
				return
			}

			log.Printf("Account %s was created in userinfo table.\n", username)

			token := RandomString()

			//preparing token activation
			stmt, err = db.Prepare("INSERT activate SET username=?, token=?, email=?, expire=?")
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 32"))
				log.Println(err)
				return
			}

			_, err = stmt.Exec(username, token, email, date)
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 33"))
				log.Println(err)
				return
			}
			//sends email to user
			go func(email, token, username, url string) {
				Sendmail(email, token, username, url)
			}(email, token, username, r.Host)

			//setting up player's rating
			stmt, err = db.Prepare("INSERT rating SET username=?, bullet=?, blitz=?, standard=?, bulletRD=?, blitzRD=?, standardRD=?")
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 34"))
				log.Println(err)
				return
			}

			_, err = stmt.Exec(username, "1500", "1500", "1500", "350.0", "350.0", "350.0")
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 35"))
				log.Println(err)
				return
			}

			// add player to rating history table
			stmt, err = db.Prepare("INSERT ratinghistory SET username=?")
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 36"))
				log.Println(err)
				return
			}

			_, err = stmt.Exec(username)
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 37"))
				log.Println(err)
				return
			}

			if r.Host == "localhost" { // handle corner case for localhost testing
				stmt, err = db.Prepare("UPDATE userinfo SET country=? WHERE username=?")
				if err != nil {
					log.Println(err)
				}

				_, err = stmt.Exec("globe", username)
				if err != nil {
					log.Println(err)
				}
			} else {
				// updates players country in database when they register for the first time
				setCountry(username, ipAddress)
			}
		}
	}
}
