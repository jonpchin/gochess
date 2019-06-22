package gostuff

import (
	"fmt"
	"sync"

	"golang.org/x/net/websocket"
)

//used to identify who the socket connection is
type Connection struct {
	username  string
	websocket *websocket.Conn
	clientIP  string
}

//stores information for a message from chat for JSON
type ChatInfo struct {
	Type string
	Name string
	Text string
}

//sends out seek matches real time in the lobby
type SeekMatch struct {
	Type        string
	Name        string
	Opponent    string
	Rating      int16  //player own rating
	GameType    string //bullet, blitz, standard, correspondence
	MatchID     int
	TimeControl int
	MinRating   int16
	MaxRating   int16
	Rated       string
}

//used to store two player's name for redirecting on the front end in JavaScript
type AcceptMatch struct {
	Type         string
	Name         string
	TargetPlayer string
}

//active users connected to chess room socket
var Active = struct {
	sync.RWMutex
	Clients map[string]*websocket.Conn
}{Clients: make(map[string]*websocket.Conn)}

//active users connected to lobby socket
var Chat = struct {
	sync.RWMutex
	Lobby map[string]*websocket.Conn
}{Lobby: make(map[string]*websocket.Conn)}

// stores information of players that are online
// used to identify what kind of message it for incoming socket messages for JSON and check authentication
type MessageType struct {
	Type    string
	Name    string
	Message string
}

//websocket handler for lobby
func EnterLobby(ws *websocket.Conn) {

	defer ws.Close()
	username, err := ws.Request().Cookie("username")
	if err == nil {
		sessionID, err := ws.Request().Cookie("sessionID")
		if err == nil {
			if sessionID.Value != "" && SessionManager[username.Value] == sessionID.Value {

				ip := ws.Request().RemoteAddr
				Client := &Connection{username.Value, ws, ip}

				//only difference between lobby and chatroom is the two lines below
				Chat.Lobby[username.Value] = ws
				Client.LobbyConnect()
			}
		}
	}
}

//websocket handler for gameroom
func EnterChess(ws *websocket.Conn) {
	defer ws.Close()
	username, err := ws.Request().Cookie("username")
	if err == nil {
		sessionID, err := ws.Request().Cookie("sessionID")
		if err == nil {
			if sessionID.Value != "" && SessionManager[username.Value] == sessionID.Value {

				ip := ws.Request().RemoteAddr
				Client := &Connection{username.Value, ws, ip}
				Active.Clients[username.Value] = ws

				Client.ChessConnect()
			}
		}
	}
}

//returns the total number of seeks that a player has pending in the lobby
func countMatches(player string) int8 { //player should have no more then 3 seeks at a time

	var total int8 = 0

	for _, match := range Pending.Matches {
		if match.Name == player {
			total++
		}
	}
	return total
}

//broadcast to chess room that player has disconnected socket
func broadCast(user string) {

	delete(Chat.Lobby, user)

	var on MessageType
	on.Type = "broadcast"
	on.Name = user
	for _, cs := range Chat.Lobby {
		if err := websocket.JSON.Send(cs, on); err != nil {
			// we could not send the message to a peer
			fmt.Println("broadcast error ", err)
		}
	}
}

//function is called when player leaves the chess room
func exitGame(user string) {
	//check if user is in PrivateChat map, delete player key's if necessary
	if _, ok := PrivateChat[user]; ok {
		if checkTable(user) == false {
			var t ChatInfo
			t.Type = "leave"
			t.Text = user + " has left the chess room."
			var otherPlayer = PrivateChat[user]

			removePendingMatches(user)

			if _, pass := Active.Clients[otherPlayer]; pass {
				if err := websocket.JSON.Send(Active.Clients[otherPlayer], &t); err != nil {
					//we could not send the message to a peer
					fmt.Println("exitgame.go error  Could not send message to ", err)
				}
				delete(PrivateChat, user)
				delete(PrivateChat, otherPlayer)
			}
		}
	}
	delete(Active.Clients, user)
}

//returns true if a player is at a given table
// function isPlayersInGame is similar to this checkTable
// this checkTable only checks one username instead of two like isPlayersInGame
func checkTable(user string) bool {
	for _, table := range All.Games {
		if table.WhitePlayer == user || table.BlackPlayer == user {
			return true
		}
	}
	return false
}

//returns true if a player or opponent has a game started
func isPlayersInGame(name, opponent string) bool {
	for _, game := range All.Games {
		if game.WhitePlayer == name || game.BlackPlayer == name {
			return true
		}
		if game.WhitePlayer == opponent || game.BlackPlayer == opponent {
			return true
		}
	}
	return false
}

//checks if a player is in the lobby
func isPlayerInLobby(player string) bool {
	for name, _ := range Chat.Lobby {
		if name == player {
			return true
		}
	}
	return false
}

//checks if player is in chess room
func isPlayerInChess(player string) bool {
	for name, _ := range Active.Clients {
		if name == player {
			return true
		}
	}
	return false
}

//remove all pending matches from a player
func removePendingMatches(name string) {
	for key, value := range Pending.Matches {

		//deletes all pending matches for either players
		if value.Name == name || value.Opponent == name {
			delete(Pending.Matches, key)
		}
	}
}

//remove player from the table that they were viewing when they leave the chess room
func removeViewer(name string, id int) []string {

	var remainingViewers []string

	for _, value := range Verify.AllTables[id].observe.Names {
		if name != value {
			remainingViewers = append(remainingViewers, value)
		}
	}
	return remainingViewers
}

//gets the game ID that a player is currently playing
//the bool indicates whether the player is currently playing a game
func GetGameID(name string) (int, bool) {
	for key, value := range All.Games {
		if name == value.WhitePlayer || name == value.BlackPlayer {
			return key, true
		}
	}
	return 0, false
}
