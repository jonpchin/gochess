package gostuff

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
				for _, name := range table.observe.Names {
					if _, ok := Active.Clients[name]; ok {
						if err := websocket.Message.Send(Active.Clients[name], reply); err != nil {
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

				for _, name := range table.observe.Names {
					if _, ok := Active.Clients[name]; ok {
						if err := websocket.Message.Send(Active.Clients[name], reply); err != nil {
							log.Println(err)
						}
					}
				}

				//rate.go
				if chessgame.Rated == "Yes" {
					ComputeRating(t.Name, game.ID, chessgame.GameType, 0.5)
				}

				wrapUpGame(game.ID)

			case "game_over":

				var game Fin

				if err := json.Unmarshal(message, &game); err != nil {
					log.Println("Just receieved a message I couldn't decode:", string(message), err)
					break
				}

				var checkMate bool
				var mater string
				var mated string

				chessgame := All.Games[game.ID]
				table := Verify.AllTables[game.ID]

				check, mate := isCheckMate(game.Fen)
				if check && mate {
					log.Println("Verified its mate")
				} else if check {
					log.Println("It's only a check, no mate.")
				} else if mate {
					log.Println("It is a stalemate.")
				} else {
					log.Println("It's not mate, check or stalemate.")
				}

				if game.Status == "White" {
					checkMate = table.isWhiteInMate()
					mater = chessgame.BlackPlayer
					mated = chessgame.WhitePlayer

				} else if game.Status == "Black" {
					mater = chessgame.WhitePlayer
					mated = chessgame.BlackPlayer
					checkMate = table.isBlackInMate()

				} else { //this should never happen, if it does most likely caused by tampering or its a bug
					fmt.Println("Invalid game status for checking mate.")
					break
				}
				//gets length of all the moves in the game
				totalMoves := (len(chessgame.GameMoves) + 1) / 2

				if checkMate {
					log.Println(mater, "has checkmated", mated, "in", totalMoves, "moves.")
				} else {
					log.Println("No Checkmate for player, could be bug or cheat attempt by", mater, "on move", totalMoves, "against", mated)
					if !mate {
						fmt.Println("Both mates cannot be verified.")
						break
					} else {
						log.Println("Conflincting checkmates", checkMate, check, mate)
					}
				}

				var result float64

				if game.Status == "White" { //then white was checkmated
					chessgame.Status = "White is checkmated"
					result = 0

				} else { // then its black that was checkmated
					chessgame.Status = "Black is checkmated"
					result = 1.0
				}

				table.gameOver <- true

				//notifying both players and spectators game is over
				for _, name := range table.observe.Names {
					if _, ok := Active.Clients[name]; ok {
						if err := websocket.Message.Send(Active.Clients[name], reply); err != nil {
							log.Println(err)
						}
					}
				}

				//update ratings
				if chessgame.Rated == "Yes" {
					ComputeRating(t.Name, game.ID, chessgame.GameType, result)
				}

				wrapUpGame(game.ID)

			case "resign":

				var game GameMove
				if err := json.Unmarshal(message, &game); err != nil {
					log.Println(err)
				}
				var result float64
				chessgame := All.Games[game.ID]

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
						if err := websocket.Message.Send(Active.Clients[name], reply); err != nil {
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

				if _, ok := PrivateChat[match.Opponent]; ok {
					t.Type = "rematch"
					if err := websocket.JSON.Send(Active.Clients[match.Opponent], &t); err != nil {
						// we could not send the message to a peer
						log.Println("Could not send message to ", match.Opponent, err)
					}
				}

			case "accept_rematch":

				var match SeekMatch
				var game ChessGame
				if err := json.Unmarshal(message, &match); err != nil {
					log.Println("Just receieved a message I couldn't decode:", string(message), err)
					break
				}
				//isPlayersInGame function is located in socket.go
				if isPlayersInGame(match.Name, match.Opponent) {
					log.Println("Player is already in a game")
					break
				}

				//checking to make sure both player's rating is in range, used as a backend rating check
				errMessage, bullet, blitz, standard, correspondence := GetRating(match.Name)
				if errMessage != "" {
					log.Println("Cannot get rating")
					break
				}

				game.Type = "chess_game"

				//bullet, blitz or standard game type
				game.GameType = Pending.Matches[match.MatchID].GameType

				//seting up the game info such as white/black player, time control, etc
				rand.Seed(time.Now().UnixNano())

				//randomly selects both players to be white or black
				if rand.Intn(2) == 0 {
					game.WhitePlayer = match.Name
					if game.GameType == "bullet" {
						game.WhiteRating = bullet

					} else if game.GameType == "blitz" {
						game.WhiteRating = blitz

					} else if game.GameType == "standard" {
						game.WhiteRating = standard

					} else {
						game.WhiteRating = correspondence
					}

					game.BlackRating = Pending.Matches[match.MatchID].Rating
					game.BlackPlayer = Pending.Matches[match.MatchID].Name

				} else {
					game.WhitePlayer = Pending.Matches[match.MatchID].Name
					if game.GameType == "bullet" {
						game.BlackRating = bullet

					} else if game.GameType == "blitz" {
						game.BlackRating = blitz

					} else if game.GameType == "standard" {
						game.BlackRating = standard

					} else {
						game.BlackRating = correspondence
					}

					game.WhiteRating = Pending.Matches[match.MatchID].Rating
					game.BlackPlayer = match.Name
				}
				//White for white to move or Black for black to move, white won, black won, stalemate or draw.
				game.Status = "White"

				//no moves yet so nill/null
				game.GameMoves = nil
				game.TimeControl = Pending.Matches[match.MatchID].TimeControl
				game.WhiteMinutes = Pending.Matches[match.MatchID].TimeControl
				game.WhiteSeconds = 0
				game.BlackMinutes = Pending.Matches[match.MatchID].TimeControl
				game.BlackSeconds = 0
				game.PendingDraw = false
				game.Rated = Pending.Matches[match.MatchID].Rated

				var start int = 0
				for {
					if _, ok := All.Games[start]; ok {
						start++

					} else {
						break
					}
				}

				game.ID = start
				//used in backend to keep track of all pending games waiting for a player to accept
				All.Games[start] = &game

				//no longer need all the pending matches as game will be started
				for key, value := range Pending.Matches {
					//deletes all pending matches for either players
					if value.Name == game.WhitePlayer || value.Name == game.BlackPlayer {
						delete(Pending.Matches, key)
					}
				}

				//intitalizes all the variables of the game
				InitGame(game.ID, game.WhitePlayer, game.BlackPlayer)

				//starting game for both players, this does NOT include spectators
				for _, name := range Verify.AllTables[game.ID].observe.Names {
					if client, ok := Active.Clients[name]; ok {
						if err := websocket.JSON.Send(client, &game); err != nil {
							log.Println(err)
						}
					}
				}

				//starting white's clock first, this goroutine will keep track of both players clock for this game
				table := Verify.AllTables[game.ID]
				go table.StartClock(game.ID, game.WhiteMinutes, game.WhiteSeconds, t.Name)

			case "draw_game":

				var game GameMove
				if err := json.Unmarshal(message, &game); err != nil {
					log.Println(err)
				}

				table := Verify.AllTables[game.ID]
				chessgame := All.Games[game.ID]

				check, mate := isCheckMate(game.Fen)
				if mate && !check {
					log.Println("It is a stalemate.")
				} else {
					//checking to see if the side whose turn it is to move is in stalemate
					if table.whiteTurn {
						if table.isWhiteStaleMate() || table.noMaterial() ||
							chessgame.threeRep() || table.fiftyMoves(game.ID) {
							log.Println("forced draw_game")
						} else {
							break
						}
					} else {

						if table.isBlackStaleMate() || table.noMaterial() ||
							chessgame.threeRep() || table.fiftyMoves(game.ID) {
							log.Println("forced draw_game")
						} else {
							break
						}
					}
				}

				table.gameOver <- true
				chessgame.Status = "Forced Draw"
				//2 means the game is a draw and stored as an int in the database
				chessgame.Result = 2

				//rate.go
				if chessgame.Rated == "Yes" {
					ComputeRating(t.Name, game.ID, chessgame.GameType, 0.5)
				}

				//closing web socket on front end for self and opponent
				for _, name := range table.observe.Names {
					if client, ok := Active.Clients[name]; ok {
						if err := websocket.Message.Send(client, reply); err != nil {
							log.Println(err)
						}
					}
				}
				wrapUpGame(game.ID)

			default:
				fmt.Println("I'm not familiar with type " + t.Type)
			}
		} else {
			log.Printf("IP %s Invalid websocket authentication in chess room.\n", c.clientIP)
			return
		}
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
