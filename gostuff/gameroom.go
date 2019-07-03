package gostuff

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

// Manages web sockets for the game room
func (c *Connection) ChessConnect() {

	defer exitGame(c.username) //remove user when they disconnect from socket
	counter := 0
	start := time.Now()

	logFile, _ := os.OpenFile("logs/chat.txt", os.O_APPEND|os.O_WRONLY, 0666)

	defer logFile.Close()
	log := log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	//PrintMemoryStats()
	for {
		var reply string

		if err := websocket.Message.Receive(c.websocket, &reply); err != nil {
			//fmt.Println("A user has drop web socket connection ", err)
			break
		}

		var t MessageType
		message := []byte(reply)
		if err := json.Unmarshal(message, &t); err != nil {
			log.Println("Just receieved a message I couldn't decode:", string(reply), err)
			break
		}
		if c.username == t.Name {
			switch t.Type {

			case "send_move":

				var game GameMove

				err := json.Unmarshal(message, &game)
				if err != nil {
					log.Println("Just receieved a message I couldn't decode:", string(message), err)
					break
				}

				//this can also be triggered when a player starts moving pieces on the board alone
				//also prevents a move from being sent if a game hasn't started
				chessgame, ok := All.Games[game.ID]
				if !ok {
					break
				}

				var white = chessgame.WhitePlayer
				var black = chessgame.BlackPlayer

				// Check if the color that is moving is the correct player
				if chessgame.Status == "White" && white != c.username {
					break
				}
				if chessgame.Status == "Black" && black != c.username {
					break
				}

				// spectators should not be able to make moves for the two chess players
				if t.Name != white && t.Name != black {
					fmt.Println(t.Name, "tried to cheat by making a move as a spectator")
					break
				}

				var result bool
				//check if its correct players turn and if move is valid before sending
				result = ChessVerify(game.Source, game.Target, game.Promotion, game.ID)

				if result == false {
					totalMoves := (len(chessgame.GameMoves) + 1) / 2
					log.Printf("Invalid chess move by %s move %s - %s in gameID %d on move %d", c.username, game.Source, game.Target, game.ID, totalMoves)
					break
				}

				chessgame.Validator.MoveStr(game.Source + game.Target + game.Promotion)

				if game.Fen == "" {
					game.Fen = chessgame.Validator.Position().String()
				}

				table := Verify.AllTables[game.ID]
				//printBoard(game.ID)

				//checkin if there is a pending draw and if so it removes it
				if chessgame.PendingDraw {
					chessgame.PendingDraw = false

					t.Type = "cancel_draw"

					//notifiy both player that the draw offer was declined
					websocket.JSON.Send(Active.Clients[PrivateChat[t.Name]], &t)
					websocket.JSON.Send(Active.Clients[t.Name], &t)
				}

				//now switch to the other players turn
				if chessgame.Status == "White" {
					chessgame.Status = "Black"
					if chessgame.GameType == "correspondence" {
						table.resetWhiteTime <- true
					}
				} else if chessgame.Status == "Black" {
					chessgame.Status = "White"
					if chessgame.GameType == "correspondence" {
						table.resetBlackTime <- true
					}
				} else {
					log.Println("Invalid game status, most likely game is over for ", t.Name)
					break
				}

				var move Move
				move.S = game.Source
				move.T = game.Target
				move.P = game.Promotion
				move.Fen = game.Fen

				//append move to back end storage for retrieval from database later
				chessgame.GameMoves = append(chessgame.GameMoves, move)

				for _, name := range table.observe.Names {
					if _, ok := Active.Clients[name]; ok {
						if err := websocket.JSON.Send(Active.Clients[name], &game); err != nil {
							log.Println("error sending chess move to", name)
						}
					} else if name != white && name != black { //remove spectator if they are no longer viewing game
						table.observe.Lock()
						table.observe.Names = removeViewer(name, game.ID)
						table.observe.Unlock()
					}
				}

				checkGameOver(t.Name, game.ID, game.Fen, chessgame.Status)

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

				gameID, exist := GetGameID(t.Name)

				// prevent invalid map access which can cause race conditions
				if _, ok := All.Games[gameID]; !ok {
					// spectators should not be able to send chat messages
					break
				}
				var white = All.Games[gameID].WhitePlayer
				var black = All.Games[gameID].BlackPlayer

				if t.Name != white && t.Name != black {
					log.Println(t.Name, " tried to send a chat message while spectating")
					break
				}

				if exist {
					table := Verify.AllTables[gameID]
					for _, name := range table.observe.Names {
						if _, ok := Active.Clients[name]; ok {
							if err := websocket.Message.Send(Active.Clients[name], reply); err != nil {
								// we could not send the message to a peer
								log.Println("Could not send message to ", name, err.Error())
							}
						} else if name != white && name != black {
							//remove spectator if they are no longer viewing game
							table.observe.Lock()
							table.observe.Names = removeViewer(name, gameID)
							table.observe.Unlock()
						}
					}
					// if game does not exist but user is still in chess room allow
					// private chat only between the two chess players
				} else if isPlayerInChess(t.Name) {
					//checking if other player has disconnected from the websocket
					if _, ok := Active.Clients[PrivateChat[t.Name]]; ok {

						//sending message to target person
						if err := websocket.Message.Send(Active.Clients[PrivateChat[t.Name]], reply); err != nil {
							// we could not send the message to a peer
							log.Println("Could not send message to ", PrivateChat[t.Name], err.Error())
						}
					}

					//sending message to self
					if err := websocket.Message.Send(Active.Clients[t.Name], reply); err != nil {
						// we could not send the message to a peer
						log.Println("Could not send message to ", t.Name, err.Error())
					}
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

						//send to self the game info
						if err := websocket.JSON.Send(c.websocket, &game); err != nil {
							log.Println(err)
						}
					}
				}

			case "abort_game":

				var game GameMove

				if err := json.Unmarshal(message, &game); err != nil {
					log.Println("Failed to unmarshal", err)
				}

				chessgame := All.Games[game.ID]

				//can only abort game before move 2
				if len(chessgame.GameMoves) > 2 {
					log.Println("You can only abort before move 2")
					break
				}

				// spectators should not be able to abort game
				if t.Name != chessgame.WhitePlayer && t.Name != chessgame.BlackPlayer {
					log.Println(t.Name, " tried to abort game while spectating")
					return
				}

				table := Verify.AllTables[game.ID]
				chessgame.Type = "abort_game"
				chessgame.Status = "Game aborted by " + t.Name
				for _, name := range table.observe.Names {
					if _, ok := Active.Clients[name]; ok {
						if err := websocket.JSON.Send(Active.Clients[name], &chessgame); err != nil {
							log.Println("error sending abort message", err)
						}
					}
				}

				table.gameOver <- true

				delete(All.Games, game.ID)
				delete(Verify.AllTables, game.ID)

			case "update_spectate":

				var spectate SpectateGame

				err := json.Unmarshal(message, &spectate)
				if err != nil {
					log.Println("Just receieved a message I couldn't decode:", string(message), err)
					break
				}
				if spectate.Spectate == "Yes" {
					All.Games[spectate.ID].Spectate = true
				} else {
					All.Games[spectate.ID].Spectate = false
				}

			case "spectate_game":

				var spectate SpectateGame

				if err := json.Unmarshal(message, &spectate); err != nil {
					log.Println("Just receieved a message I couldn't decode:", string(message), err)
					break
				}

				// search table of games for the ID in spectate and return the data back
				// to the spectator
				if table, ok := Verify.AllTables[spectate.ID]; ok {
					// only send data to spectator if spectator mode is turned on
					if All.Games[spectate.ID].Spectate {
						// register spectator to observers list
						table.observe.Lock()
						table.observe.Names = append(table.observe.Names, t.Name)
						table.observe.Unlock()
						// send data to spectator
						if err := websocket.JSON.Send(Active.Clients[t.Name], All.Games[spectate.ID]); err != nil {
							log.Println(err)
						}

						//send a message to everyone saying spectator has entered room
						for _, name := range table.observe.Names {
							if _, ok := Active.Clients[name]; ok {
								if err := websocket.Message.Send(Active.Clients[name], reply); err != nil {
									log.Println("error sending abort message", err)
								}
							}
						}
					}
				} else {
					log.Println(t.Name, " tried viewing a game that doesn't exist.")
				}

			case "offer_draw":

				var game GameMove

				if err := json.Unmarshal(message, &game); err != nil {
					log.Println("error unmarshalling data", err)
					return
				}

				chessgame := All.Games[game.ID]

				// spectators should not be able to offer draw while spectating
				if t.Name != chessgame.WhitePlayer && t.Name != chessgame.BlackPlayer {
					log.Println(t.Name, " tried to offer draw while spectating")
					return
				}
				chessgame.PendingDraw = true

				//offering draw to opponent if he is still connected
				if _, ok := Active.Clients[PrivateChat[t.Name]]; ok { // send data if other guy is still connected
					if err := websocket.Message.Send(Active.Clients[PrivateChat[t.Name]], reply); err != nil {
						log.Println(err)
					}
				}

			case "accept_draw":

				var game GameMove

				if err := json.Unmarshal(message, &game); err != nil {
					log.Println("error in unmarshalling data")
				}

				//make sure key exist in map before accessing it
				chessgame, ok := All.Games[game.ID]
				if !ok {
					break
				}

				// spectators should not be able to accept draw while spectating
				if t.Name != chessgame.WhitePlayer && t.Name != chessgame.BlackPlayer {
					log.Println(t.Name, " tried to accept draw while spectating")
					break
				}
				//if a draw was not offered then break out
				if chessgame.PendingDraw == false {
					break
				}

				table, ok := Verify.AllTables[game.ID]
				if !ok {
					break
				}

				table.gameOver <- true

				chessgame.Status = "Agreed Draw"
				//2 means the game is a draw and stored as an int in the database
				chessgame.Result = 2
				chessgame.Type = "accept_draw"

				for _, name := range table.observe.Names {
					if _, ok := Active.Clients[name]; ok {
						if err := websocket.JSON.Send(Active.Clients[name], &chessgame); err != nil {
							log.Println(err)
						}
					}
				}

				//rate.go
				if chessgame.Rated == "Yes" {
					ComputeRating(t.Name, game.ID, chessgame.GameType, 0.5)
				}

				wrapUpGame(game.ID)

			case "resign":

				var game GameMove
				if err := json.Unmarshal(message, &game); err != nil {
					log.Println(err)
				}
				var result float64
				chessgame := All.Games[game.ID]
				chessgame.Type = "resign"

				if t.Name == chessgame.WhitePlayer {
					chessgame.Status = "White Resigned"
					result = 0.0
					chessgame.Result = 0

				} else if t.Name == chessgame.BlackPlayer {
					chessgame.Status = "Black Resigned"
					result = 1.0
					chessgame.Result = 1

				} else {
					fmt.Println("Invalid resign, no player found.")
					break
				}

				table := Verify.AllTables[game.ID]
				table.gameOver <- true

				//letting both players and spectators know that a resignation occured
				for _, name := range table.observe.Names {
					if _, ok := Active.Clients[name]; ok {
						if err := websocket.JSON.Send(Active.Clients[name], &chessgame); err != nil {
							log.Println(err)
						}
					}
				}

				//rate.go
				if chessgame.Rated == "Yes" {
					ComputeRating(t.Name, game.ID, chessgame.GameType, result)
				}
				wrapUpGame(game.ID)

			case "rematch":

				var match SeekMatch
				if err := json.Unmarshal(message, &match); err != nil {
					log.Println("Just receieved a message I couldn't decode:", string(message), err)
					break
				}
				//check length of name to make sure its 3-12 characters long, prevent hack abuse
				if len(match.Opponent) < 3 || len(match.Opponent) > 12 {
					fmt.Println("Username is too long or too short")
					break
				}

				if match.assignMatchRatingType() == false {
					break
				}

				//check to make sure player only has a max of three matches seeks pending, used to prevent flood match seeking
				if countMatches(c.username) >= 3 {
					t.Type = "maxThree"
					if err := websocket.JSON.Send(Active.Clients[t.Name], &t); err != nil {
						// we could not send the message to a peer
						log.Println("Could not send message to ", t.Name, err)
					}
					break //notify user that only three matches pending max are allowed
				}

				var start int = 0
				for {
					if _, ok := Pending.Matches[start]; ok {
						start++
					} else {
						break
					}
				}

				match.MatchID = start
				//used in backend to keep track of all pending seeks waiting for a player to accept
				Pending.Matches[start] = &match

				if _, ok := Active.Clients[match.Opponent]; ok {
					t.Type = "rematch"
					if err := websocket.JSON.Send(Active.Clients[match.Opponent], &t); err != nil {
						// we could not send the message to a peer
						log.Println("Could not send message to ", match.Opponent, err)
					}
				}

			case "accept_rematch":

				var match SeekMatch
				if err := json.Unmarshal(message, &match); err != nil {
					log.Println("Just receieved a message I couldn't decode:", string(message), err)
					break
				}

				startPendingMatch(match.Name, match.MatchID)

			default:
				fmt.Println("I'm not familiar with type " + t.Type)
			}
		} else {
			log.Printf("IP %s Invalid websocket authentication in chess room.\n", c.clientIP)
			return
		}
	}
}

// Detects if checkmate, stalemate, 50 move rule, 3 repetition or insufficent material to checkmate
func checkGameOver(playerName string, gameID int, gameFen string, gameStatus string) {

	var mater string
	var mated string

	chessgame := All.Games[gameID]
	table := Verify.AllTables[gameID]

	check, mate := isCheckMate(gameFen)

	isDraw := false

	if gameStatus == "White" {
		mater = chessgame.BlackPlayer
		mated = chessgame.WhitePlayer

	} else if gameStatus == "Black" {
		mater = chessgame.WhitePlayer
		mated = chessgame.BlackPlayer

	} else { //this should never happen, if it does most likely caused by tampering or its a bug
		fmt.Println("Invalid game status for checking mate.")
		return
	}
	//gets length of all the moves in the game
	totalMoves := (len(chessgame.GameMoves) + 1) / 2

	if check && mate {
		log.Println(mater, "has checkmated", mated, "in", totalMoves, "moves.")
		var result float64

		if gameStatus == "White" { //then white was checkmated
			chessgame.Status = "White is checkmated"
			result = 0

		} else { // then its black that was checkmated
			chessgame.Status = "Black is checkmated"
			result = 1.0
		}

		table.gameOver <- true
		chessgame.Type = "game_over"

		//notifying both players and spectators game is over
		for _, name := range table.observe.Names {
			if _, ok := Active.Clients[name]; ok {
				if err := websocket.JSON.Send(Active.Clients[name], &chessgame); err != nil {
					fmt.Println(err)
				}
			}
		}

		//update ratings
		if chessgame.Rated == "Yes" {
			ComputeRating(playerName, gameID, chessgame.GameType, result)
		}

		wrapUpGame(gameID)

	} else if mate {
		log.Println(mater, "has stalemated", mated, "in", totalMoves, "moves.")
		isDraw = true
		chessgame.Status = "Stalemate"
	} else if table.noMaterial() {
		log.Println(mater, "does not have sufficient mating material to mate", mated, "in", totalMoves, "moves.")
		isDraw = true
		chessgame.Status = "Insufficent mating material"
	} else if chessgame.threeRep() {
		log.Println(mater, "has triggered three reptition draw", mated, "in", totalMoves, "moves.")
		isDraw = true
		chessgame.Status = "Three repetition draw"
	} else if table.fiftyMoves(gameID) {
		log.Println(mater, "has triggered fifty move rule", mated, "in", totalMoves, "moves.")
		isDraw = true
		chessgame.Status = "Fifty move rule draw"
	}

	if isDraw {
		table.gameOver <- true
		//2 means the game is a draw and stored as an int in the database
		chessgame.Result = 2

		//rate.go
		if chessgame.Rated == "Yes" {
			ComputeRating(playerName, gameID, chessgame.GameType, 0.5)
		}

		chessgame.Type = "draw_game"

		//closing web socket on front end for self and opponent
		for _, name := range table.observe.Names {
			if client, ok := Active.Clients[name]; ok {
				if err := websocket.JSON.Send(client, &chessgame); err != nil {
					log.Println("Can't send message for draw game isGameOver()", err)
				}
			}
		}

		wrapUpGame(gameID)
	}
}

// Cleanup function to store game in database and delete from memory
func wrapUpGame(id int) {

	chessgame := All.Games[id]
	//now store game in MySQL database
	allMoves, err := json.Marshal(chessgame.GameMoves)
	if err != nil {
		fmt.Println("Error marshalling data to store in MySQL")
	}
	//gets length of all the moves in the game
	totalMoves := (len(chessgame.GameMoves) + 1) / 2
	storeGame(totalMoves, allMoves, chessgame)

	//now delete game from memory
	delete(All.Games, id)
	delete(Verify.AllTables, id)
}
