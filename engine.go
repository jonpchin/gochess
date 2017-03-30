package gostuff

import (
	"log"
	"os"
	"runtime"

	"github.com/malbrecht/chess"
	"github.com/malbrecht/chess/engine/uci"
)

// Run starts an engine executable, with the given arguments.
// Returns the engine, make sure to call Quit() on the engine to clean up
func StartEngine(args []string) *uci.Engine {

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
	board, err := chess.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	engine.SetPosition(board)

	for info := range engine.SearchDepth(4) {
		if err := info.Err(); err != nil {
			log.Fatal(err)
		} else if move, ok := info.BestMove(); ok {
			log.Print("the best move is", move)
		} else {
			log.Print(info.Pv(), info.Stats())
		}
	}
	return engine
}
