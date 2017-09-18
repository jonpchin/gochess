package mud

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"golang.org/x/net/websocket"
)

type Authentication struct {
	Type      string
	Username  string // Go Play Chess account
	Name      string // Mud account (optional)
	SessionID string
}

func (c *MudConnection) MudConnect() {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	for {
		var reply string

		if err := websocket.Message.Receive(c.websocket, &reply); err != nil {
			break
		}
		var a Authentication
		message := []byte(reply)

		if err := json.Unmarshal(message, &a); err != nil {
			log.Println("Just receieved a message I couldn't decode:", string(reply), err)
			break
		}

		// Check to make sure player is not pretending to be someone else or changing name without permission
		if MudServer.Players[a.Username].isCredValid(a.Username, a.SessionID) == false {
			log.Println("Invalid credentials")
			break
		}

		switch a.Type {

		case "connect_mud":

			var player Player
			if err := json.Unmarshal(message, &player); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			if isNameExistForPlayer(c.username) {
				player.Type = "enter_world"
				player.enterWorld(LOAD_PLAYER, c) // Name already exists for player
				fmt.Println("Player already exists", c.username)
			} else {
				a.Type = "ask_name"
				err := websocket.JSON.Send(MudServer.Lobby[c.username], &a)
				if err != nil {
					log.Println(err)
				}
			}
		case "check_name":
			if isNameTaken(a.Name) {
				a.Type = "name_taken"
				c.sendJSONWebSocket(&a)
			} else {
				a.Type = "name_available"
				c.sendJSONWebSocket(&a)
			}
		case "enter_world_first_time":
			var player Player
			if err := json.Unmarshal(message, &player); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			player.Type = "update_player"
			//player.updateByRaceClass()
			player.Location = HOME_BASE
			player.Inventory = nil

			MudServer.Players[player.Name] = &player
			registerName(player.Name, c.username)
			player.save()
			player.enterWorld(SKIP_LOAD, c)
		case "fetch_map":

			var playerMap PlayerMap
			if err := json.Unmarshal(message, &playerMap); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			playerMap.Type = "update_map"
			playerMap.setPlayerMap()
			c.sendJSONWebSocket(&playerMap)
		default:
			log.Println("I'm not familiar with type in MUD", a.Type, " sent by ", a.Name)
		}
	}
}

// Sends message to only one person
func (c *MudConnection) sendJSONWebSocket(message interface{}) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	if err := websocket.JSON.Send(c.websocket, message); err != nil {
		log.Println(err)
	}
}

// targets could be a list of players or a map of players
// message is a struct that needs to be encoded to JSON before sending
func (c *MudConnection) broadCastJSONWebSocket(targets interface{}, message interface{}) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	switch reflect.TypeOf(targets).Kind() {
	case reflect.Slice:
		listOfPlayers := reflect.ValueOf(targets)

		for i := 0; i < listOfPlayers.Len(); i++ {
			fmt.Println(listOfPlayers.Index(i))
		}
	case reflect.Map:
		listOfPlayers := reflect.ValueOf(targets)
		for _, key := range listOfPlayers.MapKeys() {
			strct := listOfPlayers.MapIndex(key)
			fmt.Println(key.Interface(), strct.Interface())
		}

	default:
		log.Println("No reflection type in sendJSONWebSocket")
	}
}
