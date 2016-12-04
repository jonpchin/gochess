package gostuff

import (
	"crypto/rand"
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
	username := template.HTMLEscapeString(r.FormValue("user"))
	password := template.HTMLEscapeString(r.FormValue("password"))

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
		dk, err1 := scrypt.Key([]byte(password), []byte(username), 16384, 8, 1, 32)
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
		err2 := db.QueryRow("SELECT password, verify, captcha, email FROM userinfo WHERE username=?", username).Scan(&pass, &verify, &captcha, &email)
		if err2 != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/password combination."))
			log.Println("Incorrect login for ", username)
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

			addOneToCaptcha(username, captcha)
			return
		}
		if verify != "YES" {
			needToActivate(w, username)
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

		_, err = stmt.Exec(0, username)
		if err != nil {
			log.Println(err)
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Error in captcha section 4"))
			return
		}

		enterInside(w, username, ipAddress)

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
		dk, err1 := scrypt.Key([]byte(password), []byte(username), 16384, 8, 1, 32)
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
		err2 := db.QueryRow("SELECT password, verify, captcha, email FROM userinfo WHERE username=?", username).Scan(&pass, &verify, &captcha, &email)
		if err2 != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/password combination."))
			//check if there was more then one incorrect login attempt
			log.Println(err2)
			return
		}

		if captcha == 5 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> This account has been deactivated becaue too many incorrect login attempts were made. An email has been sent to you on instructions on reactivating your account."))

			addOneToCaptcha(username, captcha)

			//create activation token in database and send user notifying them that their was five incorrect login attempts
			stmt, err := db.Prepare("INSERT activate SET username=?, token=?, email=?, expire=?")
			if err != nil {
				log.Println(err)
				return
			}

			token = RandomString()
			_, err = stmt.Exec(username, token, email, time.Now())
			if err != nil {
				log.Println(err)
				return
			}
			//sends email to user with the token activation
			go func(email, token, username, address string) {
				SendAttempt(email, token, username, address)
			}(email, token, username, ipAddress)
			return

		} else if captcha > 5 { //tell user on the front end that this account has too many login attempts, resends activation token
			w.Write([]byte("<img src='img/ajax/not-available.png' /> This account has been deactivated because too many incorrect login attempts were made. An email has been sent again regarding how to reactivate your account."))

			err2 := db.QueryRow("SELECT token FROM activate WHERE username=?", username).Scan(&token)
			if err2 != nil {
				log.Println(err2)
				return
			}

			//sends email again to user with the token activation
			go func(email, token, username, address string) {
				SendAttempt(email, token, username, address)
			}(email, token, username, ipAddress)

			return
		}
		if pass != key {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/password combination."))
			log.Printf("FAILED LOGIN IP: %s  Method: %s Location: %s Agent: %s\n", ipAddress, r.Method, r.URL.Path, r.UserAgent())

			addOneToCaptcha(username, captcha)
			return
		}
		if verify != "YES" {
			needToActivate(w, username)
			return
		}
		enterInside(w, username, ipAddress)
		
		// update captcha to zero since login was a sucess
		stmt, err := db.Prepare("update userinfo set captcha=? where username=?")
		defer stmt.Close()

		if err != nil {
			log.Println(err)
		}

		_, err = stmt.Exec(0, username)
		if err != nil {
			log.Println(err)
		}
	}
}

func RandomString() string {
	size := 12 // change the length of the generated random string here 12 is actually 16

	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println("login.go RandomString 1 ", err)
	}

	token := base64.URLEncoding.EncodeToString(rb)
	return token
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

	country := getCountry(username)

	// country cookie can be modified by JS
	cookie = http.Cookie{Name: "country", Value: country, Secure: false, HttpOnly: false, Expires: expiration}
	http.SetCookie(w, &cookie)

	SessionManager[username] = sessionID

	w.Write([]byte("<script>window.location = '/memberHome'</script>"))
}

// sends an email again to reactivate an inactivated account
func needToActivate(w http.ResponseWriter, username string) {
	var tokenInDB string
	var email string

	log.Printf("%s needs to activate his account before logging in.\n", username)
	w.Write([]byte("<img src='img/ajax/not-available.png' /> You must activate your account by entering the activation token in your email at the activation page. An email has been sent again containing your activation code."))

	//checking if token matches the one entered by user
	err2 := db.QueryRow("SELECT token, email FROM activate WHERE username=?", username).Scan(&tokenInDB, &email)
	if err2 != nil {
		log.Println(err2)
	} else {
		go func(email, tokenInDB, username string) {
			Sendmail(email, tokenInDB, username)
		}(email, tokenInDB, username)
	}
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
