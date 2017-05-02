package mud

import (
	"encoding/json"
	"log"
	"os"

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

		if c.username == t.Name {
			switch t.Type {

			case "chat_all":
			default:
				log.Println("I'm not familiar with type in MUD", t.Type, " sent by ", t.Name)
			}
		}
	}
}
