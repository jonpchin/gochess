package testing

import (
	"testing"

	"github.com/jonpchin/gochess/gostuff"
)

func TestKingCastle(t *testing.T) {

	var legal bool
	const gameID = 0
	// pass in blank string for white and black players as spectators
	//don't matter for testing
	gostuff.InitGame(gameID, "", "")

	legal = gostuff.ChessVerify("e2", "e4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move two squares white")
	}
	legal = gostuff.ChessVerify("e7", "e5", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move two squares black")
	}
	legal = gostuff.ChessVerify("g1", "f3", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move knight")
	}
	legal = gostuff.ChessVerify("d8", "h4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move queen")
	}
	legal = gostuff.ChessVerify("f1", "b5", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move bishop")
	}
	legal = gostuff.ChessVerify("h4", "e4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move queen")
	}
	legal = gostuff.ChessVerify("e1", "e2", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestKingCastle() king is in check and should not be able to move")
	}
	legal = gostuff.ChessVerify("d1", "e2", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move bishop")
	}
	legal = gostuff.ChessVerify("g8", "f6", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move knight")
	}
	legal = gostuff.ChessVerify("e1", "g1", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't castle")
	}
	legal = gostuff.ChessVerify("f8", "b4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move bishop")
	}
	legal = gostuff.ChessVerify("e2", "c4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move queen")
	}
	legal = gostuff.ChessVerify("a7", "a6", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move pawn")
	}
	legal = gostuff.ChessVerify("c4", "c5", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move bishop")
	}
	legal = gostuff.ChessVerify("e8", "g8", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestKingCastle() can't move through check")
	}
	legal = gostuff.ChessVerify("h8", "g8", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move rook")
	}
	legal = gostuff.ChessVerify("g1", "f1", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestKingCastle() rook is already there")
	}
	legal = gostuff.ChessVerify("g1", "h1", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move bishop")
	}
	legal = gostuff.ChessVerify("g8", "h8", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move rook")
	}
	legal = gostuff.ChessVerify("h1", "f3", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestKingCastle() illegal king move")
	}
	legal = gostuff.ChessVerify("c5", "c3", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move queen")
	}
	legal = gostuff.ChessVerify("e8", "f8", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move king")
	}
}

func TestQueenCastle(t *testing.T) {
	var legal bool
	const gameID = 0
	// pass in blank string for white and black players as spectators
	//don't matter for testing
	gostuff.InitGame(gameID, "", "")

	legal = gostuff.ChessVerify("d2", "d4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move two squares white")
	}
	legal = gostuff.ChessVerify("d7", "d5", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move two squares black")
	}
	legal = gostuff.ChessVerify("c1", "g5", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move bishop")
	}

	legal = gostuff.ChessVerify("c8", "g4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move bishop")
	}
	legal = gostuff.ChessVerify("b1", "c3", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move knight")
	}
	legal = gostuff.ChessVerify("e8", "c8", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestQueenCastle() can't castle queenside with queen blocking")
	}
	legal = gostuff.ChessVerify("d8", "d7", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move queen")
	}
	legal = gostuff.ChessVerify("d1", "d2", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move queen")
	}
	legal = gostuff.ChessVerify("b8", "c6", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move knight")
	}
	legal = gostuff.ChessVerify("e1", "c1", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't castle king")
	}
	legal = gostuff.ChessVerify("e8", "d8", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move king")
	}
	legal = gostuff.ChessVerify("c1", "b1", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move king")
	}
	legal = gostuff.ChessVerify("d8", "e8", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move king")
	}
	legal = gostuff.ChessVerify("d1", "e1", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move rook")
	}
	legal = gostuff.ChessVerify("e8", "c8", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestQueenCastle() can't castle queenside as king already moved")
	}
}
