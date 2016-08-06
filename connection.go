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

//reads in all incoming web sockets and sends them out to the correct users
func (c *Connection) ChessConnect() {

	defer exitGame(c.username) //remove user when they disconnect from socket
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
			fmt.Println("connection.go 1 ChessConnect 1 ", err.Error())
			break
		}
		if c.username == t.Name {
			switch t.Type {

			case "send_move":

				var game GameMove

				err := json.Unmarshal(message, &game)
				if err != nil {
					fmt.Println("Just receieved a message I couldn't decode:")
					fmt.Println(string(message))
					fmt.Println("connection.go 1 ChessConnect 2 ", err.Error())
					break
				}

				//check to make sure game exists to prevent nil pointer dereference
				if _, ok := All.Games[game.ID]; !ok {
					fmt.Println("No such game exists for ", t.Name)
					break
				}

				var result bool
				//check if its correct players turn and if move is valid before sending
				result = chessVerify(game.Source, game.Target, game.ID)
				if result == false {
					totalMoves := (len(All.Games[game.ID].GameMoves) + 1) / 2
					log.Printf("Invalid chess move by %s move %s - %s in gameID %d on move %d", c.username, game.Source, game.Target, game.ID, totalMoves)
					break
				}
				Verify.AllTables[game.ID].Connection <- true
				//printBoard(game.ID)

				//checkin if there is a pending draw and if so it removes it
				if All.Games[game.ID].PendingDraw == true {
					All.Games[game.ID].PendingDraw = false

					t.Type = "cancel_draw"

					//notifiy both player that the draw offer was declined
					websocket.JSON.Send(Active.Clients[PrivateChat[t.Name]], &t)
					websocket.JSON.Send(Active.Clients[t.Name], &t)
				}

				//now switch to the other players turn
				if All.Games[game.ID].Status == "White" {
					All.Games[game.ID].Status = "Black"

					//now switch clocks
					go func() {
						All.Games[game.ID].BlackMinutes, All.Games[game.ID].BlackSeconds = StartClock(game.ID, All.Games[game.ID].BlackMinutes, All.Games[game.ID].BlackSeconds, "Black")

					}()

				} else if All.Games[game.ID].Status == "Black" {
					All.Games[game.ID].Status = "White"

					go func() {
						All.Games[game.ID].WhiteMinutes, All.Games[game.ID].WhiteSeconds = StartClock(game.ID, All.Games[game.ID].WhiteMinutes, All.Games[game.ID].WhiteSeconds, "White")

					}()

				} else {
					fmt.Println("Invalid game status, most likely game is over for ", t.Name)
					break
				}

				game.WhiteMinutes = All.Games[game.ID].WhiteMinutes
				game.WhiteSeconds = All.Games[game.ID].WhiteSeconds
				game.BlackMinutes = All.Games[game.ID].BlackMinutes
				game.BlackSeconds = All.Games[game.ID].BlackSeconds

				var move Move
				move.S = game.Source
				move.T = game.Target
				//append move to back end storage for retrieval from database later
				All.Games[game.ID].GameMoves = append(All.Games[game.ID].GameMoves, move)

				if _, ok := Active.Clients[PrivateChat[t.Name]]; !ok { //don't send move if other guy dropped connection
					//fmt.Println("connection.go 123");
					break
				}

				//sending move
				if err := websocket.JSON.Send(Active.Clients[PrivateChat[t.Name]], &game); err != nil {
					fmt.Println("error sending chess move, other player left chessConnect")
				}

			case "chat_private":

				if len(reply) > 500 {
					log.Printf("User: %s IP %s has exeeded the 500 character limit by %d byte units.\n", t.Name, c.clientIP, len(reply))
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

				//checking if other player has disconnected from the websocket
				if _, ok := Active.Clients[PrivateChat[t.Name]]; ok {

					//sending message to target person
					if err := websocket.Message.Send(Active.Clients[PrivateChat[t.Name]], reply); err != nil {
						// we could not send the message to a peer
						fmt.Println("Connection.go error 5 Could not send message to ", c.clientIP, err.Error())
					}
				}

				//sending message to self
				if err := websocket.Message.Send(Active.Clients[t.Name], reply); err != nil {
					// we could not send the message to a peer
					fmt.Println("Connection.go error 6 Could not send message to ", c.clientIP, err.Error())
				}

			case "chess_game":

				//if game match was not accepted then player name will not be stored in PrivateChat from match_accept
				if _, ok := PrivateChat[t.Name]; !ok {
					//closing socket
					break
				}

				for _, game := range All.Games {
					if game.WhitePlayer == t.Name || game.BlackPlayer == t.Name {
						game.Type = "chess_game"

						//storing the golang data structure as a string to be sent to front end
						result, _ := json.Marshal(game)

						//send to self the game info
						websocket.Message.Send(c.websocket, string(result))
					}
				}

			case "abort_game":

				var game GameMove

				json.Unmarshal(message, &game)
				//can only abort game before move 2
				if len(All.Games[game.ID].GameMoves) > 2 {
					fmt.Println("You can only abort before move 2")
					break
				}

				//closing web socket on front end for self and opponent
				websocket.Message.Send(Active.Clients[t.Name], reply)

				if _, ok := Active.Clients[PrivateChat[t.Name]]; ok { // send data if other guy is still connected
					websocket.Message.Send(Active.Clients[PrivateChat[t.Name]], reply)
				}
				Verify.AllTables[game.ID].Connection <- true
				Verify.AllTables[game.ID].gameOver <- true

				delete(All.Games, game.ID)
				delete(Verify.AllTables, game.ID)

			case "offer_draw":

				var game GameMove

				json.Unmarshal(message, &game)
				All.Games[game.ID].PendingDraw = true

				//offering draw to opponent if he is still connected
				if _, ok := Active.Clients[PrivateChat[t.Name]]; ok { // send data if other guy is still connected
					websocket.Message.Send(Active.Clients[PrivateChat[t.Name]], reply)
				}

			case "accept_draw":

				var game GameMove

				json.Unmarshal(message, &game)

				//if a draw was not offered then break out
				if All.Games[game.ID].PendingDraw == false {
					break
				}
				Verify.AllTables[game.ID].Connection <- true
				Verify.AllTables[game.ID].gameOver <- true

				All.Games[game.ID].Status = "Agreed Draw"
				//2 means the game is a draw and stored as an int in the database
				All.Games[game.ID].Result = 2

				//rate.go
				if All.Games[game.ID].Rated == "Yes"{
					ComputeRating(t.Name, game.ID, All.Games[game.ID].GameType, 0.5)
				}

				//closing web socket on front end for self and opponent
				websocket.Message.Send(Active.Clients[t.Name], reply)

				if _, ok := Active.Clients[PrivateChat[t.Name]]; ok { // send data if other guy is still connected
					websocket.Message.Send(Active.Clients[PrivateChat[t.Name]], reply)
				}

				//now store game in MySQL database
				allMoves, err := json.Marshal(All.Games[game.ID].GameMoves)
				if err != nil {
					fmt.Println("Error marshalling data to store in MySQL")
				}
				//gets length of all the moves in the game
				totalMoves := (len(All.Games[game.ID].GameMoves) + 1) / 2
				storeGame(totalMoves, allMoves, All.Games[game.ID])

				//now delete game from memory
				delete(All.Games, game.ID)
				delete(Verify.AllTables, game.ID)

			case "game_over":
				var game Fin
				var result float64

				if err := json.Unmarshal(message, &game); err != nil {
					log.Println("Just receieved a message I couldn't decode:")
					log.Println(string(reply))
					log.Println("Connection.go error 11 Exact error: " + err.Error())
					break
				}

				var checkMate bool
				var mater string
				var mated string
				if game.Status == "White" {
					checkMate = isWhiteInMate(game.ID)
					mater = All.Games[game.ID].BlackPlayer
					mated = All.Games[game.ID].WhitePlayer
				} else if game.Status == "Black" {
					mater = All.Games[game.ID].WhitePlayer
					mated = All.Games[game.ID].BlackPlayer
					checkMate = isBlackInMate(game.ID)
				} else {
					fmt.Println("Invalid game status for checking mate.")
				}
				//gets length of all the moves in the game
				totalMoves := (len(All.Games[game.ID].GameMoves) + 1) / 2
				
				if checkMate == true {
					log.Println(mater," has checkmated", mated, "in", totalMoves, "moves.")
				} else {
					log.Println("No Checkmate for player, could be bug or cheat attempt by", mater, "on move", totalMoves, "against", mated)
					break
				}

				if game.Status == "White" { //then white was checkmated
					All.Games[game.ID].Status = "White is checkmated"
					result = 0

				} else if game.Status == "Black" {
					All.Games[game.ID].Status = "Black is checkmated"
					result = 1.0

				} else {
					fmt.Println("Invalid color checkmate")
				}
				Verify.AllTables[game.ID].Connection <- true
				Verify.AllTables[game.ID].gameOver <- true
				//now store game in MySQL database
				allMoves, err := json.Marshal(All.Games[game.ID].GameMoves)
				if err != nil {
					fmt.Println("Error marshalling data to store in MySQL")
				}

				storeGame(totalMoves, allMoves, All.Games[game.ID])

				//update ratings
				if All.Games[game.ID].Rated == "Yes"{
					ComputeRating(t.Name, game.ID, All.Games[game.ID].GameType, result)
				}

				//notifying players game is over
				if err := websocket.Message.Send(Active.Clients[t.Name], reply); err != nil {
					fmt.Println("error gameover 1 connection.go error is ", err)
				}
				if _, ok := Active.Clients[PrivateChat[t.Name]]; ok { // send data if other guy is still connected
					if err := websocket.Message.Send(Active.Clients[PrivateChat[t.Name]], reply); err != nil {
						fmt.Println("connection.go gameover 2 error is ", err)
					}
				}

				delete(All.Games, game.ID)
				delete(Verify.AllTables, game.ID)

			case "resign":

				var game GameMove
				json.Unmarshal(message, &game)
				var result float64

				if t.Name == All.Games[game.ID].WhitePlayer {
					All.Games[game.ID].Status = "White Resigned"
					result = 0.0
					All.Games[game.ID].Result = 0

				} else if t.Name == All.Games[game.ID].BlackPlayer {
					All.Games[game.ID].Status = "Black Resigned"
					result = 1.0
					All.Games[game.ID].Result = 1

				} else {
					fmt.Println("Invalid resign, no player found.")
					break
				}

				Verify.AllTables[game.ID].Connection <- true
				Verify.AllTables[game.ID].gameOver <- true

				//now store game in MySQL database
				allMoves, err := json.Marshal(All.Games[game.ID].GameMoves)
				if err != nil {
					fmt.Println("Error marshalling data to store in MySQL")
				}
				//gets length of all the moves in the game
				totalMoves := (len(All.Games[game.ID].GameMoves) + 1) / 2

				storeGame(totalMoves, allMoves, All.Games[game.ID])

				//rate.go
				if All.Games[game.ID].Rated == "Yes"{
					ComputeRating(t.Name, game.ID, All.Games[game.ID].GameType, result)
				}

				//letting both players know that a resignation occured
				if _, ok := Active.Clients[PrivateChat[t.Name]]; ok { // send data if other guy is still connected
					websocket.Message.Send(Active.Clients[PrivateChat[t.Name]], reply)
				}

				websocket.Message.Send(Active.Clients[t.Name], reply)

				delete(All.Games, game.ID)
				delete(Verify.AllTables, game.ID)
			case "rematch":

				var match SeekMatch
				if err := json.Unmarshal(message, &match); err != nil {
					fmt.Println("Just receieved a message I couldn't decode:")
					fmt.Println(string(reply))
					fmt.Println("Exact error: " + err.Error())
					break
				}
				//check length of name to make sure its 3-12 characters long, prevent hack abuse
				if len(match.Opponent) < 3 || len(match.Opponent) > 12 {
					fmt.Println("Username is too long or too short")
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
					fmt.Println("Cannot get rating connection.go private_match")
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
					if err := websocket.JSON.Send(Active.Clients[t.Name], &t); err != nil {
						// we could not send the message to a peer
						log.Println("match connection.go Could not send message to ", c.clientIP, err.Error())
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

				if _, ok := PrivateChat[match.Opponent]; ok {
					log.Println("rematch 1 connection.go Player already has a game. ")

					t.Type = "rematch"
					if err := websocket.JSON.Send(Active.Clients[match.Opponent], &t); err != nil {
						// we could not send the message to a peer
						log.Println("rematch 2 connection.go Could not send message to ", c.clientIP, err.Error())
					}
				}

			case "accept_rematch":

				var match SeekMatch
				var game ChessGame
				if err := json.Unmarshal(message, &match); err != nil {
					log.Println("Just receieved a message I couldn't decode:")
					log.Println(string(reply))
					log.Println("Connection.go error 11 Exact error: " + err.Error())
					break
				}
				//isPlayerInGame function is located in socket.go
				if isPlayerInGame(match.Name, match.Opponent) {
					fmt.Println("connection.go accept rematch 12")
					break
				}

				//checking to make sure both player's rating is in range, used as a backend rating check
				errMessage, bullet, blitz, standard := GetRating(match.Name)
				if errMessage != "" {
					fmt.Println("Cannot get rating connection.go accept_match")
					break
				}

				game.Type = "chess_game"

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
				game.BlackMinutes = Pending.Matches[match.MatchID].TimeControl
				game.BlackSeconds = 0
				game.PendingDraw = false
				game.Rated = Pending.Matches[match.MatchID].Rated

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

				//intitalizes all the variables of the game
				initGame(game.ID)

				startGame, _ := json.Marshal(game)

				//starting game for both players
				if err := websocket.Message.Send(Active.Clients[game.WhitePlayer], string(startGame)); err != nil {
					fmt.Println("error accept_rematch 1 is ", err)

				}
				if err := websocket.Message.Send(Active.Clients[game.BlackPlayer], string(startGame)); err != nil {
					fmt.Println("error accept_rematch 2 is ", err)

				}

				//starting white's clock first, this goroutine will keep track of both players clock for this game
				go setClocks(game.ID, t.Name)

			case "draw_game":
				var game GameMove

				json.Unmarshal(message, &game)

				//checking to see if the side whose turn it is to move is in stalemate
				if Verify.AllTables[game.ID].whiteTurn == true {
					if isWhiteStaleMate(game.ID) == true || noMaterial(game.ID) == true || threeRep(game.ID) == true || fiftyMoves(game.ID) == true {
						fmt.Println("forced draw_game connection.go success 1")
					} else {
						break
					}
				} else {

					if isBlackStaleMate(game.ID) == true || noMaterial(game.ID) == true || threeRep(game.ID) == true || fiftyMoves(game.ID) == true {
						fmt.Println("forced draw_game connection.go success 2")
					} else {
						break
					}
				}

				Verify.AllTables[game.ID].Connection <- true
				Verify.AllTables[game.ID].gameOver <- true

				All.Games[game.ID].Status = "Forced Draw"
				//2 means the game is a draw and stored as an int in the database
				All.Games[game.ID].Result = 2

				//rate.go
				if All.Games[game.ID].Rated == "Yes"{
					ComputeRating(t.Name, game.ID, All.Games[game.ID].GameType, 0.5)
				}

				//closing web socket on front end for self and opponent
				websocket.Message.Send(Active.Clients[t.Name], reply)

				if _, ok := Active.Clients[PrivateChat[t.Name]]; ok { // send data if other guy is still connected
					websocket.Message.Send(Active.Clients[PrivateChat[t.Name]], reply)
				}

				//now store game in MySQL database
				allMoves, err := json.Marshal(All.Games[game.ID].GameMoves)
				if err != nil {
					fmt.Println("Error marshalling data to store in MySQL")
				}
				//gets length of all the moves in the game
				totalMoves := (len(All.Games[game.ID].GameMoves) + 1) / 2
				storeGame(totalMoves, allMoves, All.Games[game.ID])

				//now delete game from memory
				delete(All.Games, game.ID)
				delete(Verify.AllTables, game.ID)

			default:
				fmt.Println("I'm not familiar with type " + t.Type)
			}
		} else {
			log.Printf("IP %s Invalid websocket authentication in chess room.\n", c.clientIP)
			return
		}
	}
}
