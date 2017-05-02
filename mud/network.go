package mud

import (
	"sync"

	"golang.org/x/net/websocket"
)

type MudConnection struct {
	username  string
	websocket *websocket.Conn
	clientIP  string
	name      string
}

// Active users connected to MUD
var MudServer = struct {
	sync.RWMutex
	Lobby map[string]*websocket.Conn
}{Lobby: make(map[string]*websocket.Conn)}
