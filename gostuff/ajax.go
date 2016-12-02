package gostuff

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/dchest/captcha"
	"golang.org/x/net/websocket"
)

func UpdateCaptcha(w http.ResponseWriter, r *http.Request) {
	cap := captcha.New()
	w.Write([]byte(cap))
}

//displays player data when mouse hovers over
func GetPlayerData(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err != nil || len(username.Value) < 3 || len(username.Value) > 12 {
		return
	}

	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		return
	}

	if SessionManager[username.Value] != sessionID.Value {
		return
	}

	// the name of the player being looked up by the AJAX call
	lookupName := template.HTMLEscapeString(r.FormValue("user"))

	//getting player rating
	ratingError, bulletRating, blitzRating, standardRating := GetRating(lookupName)
	if ratingError != "" {
		w.Write([]byte("Service is down."))
		return
	}

	bullet := fmt.Sprintf("%d", bulletRating)
	blitz := fmt.Sprintf("%d", blitzRating)
	standard := fmt.Sprintf("%d", standardRating)

	//checking if the player is a game
	status := ""
	icon := ""
	url := ""
	endUrl := "" //closing the href link
	countryFlag := getCountry(lookupName)
	enemyFlag := getCountry(PrivateChat[lookupName])
	if countryFlag == "" {
		countryFlag = "globe"
	}
	if enemyFlag == "" {
		enemyFlag = "globe"
	}

	//second username is nil as it only checks one name
	if isPlayerInGame(lookupName, "") {
		status = "vs. " + PrivateChat[lookupName] + "<src='img/flags/'" +
			enemyFlag + ".png>"
		icon = "<img src='../img/icons/playing.png' alt='status'>"
		id, _ := GetGameID(lookupName)
		url = "<a href=/chess/memberChess?spectate&id=" + strconv.Itoa(id) + ">"
		endUrl = "</a>"
	}

	var result = icon + url + lookupName + "<img src='../img/flags/" + countryFlag +
		".png'>" + " " + status + endUrl +
		"<br><img src='../img/icons/bullet.png' alt='bullet'>" + bullet +
		"<img src='../img/icons/blitz.png' alt='blitz'>" + blitz +
		"<img src='../img/icons/standard.png' alt='standard'>" + standard

	w.Write([]byte(result))
}

func ResumeGame(w http.ResponseWriter, r *http.Request) {
	user, err := r.Cookie("username")
	if err != nil || len(user.Value) < 3 || len(user.Value) > 12 {
		return
	}

	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		return
	}

	if SessionManager[user.Value] != sessionID.Value {
		return
	}
	id := template.HTMLEscapeString(r.FormValue("id"))
	white := template.HTMLEscapeString(r.FormValue("white"))
	black := template.HTMLEscapeString(r.FormValue("black"))

	var chat ChatInfo
	chat.Type = "chess_game"
	var success bool
	if user.Value == white {
		if isPlayerInLobby(black) == true && !isPlayerInGame(black, "") {
			success = fetchSavedGame(id, user.Value)
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
		if isPlayerInLobby(white) == true && !isPlayerInGame(white, "") {
			success = fetchSavedGame(id, user.Value)
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

	} else {
		fmt.Println("Invalid user ajax.go ResumeGame 1")
	}
	w.Write([]byte("false"))
}

// fetches all data of a chess game by the ID
func FetchGameByID(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err != nil || len(username.Value) < 3 || len(username.Value) > 12 {
		w.Write([]byte("Wrong authentication"))
		return
	}

	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		w.Write([]byte("Wrong authentication"))
		return
	}

	if SessionManager[username.Value] != sessionID.Value {
		w.Write([]byte("Wrong authentication"))
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
	username, err := r.Cookie("username")
	if err != nil || len(username.Value) < 3 || len(username.Value) > 12 {
		w.Write([]byte("Wrong authentication"))
		return
	}

	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		w.Write([]byte("Wrong authentication"))
		return
	}

	if SessionManager[username.Value] != sessionID.Value {
		w.Write([]byte("Wrong authentication"))
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
		//check if database connection is open
		if db.Ping() != nil {
			fmt.Printf("ERROR 2 PINGING DB IP: %s \n", ipAddress)
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Please come back later."))
			return
		}

		var name string
		//checking if name exists
		checkName := db.QueryRow("SELECT username FROM userinfo WHERE username=?", username).Scan(&name)
		switch {
		case checkName == sql.ErrNoRows:
			w.Write([]byte(" <img src='img/ajax/available.png' /> Username available"))
			fmt.Printf("Username %s is available.\n", username)
		case checkName != nil:
			fmt.Printf("ERROR 3 CHECKNAME IP is %s\n", ipAddress)
		default:
			w.Write([]byte("<img src='img/ajax/not-available.png' /> Username taken"))
		}
	}
}
