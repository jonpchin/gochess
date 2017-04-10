package gostuff

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/malbrecht/chess"
	"github.com/malbrecht/chess/engine/uci"
)

// Run starts an engine executable, with the given arguments.
// Returns the engine, make sure to call Quit() on the engine to clean up
func startEngine(args []string) *uci.Engine {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	var engine *uci.Engine
	var err error

	if runtime.GOOS == "windows" {
		engine, err = uci.Run("./stockfish/stockfish_8_x64.exe", args, log)
		if err != nil {
			log.Println(err)
		}
	} else {
		engine, err = uci.Run("./stockfish/stockfish_8_x64", args, log)
		if err != nil {
			log.Println(err)
		}
	}
	return engine
}

// isCheckOrMate returns whether the side to move is in check and/or has been mated.
// Mate without check means stalemate.
func isCheckMate(fen string) (check, mate bool) {

	board, err := chess.ParseFen(fen)
	if err != nil {
		fmt.Println(err)
	}
	return board.IsCheckOrMate()
}

// Start position is "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
// Returns false if there was an error or true if a best move was found
func engineSearchDepth(fen string, engine *uci.Engine, depth int) (bool, chess.Move) {

	board, err := chess.ParseFen(fen)
	var chessMove chess.Move

	if err != nil {
		return false, chessMove
	}

	engine.SetPosition(board)

	for info := range engine.SearchDepth(depth) {
		if move, ok := info.BestMove(); ok {
			return true, move
		} else if err := info.Err(); err != nil {
			return false, chessMove
		}
	}
	return false, chessMove
}

// Returns false if there was an error or true if a best move was found
func engineSearchTime(fen string, engine *uci.Engine, t time.Duration) (bool, chess.Move) {

	board, err := chess.ParseFen(fen)
	var chessMove chess.Move

	if err != nil {
		return false, chessMove
	}

	engine.SetPosition(board)

	for info := range engine.SearchTime(t) {
		if move, ok := info.BestMove(); ok {
			return true, move
		} else if err := info.Err(); err != nil {
			return false, chessMove
		}
	}
	return false, chessMove
}

func Quit(engine *uci.Engine) {
	engine.Stop()
	engine.Quit()
}
