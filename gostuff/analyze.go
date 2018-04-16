package gostuff

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/malbrecht/chess"
)

// FEN string of played move and best move suggested by engine
type MoveAnalysis struct {
	PlayedMove string
	BestMove   string
}

type GameAnalysis struct {
	Moves []MoveAnalysis // List of actually and best moves in FEN string
	Depth int            // the depth searched
}

// Uses stockfish engine to analyze game, returns a GameAnalysis that can be marshalled and sent to front end
// that match the engine for the given depth
func (gameAnalysis *GameAnalysis) analyzeGame(chessMoves []chess.Move) {

	engine := startEngine(nil)

	// All standard chess games start with the same position
	startPosition := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	var moveAnalysis MoveAnalysis
	moveAnalysis.PlayedMove = startPosition
	moveAnalysis.BestMove = ""
	gameAnalysis.Moves = append(gameAnalysis.Moves, moveAnalysis)

	board, err := chess.ParseFen(startPosition)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentFen := board.Fen()
	// Original board will keep FEN string of the position before the next move is made
	originalBoard := board

	for _, move := range chessMoves {

		board = board.MakeMove(move)

		isOk, bestMove := engineSearchDepth(currentFen, engine, gameAnalysis.Depth)
		currentFen = board.Fen()

		if isOk == false {
			fmt.Println("Error processing move in analyze games for currentFen:", currentFen)
			break
		}

		bestMoveBoard := originalBoard.MakeMove(bestMove)
		originalBoard = board

		moveAnalysis.PlayedMove = currentFen
		moveAnalysis.BestMove = bestMoveBoard.Fen()
		gameAnalysis.Moves = append(gameAnalysis.Moves, moveAnalysis)
	}

	engine.Quit()
}

// Convert JSON list of FEN strings into malbrecht chess notation for stock fish analysis
func EngineAnalysisByJsonFen(w http.ResponseWriter, r *http.Request) {

}

// Convert PGN chess game into malbrecht chess notation for stock fish analysis
func EngineAnalysisByPgnFen(w http.ResponseWriter, r *http.Request) {

}

// Gets all moves (in engine notation) for a given game id in the database
// by converting gochess move notation into malbrecht chess notation
func GameAnalysisById(w http.ResponseWriter, r *http.Request) {

	valid := ValidateCredentials(w, r)
	if valid == false {
		return
	}

	// Get the gameID specified in the front end
	id := template.HTMLEscapeString(r.FormValue("id"))
	depth := template.HTMLEscapeString(r.FormValue("depth"))

	var moves string

	err := db.QueryRow("SELECT moves FROM games WHERE id=?", id).Scan(&moves)
	if err != nil {
		log.Println(err)
		return
	}

	var gochessMoves []Move

	temp := []byte(moves)
	err = json.Unmarshal(temp, &gochessMoves)
	if err != nil {
		log.Println(err)
		return
	}
	color := 0 // 0 is black, 1 is white
	var engineMoves []chess.Move
	var chessMove chess.Move

	for index, move := range gochessMoves {
		chessMove.From = getEngineSquare(move.S)
		chessMove.To = getEngineSquare(move.T)

		if index%2 == 0 { // then its black's move
			color = 0
		} else { // then its white's move
			color = 1
		}
		chessMove.Promotion = getPromotionPiece(move.P, color)
		engineMoves = append(engineMoves, chessMove)
	}

	var gameAnalysis GameAnalysis
	convertedDepth, err := strconv.Atoi(depth)
	if err != nil {
		fmt.Println("Could not convert string to int in GetEngineAnalysisById", err)
		return
	}

	gameAnalysis.Depth = convertedDepth
	gameAnalysis.analyzeGame(engineMoves)

	// JSON marshal and send game to front end
	jsonGameAnalysis, err := json.Marshal(gameAnalysis)
	if err != nil {
		fmt.Println("Could not marshal gameAnalysis", err)
		return
	}
	w.Write([]byte((jsonGameAnalysis)))
}
