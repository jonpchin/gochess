package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dchest/captcha"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jonpchin/GoChess/gostuff"
	"golang.org/x/net/websocket"
)

var (
	domain = "localhost"
)

func main() {

	http.HandleFunc("/", mainPage)
	http.HandleFunc("/memberHome", memberHome)
	http.HandleFunc("/login", login)
	http.HandleFunc("/server/lobby", lobby)
	http.HandleFunc("/chess/memberChess", memberChess)
	http.HandleFunc("/database", database)
	http.HandleFunc("/profile", playerProfile)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/help", help)
	http.HandleFunc("/screenshots", screenshots)
	http.HandleFunc("/activate", activate)
	http.HandleFunc("/register", register)
	http.HandleFunc("/forgot", forgot)
	http.HandleFunc("/processForgot", gostuff.ProcessForgot)
	http.HandleFunc("/resetpass", resetpass)
	http.HandleFunc("/processResetPass", gostuff.ProcessResetPass)
	http.HandleFunc("/processRegister", gostuff.ProcessRegister)
	http.HandleFunc("/processLogin", gostuff.ProcessLogin)
	http.HandleFunc("/processActivate", gostuff.ProcessActivate)
	http.HandleFunc("/settings", settings)
	http.HandleFunc("/robots.txt", robot)
	http.HandleFunc("/saved", saved)
	http.HandleFunc("/highscores", highscores)
	http.HandleFunc("/server/getPlayerData", gostuff.GetPlayerData)

	http.HandleFunc("/updateCaptcha", gostuff.UpdateCaptcha)
	http.HandleFunc("/checkname", gostuff.CheckUserName)
	http.HandleFunc("/resumeGame", gostuff.ResumeGame)
	http.HandleFunc("/fetchgameID", gostuff.FetchGameByID)
	http.HandleFunc("/fetchgameByECO", gostuff.FetchGameByECO)

	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))

	http.Handle("/css/", http.FileServer(http.Dir("")))
	http.Handle("/img/", http.FileServer(http.Dir("")))
	http.Handle("/js/", http.FileServer(http.Dir("")))
	http.Handle("/data/", http.FileServer(http.Dir("")))
	http.Handle("/sound/", http.FileServer(http.Dir("")))

	http.Handle("/server", websocket.Handler(gostuff.EnterLobby))
	http.Handle("/chess", websocket.Handler(gostuff.EnterChess))

	var certPath = "secret/device.crt"
	var keyPath = "secret/device.key"

	//parse console arguments to determine OS environment to use localhost or goplaychess.com
	//default is localhost if no argument is passed
	if len(os.Args) > 1 {
		domain = os.Args[1]
		certPath = "secret/device.crt"
		keyPath = "secret/device.key"
	}

	go func() {

		//setting up database, the directory location of database backups is passed in
		proceed := gostuff.DbSetup("./backup")

		//setting up cron job
		gostuff.StartCron()
		//removes games older then 30 days from database
		if proceed == true {
			// Number of days used to remove old games, forgot and activate tokens
			days := "30"
			gostuff.RemoveOldGames(days)
			gostuff.RemoveOldActivate(days)
			gostuff.RemoveOldForgot(days)
			//fetch high score data from database
			gostuff.UpdateHighScore()
			gostuff.UpdateTotalGrandmasterGames()
			//gostuff.ResizeImages()
			//pass := gostuff.VerifyGrandmasterGames(100000)
			//if pass == true {
			//	fmt.Println("All games are accurate!")
			//}
			//gostuff.ExportDatabase()
		}

		//gostuff.SpawnProcess()

		//setting up clean up function for graceful shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
			<-c
			gostuff.Cleanup()
			os.Exit(1)
		}()
	}()

	go func() {
		if err := http.ListenAndServeTLS(":443", certPath, keyPath, nil); err != nil {
			fmt.Printf("ListenAndServeTLS error: %v\n", err)
		}
	}()
	//	gostuff.ConvertAllPGN()
	fmt.Println("Web server is now running.")

	if err := http.ListenAndServe(":80", http.HandlerFunc(redir)); err != nil {
		fmt.Printf("ListenAndServe error: %v\n", err)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		http.ServeFile(w, r, "404.html")
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=432000")
	http.ServeFile(w, r, "index.html")
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	var login = template.Must(template.ParseFiles("login.html"))

	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	if err := login.Execute(w, &d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func help(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=432000")
	http.ServeFile(w, r, "help.html")
}

func screenshots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=432000")
	http.ServeFile(w, r, "screenshots.html")
}
func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	var register = template.Must(template.ParseFiles("register.html"))

	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	if err := register.Execute(w, &d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func activate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	var activate = template.Must(template.ParseFiles("activate.html"))

	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	if err := activate.Execute(w, &d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func forgot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	var formTemplate = template.Must(template.ParseFiles("forgot.html"))

	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	if err := formTemplate.Execute(w, &d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func resetpass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	var formTemplate = template.Must(template.ParseFiles("resetpass.html"))

	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	if err := formTemplate.Execute(w, &d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func lobby(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err == nil {

		sessionID, err := r.Cookie("sessionID")
		if err == nil {

			if gostuff.SessionManager[username.Value] == sessionID.Value {

				w.Header().Set("Cache-Control", "private, max-age=86400")
				var lobby = template.Must(template.ParseFiles("lobby.html"))
				var bulletRating, blitzRating, standardRating int16
				var errMessage string

				errMessage, bulletRating, blitzRating, standardRating = gostuff.GetRating(username.Value)
				if errMessage != "" {
					fmt.Println("Problem fetching rating lobby main.go")
				}
				p := gostuff.Person{User: username.Value, Bullet: bulletRating, Blitz: blitzRating, Standard: standardRating}

				if err := lobby.Execute(w, &p); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}
	w.WriteHeader(404)
	http.ServeFile(w, r, "404.html")
}

func memberChess(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err == nil {

		sessionID, err := r.Cookie("sessionID")
		if err == nil {

			if gostuff.SessionManager[username.Value] == sessionID.Value {

				w.Header().Set("Cache-Control", "private, max-age=86400")
				var memberChess = template.Must(template.ParseFiles("memberchess.html"))
				p := gostuff.Person{User: username.Value}

				if err := memberChess.Execute(w, &p); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}
	w.WriteHeader(404)
	http.ServeFile(w, r, "404.html")
}

func memberHome(w http.ResponseWriter, r *http.Request) {

	username, err := r.Cookie("username")
	if err == nil {
		sessionID, err := r.Cookie("sessionID")
		if err == nil {

			if gostuff.SessionManager[username.Value] == sessionID.Value {
				w.Header().Set("Cache-Control", "private, max-age=432000")
				var memberHome = template.Must(template.ParseFiles("memberHome.html"))
				p := gostuff.Person{User: username.Value}

				if err := memberHome.Execute(w, &p); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}
	w.WriteHeader(404)
	http.ServeFile(w, r, "404.html")
}

func database(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err == nil {
		sessionID, err := r.Cookie("sessionID")
		if err == nil {

			if gostuff.SessionManager[username.Value] == sessionID.Value {
				w.Header().Set("Cache-Control", "private, max-age=432000")
				var memberHome = template.Must(template.ParseFiles("database.html"))
				p := gostuff.Person{User: username.Value}

				if err := memberHome.Execute(w, &p); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}
	w.WriteHeader(404)
	http.ServeFile(w, r, "404.html")
}

func playerProfile(w http.ResponseWriter, r *http.Request) {

	username, err := r.Cookie("username")
	if err == nil {

		sessionID, err := r.Cookie("sessionID")
		if err == nil {

			if gostuff.SessionManager[username.Value] == sessionID.Value {

				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				var all []gostuff.GoGame
				var ratErr string
				var bulletRating, blitzRating, standardRating, bulletRD, blitzRD, standardRD float64

				name := r.URL.Query().Get("name")
				var inputName string                 //used to pass to template to specify what profile name is being viewed
				if r.URL.Query().Get("name") == "" { //then look at own profile
					_, bulletRating, blitzRating, standardRating, bulletRD, blitzRD, standardRD = gostuff.GetRatingAndRD(username.Value)
					all = gostuff.GetGames(username.Value)
					inputName = username.Value
				} else { //otherwise look at specified player's profile
					ratErr, bulletRating, blitzRating, standardRating, bulletRD, blitzRD, standardRD = gostuff.GetRatingAndRD(name)
					if ratErr != "" { //this means someone typed a profile url which no player exists in database
						http.ServeFile(w, r, "nouser.html")
						return
					}
					all = gostuff.GetGames(name)
					inputName = name
				}

				var playerProfile = template.Must(template.ParseFiles("profile.html"))

				//rounding floats
				bulletN := gostuff.Round(bulletRating)
				blitzN := gostuff.Round(blitzRating)
				standardN := gostuff.Round(standardRating)
				bulletR := gostuff.RoundPlus(bulletRD, 2)
				blitzR := gostuff.RoundPlus(blitzRD, 2)
				standardR := gostuff.RoundPlus(standardRD, 2)
				gameID, exist := gostuff.GetGameID(inputName)
				opponent := ""

				// if a player is not playing a game use -1 for the gameID
				if exist == false {
					gameID = -1
				} else {
					opponent = gostuff.PrivateChat[inputName]
				}

				p := gostuff.ProfileGames{User: inputName, Bullet: bulletN, Blitz: blitzN, Standard: standardN,
					BulletRD: bulletR, BlitzRD: blitzR, StandardRD: standardR,
					Games: all, GameID: gameID, Opponent: opponent}

				if err := playerProfile.Execute(w, &p); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}
	w.WriteHeader(404)
	http.ServeFile(w, r, "404.html")
}

//logs the user out by deleting the cookies and back end session and redirecting them to the homepage
func logout(w http.ResponseWriter, r *http.Request) {

	username, err := r.Cookie("username")
	if err == nil && len(username.Value) >= 3 && len(username.Value) <= 12 {

		sessionID, err := r.Cookie("sessionID")
		if err == nil {

			if gostuff.SessionManager[username.Value] == sessionID.Value {

				delete(gostuff.SessionManager, username.Value)
				cookie := http.Cookie{Name: "username", Value: "0", MaxAge: -1}
				http.SetCookie(w, &cookie)
				cookie = http.Cookie{Name: "sessionID", Value: "0", MaxAge: -1}
				http.SetCookie(w, &cookie)
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				http.ServeFile(w, r, "index.html")
				return
			}
		}
	}
	w.WriteHeader(404)
	http.ServeFile(w, r, "404.html")
}

func settings(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err == nil {
		sessionID, err := r.Cookie("sessionID")
		if err == nil {

			if gostuff.SessionManager[username.Value] == sessionID.Value {
				w.Header().Set("Cache-Control", "private, max-age=86400")
				http.ServeFile(w, r, "settings.html")
				return
			}
		}
	}
	w.WriteHeader(404)
	http.ServeFile(w, r, "404.html")
}

func highscores(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err == nil {
		sessionID, err := r.Cookie("sessionID")
		if err == nil {
			if gostuff.SessionManager[username.Value] == sessionID.Value {

				w.Header().Set("Cache-Control", "private, max-age=3600")
				var highscores = template.Must(template.ParseFiles("highscore.html"))

				var p gostuff.ScoreBoard
				p = gostuff.LeaderBoard.Scores

				if err := highscores.Execute(w, &p); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}
	w.WriteHeader(404)
	http.ServeFile(w, r, "404.html")
}

func saved(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err == nil {
		sessionID, err := r.Cookie("sessionID")
		if err == nil {

			if gostuff.SessionManager[username.Value] == sessionID.Value {

				w.Header().Set("Cache-Control", "private, max-age=86400")
				var all []gostuff.GoGame

				name := r.URL.Query().Get("user")
				all = gostuff.GetSaved(name)

				var saved = template.Must(template.ParseFiles("saved.html"))

				p := gostuff.ProfileGames{User: name, Games: all}

				if err := saved.Execute(w, &p); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}
	w.WriteHeader(404)
	http.ServeFile(w, r, "404.html")
}

func robot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "robots.txt")
}

func redir(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+domain+r.RequestURI, http.StatusMovedPermanently)
}
