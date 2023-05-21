package gostuff

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/dchest/captcha"
	"golang.org/x/net/websocket"
)

func UpdateCaptcha(w http.ResponseWriter, r *http.Request) {
	cap := captcha.New()
	w.Write([]byte(cap))
}

// displays player data when mouse hovers over
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
