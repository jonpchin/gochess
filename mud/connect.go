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
			if gostuff.SessionManager[username.Value] == sessionID.Value {

				ip := ws.Request().RemoteAddr

				Client := &MudConnection{username.Value, ws, ip, ""}

				//only difference between lobby and chatroom is the two lines below
				MudServer.Lobby[username.Value] = ws
				Client.MudConnect()
			}
		}
	}
}
