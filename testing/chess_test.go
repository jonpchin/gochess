package testing

import (
	"testing"

	"github.com/jonpchin/gochess/gostuff"
)

func TestPawn(t *testing.T) {

	var legal bool
	const gameID = 0
	// pass in blank string for white and black players as spectators
	//don't matter for testing
	gostuff.InitGame(gameID, "", "")

	legal = gostuff.ChessVerify("e2", "e4", "", gameID)
	if legal == false {
		t.Error("pawn_test.go TestPawn() 1 can't move two squares white")
	}
	legal = gostuff.ChessVerify("e7", "e5", "", gameID)
	if legal == false {
		t.Error("pawn_test.go TestPawn() 2 can't move two squares black")
	}
	legal = gostuff.ChessVerify("d2", "d5", "", gameID)
	if legal == true {
		t.Error("pawn_test.go TestPawn() 3, illegal three square move")
	}
	legal = gostuff.ChessVerify("d2", "d3", "", gameID)
	if legal == false {
		t.Error("pawn_test.go TestPawn() 4 can't move one square white")
	}
	legal = gostuff.ChessVerify("d2", "d3", "", gameID)
	if legal == true {
		t.Error("pawn_test.go TestPawn() 5, no pawn on d2")
	}
	legal = gostuff.ChessVerify("e5", "e4", "", gameID)
	if legal == true {
		t.Error("pawn_test.go TestPawn() 6 pawn is blocking square")
	}
	legal = gostuff.ChessVerify("f7", "f5", "", gameID)
	if legal == false {
		t.Error("pawn_test.go TestPawn() 7 black moves pawn up two squares")
	}
	legal = gostuff.ChessVerify("e4", "f6", "", gameID)
	if legal == true {
		t.Error("pawn_test.go TestPawn() 8 illegal white move")
	}
	legal = gostuff.ChessVerify("e4", "f5", "", gameID)
	if legal == false {
		t.Error("pawn_test.go TestPawn() 9 white pawn captures")
	}
	legal = gostuff.ChessVerify("g7", "g5", "", gameID)
	if legal == false {
		t.Error("pawn_test.go TestPawn() 10 black pawn two squares")
	}
	legal = gostuff.ChessVerify("f5", "g6", "", gameID)
	if legal == false {
		t.Error("pawn_test.go TestPawn() 11 enpassent")
	}
	legal = gostuff.ChessVerify("h7", "g6", "", gameID)
	if legal == false {
		t.Error("pawn_test.go TestPawn() black captures")
	}
	legal = gostuff.ChessVerify("d3", "d2", "", gameID)
	if legal == true {
		t.Error("pawn_test.go TestPawn() pawn cannot go backwards")
	}
}

func TestKnght(t *testing.T) {

	var legal bool
	const gameID = 0
	gostuff.InitGame(gameID, "", "")

	legal = gostuff.ChessVerify("g1", "e2", "", gameID)
	if legal == true {
		t.Error("white knight capturing ally piece was deemed legal")
	}
	legal = gostuff.ChessVerify("g1", "f3", "", gameID)
	if legal == false {
		t.Error("legal white knight move deemed illegal")
	}
	legal = gostuff.ChessVerify("g8", "h6", "", gameID)
	if legal == false {
		t.Error("legal black knight move deemed illegal")
	}
	legal = gostuff.ChessVerify("f3", "g4", "", gameID)
	if legal == true {
		t.Error("knight_test.go 4 testing illegal move")
	}
	legal = gostuff.ChessVerify("f3", "h4", "", gameID)
	if legal == false {
		t.Error("knight_test.go 5")
	}
	legal = gostuff.ChessVerify("b8", "c6", "", gameID)
	if legal == false {
		t.Error("knight_test.go 6")
	}
	legal = gostuff.ChessVerify("h4", "f5", "", gameID)
	if legal == false {
		t.Error("knight_test.go 7")
	}
	legal = gostuff.ChessVerify("h6", "f5", "", gameID)
	if legal == false {
		t.Error("knight_test.go 8 black captures white knight")
	}
}

func TestBishop(t *testing.T) {

	var legal bool
	const gameID = 0
	gostuff.InitGame(gameID, "", "")

	legal = gostuff.ChessVerify("f1", "b5", "", gameID)
	if legal == true {
		t.Error("bishop cannot move over ally pawn")
	}
	legal = gostuff.ChessVerify("e2", "e3", "", gameID)
	if legal == false {
		t.Error("legal white pawn move deemed illegal")
	}
	legal = gostuff.ChessVerify("a7", "a6", "", gameID)
	if legal == false {
		t.Error("legal black pawn move deemed illegal")
	}
	legal = gostuff.ChessVerify("f1", "b5", "", gameID)
	if legal == false {
		t.Error("legal white bishop move deemed illegal")
	}
	legal = gostuff.ChessVerify("g7", "g6", "", gameID)
	if legal == false {
		t.Error("legal black pawn move deemed illegal")
	}
	legal = gostuff.ChessVerify("b5", "d7", "", gameID)
	if legal == false {
		t.Error("legal white bishop capture deemed illegal")
	}
	legal = gostuff.ChessVerify("c8", "d7", "", gameID)
	if legal == false {
		t.Error("legal black biship capture deemed illegal")
	}
	legal = gostuff.ChessVerify("b2", "b4", "", gameID)
	if legal == false {
		t.Error("legal white pawn move deemed illegal")
	}
	legal = gostuff.ChessVerify("f8", "g7", "", gameID)
	if legal == false {
		t.Error("legal black bishop move deemed illegal")
	}
	legal = gostuff.ChessVerify("c1", "a3", "", gameID)
	if legal == false {
		t.Error("legal white bishop move deemed illegal")
	}
	legal = gostuff.ChessVerify("g7", "h8", "", gameID)
	if legal == true {
		t.Error("black bishop is not allowed to capture ally piece")
	}
	legal = gostuff.ChessVerify("a3", "a4", "", gameID)
	if legal == true {
		t.Error("white pawn cannot jump over ally piece")
	}
	legal = gostuff.ChessVerify("a3", "b4", "", gameID)
	if legal == true {
		t.Error("white pawn cannot move there")
	}
	legal = gostuff.ChessVerify("g7", "a1", "", gameID)
	if legal == false {
		t.Error("legal black bishop capture deemed illegal")
	}
	legal = gostuff.ChessVerify("a3", "c5", "", gameID)
	if legal == false {
		t.Error("legal white bishop move deemed illegal")
	}
}

func TestRook(t *testing.T) {
	var legal bool
	const gameID = 0
	gostuff.InitGame(gameID, "", "")

	legal = gostuff.ChessVerify("a2", "a4", "", gameID)
	if legal == false {
		t.Error("TestRook 1")
	}
	legal = gostuff.ChessVerify("h8", "h3", "", gameID)
	if legal == true {
		t.Error("TestRook 2")
	}
	legal = gostuff.ChessVerify("h7", "h5", "", gameID)
	if legal == false {
		t.Error("TestRook 3")
	}
	legal = gostuff.ChessVerify("a1", "a4", "", gameID)
	if legal == true {
		t.Error("TestRook 4")
	}
	legal = gostuff.ChessVerify("a1", "a3", "", gameID)
	if legal == false {
		t.Error("TestRook 5")
	}
	legal = gostuff.ChessVerify("h8", "c3", "", gameID)
	if legal == true {
		t.Error("TestRook 6")
	}
	legal = gostuff.ChessVerify("h8", "h6", "", gameID)
	if legal == false {
		t.Error("TestRook 7")
	}
	legal = gostuff.ChessVerify("a3", "b4", "", gameID)
	if legal == true {
		t.Error("TestRook 8")
	}
	legal = gostuff.ChessVerify("a3", "e3", "", gameID)
	if legal == false {
		t.Error("TestRook 9")
	}
	legal = gostuff.ChessVerify("a7", "a6", "", gameID)
	if legal == false {
		t.Error("TestRook 10")
	}
	legal = gostuff.ChessVerify("e3", "e7", "", gameID)
	if legal == false {
		t.Error("TestRook 11")
	}
	legal = gostuff.ChessVerify("h6", "f6", "", gameID)
	if legal == true {
		t.Error("TestRook 12")
	}
	legal = gostuff.ChessVerify("d8", "e7", "", gameID)
	if legal == false {
		t.Error("TestRook 13")
	}
	legal = gostuff.ChessVerify("a4", "a5", "", gameID)
	if legal == false {
		t.Error("TestRook 14")
	}
	legal = gostuff.ChessVerify("h6", "a6", "", gameID)
	if legal == true {
		t.Error("TestRook 15")
	}
	legal = gostuff.ChessVerify("h6", "b6", "", gameID)
	if legal == false {
		t.Error("TestRook 16")
	}
}
