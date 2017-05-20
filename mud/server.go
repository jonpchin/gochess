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
			//fmt.Println("A user has drop web socket connection ", err)
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
			if checkNameExist(c.username) {
				//enterWorld(c.username)
			} else {
				t.Type = "askName"
				c.sendJSONWebSocket(All.Players, &t)
			}
		default:
			log.Println("I'm not familiar with type in MUD", t.Type, " sent by ", t.Name)
		}
	}
}

// targets could be a list of players or a map of players
// message is a struct that needs to be encoded to JSON before sending
func (c *MudConnection) sendJSONWebSocket(targets interface{}, message interface{}) {

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
