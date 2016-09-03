package gostuff

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"os"
	"time"
)

//this function gets executed on ctrl-c
func Cleanup() {

	fmt.Println("Web server is shutting down...saving games please wait...")

	var message ChatInfo
	message.Type = "massMessage"
	message.Text = "ATTENTION: Web server is shutting down NOW for maintenance, brace for impact..."

	go func() {
		for _, value := range Active.Clients {
			if err := websocket.JSON.Send(value, &message); err != nil {
				// we could not send the message to a peer
				fmt.Println("cleanup.go CleanUp() error 1  Could not send message to ", err)
			}
		}

		for _, value := range Chat.Lobby {
			if err := websocket.JSON.Send(value, &message); err != nil {
				// we could not send the message to a peer
				fmt.Println("cleanup.go CleanUp() error 2  Could not send message to ", err)
			}
		}
	}()

	for _, game := range All.Games {
		//now store game in MySQL database
		allMoves, err := json.Marshal(game.GameMoves)
		if err == nil {
			//gets length of all the moves in the game
			totalMoves := (len(All.Games[game.ID].GameMoves) + 1) / 2
			saveGame(totalMoves, allMoves, game)
		} else {
			fmt.Println("Error in Cleanup.go cleanup 1")
		}

	}
	fmt.Println("All games are saved. Web server is shutting down in 1 second.")
}

//used when web server is shutting down to save all current games into database
func saveGame(totalMoves int, allMoves []byte, game *ChessGame) {

	moves := string(allMoves)
	//fmt.Println("The game moves are ", moves)

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("cleanup.go saveGame 0 DATABASE DOWN!")
		return
	}

	stmt, err := db.Prepare("INSERT saved SET white=?, black=?, gametype=?, rated=?," +
		" whiterating=?, blackrating=?, blackminutes=?, blackseconds=?, whiteminutes=?," +
		" whiteseconds=?, timecontrol=?, moves=?, totalmoves=?, status=?, date=?, time=?")
	if err != nil {
		log.Println("Cleanup.go saveGame 1 ", err)
		return
	}
	date := time.Now()
	res, err := stmt.Exec(game.WhitePlayer, game.BlackPlayer, game.GameType, game.Rated,
		game.WhiteRating, game.BlackRating, game.BlackMinutes, game.BlackSeconds,
		game.WhiteMinutes, game.WhiteSeconds, game.TimeControl, moves, totalMoves,
		game.Status, date, date)
	if err != nil {
		log.Println("Cleanup.go saveGame 2 ", err)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("Cleanup.go saveGame 3 ", err)
		return
	}
	log.Printf("Game ID %d has been updated in the saved table.\n", id)
}
