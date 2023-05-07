package gostuff

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/jonpchin/chess"
	"github.com/jonpchin/chess/engine/uci"
	pgn "gopkg.in/freeeve/pgn.v1"
)

// malbrecht chess engine uses notation where a1 is 0, b1 is 1 c1 is 2...h8 is 63
// while gochess notation uses standard chess notation such as e4, c5
// engine notation is used to map gochess notation into malbrecht chess engine notation
var engineBoard = [64]string{
	"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
	"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
	"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
	"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
	"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
	"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
	"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
	"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
}

const (
	White = iota
	Black
)

// Converts pgn to malbrecht promote for piece promotion
func pgnPromoteToMalbrechtPromote(piece pgn.Piece) chess.Piece {
	result := chess.Piece(pgn.NoPiece)
	switch piece {
	case pgn.WhiteKnight:
		result = chess.WN
	case pgn.WhiteBishop:
		result = chess.WB
	case pgn.WhiteRook:
		result = chess.WR
	case pgn.WhiteQueen:
		result = chess.WQ
	case pgn.BlackBishop:
		result = chess.BB
	case pgn.BlackKnight:
		result = chess.BN
	case pgn.BlackRook:
		result = chess.BR
	case pgn.BlackQueen:
		result = chess.BQ
	default:
		result = chess.NoPiece
	}
	return result
}

// Converts gochess square notation to engine square notation
// Returns -1 if no notation matches
func getEngineSquare(gochessNotation string) chess.Sq {

	for index, square := range engineBoard {
		if gochessNotation == square {
			return chess.Sq(index)
		}
	}
	return -1
}

func getGochessSquare(engineNotation chess.Sq) string {
	for index, square := range engineBoard {
		if engineNotation == chess.Sq(index) {
			return square
		}
	}
	return ""
}

func getGoChessPromotionPiece(enginePiece chess.Piece) string {
	if enginePiece == chess.WN || enginePiece == chess.BN {
		return "n"
	} else if enginePiece == chess.WB || enginePiece == chess.BB {
		return "b"
	} else if enginePiece == chess.WQ || enginePiece == chess.BQ {
		return "q"
	} else if enginePiece == chess.WR || enginePiece == chess.BR {
		return "r"
	} else {
		return ""
	}
}

// Engine promotion piece
func getPromotionPiece(piece string, color int) chess.Piece {
	enginePiece := chess.NoPiece
	switch piece {
	case "q":
		enginePiece = chess.Queen
	case "r":
		enginePiece = chess.Rook
	case "b":
		enginePiece = chess.Bishop
	case "n":
		enginePiece = chess.Knight
	default: // Can't promote to anything other then queen, rook, bishop or knight
		fmt.Println("Invalid piece in getPromotionPiece", piece, color)
	}
	return chess.Piece(enginePiece | color)
}

// Run starts an engine executable, with the given arguments.
// Returns the engine, make sure to call Quit() on the engine to clean up
func StartEngine(args []string) *uci.Engine {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	var engine *uci.Engine
	var err error

	if runtime.GOOS == "windows" {
		engine, err = uci.Run("stockfish_8_x64.exe", args, log)
		if err != nil {
			log.Println(err)
		}
	} else {
		engine, err = uci.Run("./stockfish", args, log)
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

func EngineSearchTimeRaw(fen string, engine *uci.Engine, t time.Duration) (bool, string) {

	board, err := chess.ParseFen(fen)

	if err != nil {
		return false, ""
	}

	engine.SetPosition(board)

	for info := range engine.SearchTime(t) {
		if move, ok := info.BestMoveRaw(); ok {
			return ok, move
		}
	}
	return false, ""
}

func Quit(engine *uci.Engine) {
	engine.Stop()
	engine.Quit()
}
