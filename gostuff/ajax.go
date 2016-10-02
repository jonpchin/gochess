package gostuff

import (
	"fmt"
	"github.com/dchest/captcha"
	"golang.org/x/net/websocket"
	"html/template"
	"net/http"
)

func UpdateCaptcha(w http.ResponseWriter, r *http.Request) {
	cap := captcha.New()
	w.Write([]byte(cap))
}

func GetPlayerData(w http.ResponseWriter, r *http.Request) { //displays player data when mouse hovers over
	username, err := r.Cookie("username")
	if err != nil{
		return
	}
	
	sessionID, err := r.Cookie("sessionID")
	if err != nil{
		return
	}
	
	if SessionManager[username.Value] != sessionID.Value {
		return
	}
	
	if len(username.Value) < 3 || len(username.Value) > 12 {
		w.Write([]byte("Invalid name"))
		return
	}

	//getting player raating
	ratingError, bulletRating, blitzRating, standardRating := GetRating(username.Value)
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
	//second username is nil as it only checks one name
	if isPlayerInGame(username.Value, "") == true {
		status = "vs. " + PrivateChat[username.Value]
		icon = "playing"
	}

	var result = "<img src='../img/icons/" + icon + ".png' alt='status'>" + username.Value + " " + status +
		"<br><img src='../img/icons/bullet.png' alt='bullet'>" + bullet +
		"<img src='../img/icons/blitz.png' alt='blitz'>" + blitz +
		"<img src='../img/icons/standard.png' alt='standard'>" + standard

	w.Write([]byte(result))
}

func ResumeGame(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Cookie("username")
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
