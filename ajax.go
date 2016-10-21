package gostuff

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
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

	//getting player raating
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
	icon := "ready"
	url := ""
	endUrl := "" //closing the href link
	//second username is nil as it only checks one name
	if isPlayerInGame(lookupName, "") {
		status = "vs. " + PrivateChat[lookupName]
		icon = "playing"
		id, _ := GetGameID(lookupName)
		url = "<a href=/chess/memberChess?spectate&id=" + strconv.Itoa(id) + ">"
		endUrl = "</a>"
	}

	var result = "<img src='../img/icons/" + icon + ".png' alt='status'>" +
		url + lookupName + " " + status + endUrl +
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

// fetches names of players and ID of games from database with the param being the range of the ID inclusive
// returns JSON string of all games in range, returning blank string means there was an error
func FetchPlayersInRange(w http.ResponseWriter, r *http.Request) {

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
	start := template.HTMLEscapeString(r.FormValue("start"))
	last := template.HTMLEscapeString(r.FormValue("last"))

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		w.Write([]byte(""))
		return
	}

	//looking up players rating
	rows, err := db.Query("SELECT id, white, black FROM grandmaster WHERE id >= ? AND id <= ?", start, last)
	if err != nil {
		log.Println(err)
		w.Write([]byte(""))
		return
	}

	defer rows.Close()
	var all NamesAndID
	var storage []NamesAndID

	for rows.Next() {

		err = rows.Scan(&all.ID, &all.White, &all.Black)

		if err != nil {
			log.Println(err)
			w.Write([]byte(""))
			return
		}
		storage = append(storage, all)
	}
	allNamesAndID, err := json.Marshal(storage)
	if err != nil {
		log.Println(err)
	}

	w.Write([]byte(string(allNamesAndID)))
}

// fetches all data of a chess game by the ID
func FetchGameByID(w http.ResponseWriter, r *http.Request) {
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

	id := template.HTMLEscapeString(r.FormValue("gameID"))

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		w.Write([]byte(""))
		return
	}

	var all GrandMasterGame

	//looking up players rating
	err = db.QueryRow("SELECT * FROM grandmaster WHERE id=?", id).Scan(&all.ID, &all.Event, &all.Site,
		&all.Date, &all.Round, &all.White,
		&all.Black, &all.Result, &all.WhiteElo, &all.BlackElo,
		&all.ECO, &all.Moves, &all.EventDate)

	if err != nil {
		log.Println(err)
		w.Write([]byte(""))
		return
	}

	allGames, err := json.Marshal(all)
	if err != nil {
		log.Println(err)
		w.Write([]byte(""))
		return
	}
	w.Write([]byte(string(allGames)))
}
