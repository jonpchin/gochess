package mud

import (
	"sync"

	"golang.org/x/net/websocket"
)

type Player struct {
	Name      string
	Inventory []interface{} // What the player is carrying
	Equipment Equipment
	Stats     PlayerStats
	Status    []string // List of afflictions or buffs affecting player
	Bleed     int      // Amount of health the player will lose every tick
	Level     int
	Location  Coordinate
	Area      Area
}

type PlayerStats struct {
	Name         string
	Health       int
	Mana         int
	Energy       int
	Strength     int
	Speed        int
	Dexterity    int
	Intelligence int
	Wisdom       int
}

var All = struct {
	sync.RWMutex
	Players map[string]*websocket.Conn
}{Players: make(map[string]*websocket.Conn)}
