package gostuff

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/jonpchin/chess"
	"github.com/jonpchin/chess/engine/uci"
	pgn "gopkg.in/freeeve/pgn.v1"
)

// On a brand new installation on a server make sure to check if castle fix is incorporated
// in third party library at https://github.com/malbrecht/chess/pull/3

// FEN string of played move and best move suggested by engine
type MoveAnalysis struct {
	PlayedMoveFen       string
	BestMoveFen         string
	PlayedMoveSrc       string //gochess notation
	PlayedMoveTar       string //gochess notation
	PlayedMovePromotion string //gochess notation
	BestMoveSrc         string //gochess notation
	BestMoveTar         string //gochess notation
	BestMovePromotion   string //gochess notation
}

type GameAnalysis struct {
	Moves []MoveAnalysis // List of actually and best moves in FEN string
	Depth int            // the depth searched
}

// Uses stockfish engine to analyze game, returns a GameAnalysis that can be marshalled and sent to front end
// that match the engine for the given depth
func (gameAnalysis *GameAnalysis) analyzeGame(chessMoves []chess.Move, gochessMoves []GameMove) {

	engine := StartEngine(nil)

	// All standard chess games start with the same position
	startPosition := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	var moveAnalysis MoveAnalysis
	moveAnalysis.PlayedMoveFen = startPosition
	moveAnalysis.BestMoveFen = ""
	gameAnalysis.Moves = append(gameAnalysis.Moves, moveAnalysis)

	board, err := chess.ParseFen(startPosition)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentFen := board.Fen()
	// Original board will keep FEN string of the position before the next move is made
	originalBoard := board

	for index, move := range chessMoves {

		board = board.MakeMove(move)

		isOk, bestMove := engineSearchDepth(currentFen, engine, gameAnalysis.Depth)
		currentFen = board.Fen()

		if isOk == false {
			fmt.Println("Error processing move in analyze games for currentFen:", currentFen)
			break
		}

		bestMoveBoard := originalBoard.MakeMove(bestMove)
		originalBoard = board

		moveAnalysis.PlayedMoveFen = currentFen
		moveAnalysis.BestMoveFen = bestMoveBoard.Fen()

		moveAnalysis.PlayedMoveSrc = gochessMoves[index].S
		moveAnalysis.PlayedMoveTar = gochessMoves[index].T
		moveAnalysis.PlayedMovePromotion = gochessMoves[index].P
		moveAnalysis.BestMoveSrc = engineBoard[bestMove.From]
		moveAnalysis.BestMoveTar = engineBoard[bestMove.To]
		moveAnalysis.BestMovePromotion = string(bestMove.Promotion)

		gameAnalysis.Moves = append(gameAnalysis.Moves, moveAnalysis)
	}

	engine.Quit()
}

func (allGames *allPgnGames) analyzePgnGames(pgnMoves []pgn.Move, engine *uci.Engine) {

	// All standard chess games start with the same position
	startPosition := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	var moveAnalysis MoveAnalysis
	var gameAnalysis GameAnalysis
	moveAnalysis.PlayedMoveFen = startPosition
	moveAnalysis.BestMoveFen = ""
	gameAnalysis.Moves = append(gameAnalysis.Moves, moveAnalysis)

	board, err := chess.ParseFen(startPosition)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentFen := board.Fen()
	// Original board will keep FEN string of the position before the next move is made
	originalBoard := board
	var malbrechtMove chess.Move

	for _, move := range pgnMoves {

		malbrechtMove.From = getEngineSquare(move.From.String())
		malbrechtMove.To = getEngineSquare(move.To.String())
		malbrechtMove.Promotion = pgnPromoteToMalbrechtPromote(move.Promote)
		board = board.MakeMove(malbrechtMove)

		isOk, bestMove := engineSearchDepth(currentFen, engine, gameAnalysis.Depth)
		currentFen = board.Fen()

		if isOk == false {
			fmt.Println("Error processing move in analyze games for currentFen:", currentFen)
			break
		}

		bestMoveBoard := originalBoard.MakeMove(bestMove)
		originalBoard = board

		moveAnalysis.PlayedMoveFen = currentFen
		moveAnalysis.BestMoveFen = bestMoveBoard.Fen()

		moveAnalysis.PlayedMoveSrc = move.From.String()
		moveAnalysis.PlayedMoveTar = move.To.String()
		moveAnalysis.PlayedMovePromotion = string(move.Promote)

		moveAnalysis.BestMoveSrc = engineBoard[bestMove.From]
		moveAnalysis.BestMoveTar = engineBoard[bestMove.To]
		moveAnalysis.BestMovePromotion = string(bestMove.Promotion)

		gameAnalysis.Moves = append(gameAnalysis.Moves, moveAnalysis)
	}

	*allGames = append(*allGames, gameAnalysis)
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

	username, err := r.Cookie("username")
	if err != nil {
		fmt.Println(err)
		return
	}

	// If user is an admin or a mod then they allowed to analyze games in back end
	if IsAdmin(username.Value) == false && IsMod(username.Value) == false {
		fmt.Println(username.Value, "tried analyzing a game when they shouldn't have GameAnalysisById")
		return
	}

	// Get the gameID specified in the front end
	id := template.HTMLEscapeString(r.FormValue("id"))
	depth, err := strconv.Atoi(template.HTMLEscapeString(r.FormValue("depth")))

	if err != nil {
		fmt.Println("Could not convert string to int in GetEngineAnalysisById", err)
		return
	}

	if depth < 1 || depth > 7 {
		fmt.Println("Depth is not in a valid range: ", depth)
		return
	}

	var moves string

	err = db.QueryRow("SELECT moves FROM games WHERE id=?", id).Scan(&moves)
	if err != nil {
		log.Println(err)
		return
	}

	var gochessMoves []GameMove

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

	gameAnalysis.Depth = depth
	gameAnalysis.analyzeGame(engineMoves, gochessMoves)

	// JSON marshal and send game to front end
	jsonGameAnalysis, err := json.Marshal(gameAnalysis)
	if err != nil {
		fmt.Println("Could not marshal gameAnalysis", err)
		return
	}
	w.Write([]byte((jsonGameAnalysis)))
}
