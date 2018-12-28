package mud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
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
	Lobby   map[string]*websocket.Conn
	Players map[string]*Player // Active users connected to MUD
}{Lobby: make(map[string]*websocket.Conn),
	Players: make(map[string]*Player)}

type KnownCommands struct {
	Commands []string
}

// Parses known commands from commands.json
func ParseKnownCommands() []string {
	data, err := ioutil.ReadFile("data/mud/commands.json")
	if err != nil {
		log.Fatal(err)
	}

	var commands KnownCommands

	err = json.Unmarshal(data, &commands)
	if err != nil {
		log.Fatal(err)
	}

	return commands.Commands
}

func (player *Player) processCommand(command string, connection *MudConnection) {
	knownCommands := ParseKnownCommands()

	for _, tempCommand := range knownCommands {
		if strings.HasPrefix(tempCommand, command) {
			switch tempCommand {
			case "north":
				if player.Location.Row > 0 {
					player.Location.Row = player.Location.Row - 1
					player.setMapVision(world)
					player.Type = "update_map"
					connection.sendJSONWebSocket(&player)
				}
			case "east":
				if player.Location.Col < world.Floors[player.Location.Level].Width {
					player.Location.Col = player.Location.Col + 1
					player.setMapVision(world)
					player.Type = "update_map"
					fmt.Println("player map is ", player.Map)
					connection.sendJSONWebSocket(&player)
				}
			case "south":
				if player.Location.Row < world.Floors[player.Location.Level].Length {
					player.Location.Row = player.Location.Row + 1
					player.setMapVision(world)
					player.Type = "update_map"
					connection.sendJSONWebSocket(&player)
				}
			case "west":
				if player.Location.Col > 0 {
					player.Location.Col = player.Location.Col - 1
					player.setMapVision(world)
					player.Type = "update_map"
					connection.sendJSONWebSocket(&player)
				}
			default:
				fmt.Println("Unknown process Command", command, tempCommand)
			}
		}
	}
}
