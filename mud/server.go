package mud

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/jonpchin/gochess/gostuff"

	"golang.org/x/net/websocket"
)

func (c *MudConnection) MudConnect() {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	for {
		var reply string

		if err := websocket.Message.Receive(c.websocket, &reply); err != nil {
			break
		}
		var t gostuff.MessageType
		message := []byte(reply)

		if err := json.Unmarshal(message, &t); err != nil {
			log.Println("Just receieved a message I couldn't decode:", string(reply), err)
			break
		}

		switch t.Type {

		case "connect_mud":
			if isNameExistForPlayer(c.username) {
				var player Player
				if err := json.Unmarshal(message, &player); err != nil {
					log.Println("Just receieved a message I couldn't decode:", string(reply), err)
					break
				}
				player.Type = "enter_world"
				// Check to make sure player is not pretending to be someone else or changing name without permission
				if MudServer.Players[player.Username].isCredValid(player.Username, player.Name, player.SessionID) == false {
					log.Println("Invalid credentials")
					break
				}
				player.enterWorld(LOAD_PLAYER, c) // Name already exists for player
				fmt.Println("Player already exists", c.username)
			} else {
				t.Type = "ask_name"
				err := websocket.JSON.Send(MudServer.Lobby[c.username], &t)
				if err != nil {
					log.Println(err)
				}
			}
		case "check_name":
			if isNameTaken(t.Name) {
				t.Type = "name_taken"
				c.sendJSONWebSocket(&t)
			} else {
				t.Type = "name_available"
				c.sendJSONWebSocket(&t)
			}
		case "enter_world_first_time":
			var player Player
			if err := json.Unmarshal(message, &player); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			if MudServer.Players[player.Username].isCredValidFirstTime(player.Username, player.SessionID) == false {
				log.Println("Invalid credentials")
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

			if MudServer.Players[playerMap.Creds.Username].isCredValid(
				playerMap.Creds.Username, playerMap.Creds.Name, playerMap.Creds.SessionID) == false {
				log.Println("Invalid credentials")
				break
			}

			playerMap.Type = "update_map"
			playerMap.setPlayerMap()
			c.sendJSONWebSocket(&playerMap)
		default:
			log.Println("I'm not familiar with type in MUD", t.Type, " sent by ", t.Name)
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
