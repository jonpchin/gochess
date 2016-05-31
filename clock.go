package gostuff

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"time"
)

//starting white's clock first, this goroutine will keep track of both players clock for this game
func setClocks(gameID int16, name string) {

	var result float64
	go func() {

		All.Games[gameID].WhiteMinutes, All.Games[gameID].WhiteSeconds = StartClock(gameID, All.Games[gameID].WhiteMinutes, All.Games[gameID].WhiteSeconds, "White")
		
	}()

	//checks whether or not clock has timed out
	for {
		select{
			case <-Verify.AllTables[gameID].whiteTimeOut:
				
				result = 0.0
				//update ratings
				ComputeRating(name, gameID, All.Games[gameID].GameType, result)

				//Black won as white ran out of time
				All.Games[gameID].Status = "Black won on time"
				All.Games[gameID].Type = "game_over"
				All.Games[gameID].Result = 0
				//now store game in MySQL database
				allMoves, err := json.Marshal(All.Games[gameID].GameMoves)
				if err != nil {
					fmt.Println("Error marshalling data to store in MySQL")
				}
				//gets length of all the moves in the game
				totalMoves := (len(All.Games[gameID].GameMoves) + 1) / 2
				//save game to database before deleting it from memory
				storeGame(totalMoves, allMoves, All.Games[gameID])

				//notifiy both player that black won on time
				if _, ok := Active.Clients[PrivateChat[name]]; ok { // send data if other guy is still connected
					websocket.JSON.Send(Active.Clients[PrivateChat[name]], All.Games[gameID])
				}

				if _, ok := Active.Clients[name]; ok { // send data if other guy is still connected
					websocket.JSON.Send(Active.Clients[name], All.Games[gameID])
				}

				delete(All.Games, gameID)
				delete(Verify.AllTables, gameID)

				return
			case <-Verify.AllTables[gameID].blackTimeOut:
				
				//White won as black ran out of time
				All.Games[gameID].Status = "White won on time"
				result = 1.0
				All.Games[gameID].Type = "game_over"
				All.Games[gameID].Result = 1

				//now store game in MySQL database
				allMoves, err := json.Marshal(All.Games[gameID].GameMoves)
				if err != nil {
					fmt.Println("Error marshalling data to store in MySQL")
				}
				//gets length of all the moves in the game
				totalMoves := (len(All.Games[gameID].GameMoves) + 1) / 2
				//save game to database before deleting it from memory
				storeGame(totalMoves, allMoves, All.Games[gameID])

				//update ratings
				ComputeRating(name, gameID, All.Games[gameID].GameType, result)

				//notifiy both player that white won on time
				if _, ok := Active.Clients[PrivateChat[name]]; ok { // send data if other guy is still connected
					websocket.JSON.Send(Active.Clients[PrivateChat[name]], All.Games[gameID])
				}

				if _, ok := Active.Clients[name]; ok { // send data if other guy is still connected
					websocket.JSON.Send(Active.Clients[name], All.Games[gameID])
				}

				delete(All.Games, gameID)
				delete(Verify.AllTables, gameID)
				
				return
			case <-Verify.AllTables[gameID].gameOver:
//				fmt.Println("Game is over but clocks have not ran out so break out of clocks")
				return
		} 
		
	
	}

}

//returns the remaining time of players's clock
func StartClock(gameID int16, minutes int, seconds int, color string) (int, int) {
	
	timerChan := time.NewTicker(time.Second).C

	clock := (60 * minutes) + seconds 
	if clock <= 0{
		return 0, 0 
	}
	
	Verify.AllTables[gameID].Connection = make(chan bool)
	
	for{
		select{
			case <-Verify.AllTables[gameID].Connection:
				
				remainingMinutes := clock / 60
				remainingSeconds := clock  % 60
	
//				fmt.Printf("Clock here is %d %d color is %s\n",  remainingMinutes, remainingSeconds, color)
				return remainingMinutes, remainingSeconds
				
			case <-timerChan:
				
				clock--
				if clock <= 0 && color == "White"{
					Verify.AllTables[gameID].whiteTimeOut <- true
					return 0, 0
				}else if clock <= 0 && color == "Black"{
					Verify.AllTables[gameID].blackTimeOut <- true
					return 0, 0
				}
			
		}
	}
	
	return 0, 0
}
