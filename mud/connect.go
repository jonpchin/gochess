package mud

import (
	"fmt"

	"github.com/jonpchin/gochess/gostuff"
	"golang.org/x/net/websocket"
)

func EnterMud(ws *websocket.Conn) {

	defer ws.Close()
	username, err := ws.Request().Cookie("username")
	if err == nil {
		sessionID, err := ws.Request().Cookie("sessionID")
		if err == nil {
			if gostuff.SessionManager[username.Value] == sessionID.Value {

				ip := ws.Request().RemoteAddr

				Client := &MudConnection{username.Value, ws, ip, ""}

				MudServer.Lobby[username.Value] = ws
				// Ensures username is registered in mud table
				err := lookupName(username.Value)
				if err == nil {
					var player Player
					player.Username = username.Value
					player.SessionID = sessionID.Value
					MudServer.Players[username.Value] = &player
					Client.MudConnect()
				} else {
					fmt.Println("Could not get username from mud", err)
				}
			}
		}
	}
}
