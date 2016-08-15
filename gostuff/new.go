// processRegister.go
package gostuff

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/dchest/captcha"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/scrypt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// global SessionManager["username"] = sessionID
var SessionManager = make(map[string]string)

//stores rating inforation about user in memory
type Person struct {
	User     string
	Bullet   int16
	Blitz    int16
	Standard int16
	//email   string
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
	log.SetOutput(problems)

	capID := template.HTMLEscapeString(r.FormValue("captchaId"))
	capSol := template.HTMLEscapeString(r.FormValue("captchaSolution"))

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
			log.Println("new.go ProcessLogin 1 ", err1)
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
			browser := r.UserAgent()
			log.Printf("FAILED LOGIN IP: %s  Method: %s Location: %s Agent: %s\n", r.RemoteAddr, r.Method, r.URL.Path, browser)

			if captcha == 1 {
				w.Write([]byte("<script>$('#hiddenCap').show();</script><img src='img/ajax/not-available.png' />  You entered password incorrectly too many times. Now you need to enter captcha."))
			} else {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/password combination."))
			}

			//add 1 to captcha if password was incorrect
			stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
			if err != nil {
				log.Println("Error in captcha section 3")
				return
			}
			captcha = captcha + 1

			_, err = stmt.Exec(captcha, userName)
			if err != nil {
				log.Println("Error in captcha section 4")
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
				log.Println("new.go ProcessLogin 3 ", err2)
			} else {
				go func(email, tokenInDB, userName string){
					Sendmail(email, tokenInDB, userName)
				}(email, tokenInDB, userName)
			}
			return
		}
		// update captcha to zero since login was a sucess
		stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
		if err != nil {
			log.Println("new.go ProcessLogin 4 ", err)
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Error in captcha section 3"))
			return
		}

		_, err = stmt.Exec(0, userName)
		if err != nil {
			log.Println("new.go ProcessLogin 5 ", err)
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
			log.Println("new.go DATABASE DOWN!")
			return
		}

		//hashing password
		dk, err1 := scrypt.Key([]byte(passWord), []byte(userName), 16384, 8, 1, 32)
		key := hex.EncodeToString(dk)
		if err1 != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 25"))
			log.Println("error in hasing password new.go ", err1)
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
			log.Println("new.go ProcessLogin error 6", err2)
			return
		}

		if captcha == 5 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> This account has been deactivated becaue too many incorrect login attempts were made. An email has been sent to you on instructions on reactivating your account."))

			//increment captcha counter
			stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
			if err != nil {
				log.Println("new.go ProcessLogin error 7")
				return
			}
			captcha = captcha + 1

			_, err = stmt.Exec(captcha, userName)
			if err != nil {
				log.Println("new.go ProcessLogin error 8")
				return
			}

			//create activation token in database and send user notifying them that their was five incorrect login attempts
			stmt, err = db.Prepare("INSERT activate SET username=?, token=?, email=?, expire=?")
			if err != nil {
				log.Println("new.go ProcessLogin 9 ", err)
				return
			}
			date := time.Now()
			token = RandomString()
			_, err = stmt.Exec(userName, token, email, date)
			if err != nil {
				log.Println("new.go ProcessLogin 10 ", err)
				return
			}
			//sends email to user with the token activation
			go func(email, token, userName, address string){
				SendAttempt(email, token, userName, address)
			}(email, token, userName, r.RemoteAddr)
			return

		} else if captcha > 5 { //tell user on the front end that this account has too many login attempts, resends activation token
			w.Write([]byte("<img src='img/ajax/not-available.png' /> This account has been deactivated because too many incorrect login attempts were made. An email has been sent again regarding how to reactivate your account."))

			err2 := db.QueryRow("SELECT token FROM activate WHERE username=?", userName).Scan(&token)
			if err2 != nil {
				log.Println("new.go ProcessLogin 11 ", err2)
				return
			}

			//sends email again to user with the token activation
			go func(email, token, userName, address string){
				SendAttempt(email, token, userName, address)
			}(email, token, userName, r.RemoteAddr)
			
			return
		}
		if pass != key {
			browser := r.UserAgent()
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/password combination."))
			log.Printf("FAILED LOGIN IP: %s  Method: %s Location: %s Agent: %s\n", r.RemoteAddr, r.Method, r.URL.Path, browser)

			//add 1 to captcha if password was incorrect
			stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
			if err != nil {
				log.Println("new.go ProcessLogin 12 ", err)
			}
			captcha = captcha + 1

			_, err = stmt.Exec(captcha, userName)
			if err != nil {
				log.Println("new.go ProcessLogin 13 ", err)
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
				log.Println("new.go ProcessLogin 14 ", err2)
			} else {
				
				go func(email, tokenInDB, userName string){
					Sendmail(email, tokenInDB, userName)
				}(email, tokenInDB, userName)		
			}
			return
		}

		// update captcha to zero since login was a sucess
		stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
		if err != nil {
			log.Println("new.go ProcessLogin 15 ", err)
		}

		_, err = stmt.Exec(0, userName)
		if err != nil {
			log.Println("new.go ProcessLogin 16 ", err)
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

			log.SetOutput(problems)

			//check if database connection is open
			if db.Ping() != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 27"))
				log.Println("new.go processRegister 1 DATABASE DOWN!")
				return
			}

			//check if username already exists, if it does then break out and inform user
			//use javascript for as well as check this in backend in Golang
			var name string
			//checking if name exists
			checkName := db.QueryRow("SELECT username FROM userinfo WHERE username=?", userName).Scan(&name)

			if checkName != sql.ErrNoRows {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> Username already exist. Please choose another username"))
				log.Printf("Prevented host %s from choosing duplicate username %s\n", r.RemoteAddr, userName)
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
				log.Println("new.go processRegister 2 ", err)
				return
			}

			//inserting into database
			stmt, err := db.Prepare("INSERT userinfo SET username=?, password=?, email=?, date=?, time=?, host=?, verify=?, captcha=?")
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 29"))
				log.Println("new.go processRegister 3 ", err)
				return
			}

			date := time.Now()
			res, err := stmt.Exec(userName, key, email, date, date, r.RemoteAddr, "NO", 0)
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 30"))
				log.Println("new.go processRegister 4 ", err)
				return
			}

			id, err := res.LastInsertId()
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 31"))
				log.Println("new.go processRegister 5 ", err)
				return
			}
			log.Printf("Account %s and id %d was created in userinfo table.\n", userName, id)

			token := RandomString()

			//preparing token activation
			stmt, err = db.Prepare("INSERT activate SET username=?, token=?, email=?, expire=?")
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 32"))
				log.Println("new.go processRegister 6 ", err)
				return
			}

			res, err = stmt.Exec(userName, token, email, date)
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 33"))
				log.Println("error in executing token activation new.go", err)
				return
			}
			//sends email to user
			go func(email, token, userName string){
				Sendmail(email, token, userName)
			}(email, token, userName)
			
			//setting up player's rating
			stmt, err = db.Prepare("INSERT rating SET username=?, bullet=?, blitz=?, standard=?, bulletRD=?, blitzRD=?, standardRD=?")
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 34"))
				log.Println("new.go processRegister 7 ", err)
				return

			}

			res, err = stmt.Exec(userName, "1500", "1500", "1500", "350.0", "350.0", "350.0")
			if err != nil {
				w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Report to admin Error 35"))
				log.Println("error in setting up player's rating new.go", err)
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

		//check if database connection is open
		if db.Ping() != nil {
			fmt.Printf("ERROR 2 PINGING DB IP: %s \n", r.RemoteAddr)
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
			fmt.Printf("ERROR 3 CHECKNAME IP is %s\n", r.RemoteAddr)
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
