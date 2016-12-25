package gostuff

import (
	"encoding/hex"
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

//checks database to see if username and token is correct and then updates the new password
//with the encrypted password
func ProcessResetPass(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(404)
		http.ServeFile(w, r, "404.html")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong captcha solution! Please try again."))
	} else {
		username := template.HTMLEscapeString(r.FormValue("user"))
		token := template.HTMLEscapeString(r.FormValue("token"))
		password := template.HTMLEscapeString(r.FormValue("pass"))
		confirm := template.HTMLEscapeString(r.FormValue("confirm"))
		ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)

		problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
		defer problems.Close()
		log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

		//check if password and confirm match
		if password != confirm {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> The password you entered does not match. Please try again."))
			return
		}
		if len(password) < 5 || len(password) > 32 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> The password you entered must be 5-32 characters long."))
			return
		}

		//check if database connection is open
		if db.Ping() != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later."))
			log.Println("account.go processResetPass 0 DATABASE DOWN!")
			return
		}
		var tokenInDB string

		//checking if token matches the one entered by user
		err2 := db.QueryRow("SELECT token FROM forgot WHERE username=?", username).Scan(&tokenInDB)
		if err2 != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/token combination."))
			log.Println(err2)
			return
		}
		if tokenInDB != token {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/token combination."))
			log.Printf("FAILED PASSWORD RESET IP: %s  Method: %s Location: %s Agent: %s\n", ipAddress, r.Method, r.URL.Path, r.UserAgent())
			return
		}

		//setting password for user and deleting a row from the forgot table
		stmt, err := db.Prepare("UPDATE userinfo SET password=? where username=?")
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 3."))
			log.Println(err)
			return
		}

		//hashing password
		dk, err1 := scrypt.Key([]byte(password), []byte(username), 16384, 8, 1, 32)
		key := hex.EncodeToString(dk)
		if err1 != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 4."))
			log.Println(err)
			return
		}

		_, err = stmt.Exec(key, username)
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 5."))
			log.Println(err)
			return
		}

		// delete row from forgot table
		stmt, err = db.Prepare("DELETE FROM forgot where username=?")
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 6."))
			log.Println(err)
			return
		}

		res, err := stmt.Exec(username)
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 7"))
			log.Println(err)
			return
		}
		stmt.Close()
		affect, err := res.RowsAffected()
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 8."))
			log.Println(err)
			return
		}

		log.Printf("%d row was deleted from the activate table by user %s\n", affect, username)
		w.Write([]byte("<img src='img/ajax/available.png' /> Your password is now changed!"))
	}
}

//activates user account and deletes entry from database
func ProcessActivate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(404)
		http.ServeFile(w, r, "404.html")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !captcha.VerifyString(template.HTMLEscapeString(r.FormValue("captchaId")), template.HTMLEscapeString(r.FormValue("captchaSolution"))) {
		w.Write([]byte("<img src='img/ajax/not-available.png'/> Wrong captcha solution"))

	} else {
		username := template.HTMLEscapeString(r.FormValue("user"))
		token := template.HTMLEscapeString(r.FormValue("token"))
		ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)

		problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
		defer problems.Close()
		log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

		//check if database connection is open
		if db.Ping() != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png'/> We are having trouble with our server. Report to admin Error 11"))
			log.Println("DATABASE DOWN!")
			return
		}
		var tokenInDB string

		//checking if token matches the one entered by user
		err2 := db.QueryRow("SELECT token FROM activate WHERE username=?", username).Scan(&tokenInDB)
		if err2 != nil || tokenInDB != token {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/token combination"))
			log.Printf("FAILED ACTIVATION Host: %s  Method: %s Location: %s Agent: %s\n", ipAddress, r.Method, r.URL.Path, r.UserAgent())
			return
		}
		//setting verify to yes and deleting row from activate table as well as captcha to zero to signfy user unlocked account
		stmt, err := db.Prepare("UPDATE userinfo SET verify=?, captcha=? where username=?")
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 12."))
			log.Println(err)
			return
		}

		_, err = stmt.Exec("YES", 0, username)
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 13."))
			log.Println(err)
			return
		}

		//now user may login so we can redirect while token deletion proceeds in the background
		message := "<script>window.location = 'login?user=" + username + "';</script>"
		w.Write([]byte(message))

		stmt, err = db.Prepare("DELETE FROM activate where username=?")
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 15."))
			log.Println(err)
			return
		}

		res, err := stmt.Exec(username)
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 16."))
			log.Println(err)
			return
		}
		stmt.Close()
		affect, err := res.RowsAffected()
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 17."))
			log.Println(err)
			return
		}
		log.Printf("%d row was deleted from the activate table by user %s\n", affect, username)

		//once a user activates his acccount update the highscore board so it shows him as new user
		UpdateHighScore()
	}
}

//connects to database and adds token to the forgot table
func ProcessForgot(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(404)
		http.ServeFile(w, r, "404.html")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)

	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong captcha solution! Please try again."))

	} else {
		username := template.HTMLEscapeString(r.FormValue("user"))
		email := template.HTMLEscapeString(r.FormValue("email"))

		problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
		defer problems.Close()
		log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

		//check if database connection is open
		if db.Ping() != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Error 19"))
			return
		}
		var match string
		//checking if email and username entered matches what is in database
		err2 := db.QueryRow("SELECT email FROM userinfo WHERE username=?", username).Scan(&match)
		if err2 != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong email/username combination."))
			log.Println(err2)
			return
		}
		if match != email {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong email/username combination."))
			log.Printf("FAILED SEND PASSWORD RESET TO EMAIL Host: %s  Method: %s Location: %s Agent: %s\n", ipAddress, r.Method, r.URL.Path, r.UserAgent())
			return
		}

		token := RandomString()
		//check for duplicate entry in forgot table
		var found string
		_ = db.QueryRow("SELECT token FROM forgot WHERE username=?", username).Scan(&found)

		if found != "" {
			w.Write([]byte("<img src='img/ajax/available.png' /> Activation token resent to your email."))
			SendForgot(email, found, r.Host)
			return
		}

		//preparing token activation
		stmt, err := db.Prepare("INSERT forgot SET username=?, token=?, expire=?")
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Report to admin Error 20"))
			log.Println(err)
			return
		}

		res, err := stmt.Exec(username, token, time.Now())
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Report to admin Error 21"))
			log.Println(err)
			return
		}
		stmt.Close()
		affect, err := res.RowsAffected()
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Report to admin Error 22"))
			log.Println(err)
			return
		}

		log.Printf("%d rows were affected by the the token activation check.\n", affect)

		//sends pasword reset information to email of user
		go func(email, token, url string) {
			SendForgot(email, token, url)
		}(email, token, r.Host)

		w.Write([]byte("<img src='img/ajax/available.png' /> Your password reset information has been sent your email."))
	}
}
