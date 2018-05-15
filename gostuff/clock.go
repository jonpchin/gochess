package gostuff

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

// keep track of both sides of the clock, when function exits the game is over
func (table *Table) StartClock(gameID int, minutes int, seconds int, name string) {

	timerChan := time.NewTicker(time.Second).C

	whiteClock := (minutes * 60) + seconds
	blackClock := (minutes * 60) + seconds

	startTime := whiteClock

	chessgame := All.Games[gameID]

	for {
		select {

		case <-table.gameOver:
			return

		case <-timerChan:
			if table.whiteToMove {

				whiteClock--
				remainingMinutes := whiteClock / 60
				remainingSeconds := whiteClock % 60

				chessgame.WhiteMinutes = remainingMinutes
				chessgame.WhiteSeconds = remainingSeconds
				//fmt.Printf("White clock %d %d \n", remainingMinutes, remainingSeconds)

				if whiteClock <= 0 {
					chessgame.whiteTimeout(name)
					return
				}
			} else {

				blackClock--
				remainingMinutes := blackClock / 60
				remainingSeconds := blackClock % 60

				chessgame.BlackMinutes = remainingMinutes
				chessgame.BlackSeconds = remainingSeconds
				//fmt.Printf("Black clock %d %d \n", remainingMinutes, remainingSeconds)

				if blackClock <= 0 {
					chessgame.blackTimeout(name)
					return
				}
			}
		case <-table.resetBlackTime:
			blackClock = startTime

		case <-table.resetWhiteTime:
			whiteClock = startTime
		}
	}
}

// when white's clock runs out save database information and cleanup game
func (game *ChessGame) whiteTimeout(name string) {

	var result float64
	result = 0.0
	//update ratings
	if game.Rated == "Yes" {
		ComputeRating(name, game.ID, game.GameType, result)
	}

	//Black won as white ran out of time
	game.Status = "Black won on time"
	game.Type = "game_over"
	game.Result = 0

	//now store game in MySQL database
	allMoves, err := json.Marshal(game.GameMoves)
	if err != nil {
		fmt.Println("Error marshalling data to store in MySQL")
	}
	//gets length of all the moves in the game
	totalMoves := (len(game.GameMoves) + 1) / 2
	//save game to database before deleting it from memory
	storeGame(totalMoves, allMoves, game)

	//notifiy both player that black won on time
	if _, ok := Active.Clients[PrivateChat[name]]; ok { // send data if other guy is still connected
		websocket.JSON.Send(Active.Clients[PrivateChat[name]], game)
	}

	if _, ok := Active.Clients[name]; ok { // send data if other guy is still connected
		websocket.JSON.Send(Active.Clients[name], game)
	}

	delete(All.Games, game.ID)
	delete(Verify.AllTables, game.ID)
}

// when black's clock runs out save database information and cleanup game
func (game *ChessGame) blackTimeout(name string) {

	var result float64

	//White won as black ran out of time
	game.Status = "White won on time"
	result = 1.0
	game.Type = "game_over"
	game.Result = 1

	//now store game in MySQL database
	allMoves, err := json.Marshal(game.GameMoves)
	if err != nil {
		fmt.Println("Error marshalling data to store in MySQL")
	}
	//gets length of all the moves in the game
	totalMoves := (len(game.GameMoves) + 1) / 2
	//save game to database before deleting it from memory
	storeGame(totalMoves, allMoves, game)

	//update ratings
	if game.Rated == "Yes" {
		ComputeRating(name, game.ID, game.GameType, result)
	}

	//notifiy both player that white won on time
	if _, ok := Active.Clients[PrivateChat[name]]; ok { // send data if other guy is still connected
		websocket.JSON.Send(Active.Clients[PrivateChat[name]], game)
	}

	if _, ok := Active.Clients[name]; ok { // send data if other guy is still connected
		websocket.JSON.Send(Active.Clients[name], game)
	}

	delete(All.Games, game.ID)
	delete(Verify.AllTables, game.ID)
}
