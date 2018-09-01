package testing

import (
	"testing"

	"github.com/jonpchin/gochess/gostuff"
)

// Tests king castle
func TestKingCastle(t *testing.T) {

	var legal bool
	const gameID = 0
	// pass in blank string for white and black players as spectators
	//don't matter for testing
	gostuff.InitGame(gameID, "", "")

	cases := []struct {
		src   string
		dst   string
		legal bool
		msg   string
	}{
		{"e2", "e4", false, "e2 e4 can't move two squares white"},
		{"e7", "e5", false, "e7 e5 can't move two squares black"},
		{"g1", "f3", false, "g1 f3 can't move knight"},
		{"d8", "h4", false, "d8 h4 can't move queen"},
		{"f1", "b5", false, "f1 b5 can't move bishop"},
		{"h4", "e4", false, "h4 e4  can't move queen"},
		{"e1", "e2", true, "e1 e2 king is in check and should not be able to move"},
		{"d1", "e2", false, "d1 e2 can't move bishop"},
		{"g8", "f6", false, "g8 f6 can't move knight"},
		{"e1", "g1", false, "e1 g1 can't castle"},
		{"f8", "b4", false, "f8 b4 can't move bishop"},
		{"e2", "c4", false, "e2 c4 can't move queen"},
		{"a7", "a6", false, "a7 a6 can't move pawn"},
		{"c4", "c5", false, "c4 c5 can't move bishop"},
		{"e8", "g8", true, "e8 g8 can't move through check"},
		{"h8", "g8", false, "h8, g8  can't move rook"},
		{"g1", "f1", true, "g1 f1  rook is already there"},
		{"g1", "h1", false, "g1 h1 can't move bishop"},
		{"g8", "h8", false, "g8 h8 can't move rook"},
		{"h1", "f3", true, "h1 f3  illegal king move"},
		{"c5", "c3", false, "c5 c3 can't move queen"},
		{"e8", "f8", false, "e8 f8 can't move king"},
	}

	for _, c := range cases {
		legal = gostuff.ChessVerify(c.src, c.dst, "", gameID)
		if legal == c.legal {
			t.Error(c.msg)
		}
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
