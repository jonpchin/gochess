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
			fmt.Println("error 1", err)
			//fmt.Println("A user has drop web socket connection ", err)
			break
		}
		var t gostuff.MessageType
		message := []byte(reply)

		if err := json.Unmarshal(message, &t); err != nil {
			fmt.Println("Just receieved a message I couldn't decode:", string(reply), err)
			break
		}

		switch t.Type {

		case "connect_mud":
			if isNameExistForPlayer(c.username) {
				var player Player
				player.Username = c.username
				player.enterWorld(LOAD_PLAYER, c)
				fmt.Println("Name already exists for player", c.username)
			} else {
				t.Type = "get_player_data"
				err := websocket.JSON.Send(MudServer.Lobby[c.username], &t)
				if err != nil {
					fmt.Println(err)
				}
			}
		case "check_name":
			if isNameTaken(t.Name) {
				fmt.Println("Name already exists for", t.Name)
				t.Type = "name_taken"
				c.sendJSONWebSocket(&t)
			} else {
				t.Type = "name_available"
				c.sendJSONWebSocket(&t)
			}
		case "enter_world_first_time":
			var player Player
			if err := json.Unmarshal(message, &player); err != nil {
				fmt.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			player.Type = "update_player"
			player.updateByRaceClass()
			player.Location = HOME_BASE
			player.Inventory = nil

			MudServer.Players[player.Name] = &player
			registerName(player.Name, c.username)
			player.save()
			player.enterWorld(SKIP_LOAD, c)
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
