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
	"github.com/jonpchin/gochess/goforum"
	"github.com/jonpchin/gochess/gostuff"

	"golang.org/x/net/websocket"
)

const (
	days = "180" // Number of days used to remove old games, forgot and activate tokens
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

type neuteredReaddirFile struct {
	http.File
}

type fullPath struct {
	Url string
}

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
	http.HandleFunc("/engine", engine)
	http.HandleFunc("/news", news)
	http.HandleFunc("/runtest", runJsTests)
	http.HandleFunc("/forum", forum)
	http.HandleFunc("/createthread", createThread)
	http.HandleFunc("/sendFirstForumPost", goforum.SendFirstForumPost)
	http.HandleFunc("/server/getPlayerData", gostuff.GetPlayerData)

	http.HandleFunc("/updateCaptcha", gostuff.UpdateCaptcha)
	http.HandleFunc("/checkname", gostuff.CheckUserName)
	http.HandleFunc("/resumeGame", gostuff.ResumeGame)
	http.HandleFunc("/fetchgameID", gostuff.FetchGameByID)
	http.HandleFunc("/fetchgameByECO", gostuff.FetchGameByECO)
	http.HandleFunc("/fetchBulletHistory", gostuff.FetchBulletHistory)
	http.HandleFunc("/fetchBlitzHistory", gostuff.FetchBlitzHistory)
	http.HandleFunc("/fetchStandardHistory", gostuff.FetchStandardHistory)
	http.HandleFunc("/fetchCorrespondenceHistory", gostuff.FetchCorrespondenceHistory)
	http.HandleFunc("/checkInGame", gostuff.CheckInGame)

	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))

	// prevent directory listing
	currentDir := justFilesFilesystem{http.Dir("")}

	http.Handle("/css/", cacheControl(http.FileServer(currentDir), "259200"))
	http.Handle("/img/", http.FileServer(currentDir))
	http.Handle("/js/", cacheControl(http.FileServer(currentDir), "86400"))
	http.Handle("/third-party-js/", cacheControl(http.FileServer(currentDir), "432000"))
	http.Handle("/data/", http.FileServer(currentDir))
	http.Handle("/sound/", http.FileServer(currentDir))

	http.Handle("/server", websocket.Handler(gostuff.EnterLobby))
	http.Handle("/chess", websocket.Handler(gostuff.EnterChess))

	var certPath = "secret/device.crt"
	var keyPath = "secret/device.key"

	//parse console arguments to determine OS environment to use localhost or goplaychess.com
	//default is localhost if no argument is passed
	if len(os.Args) > 1 {
		certPath = "secret/combine2017.crt"
		keyPath = "secret/go2017.key"
	} else if gostuff.IsEnvironmentTravis() {
		certPath = "_travis/data/device.crt"
		keyPath = "_travis/data/device.key"
	}

	//gostuff.OneTimeParseTemplates()

	go func() {
		//setting up database, the directory location of database backups is passed in
		proceed := gostuff.DbSetup("./backup")

		//removes games older then 180 days from database
		// only proceed if not in App Veyor or Travis environments
		// this is only temporarily until all the tables are imported
		if proceed == true && gostuff.IsEnvironmentTravis() == false &&
			gostuff.IsEnvironmentAppVeyor() == false {

			//setting up cron job
			gostuff.StartCron()

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
			// pass in true to export template(No grandmaster) without data in the tables

			gostuff.ExportDatabase(false)
			//gostuff.CompressDatabase()
			goforum.ConnectForumDb()
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
		//gostuff.CheckNullInTable("rating")
	}()
	//gostuff.FetchNewsSources()
	//gostuff.ReadAllNews()
	//gostuff.UpdateNewsFromConfig()
	//gostuff.FakeDataForTravis()

	//go func() {
	//	engine := gostuff.StartEngine(nil)
	//	engine.Quit()
	//}()

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
		var doesNotExist = template.Must(template.ParseFiles("404.html"))

		p := fullPath{Url: r.Host}

		if err := doesNotExist.Execute(w, &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.ServeFile(w, r, "index.html")
	}
}

func login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	d := struct {
		CaptchaId string
		PageTitle string
	}{
		captcha.New(),
		"Login",
	}

	gostuff.ParseTemplates(d, w, "login.html", []string{"templates/loginTemplate.html",
		"templates/guestHeader.html"}...)
}

func help(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=432000")
	http.ServeFile(w, r, "help.html")
}

func news(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeFile(w, r, "news.html")
}

func screenshots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=432000")
	http.ServeFile(w, r, "screenshots.html")
}
func register(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	d := struct {
		CaptchaId string
		PageTitle string
	}{
		captcha.New(),
		"Register",
	}
	gostuff.ParseTemplates(d, w, "register.html", []string{"templates/registerTemplate.html",
		"templates/guestHeader.html"}...)
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

	if isAuthorized(w, r) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		var lobby = template.Must(template.ParseFiles("lobby.html"))
		var bulletRating, blitzRating, standardRating, correspondenceRating int16
		var errMessage string
		username, _ := r.Cookie("username")

		errMessage, bulletRating, blitzRating, standardRating,
			correspondenceRating = gostuff.GetRating(username.Value)

		if errMessage != "" {
			fmt.Println("Problem fetching rating lobby main.go")
		}

		p := gostuff.Person{User: username.Value, Bullet: bulletRating, Blitz: blitzRating,
			Standard: standardRating, Correspondence: correspondenceRating}

		if err := lobby.Execute(w, &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func memberChess(w http.ResponseWriter, r *http.Request) {

	if isAuthorized(w, r) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		var memberChess = template.Must(template.ParseFiles("memberchess.html"))
		username, _ := r.Cookie("username")
		p := gostuff.Person{User: username.Value}

		if err := memberChess.Execute(w, &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func memberHome(w http.ResponseWriter, r *http.Request) {

	if isAuthorized(w, r) {
		w.Header().Set("Cache-Control", "private, max-age=432000")

		username, _ := r.Cookie("username")
		p := struct {
			User      string
			PageTitle string // Title of the web page
		}{
			username.Value,
			"Welcome",
		}

		gostuff.ParseTemplates(p, w, "memberHome.html", []string{"templates/memberHomeTemplate.html",
			"templates/memberHeader.html"}...)
	}
}

func database(w http.ResponseWriter, r *http.Request) {

	if isAuthorized(w, r) {
		w.Header().Set("Cache-Control", "private, max-age=432000")
		var memberHome = template.Must(template.ParseFiles("database.html"))
		username, _ := r.Cookie("username")

		p := gostuff.Person{User: username.Value}

		if err := memberHome.Execute(w, &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func playerProfile(w http.ResponseWriter, r *http.Request) {

	if isAuthorized(w, r) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		var all []gostuff.GoGame
		var ratErr string
		var bulletRating, blitzRating, standardRating, correspondenceRating,
			bulletRD, blitzRD, standardRD, correspondenceRD float64

		name := r.URL.Query().Get("name")
		username, _ := r.Cookie("username")

		var inputName string                 //used to pass to template to specify what profile name is being viewed
		if r.URL.Query().Get("name") == "" { //then look at own profile
			_, bulletRating, blitzRating, standardRating, correspondenceRating,
				bulletRD, blitzRD, standardRD, correspondenceRD = gostuff.GetRatingAndRD(username.Value)
			all = gostuff.GetGames(username.Value)
			inputName = username.Value
		} else { //otherwise look at specified player's profile
			ratErr, bulletRating, blitzRating, standardRating, correspondenceRating,
				bulletRD, blitzRD, standardRD, correspondenceRD = gostuff.GetRatingAndRD(name)
			if ratErr != "" { //this means someone typed a profile url which no player exists in database
				http.ServeFile(w, r, "nouser.html")
				return
			}
			all = gostuff.GetGames(name)
			inputName = name
		}

		country := gostuff.GetCountry(name)

		var playerProfile = template.Must(template.ParseFiles("profile.html"))

		//rounding floats
		bulletN := gostuff.Round(bulletRating)
		blitzN := gostuff.Round(blitzRating)
		standardN := gostuff.Round(standardRating)
		correspondenceN := gostuff.Round(correspondenceRating)
		bulletR := gostuff.RoundPlus(bulletRD, 2)
		blitzR := gostuff.RoundPlus(blitzRD, 2)
		standardR := gostuff.RoundPlus(standardRD, 2)
		correspondenceR := gostuff.RoundPlus(correspondenceRD, 2)
		gameID, exist := gostuff.GetGameID(inputName)
		opponent := ""

		// if a player is not playing a game use -1 for the gameID
		if exist == false {
			gameID = -1
		} else {
			opponent = gostuff.PrivateChat[inputName]
		}

		p := gostuff.ProfileGames{User: inputName, Bullet: bulletN, Blitz: blitzN, Standard: standardN,
			Correspondence: correspondenceN, BulletRD: bulletR, BlitzRD: blitzR, StandardRD: standardR,
			CorrespondenceRD: correspondenceR, Games: all, GameID: gameID, Opponent: opponent,
			Days: days, Country: country}

		if err := playerProfile.Execute(w, &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

//logs the user out by deleting the cookies and back end session and redirecting them to the homepage
func logout(w http.ResponseWriter, r *http.Request) {

	if isAuthorized(w, r) {
		username, _ := r.Cookie("username")
		delete(gostuff.SessionManager, username.Value)
		cookie := http.Cookie{Name: "username", Value: "0", MaxAge: -1}
		http.SetCookie(w, &cookie)
		cookie = http.Cookie{Name: "sessionID", Value: "0", MaxAge: -1}
		http.SetCookie(w, &cookie)
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.ServeFile(w, r, "index.html")
	}
}

func settings(w http.ResponseWriter, r *http.Request) {

	if isAuthorized(w, r) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.ServeFile(w, r, "settings.html")
	}
}

func highscores(w http.ResponseWriter, r *http.Request) {

	if isAuthorized(w, r) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		var highscores = template.Must(template.ParseFiles("highscore.html"))

		var p gostuff.ScoreBoard
		p = gostuff.LeaderBoard.Scores

		if err := highscores.Execute(w, &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func engine(w http.ResponseWriter, r *http.Request) {

	if isAuthorized(w, r) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		var engine = template.Must(template.ParseFiles("engine.html"))
		username, _ := r.Cookie("username")

		p := gostuff.Person{User: username.Value}

		if err := engine.Execute(w, &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func saved(w http.ResponseWriter, r *http.Request) {

	if isAuthorized(w, r) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		var all []gostuff.GoGame

		name := r.URL.Query().Get("user")
		all = gostuff.GetSaved(name)

		var saved = template.Must(template.ParseFiles("saved.html"))

		p := gostuff.ProfileGames{User: name, Games: all, Days: days}

		if err := saved.Execute(w, &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func runJsTests(w http.ResponseWriter, r *http.Request) {
	if isAuthorized(w, r) {
		http.ServeFile(w, r, "runJsTests.html")
	}
}

func forum(w http.ResponseWriter, r *http.Request) {

	// goplaychess.com/forum?forumId=2
	forumId := r.URL.Query().Get("forumid")
	threadId := r.URL.Query().Get("threadid")

	var output string
	var templatePath string
	var p interface{}

	if forumId == "" && threadId == "" { // show main forum
		p = struct {
			PageTitle string
			Forums    []goforum.Forum
		}{
			"Forums",
			goforum.GetForums(),
		}

		output = "forum.html"
		templatePath = "templates/forumTemplate.html"

	} else if threadId == "" { //  show all threads in a section

		p = struct {
			PageTitle string
			Threads   goforum.ThreadSection
		}{
			"Threads",
			goforum.GetThreads(forumId),
		}

		output = "threads.html"
		templatePath = "templates/threadsTemplate.html"

	} else { // show all posts in a thread

		p = struct {
			PageTitle string
			Posts     []goforum.Post
		}{
			goforum.GetForumTitle(forumId),
			goforum.GetPosts(threadId),
		}

		output = "posts.html"
		templatePath = "templates/postsTemplate.html"
	}

	gostuff.ParseTemplates(p, w, output, []string{templatePath,
		"templates/guestHeader.html"}...)
}

func createThread(w http.ResponseWriter, r *http.Request) {

	if isAuthorized(w, r) {
		username, _ := r.Cookie("username")
		p := struct {
			User      string
			ForumName string
			PageTitle string
		}{
			username.Value,
			r.URL.Query().Get("forumname"),
			"Create Thread",
		}
		gostuff.ParseTemplates(p, w, "createthread.html", []string{"templates/createthreadTemplate.html",
			"templates/guestHeader.html"}...)
	}
}

func robot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "robots.txt")
}

func redir(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

// returns false if a user credentials are invalid
func isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	username, err := r.Cookie("username")
	if err == nil {
		sessionID, err := r.Cookie("sessionID")
		if err == nil {
			if gostuff.SessionManager[username.Value] == sessionID.Value {
				return true
			}
		}
	}
	w.WriteHeader(404)
	http.ServeFile(w, r, "404.html")
	return false
}

// used to cache static assets for specified seconds passed in function parameter
func cacheControl(h http.Handler, seconds string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// max age is the number of seconds to cache
		w.Header().Set("Cache-Control", "private, max-age="+seconds)
		h.ServeHTTP(w, r)
	}
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}
