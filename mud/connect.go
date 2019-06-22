package mud

import (
	"github.com/jonpchin/gochess/gostuff"
	"golang.org/x/net/websocket"
)

func EnterMud(ws *websocket.Conn) {

	defer ws.Close()
	username, err := ws.Request().Cookie("username")
	if err == nil {
		sessionID, err := ws.Request().Cookie("sessionID")
		if err == nil {
			if sessionID.Value != "" && gostuff.SessionManager[username.Value] == sessionID.Value {

				ip := ws.Request().RemoteAddr

				Client := &MudConnection{username.Value, ws, ip, ""}

				MudServer.Lobby[username.Value] = ws
				// Ensures username is registered in mud table, if not adds it
				var player Player
				player.Username = username.Value
				player.SessionID = sessionID.Value
				MudServer.Players[username.Value] = &player
				Client.MudConnect()

			}
		}
	}
}
