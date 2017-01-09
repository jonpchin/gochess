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
	go func() {

		game.WhiteMinutes, game.WhiteSeconds, game.WhiteMilli = StartClock(game.ID, game.WhiteMinutes, game.WhiteSeconds, game.WhiteMilli, "White")

	}()

	//checks whether or not clock has timed out
	for {
		select {
		case <-Verify.AllTables[game.ID].whiteTimeOut:

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
		case <-Verify.AllTables[game.ID].blackTimeOut:

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
		case <-Verify.AllTables[game.ID].gameOver:
			//fmt.Println("Game is over but clocks have not ran out so break out of clocks")
			return
		}
	}
}

//returns the remaining time of players's clock
func StartClock(gameID int, minutes int, seconds int, milliseconds int, color string) (int, int, int) {

	timerChan := time.NewTicker(time.Millisecond).C

	clock := (60000 * minutes) + (seconds * 1000) + milliseconds
	if clock <= 0 {
		return 0, 0, 0
	}

	Verify.AllTables[gameID].Connection = make(chan bool)

	for {
		select {
		case <-Verify.AllTables[gameID].Connection:

			remainingMinutes := (clock / 60000)
			remainingSeconds := (clock / 1000) % 60
			remainingMilli := clock % 1000
			// fmt.Printf("Clock here is %d %d color is %s\n",  remainingMinutes, remainingSeconds, color)
			return remainingMinutes, remainingSeconds, remainingMilli

		case <-timerChan:

			clock--
			if clock <= 0 && color == "White" {
				Verify.AllTables[gameID].whiteTimeOut <- true
				return 0, 0, 0
			} else if clock <= 0 && color == "Black" {
				Verify.AllTables[gameID].blackTimeOut <- true
				return 0, 0, 0
			}
		}
	}

	return 0, 0, 0
}
