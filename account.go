package gostuff

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dchest/captcha"
	_ "github.com/go-sql-driver/mysql"
)

//checks database to see if username and token is correct and then updates the new password
//with the encrypted password
func ProcessResetPass(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		Show404Page(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong captcha solution! Please try again."))
	} else {

		var userInfo UserInfo
		userInfo.Username = template.HTMLEscapeString(r.FormValue("user"))
		userInfo.Token = template.HTMLEscapeString(r.FormValue("token"))
		userInfo.Password = template.HTMLEscapeString(r.FormValue("pass"))
		confirm := template.HTMLEscapeString(r.FormValue("confirm"))
		userInfo.IpAddress, _, _ = net.SplitHostPort(r.RemoteAddr)

		//check if password and confirm match
		if userInfo.Password != confirm {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> The password you entered does not match. Please try again."))
			return
		}
		if len(userInfo.Password) < 5 || len(userInfo.Password) > 32 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> The password you entered must be 5-32 characters long."))
			return
		}
		userInfo.resetPass(w, r)
	}
}

func (userInfo *UserInfo) resetPass(w http.ResponseWriter, r *http.Request) {

	problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later."))
		log.Println("DATABASE DOWN!")
		return
	}
	var tokenInDB string

	//checking if token matches the one entered by user
	err2 := db.QueryRow("SELECT token FROM forgot WHERE username=?", userInfo.Username).Scan(&tokenInDB)
	if err2 != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/token combination."))
		log.Println(err2)
		return
	}
	if tokenInDB != userInfo.Token {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong username/token combination."))
		log.Printf("FAILED PASSWORD RESET IP: %s  Method: %s Location: %s Agent: %s\n", userInfo.IpAddress,
			r.Method, r.URL.Path, r.UserAgent())
		return
	}

	//setting password for user and deleting a row from the forgot table
	stmt, err := db.Prepare("UPDATE userinfo SET password=? where username=?")
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 3."))
		log.Println(err)
		return
	}
	defer stmt.Close()

	//hashing password
	key, err := hashPass(userInfo.Username, userInfo.Password)
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 4."))
		log.Println(err)
		return
	}

	_, err = stmt.Exec(key, userInfo.Username)
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

	res, err := stmt.Exec(userInfo.Username)
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 7"))
		log.Println(err)
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 8."))
		log.Println(err)
		return
	}

	log.Printf("%d row was deleted from the activate table by user %s\n", affect, userInfo.Username)
	w.Write([]byte("<img src='img/ajax/available.png' /> Your password is now changed!"))
}

//activates user account and deletes entry from database
func ProcessActivate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Show404Page(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !captcha.VerifyString(template.HTMLEscapeString(r.FormValue("captchaId")), template.HTMLEscapeString(r.FormValue("captchaSolution"))) {
		w.Write([]byte("<img src='img/ajax/not-available.png'/> Wrong captcha solution"))

	} else {
		var userInfo UserInfo
		userInfo.Username = template.HTMLEscapeString(r.FormValue("user"))
		userInfo.Token = template.HTMLEscapeString(r.FormValue("token"))
		userInfo.IpAddress, _, _ = net.SplitHostPort(r.RemoteAddr)

		message, result := userInfo.Activate(r.Method, r.URL.Path, r.UserAgent())

		// There should always be a message returned
		w.Write([]byte(message))
		if result {
			//once a user activates his acccount update the highscore board so it shows him as new user
			UpdateHighScore()
		}
	}
}

// returns true if account is succesfully activated
// these parameters can be replaced with blank string when using it on a unit test
// as these parameters are only used for logging purposes
func (userInfo *UserInfo) Activate(method string, url string, agent string) (string, bool) {

	problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("Can't ping database")
		return "<img src='img/ajax/not-available.png'/> Database is down. ", false
	}
	var tokenInDB string

	//checking if token matches the one entered by user
	err2 := db.QueryRow("SELECT token FROM activate WHERE username=?", userInfo.Username).Scan(&tokenInDB)
	if err2 != nil || tokenInDB != userInfo.Token {
		log.Printf("FAILED ACTIVATION Host: %s  Method: %s Location: %s Agent: %s\n", userInfo.IpAddress, method, url, agent)
		return "<img src='img/ajax/not-available.png' /> Wrong username/token combination", false
	}
	//setting verify to yes and deleting row from activate table as well as captcha to zero to signfy user unlocked account
	stmt, err := db.Prepare("UPDATE userinfo SET verify=?, captcha=? where username=?")
	if err != nil {
		log.Println(err)
		return "<img src='img/ajax/not-available.png' /> Can't activate account.", false
	}
	defer stmt.Close()

	_, err = stmt.Exec("YES", 0, userInfo.Username)
	if err != nil {
		log.Println(err)
		return "<img src='img/ajax/not-available.png' /> Can't activate account 2", false
	}

	stmt, err = db.Prepare("DELETE FROM activate where username=?")
	if err != nil {
		log.Println(err)
		return "<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 15.", false
	}

	res, err := stmt.Exec(userInfo.Username)
	if err != nil {
		log.Println(err)
		return "<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 16.", false
	}

	affect, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return "<img src='img/ajax/not-available.png' /> Something is wrong with the server. Tell admin error 17.", false
	}
	log.Printf("%d row was deleted from the activate table by user %s\n", affect, userInfo.Username)

	//now user may login so we can redirect while token deletion proceeds in the background
	message := "<script>window.location = 'login?user=" + userInfo.Username + "';</script>"
	return message, true
}

//connects to database and adds token to the forgot table
func ProcessForgot(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Show404Page(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong captcha solution! Please try again."))

	} else {

		var userInfo UserInfo
		userInfo.Username = template.HTMLEscapeString(r.FormValue("user"))
		userInfo.Email = template.HTMLEscapeString(r.FormValue("email"))

		success := userInfo.forgot(w, r)

		if success {

			//sends pasword reset information to email of user
			go SendForgot(userInfo.Email, userInfo.Token, r.Host)
			w.Write([]byte("<img src='img/ajax/available.png' /> Your password reset information has been sent your email."))
		}
	}
}

// if true is returned then password reset information will be emailed to user
func (userInfo *UserInfo) forgot(w http.ResponseWriter, r *http.Request) bool {

	problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Please come back later. Error 19"))
		return false
	}
	var match string
	//checking if email and username entered matches what is in database
	err2 := db.QueryRow("SELECT email FROM userinfo WHERE username=?", userInfo.Username).Scan(&match)
	if err2 != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong email/username combination."))
		log.Println(err2)
		return false
	}
	if match != userInfo.Email {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> Wrong email/username combination."))
		log.Printf("FAILED SEND PASSWORD RESET TO EMAIL Host: %s  Method: %s Location: %s Agent: %s\n", userInfo.IpAddress, r.Method, r.URL.Path, r.UserAgent())
		return false
	}

	userInfo.Token = RandomString()
	//check for duplicate entry in forgot table
	var found string
	_ = db.QueryRow("SELECT token FROM forgot WHERE username=?", userInfo.Username).Scan(&found)

	if found != "" {
		w.Write([]byte("<img src='img/ajax/available.png' /> Activation token resent to your email."))
		go SendForgot(userInfo.Email, found, r.Host)
		return false
	}

	//preparing token activation
	stmt, err := db.Prepare("INSERT forgot SET username=?, token=?, expire=?")
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Report to admin Error 20"))
		log.Println(err)
		return false
	}
	defer stmt.Close()

	res, err := stmt.Exec(userInfo.Username, userInfo.Token, time.Now())
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Report to admin Error 21"))
		log.Println(err)
		return false
	}

	affect, err := res.RowsAffected()
	if err != nil {
		w.Write([]byte("<img src='img/ajax/not-available.png' /> We are having trouble with our server. Report to admin Error 22"))
		log.Println(err)
		return false
	}

	log.Printf("%d rows were affected by the the token activation check.\n", affect)
	return true
}
