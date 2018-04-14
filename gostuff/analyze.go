package gostuff

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/malbrecht/chess"
)

// Uses stockfish engine to analyze game, returns percentage of moves
// that match the engine for the given depth
func AnalyzeGame(chessMoves []chess.Move, depth int) {

	engine := startEngine(nil)
	defaultDepth := 3 // For now set the default depth to 3

	// All standard chess games start with the same position
	startPostion := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	board, err := chess.ParseFen(startPostion)
	if err != nil {
		fmt.Println(err)
		return
	}

	// total strong moves that match the chess engine
	totalStrongMoves := 0
	var updatedBoard *chess.Board

	for _, move := range chessMoves {

		board = board.MakeMove(move)
		currentFen := board.Fen()

		isOk, bestMove := engineSearchDepth(currentFen, engine, defaultDepth)

		if isOk == false {
			fmt.Println("Error processing move in analyze games for board.Fen()", currentFen)
			return
		}

		updatedBoard = board.MakeMove(bestMove)
		if updatedBoard.Fen() == currentFen {
			totalStrongMoves += 1
		}
	}

	engine.Quit()
}

// Gets all moves (in engine notation) for a given game id in the database
// If there was an error in getting the games or if there is no games for the given id
// then an empty slice is returned
func getEngineMovesInGames(id string) []chess.Move {

	var moves string
	var engineMoves []chess.Move

	err := db.QueryRow("SELECT moves FROM games WHERE id=?", id).Scan(&moves)
	if err != nil {
		log.Println(err)
		return engineMoves
	}

	var gochessMoves []Move

	temp := []byte(moves)
	err = json.Unmarshal(temp, &gochessMoves)
	if err != nil {
		log.Println(err)
		return engineMoves
	}
	color := 0 // 0 is black, 1 is white
	for index, move := range gochessMoves {
		engineMoves[index].From = getEngineSquare(move.S)
		engineMoves[index].To = getEngineSquare(move.T)

		if index%2 == 0 { // then its black's move
			color = 0
		} else { // then its white's move
			color = 1
		}
		engineMoves[index].Promotion = getPromotionPiece(move.P, color)
	}

	return engineMoves
}
