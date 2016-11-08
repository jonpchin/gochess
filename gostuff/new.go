// processRegister.go
package gostuff

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dchest/captcha"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/scrypt"
)

// global SessionManager["username"] = sessionID
var SessionManager = make(map[string]string)

//stores rating inforation about user in memory
type Person struct {
	User     string
	Bullet   int16
	Blitz    int16
	Standard int16
}

//process user input when signing in
func ProcessLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(404)
		http.ServeFile(w, r, "404.html")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	userName := template.HTMLEscapeString(r.FormValue("user"))
	passWord := template.HTMLEscapeString(r.FormValue("password"))

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	capID := template.HTMLEscapeString(r.FormValue("captchaId"))
	capSol := template.HTMLEscapeString(r.FormValue("captchaSolution"))
	ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)

	if capSol == "" { //then assume user was not displayed captcha

		//check if database connection is open
		if db.Ping() != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Error 24"))
			log.Println("DATABASE DOWN!")
			return
		}

		//hashing password
		dk, err1 := scrypt.Key([]byte(passWord), []byte(userName), 16384, 8, 1, 32)
		key := hex.EncodeToString(dk)
		if err1 != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 25"))
			log.Println(err1)
			return
		}

		var pass string
		var verify string
		var captcha int
		var email string

		//getting password, verify, captcha and email
		err2 := db.QueryRow("SELECT password, verify, captcha, email FROM userinfo WHERE username=?", userName).Scan(&pass, &verify, &captcha, &email)
		if err2 != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/password combination."))
			log.Println("Incorrect login for ", userName)
			return
		}

		//if user entered password incorrect two times or more then they need to enter captcha to login
		if captcha >= 2 {
			w.Write([]byte("<script>$('#hiddenCap').show();</script><img src='img/ajax/not-available.png' /> You entered password incorrecty too many times. Now you need to enter captcha."))
			return
		}
		//checking if password entered by user matches encrypted key
		if pass != key {
			log.Printf("FAILED LOGIN IP: %s  Method: %s Location: %s Agent: %s\n", ipAddress, r.Method, r.URL.Path, r.UserAgent())

			if captcha == 1 {
				w.Write([]byte("<script>$('#hiddenCap').show();</script><img src='img/ajax/not-available.png' />  You entered password incorrectly too many times. Now you need to enter captcha."))
			} else {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/password combination."))
			}

			//add 1 to captcha if password was incorrect
			stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
			defer stmt.Close()

			if err != nil {
				log.Println(err)
				return
			}
			captcha = captcha + 1

			_, err = stmt.Exec(captcha, userName)
			if err != nil {
				log.Println(err)
			}
			return
		}
		if verify != "YES" {
			var tokenInDB string
			var email string

			log.Printf("%s needs to activate his account before logging in.\n", userName)
			w.Write([]byte("<img src='img/ajax/not-available.png' /> You must activate your account by entering the activation token in your email at the activation page. An email has been sent again containing your activation code."))

			//checking if token matches the one entered by user
			err2 := db.QueryRow("SELECT token, email FROM activate WHERE username=?", userName).Scan(&tokenInDB, &email)
			if err2 != nil {
				log.Println(err2)
			} else {
				go func(email, tokenInDB, userName string) {
					Sendmail(email, tokenInDB, userName)
				}(email, tokenInDB, userName)
			}
			return
		}
		// update captcha to zero since login was a sucess
		stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
		defer stmt.Close()
		if err != nil {
			log.Println(err)
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Error in captcha section 3"))
			return
		}

		_, err = stmt.Exec(0, userName)
		if err != nil {
			log.Println(err)
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Error in captcha section 4"))
			return
		}

		expiration := time.Now().Add(3 * 24 * time.Hour)
		cookie := http.Cookie{Name: "username", Value: userName, Secure: true, HttpOnly: true, Expires: expiration}
		http.SetCookie(w, &cookie)

		//generating random session ID to be stored in the backend
		sessionID := RandomString()

		cookie = http.Cookie{Name: "sessionID", Value: sessionID, Secure: true, HttpOnly: true, Expires: expiration}
		http.SetCookie(w, &cookie)

		SessionManager[userName] = sessionID

		w.Write([]byte("<script>window.location = '/memberHome'</script>"))

	} else if !captcha.VerifyString(capID, capSol) {
		w.Write([]byte("<script>document.getElementById('captchaSolution').value = '';</script><img src='img/ajax/not-available.png' /> Wrong captcha solution! Please try again."))
	} else {
		//check if database connection is open
		if db.Ping() != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Error 24"))
			log.Println("DATABASE DOWN!")
			return
		}

		//hashing password
		dk, err1 := scrypt.Key([]byte(passWord), []byte(userName), 16384, 8, 1, 32)
		key := hex.EncodeToString(dk)
		if err1 != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 25"))
			log.Println(err1)
			return
		}

		var pass string
		var verify string
		var captcha int
		var email string
		var token string

		//getting password, verify, captcha, email from database
		err2 := db.QueryRow("SELECT password, verify, captcha, email FROM userinfo WHERE username=?", userName).Scan(&pass, &verify, &captcha, &email)
		if err2 != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/password combination."))
			//check if there was more then one incorrect login attempt
			log.Println(err2)
			return
		}

		if captcha == 5 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> This account has been deactivated becaue too many incorrect login attempts were made. An email has been sent to you on instructions on reactivating your account."))

			//increment captcha counter
			stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
			defer stmt.Close()
			if err != nil {
				log.Println(err)
				return
			}
			captcha = captcha + 1

			_, err = stmt.Exec(captcha, userName)
			if err != nil {
				log.Println(err)
				return
			}

			//create activation token in database and send user notifying them that their was five incorrect login attempts
			stmt, err = db.Prepare("INSERT activate SET username=?, token=?, email=?, expire=?")
			if err != nil {
				log.Println(err)
				return
			}

			token = RandomString()
			_, err = stmt.Exec(userName, token, email, time.Now())
			if err != nil {
				log.Println(err)
				return
			}
			//sends email to user with the token activation
			go func(email, token, userName, address string) {
				SendAttempt(email, token, userName, address)
			}(email, token, userName, ipAddress)
			return

		} else if captcha > 5 { //tell user on the front end that this account has too many login attempts, resends activation token
			w.Write([]byte("<img src='img/ajax/not-available.png' /> This account has been deactivated because too many incorrect login attempts were made. An email has been sent again regarding how to reactivate your account."))

			err2 := db.QueryRow("SELECT token FROM activate WHERE username=?", userName).Scan(&token)
			if err2 != nil {
				log.Println(err2)
				return
			}

			//sends email again to user with the token activation
			go func(email, token, userName, address string) {
				SendAttempt(email, token, userName, address)
			}(email, token, userName, ipAddress)

			return
		}
		if pass != key {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/password combination."))
			log.Printf("FAILED LOGIN IP: %s  Method: %s Location: %s Agent: %s\n", ipAddress, r.Method, r.URL.Path, r.UserAgent())

			//add 1 to captcha if password was incorrect
			stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
			defer stmt.Close()

			if err != nil {
				log.Println(err)
			}
			captcha = captcha + 1

			_, err = stmt.Exec(captcha, userName)
			if err != nil {
				log.Println(err)
			}
			return
		}
		if verify != "YES" {
			var tokenInDB string
			var email string

			log.Printf("%s needs to activate his account before logging in.\n", userName)
			w.Write([]byte("<img src='img/ajax/not-available.png' /> You must activate your account by entering the activation token in your email at the activation page. An email has been sent again containing your activation code."))

			//checking if token matches the one entered by user
			err2 := db.QueryRow("SELECT token, email FROM activate WHERE username=?", userName).Scan(&tokenInDB, &email)
			if err2 != nil {
				log.Println(err2)
			} else {
				go func(email, tokenInDB, userName string) {
					Sendmail(email, tokenInDB, userName)
				}(email, tokenInDB, userName)
			}
			return
		}

		// update captcha to zero since login was a sucess
		stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
		defer stmt.Close()

		if err != nil {
			log.Println(err)
		}

		_, err = stmt.Exec(0, userName)
		if err != nil {
			log.Println(err)
		}

		expiration := time.Now().Add(3 * 24 * time.Hour)
		cookie := http.Cookie{Name: "username", Value: userName, Secure: true, HttpOnly: true, Expires: expiration}
		http.SetCookie(w, &cookie)

		//generating random session ID to be stored in the backend
		sessionID := RandomString()

		cookie = http.Cookie{Name: "sessionID", Value: sessionID, Secure: true, HttpOnly: true, Expires: expiration}
		http.SetCookie(w, &cookie)

		SessionManager[userName] = sessionID

		w.Write([]byte("<script>window.location = '/memberHome'</script>"))
	}
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

		userName := template.HTMLEscapeString(r.FormValue("username"))
		passWord := template.HTMLEscapeString(r.FormValue("pass"))
		confirm := template.HTMLEscapeString(r.FormValue("confirm"))
		email := template.HTMLEscapeString(r.FormValue("email"))
		ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)

		if len(userName) < 3 || len(userName) > 12 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Please choose a username between 3 and 12 characters long."))

		} else if passWord != confirm {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Your password and confirm password did not match"))

		} else if len(passWord) < 5 || len(passWord) > 32 {
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
			checkName := db.QueryRow("SELECT username FROM userinfo WHERE username=?", userName).Scan(&name)

			if checkName != sql.ErrNoRows {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> Username already exist. Please choose another username"))
				log.Printf("Prevented host %s from choosing duplicate username %s\n", ipAddress, userName)
				return

			}
			message := "<script>$('#register').hide();</script><img src='img/ajax/available.png' /> Hello " + userName + "! Please check email for instructions to verify your account."
			//if reached here just notify user to check his email and continue on with the account creation
			w.Write([]byte(message))

			//hashing password
			dk, err1 := scrypt.Key([]byte(passWord), []byte(userName), 16384, 8, 1, 32)
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

			_, err = stmt.Exec(userName, key, email, date, date, ipAddress, "NO", 0)
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 30"))
				log.Println(err)
				return
			}

			log.Printf("Account %s was created in userinfo table.\n", userName)

			token := RandomString()

			//preparing token activation
			stmt, err = db.Prepare("INSERT activate SET username=?, token=?, email=?, expire=?")
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 32"))
				log.Println(err)
				return
			}

			_, err = stmt.Exec(userName, token, email, date)
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 33"))
				log.Println(err)
				return
			}
			//sends email to user
			go func(email, token, userName string) {
				Sendmail(email, token, userName)
			}(email, token, userName)

			//setting up player's rating
			stmt, err = db.Prepare("INSERT rating SET username=?, bullet=?, blitz=?, standard=?, bulletRD=?, blitzRD=?, standardRD=?")
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 34"))
				log.Println(err)
				return
			}

			_, err = stmt.Exec(userName, "1500", "1500", "1500", "350.0", "350.0", "350.0")
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 35"))
				log.Println(err)
				return
			}
		}
	}
}

func CheckUserName(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userName := template.HTMLEscapeString(r.FormValue("username"))

		//making sure username fits length requirement
		if len(userName) < 3 || len(userName) > 12 {
			return
		}
		ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)
		//check if database connection is open
		if db.Ping() != nil {
			fmt.Printf("ERROR 2 PINGING DB IP: %s \n", ipAddress)
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Please come back later."))
			return
		}

		var name string
		//checking if name exists
		checkName := db.QueryRow("SELECT username FROM userinfo WHERE username=?", userName).Scan(&name)
		switch {
		case checkName == sql.ErrNoRows:
			w.Write([]byte(" <img src='img/ajax/available.png' /> Username available"))
			fmt.Printf("Username %s is available.\n", userName)
		case checkName != nil:
			fmt.Printf("ERROR 3 CHECKNAME IP is %s\n", ipAddress)
		default:
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Username taken"))
		}
	}
}

func RandomString() string {
	size := 12 // change the length of the generated random string here 12 is actually 16

	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println("new.go RandomString 1 ", err)
	}

	token := base64.URLEncoding.EncodeToString(rb)
	return token
}
