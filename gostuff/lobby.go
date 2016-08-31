package gostuff

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"math/rand"
	"os"
	"time"
)

func (c *Connection) LobbyConnect() {

	defer broadCast(c.username) //remove user when they disconnect from socket
	counter := 0
	start := time.Now()

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
			fmt.Println("lobby.go 1 Reader 1 ", err.Error())
			break
		}

		if c.username == t.Name {
			switch t.Type {

			case "chat_all":

				if len(reply) > 225 {
					log.Printf("User: %s IP %s has exeeded the 225 character limit by %d byte units.\n", t.Name, c.clientIP, len(reply))
					return
				}
				//keeps track of messages are sent in a given interval
				counter++

				if counter > 4 {
					elapsed := time.Since(start)
					if elapsed < time.Second*10 {
						log.Printf("User: %s IP: %s was spamming chat.\n", t.Name, c.clientIP)
						return
					}
					start = time.Now()
					counter = 0

				}
				go func() {
					for _, cs := range Chat.Lobby {
						if err := websocket.Message.Send(cs, reply); err != nil {

							// we could not send the message to a peer
							fmt.Println("lobby.go error 2 Could not send message to ", c.clientIP, err.Error())
						}
					}
				}()
			case "fetch_matches":
				//send in array instead of sending individual
				for _, value := range Pending.Matches {
					match, err := json.Marshal(value)
					if err != nil {
						fmt.Println("Just receieved a message I couldn't encode: on fetch_matches")
						fmt.Println("Exact error: " + err.Error())
						break
					}
					result := string(match)
					websocket.Message.Send(c.websocket, result)

				}

			case "fetch_players":
				//send in array instead of sending individual
				var player Online

				player.Type = "fetch_players"

				for key, _ := range Chat.Lobby {
					player.Name = key
					websocket.JSON.Send(c.websocket, player)
				}

			case "match_seek":

				//check to make sure player only has a max of three matches seeks pending, used to prevent flood match seeking
				if c.totalMatches >= 3 {
					t.Type = "maxThree"
					if err := websocket.JSON.Send(Chat.Lobby[t.Name], &t); err != nil {
						// we could not send the message to a peer
						log.Println("match lobby.go Could not send message to ", c.clientIP, err.Error())
					}
					break //notify user that only three matches pending max are allowed
				} else {
					c.totalMatches++
				}

				var match SeekMatch
				if err := json.Unmarshal(message, &match); err != nil {
					fmt.Println("Just receieved a message I couldn't decode:")
					fmt.Println(string(reply))
					fmt.Println("Exact error: " + err.Error())
					break
				}

				//check if player already has a game started, if there is a game in progress alert player
				if isPlayerInGame(t.Name, match.Opponent) == true {
					fmt.Println("Player is already in a game!")
					t.Type = "alert"
					if err := websocket.JSON.Send(Chat.Lobby[t.Name], &t); err != nil {
						// we could not send the message to a peer
						fmt.Println("lobby.go error 7 Could not send message to ", c.clientIP, err.Error())
					}
					break
				}

				//verify.go
				if checkTime(match.TimeControl) == false {
					fmt.Println("An invalid time control has been selected.")
					break
				}

				//fetching rating from back end
				errRate, bullet, blitz, standard := GetRating(match.Name)
				if errRate != "" {
					fmt.Println("Cannot get rating lobby.go match_seek")
					break
				}
				switch match.TimeControl {
				case 1:
					match.Rating = bullet
					match.GameType = "bullet"
				case 2:
					match.Rating = bullet
					match.GameType = "bullet"
				case 3:
					match.Rating = blitz
					match.GameType = "blitz"
				case 4:
					match.Rating = blitz
					match.GameType = "blitz"
				case 5:
					match.Rating = blitz
					match.GameType = "blitz"
				case 10:
					match.Rating = blitz
					match.GameType = "blitz"
				default: //for 15, 20, 30 or 45 minute game defaults to standard
					match.Rating = standard
					match.GameType = "standard"
				}

				var start int16 = 0
				for {
					if _, ok := Pending.Matches[start]; ok {
						start++
					} else {
						break
					}
				}
				//value := fmt.Sprintf("%d", start)
				match.MatchID = start
				//used in backend to keep track of all pending games waiting for a player to accept

				Pending.Matches[start] = &match

				result, err := json.Marshal(match)
				if err != nil {
					fmt.Println("Just receieved a message I couldn't encode on error 8", err)
					break
				}

				finalMessage := string(result)
				go func() {
					for _, cs := range Chat.Lobby {
						if err := websocket.Message.Send(cs, finalMessage); err != nil {
							// we could not send the message to a peer
							fmt.Println("lobby.go error 9 Could not send message to ", c.clientIP, err.Error())
						}
					}
				}()
			case "cancel_match":

				var match SeekMatch
				if err := json.Unmarshal(message, &match); err != nil {
					fmt.Println("Just receieved a message I couldn't decode in lobby.go cancel_match:")
					fmt.Println(string(reply))
					fmt.Println("Exact error: " + err.Error())
					break
				}

				//number, _ := strconv.ParseInt(match.MatchID, 10, 0)
				//deletes key from hash table
				delete(Pending.Matches, match.MatchID)

				//deletes pending match counter
				c.totalMatches--
				//check if its a private match, if so then delete it and break out
				if match.Opponent != "" {
					fmt.Println("Private match deleted")
					break //no need to continue as this is a private match
				}

				go func() {
					for _, cs := range Chat.Lobby {
						websocket.Message.Send(cs, reply)
					}
				}()

			case "accept_match":

				var match SeekMatch
				var game ChessGame
				if err := json.Unmarshal(message, &match); err != nil {
					log.Println("Just receieved a message I couldn't decode:")
					log.Println(string(reply))
					log.Println("lobby.go error 11 Exact error: " + err.Error())
					break
				}

				//check if player already has a game started, if there is a game in progress alert player
				if isPlayerInGame(match.Name, match.Opponent) == true {
					log.Println("lobby.go Player already has a game. ")
					//alerting player
					t.Type = "alert"
					if err := websocket.JSON.Send(Chat.Lobby[t.Name], &t); err != nil {
						// we could not send the message to a peer
						log.Println("error 10 lobby.go Could not send message to ", c.clientIP, err.Error())
					}
					break
				}

				//checking to make sure both player's rating is in range, used as a backend rating check
				errMessage, bullet, blitz, standard := GetRating(match.Name)
				if errMessage != "" {
					fmt.Println("Cannot get rating connection.go accept_match")
					break
				}

				if Pending.Matches[match.MatchID].Opponent == "" { //only use this case for public matches
					if Pending.Matches[match.MatchID].GameType == "bullet" && (bullet < Pending.Matches[match.MatchID].MinRating || bullet > Pending.Matches[match.MatchID].MaxRating) {
						fmt.Println("Bullet Rating not in range.")
						break
					} else if Pending.Matches[match.MatchID].GameType == "blitz" && (blitz < Pending.Matches[match.MatchID].MinRating || blitz > Pending.Matches[match.MatchID].MaxRating) {
						fmt.Println("Blitz Rating not in range.")
						break
					} else if Pending.Matches[match.MatchID].GameType == "standard" && (standard < Pending.Matches[match.MatchID].MinRating || standard > Pending.Matches[match.MatchID].MaxRating) {
						fmt.Println("Standard Rating not in range.")
						break
					}
				}

				//bullet, blitz or standard game type
				game.GameType = Pending.Matches[match.MatchID].GameType

				//seting up the game info such as white/black player, time control, etc
				rand.Seed(time.Now().UnixNano())
				randomNum := rand.Intn(2)

				//randomly selects both players to be white or black
				if randomNum == 0 {
					game.WhitePlayer = match.Name
					if game.GameType == "bullet" {
						game.WhiteRating = bullet

					} else if game.GameType == "blitz" {
						game.WhiteRating = blitz

					} else {
						game.WhiteRating = standard

					}

					game.BlackRating = Pending.Matches[match.MatchID].Rating
					game.BlackPlayer = Pending.Matches[match.MatchID].Name

				} else {
					game.WhitePlayer = Pending.Matches[match.MatchID].Name
					if game.GameType == "bullet" {
						game.BlackRating = bullet

					} else if game.GameType == "blitz" {
						game.BlackRating = blitz
					} else {
						game.BlackRating = standard
					}

					game.WhiteRating = Pending.Matches[match.MatchID].Rating
					game.BlackPlayer = match.Name
				}
				//White for white to move or Black for black to move, white won, black won, stalemate or draw.
				game.Status = "White"

				//no moves yet so nill/null
				game.GameMoves = nil

				game.TimeControl = Pending.Matches[match.MatchID].TimeControl
				//for simplicity we will only allow minutes
				game.WhiteMinutes = Pending.Matches[match.MatchID].TimeControl
				game.WhiteSeconds = 0
				game.WhiteMilli   = 0
				game.BlackMinutes = Pending.Matches[match.MatchID].TimeControl
				game.BlackSeconds = 0
				game.BlackMilli   = 0
				game.PendingDraw = false
				game.Rated = Pending.Matches[match.MatchID].Rated;

				var start int16 = 0
				for {
					if _, ok := All.Games[start]; ok {
						start++
					} else {
						break
					}
				}
				//value := fmt.Sprintf("%d", start)
				game.ID = start
				//used in backend to keep track of all pending games waiting for a player to accept
				All.Games[start] = &game

				//number, _ := strconv.ParseInt(match.MatchID, 10, 0)
				//no longer need all the pending matches as game will be started
				for key, value := range Pending.Matches {
					//deletes all pending matches for either players
					if value.Name == game.WhitePlayer || value.Name == game.BlackPlayer {
						delete(Pending.Matches, key)
					}
				}

				//sending to front end for url redirection
				var acceptmatch AcceptMatch
				acceptmatch.Type = "accept_match"
				acceptmatch.Name = game.WhitePlayer
				acceptmatch.TargetPlayer = game.BlackPlayer

				//setting up the private chat between two players and send move connection
				PrivateChat[acceptmatch.Name] = acceptmatch.TargetPlayer
				PrivateChat[acceptmatch.TargetPlayer] = acceptmatch.Name

				//intitalizes all the variables of the game
				initGame(game.ID)

				startGame, _ := json.Marshal(acceptmatch)

				for _, cs := range Chat.Lobby {
					if err := websocket.Message.Send(cs, string(startGame)); err != nil {
						fmt.Println("lobby.go error 12 error is ", err)
					}
				}

				//starting white's clock first, this goroutine will keep track of both players clock for this game
				go setClocks(game.ID, t.Name)

			case "private_match":

				var match SeekMatch
				if err := json.Unmarshal(message, &match); err != nil {
					fmt.Println("Just receieved a message I couldn't decode:")
					fmt.Println(string(reply))
					fmt.Println("Exact error: " + err.Error())
					break
				}
				//check if player already has a game started, if there is a game in progress alert player
				if isPlayerInGame(match.Name, match.Opponent) == true {
					fmt.Println("Player already has a game.")
					//alerting player
					t.Type = "alert"
					if err := websocket.JSON.Send(Chat.Lobby[t.Name], &t); err != nil {
						// we could not send the message to a peer
						fmt.Println("lobby.go error 7 Could not send message to ", c.clientIP, err.Error())
					}
					break
				}

				//check length of name to make sure its 3-12 characters long
				if len(match.Opponent) < 3 || len(match.Opponent) > 12 {
					fmt.Println("Username is too long or too short")
					break
				}
				//a player should not be able to match himself
				if t.Name == match.Opponent {
					fmt.Println("You can't match yourself!")
					break
				}

				//check if opponent is in the lobby or not
				if _, ok := Chat.Lobby[match.Opponent]; !ok {
					//alerting player
					t.Type = "absent"
					if err := websocket.JSON.Send(Chat.Lobby[t.Name], &t); err != nil {
						// we could not send the message to a peer
						fmt.Println("lobby.go error 7 Could not send message to ", c.clientIP, err.Error())
					}
					break
				}

				//verify.go
				if checkTime(match.TimeControl) == false {
					fmt.Println("An invalid time control has been selected.")
					break
				}

				//fetching rating from back end
				errMessage, bullet, blitz, standard := GetRating(match.Name)
				if errMessage != "" {
					fmt.Println("Cannot get rating lobby.go private_match")
					break
				}
				switch match.TimeControl {
				case 1:
					match.Rating = bullet
					match.GameType = "bullet"
				case 2:
					match.Rating = bullet
					match.GameType = "bullet"
				case 3:
					match.Rating = blitz
					match.GameType = "blitz"
				case 4:
					match.Rating = blitz
					match.GameType = "blitz"
				case 5:
					match.Rating = blitz
					match.GameType = "blitz"
				case 10:
					match.Rating = blitz
					match.GameType = "blitz"
				default: //for 15, 20, 30 or 45 minute game defaults to standard
					match.Rating = standard
					match.GameType = "standard"
				}

				//check to make sure player only has a max of three matches seeks pending, used to prevent flood match seeking
				if c.totalMatches >= 3 {
					t.Type = "maxThree"
					if err := websocket.JSON.Send(Chat.Lobby[t.Name], &t); err != nil {
						// we could not send the message to a peer
						log.Println("match lobby.go Could not send message to ", c.clientIP, err.Error())
					}
					break //notify user that only three matches pending max are allowed
				} else {
					c.totalMatches++
				}

				var start int16 = 0
				for {
					if _, ok := Pending.Matches[start]; ok {
						start++
					} else {
						break
					}
				}
				//value := fmt.Sprintf("%d", start)
				match.MatchID = start
				//used in backend to keep track of all pending seeks waiting for a player to accept
				Pending.Matches[start] = &match

				result, err := json.Marshal(match)
				if err != nil {
					fmt.Println("Just receieved a message I couldn't encode on error 8", err)
					break
				}

				finalMessage := string(result)
				go func() {
					for name, _ := range Chat.Lobby {
						if name == match.Opponent || name == match.Name { //send to self and opponent
							if err := websocket.Message.Send(Chat.Lobby[name], finalMessage); err != nil {
								// we could not send the message to a peer
								fmt.Println("lobby.go error 9 Could not send message to ", c.clientIP, err.Error())
							}
						}
					}
				}()

			default:
				fmt.Println("I'm not familiar with type " + t.Type)
			}
		} else {
			log.Printf("IP %s Invalid websocket authentication in lobby.\n", c.clientIP)
			return
		}
	}
}
