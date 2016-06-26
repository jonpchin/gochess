package gostuff

import(
	"fmt"
	"os"
	"log"
	"golang.org/x/net/websocket"
	"encoding/json"
)

func (c *ConnectionEngine) EngineSetup(){
	
	defer exitGame(c.username) //remove user when they disconnect from socket

	logFile, _ := os.OpenFile("logs/chat.txt", os.O_APPEND|os.O_WRONLY, 0666)

	defer logFile.Close()

	// direct all log messages to log.txt
	log.SetOutput(logFile)

	for {
		var reply string

		if err := websocket.Message.Receive(c.websocket, &reply); err != nil {
			//fmt.Println("A user has drop web socket connection ", err)
			break
		}

		var t APITypeOnly
		message := []byte(reply)
		if err := json.Unmarshal(message, &t); err != nil {
			fmt.Println("Just receieved a message I couldn't decode:")
			fmt.Println(string(reply))
			fmt.Println("connection.go 1 ChessConnect 1 ", err.Error())
			break
		}
		if c.username == t.Name {
			switch t.Type {

			case "send_move":


			case "chat_private":
			

			case "chess_game":
				

			case "abort_game":

				

			case "game_over":

			
			case "rematch":

	
			case "accept_rematch":


			case "draw_game":
			
			
			default:
				fmt.Println("I'm not familiar with type " + t.Type)
			}
		} else {
			log.Printf("IP %s Invalid websocket authentication in chess room.\n", c.clientIP)
			return
		}
	}
}