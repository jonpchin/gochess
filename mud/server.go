package mud

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"golang.org/x/net/websocket"
)

type CommandMessage struct {
	Type    string
	Command string
}

func (c *MudConnection) MudConnect() {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	for {
		var reply string

		if err := websocket.Message.Receive(c.websocket, &reply); err != nil {
			break
		}
		var player Player
		message := []byte(reply)

		if err := json.Unmarshal(message, &player); err != nil {
			log.Println("Just receieved a message I couldn't decode:", string(reply), err)
			break
		}

		// Check to make sure player is not pretending to be someone else or changing name without permission
		if MudServer.Players[player.Username].isCredValid(player.Username, player.SessionID) == false {
			log.Println("Invalid credentials")
			break
		}

		switch player.Type {
		case "command":

			var command CommandMessage
			if err := json.Unmarshal(message, &command); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			//tempPlayer := MudServer.Players[player.Username]
			player.processCommand(command.Command, c)
			MudServer.Players[player.Username].Location = player.Location

		case "connect_mud":

			var player Player
			if err := json.Unmarshal(message, &player); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			if isNameExistForPlayer(c.username) {
				player.Type = "enter_world"
				player.enterWorld(LOAD_PLAYER, c) // Name already exists for player

				//TODO: Make sure all other data is set for players
				MudServer.Players[player.Username].Location = player.Location

				fmt.Println("Player already exists", c.username)
			} else {
				player.Type = "ask_name"
				err := websocket.JSON.Send(MudServer.Lobby[c.username], &player)
				if err != nil {
					log.Println(err)
				}
			}
		case "check_name":
			if isNameTaken(player.Name) {
				player.Type = "name_taken"
				c.sendJSONWebSocket(&player)
			} else {
				player.Type = "name_available"
				c.sendJSONWebSocket(&player)
			}
		case "enter_world_first_time":
			if err := json.Unmarshal(message, &player); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			player.Type = "enter_world"
			//player.updateByRaceClass()
			player.Location = HOME_BASE
			player.Inventory = nil

			MudServer.Players[player.Name] = &player

			player.registerPlayer()
			player.enterWorld(SKIP_LOAD, c)
		case "fetch_map":

			var player Player
			if err := json.Unmarshal(message, &player); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			player.Type = "update_map"
			player.setPlayerMap()
			c.sendJSONWebSocket(&player)
		default:
			log.Println("I'm not familiar with type in MUD", player.Type, " sent by ", player.Name)
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
