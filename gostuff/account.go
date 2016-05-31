package gostuff

import (
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
		userName := template.HTMLEscapeString(r.FormValue("user"))
		token := template.HTMLEscapeString(r.FormValue("token"))
		passWord := template.HTMLEscapeString(r.FormValue("pass"))
		confirm := template.HTMLEscapeString(r.FormValue("confirm"))

		problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
		defer problems.Close()
		log.SetOutput(problems)

		//check if password and confirm match
		if passWord != confirm {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> The password you entered does not match. Please try again."))
			return
		}
		if len(passWord) < 5 || len(passWord) > 32 {
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
		err2 := db.QueryRow("SELECT token FROM forgot WHERE username=?", userName).Scan(&tokenInDB)
		if err2 != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/token combination."))
			log.Println("account.go processResetPass 1 ", err2)
			return  
		}
		if tokenInDB != token {
			browser := r.UserAgent()
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/token combination."))
			log.Printf("FAILED PASSWORD RESET IP: %s  Method: %s Location: %s Agent: %s\n", r.RemoteAddr, r.Method, r.URL.Path, browser)
			return
		}

		//setting password for user and deleting a row from the forgot table
		stmt, err := db.Prepare("UPDATE userinfo SET password=? where username=?")
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 3."))
			log.Println("account.go processResetPass 2 ", err)
			return
		}

		//hashing password
		dk, err1 := scrypt.Key([]byte(passWord), []byte(userName), 16384, 8, 1, 32)
		key := hex.EncodeToString(dk)
		if err1 != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 4."))
			log.Println("account.go processResetPass 3 ", err)
			return
		}

		res, err := stmt.Exec(key, userName)
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 5."))
			log.Println("account.go processResetPass 4 ", err)
			return
		}
		affect, err := res.RowsAffected()
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 6."))
			log.Println("account.go processResetPass 5 ", err)
			return
		}

		log.Printf("%d rows were affected by the the token activation check.\n", affect)

		// delete row from forgot table
		stmt, err = db.Prepare("DELETE FROM forgot where username=?")
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 7."))
			log.Println("account.go processResetPass 6 ", err)
			return
		}

		res, err = stmt.Exec(userName)
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 8"))
			log.Println("account.go processResetPass 7 ", err)
			return
		}
		stmt.Close()
		affect, err = res.RowsAffected()
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 9."))
			log.Println("account.go processResetPass 8 ", err)
			return
		}

		log.Printf("%d row was deleted from the activate table by user %s\n", affect, userName)
		w.Write([]byte("<img src='img/ajax/available.png' /> Your password is now changed!"))

	}
}

//activates user account and deletes enty from database
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
		userName := template.HTMLEscapeString(r.FormValue("user"))
		token := template.HTMLEscapeString(r.FormValue("token"))

		problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
		defer problems.Close()
		log.SetOutput(problems)

		//check if database connection is open
		if db.Ping() != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png'/> We are having trouble with our server. Report to admin Error 11"))
			log.Println("DATABASE DOWN! account.go ProcessActivate 0")
			return
		}
		var tokenInDB string

		//checking if token matches the one entered by user
		err2 := db.QueryRow("SELECT token FROM activate WHERE username=?", userName).Scan(&tokenInDB)
		if err2 != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/token combination"))
			log.Println("account.go processActivate 1", err2)
		}
		if tokenInDB != token {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/token combination"))
			browser := r.UserAgent()
			log.Printf("FAILED ACTIVATION Host: %s  Method: %s Location: %s Agent: %s\n", r.RemoteAddr, r.Method, r.URL.Path, browser)
			return
		}
		//setting verify to yes and deleting row from activate table as well as captcha to zero to signfy user unlocked account
		stmt, err := db.Prepare("UPDATE userinfo SET verify=?, captcha=? where username=?")
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 12."))
			log.Println("account.go processActivate 2 ", err)
			return
		}

		res, err := stmt.Exec("YES", 0, userName)
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 13."))
			log.Println("account.go processActivate 3 ", err)
			return
		}
		affect, err := res.RowsAffected()
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 14."))
			log.Println("account.go processActivate 4 ", err)
			return
		}

		log.Printf("%s is now verified and %d row was updated.\n", userName, affect)

		//now user may login so we can redirect while token deletion proceeds in the background
		message := "<script>window.location = 'login?user=" + userName + "';</script>"
		w.Write([]byte(message))

		// this is safe
		//db.Query("SELECT name FROM users WHERE age=?", req.FormValue("age"))

		// delete
		stmt, err = db.Prepare("DELETE FROM activate where username=?")
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 15."))
			log.Println("account.go processActivate 5 ", err)
			return
		}

		res, err = stmt.Exec(userName)
		if err != nil {
			fmt.Println("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 16.")
			log.Println("account.go processActivate 6 ", err)

			return
		}
		stmt.Close()
		affect, err = res.RowsAffected()
		if err != nil {
			fmt.Println("Something is wrong with the server. Tell admin error 17.")
			log.Println("account.go processActivate 7 ", err)
			return
		}
		log.Printf("%d row was deleted from the activate table by user %s\n", affect, userName)

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
	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong captcha solution! Please try again."))

	} else {
		userName := template.HTMLEscapeString(r.FormValue("user"))
		email := template.HTMLEscapeString(r.FormValue("email"))

		problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
		defer problems.Close()
		log.SetOutput(problems)

		//check if database connection is open
		if db.Ping() != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Error 19"))
			return
		}
		var match string
		//checking if email and username entered matches what is in database
		err2 := db.QueryRow("SELECT email FROM userinfo WHERE username=?", userName).Scan(&match)
		if err2 != nil {

			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong email/username combination."))
			log.Println("account.go processForgot 1 ", err2)
			return
		}
		if match != email {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong email/username combination."))
			browser := r.UserAgent()
			log.Printf("FAILED SEND PASSWORD RESET TO EMAIL Host: %s  Method: %s Location: %s Agent: %s\n", r.RemoteAddr, r.Method, r.URL.Path, browser)
			return
		}

		token := RandomString()
		//check for duplicate entry in forgot table
		var found string
		_ = db.QueryRow("SELECT token FROM forgot WHERE username=?", userName).Scan(&found)

		if found != "" {
			w.Write([]byte("<img src='img/ajax/available.png' /> Activation token resent to your email."))
			SendForgot(email, found)
			return
		}

		//preparing token activation
		stmt, err := db.Prepare("INSERT forgot SET username=?, token=?, expire=?")
		if err != nil {

			w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Report to admin Error 20"))
			log.Println("account.go processForgot 2 ", err)
			return
		}
		date := time.Now()
		res, err := stmt.Exec(userName, token, date)
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Report to admin Error 21"))
			log.Println("account.go processForgot 3 ", err)
			return
		}
		stmt.Close()
		affect, err := res.RowsAffected()
		if err != nil {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Report to admin Error 22"))
			log.Println("account.go processForgot 4 ", err)
			return
		}

		log.Printf("%d rows were affected by the the token activation check.\n", affect)

		//sends pasword reset information to email of user
		SendForgot(email, token)

		w.Write([]byte("<img src='img/ajax/available.png' /> Your password reset information has been sent your email."))

	}
}
