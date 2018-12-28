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
	Username  string // Go Play Chess account
	Name      string // Mud account (optional)
	SessionID string
}

type Credentials struct {
	Type  string
	Creds Authentication
}

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
		var a Credentials
		message := []byte(reply)

		if err := json.Unmarshal(message, &a); err != nil {
			log.Println("Just receieved a message I couldn't decode:", string(reply), err)
			break
		}

		// Check to make sure player is not pretending to be someone else or changing name without permission
		if MudServer.Players[a.Creds.Username].isCredValid(a.Creds.Username, a.Creds.SessionID) == false {
			log.Println("Invalid credentials")
			break
		}

		switch a.Type {
		case "command":

			var command CommandMessage
			if err := json.Unmarshal(message, &command); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			player := MudServer.Players[a.Creds.Username]
			player.processCommand(command.Command, c)
			MudServer.Players[a.Creds.Username].Location = player.Location

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
				MudServer.Players[a.Creds.Username].Location = player.Location

				fmt.Println("Player already exists", c.username)
			} else {
				a.Type = "ask_name"
				err := websocket.JSON.Send(MudServer.Lobby[c.username], &a)
				if err != nil {
					log.Println(err)
				}
			}
		case "check_name":
			if isNameTaken(a.Creds.Name) {
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

			var player Player
			if err := json.Unmarshal(message, &player); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			player.Type = "update_map"
			player.setPlayerMap()
			c.sendJSONWebSocket(&player)
		default:
			log.Println("I'm not familiar with type in MUD", a.Type, " sent by ", a.Creds.Name)
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
