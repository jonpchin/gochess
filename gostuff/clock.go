package gostuff

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

//starting white's clock first, this goroutine will keep track of both players clock for this game
func (game *ChessGame) setClocks(name string) {

	var result float64
	table := Verify.AllTables[game.ID]
	go func() {
		game.WhiteMinutes, game.WhiteSeconds = table.startClock(game.ID, game.WhiteMinutes, game.WhiteSeconds, "White")
	}()

	//checks whether or not clock has timed out
	for {
		select {
		case <-table.whiteTimeOut:

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

			return
		case <-table.blackTimeOut:

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

			return
		case <-table.gameOver:
			//fmt.Println("Game is over but clocks have not ran out so break out of clocks")
			return
		}
	}
}

//returns the remaining time of players's clock
func (table *Table) startClock(gameID int, minutes int, seconds int, color string) (int, int) {

	timerChan := time.NewTicker(time.Second).C

	clock := (minutes * 60) + seconds
	if clock <= 0 {
		return 0, 0
	}

	table.Connection = make(chan bool)
	chessgame := All.Games[gameID]

	for {
		select {
		case <-table.Connection:

			remainingMinutes := clock / 60
			remainingSeconds := clock % 60
			//fmt.Printf("Clock here is %d %d color is %s\n", remainingMinutes, remainingSeconds, color)
			return remainingMinutes, remainingSeconds

		case <-timerChan:

			clock--
			remainingMinutes := clock / 60
			remainingSeconds := clock % 60
			if color == "White" {
				chessgame.WhiteMinutes = remainingMinutes
				chessgame.WhiteSeconds = remainingSeconds
			} else {
				chessgame.BlackMinutes = remainingMinutes
				chessgame.BlackSeconds = remainingSeconds
			}
			if clock <= 0 && color == "White" {
				table.whiteTimeOut <- true
				return 0, 0
			} else if clock <= 0 && color == "Black" {
				table.blackTimeOut <- true
				return 0, 0
			}
		}
	}
}
