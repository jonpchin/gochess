package gostuff

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dchest/captcha"
	"golang.org/x/net/websocket"
)

func UpdateCaptcha(w http.ResponseWriter, r *http.Request) {
	cap := captcha.New()
	w.Write([]byte(cap))
}

//displays player data when mouse hovers over
func GetPlayerData(w http.ResponseWriter, r *http.Request) {

	valid := ValidateCredentials(w, r)
	if valid == false {
		return
	}

	// the name of the player being looked up by the AJAX call
	lookupName := template.HTMLEscapeString(r.FormValue("user"))

	//checking if the player is a game
	status := ""
	icon := ""
	url := ""
	endUrl := "" //closing the href link

	countryFlag := GetCountry(lookupName)
	enemyFlag := GetCountry(PrivateChat[lookupName])
	if countryFlag == "" {
		countryFlag = "globe"
	}
	if enemyFlag == "" {
		enemyFlag = "globe"
	}

	device := ""

	if DeviceManager[lookupName] == "android" {
		device = "<img src='../img/icons/android.ico'>"
	}

	//second username is nil as it only checks one name
	if checkTable(lookupName) {
		status = "vs. " + PrivateChat[lookupName] + "<src='img/flags/'" +
			enemyFlag + ".png>"
		icon = "<img src='../img/icons/playing.png' alt='status'>"
		id, _ := GetGameID(lookupName)
		url = "<a href=/chess/memberChess?spectate&id=" + strconv.Itoa(id) + ">"
		endUrl = "</a>"
	}

	if strings.Contains(lookupName, "guest") {
		var result = device + icon + url + lookupName + " " + status + endUrl +
			"<br><img src='../img/icons/bullet.png' alt='bullet'>1500" +
			"<img src='../img/icons/blitz.png' alt='blitz'>1500" +
			"<img src='../img/icons/standard.png' alt='standard'>1500"
		w.Write([]byte(result))
		return
	}

	//getting player rating
	ratingError, bulletRating, blitzRating, standardRating,
		_ := GetRating(lookupName)

	if ratingError != "" {
		w.Write([]byte("Could not find player, " + lookupName))
		return
	}

	bullet := fmt.Sprintf("%d", bulletRating)
	blitz := fmt.Sprintf("%d", blitzRating)
	standard := fmt.Sprintf("%d", standardRating)
	//correspondence := fmt.Sprintf("%d", correspondenceRating)

	var result = device + icon + url + lookupName + "<img src='../img/flags/" + countryFlag +
		".png'>" + " " + status + endUrl +
		"<br><img src='../img/icons/bullet.png' alt='bullet'>" + bullet +
		"<img src='../img/icons/blitz.png' alt='blitz'>" + blitz +
		"<img src='../img/icons/standard.png' alt='standard'>" + standard

	w.Write([]byte(result))
}

func ResumeGame(w http.ResponseWriter, r *http.Request) {

	valid := ValidateCredentials(w, r)
	if valid == false {
		return
	}

	id := template.HTMLEscapeString(r.FormValue("id"))
	white := template.HTMLEscapeString(r.FormValue("white"))
	black := template.HTMLEscapeString(r.FormValue("black"))

	user, _ := r.Cookie("username")

	var chat ChatInfo
	chat.Type = "chess_game"
	var success bool
	var game ChessGame
	if user.Value == white {
		if isPlayerInLobby(black) && !checkTable(black) {
			success = game.fetchSavedGame(id, user.Value)
			if success == false {
				w.Write([]byte("false"))
				return
			}
			if err := websocket.JSON.Send(Chat.Lobby[black], &chat); err != nil {
				fmt.Println("error ajax.go ResumeGame 1 is ", err)
			}
			w.Write([]byte("true"))
			return
		}

	} else if user.Value == black {
		if isPlayerInLobby(white) && !checkTable(white) {
			success = game.fetchSavedGame(id, user.Value)
			if success == false {
				w.Write([]byte("false"))
				return
			}
			if err := websocket.JSON.Send(Chat.Lobby[white], &chat); err != nil {
				fmt.Println("error ajax.go ResumeGame 3 is ", err)
			}
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false"))
}

// fetches all data of a chess game by the ID
func FetchGameByID(w http.ResponseWriter, r *http.Request) {

	valid := ValidateCredentials(w, r)
	if valid == false {
		return
	}

	// a player shouldn't be using the database if they are in a game playing another person
	user, _ := r.Cookie("username")
	if checkTable(user.Value) {
		w.Write([]byte("Database use is not allowed when you are playing a game against a real person!"))
		return
	}

	id := template.HTMLEscapeString(r.FormValue("gameID"))

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	num, err := strconv.Atoi(id)
	if err != nil {
		w.Write([]byte("Not a valid number. Please enter only digits."))
		return
	}
	if num > TotalGrandmasterGames-1 || num <= 0 {
		w.Write([]byte("Please search a game ID between 1 and " + strconv.Itoa(TotalGrandmasterGames-1)))
		return
	}

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		w.Write([]byte("The database is offline"))
		return
	}

	var all GrandMasterGame

	err = db.QueryRow("SELECT * FROM grandmaster WHERE id=?", id).Scan(&all.ID, &all.Event, &all.Site,
		&all.Date, &all.Round, &all.White,
		&all.Black, &all.Result, &all.WhiteElo, &all.BlackElo,
		&all.ECO, &all.Moves, &all.EventDate)

	if err != nil {
		log.Println(err)
		w.Write([]byte("Error in processing request"))
		return
	}

	allGames, err := json.Marshal(all)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Unable to serialize data"))
		return
	}
	w.Write([]byte(string(allGames)))
}

func FetchGameByECO(w http.ResponseWriter, r *http.Request) {
	valid := ValidateCredentials(w, r)
	if valid == false {
		return
	}

	user, _ := r.Cookie("username")
	if checkTable(user.Value) {
		w.Write([]byte("Database use is not allowed when you are playing a game against a real person!"))
		return
	}

	eco := template.HTMLEscapeString(r.FormValue("ECO"))
	ecoIndex := template.HTMLEscapeString(r.FormValue("ECOIndex"))

	// for now do not allow players to scroll past 100th game in a specific opening
	ecoValue, err := strconv.Atoi(ecoIndex)
	if err != nil {
		w.Write([]byte("An invalid game index was searched"))
		log.Println("Invalid ecoIndex")
	}
	if ecoValue < 0 || ecoValue > 100 {
		log.Println("ecoIndex is out of bounds")
		w.Write([]byte("Search index out of bounds"))
		return
	}

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		w.Write([]byte("The database is offline"))
		return
	}

	var all GrandMasterGame

	err = db.QueryRow("select * from grandmaster where `ECO`=? ORDER BY ID LIMIT ?, 1", eco, ecoValue).Scan(&all.ID, &all.Event, &all.Site,
		&all.Date, &all.Round, &all.White,
		&all.Black, &all.Result, &all.WhiteElo, &all.BlackElo,
		&all.ECO, &all.Moves, &all.EventDate)

	if err != nil {
		log.Println(err)
		w.Write([]byte("Error in processing request"))
		return
	}

	game, err := json.Marshal(all)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Unable to serialize data"))
		return
	}
	w.Write([]byte(string(game)))
}

func CheckUserName(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := template.HTMLEscapeString(r.FormValue("username"))

		//making sure username fits length requirement
		if len(username) < 3 || len(username) > 12 {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Choose a name between 3 to 12 characters."))
			return
		}
		ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)

		if db.Ping() != nil {
			fmt.Printf("ERROR 2 PINGING DB IP: %s \n", ipAddress)
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Please come back later."))
			return
		}

		//check if database connection is open
		found := CheckUserNameInDb(username)
		if found {
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Username taken"))
		} else {
			w.Write([]byte(" <img src='img/ajax/available.png' /> Username available"))
		}
	}
}

// fetches players rating bullet history from database
func FetchBulletHistory(w http.ResponseWriter, r *http.Request) {
	valid := ValidateCredentials(w, r)
	if valid == false {
		return
	}
	user := template.HTMLEscapeString(r.FormValue("user"))
	bullet, isHistory, _ := GetRatingHistory(user, "bullet")
	if isHistory {
		w.Write([]byte(bullet))
	} else {
		w.Write([]byte("")) // blank string will be checked if history was succesfully fetched
	}
}

// fetches players blitz history rating from database
func FetchBlitzHistory(w http.ResponseWriter, r *http.Request) {
	valid := ValidateCredentials(w, r)
	if valid == false {
		return
	}
	user := template.HTMLEscapeString(r.FormValue("user"))
	blitz, isHistory, _ := GetRatingHistory(user, "blitz")
	if isHistory {
		w.Write([]byte(blitz))
	} else {
		w.Write([]byte("")) // blank string will be checked if history was succesfully fetched
	}
}

// fetches players standard rating history from database
func FetchStandardHistory(w http.ResponseWriter, r *http.Request) {
	valid := ValidateCredentials(w, r)
	if valid == false {
		return
	}
	user := template.HTMLEscapeString(r.FormValue("user"))
	standard, isHistory, _ := GetRatingHistory(user, "standard")
	if isHistory {
		w.Write([]byte(standard))
	} else {
		w.Write([]byte("")) // blank string will be checked if history was succesfully fetched
	}
}

// fetches players rating blitz history from database
func FetchCorrespondenceHistory(w http.ResponseWriter, r *http.Request) {
	valid := ValidateCredentials(w, r)
	if valid == false {
		return
	}
	user := template.HTMLEscapeString(r.FormValue("user"))
	correspondence, isHistory, _ := GetRatingHistory(user, "correspondence")
	if isHistory {
		w.Write([]byte(correspondence))
	} else {
		w.Write([]byte("")) // blank string will be checked if history was succesfully fetched
	}
}

// returns true if user is an registered user that is logged in
func ValidateCredentials(w http.ResponseWriter, r *http.Request) bool {
	username, err := r.Cookie("username")
	if err != nil || len(username.Value) < 3 || len(username.Value) > 12 {
		w.Write([]byte("Could not authenticate user"))
		return false
	}

	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		w.Write([]byte("Could not authenticate user"))
		return false
	}

	if sessionID.Value == "" || SessionManager[username.Value] != sessionID.Value {
		w.Write([]byte("Could not authenticate user"))
		return false
	}
	return true
}

// checks if a player is in a game
func CheckInGame(w http.ResponseWriter, r *http.Request) {
	valid := ValidateCredentials(w, r)
	if valid == false {
		return
	}
	user := template.HTMLEscapeString(r.FormValue("user"))
	if checkTable(user) {
		w.Write([]byte("inGame"))
	} else {
		w.Write([]byte("Safe"))
	}
}

// Show logs only to admins
func FetchLogs(w http.ResponseWriter, r *http.Request) {
	valid := ValidateCredentials(w, r)
	if valid == false {
		return
	}
	username, err := r.Cookie("username")
	if err != nil {
		fmt.Println(err)
		return
	}
	if IsAdmin(username.Value) == false {
		return
	}

	logType := template.HTMLEscapeString(r.FormValue("logType"))

	if logType == "chat" {
		data, err := ioutil.ReadFile("logs/chat.txt")
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(data)
		}

	} else if logType == "errors" {
		data, err := ioutil.ReadFile("logs/errors.txt")
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(data)
		}
	} else if logType == "main" {
		data, err := ioutil.ReadFile("nohup.out")
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(data)
		}
	} else {
		w.Write([]byte("Invalid log type showLogs"))
	}
}
