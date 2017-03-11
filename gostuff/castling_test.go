package gostuff

import "testing"

func TestKingCastle(t *testing.T) {

	var legal bool
	const gameID = 0
	// pass in blank string for white and black players as spectators
	//don't matter for testing
	initGame(gameID, "", "")

	legal = chessVerify("e2", "e4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move two squares white")
	}
	legal = chessVerify("e7", "e5", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move two squares black")
	}
	legal = chessVerify("g1", "f3", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move knight")
	}
	legal = chessVerify("d8", "h4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move queen")
	}
	legal = chessVerify("f1", "b5", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move bishop")
	}
	legal = chessVerify("h4", "e4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move queen")
	}
	legal = chessVerify("e1", "e2", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestKingCastle() king is in check and should not be able to move")
	}
	legal = chessVerify("d1", "e2", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move bishop")
	}
	legal = chessVerify("g8", "f6", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move knight")
	}
	legal = chessVerify("e1", "g1", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't castle")
	}
	legal = chessVerify("f8", "b4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move bishop")
	}
	legal = chessVerify("e2", "c4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move queen")
	}
	legal = chessVerify("a7", "a6", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move pawn")
	}
	legal = chessVerify("c4", "c5", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move bishop")
	}
	legal = chessVerify("e8", "g8", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestKingCastle() can't move through check")
	}
	legal = chessVerify("h8", "g8", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move rook")
	}
	legal = chessVerify("g1", "f1", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestKingCastle() rook is already there")
	}
	legal = chessVerify("g1", "h1", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move bishop")
	}
	legal = chessVerify("g8", "h8", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move rook")
	}
	legal = chessVerify("h1", "f3", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestKingCastle() illegal king move")
	}
	legal = chessVerify("c5", "c3", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move queen")
	}
	legal = chessVerify("e8", "f8", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() can't move king")
	}
}

func TestQueenCastle(t *testing.T) {
	var legal bool
	const gameID = 0
	// pass in blank string for white and black players as spectators
	//don't matter for testing
	initGame(gameID, "", "")

	legal = chessVerify("d2", "d4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move two squares white")
	}
	legal = chessVerify("d7", "d5", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move two squares black")
	}
	legal = chessVerify("c1", "g5", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move bishop")
	}

	legal = chessVerify("c8", "g4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move bishop")
	}
	legal = chessVerify("b1", "c3", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move knight")
	}
	legal = chessVerify("e8", "c8", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestQueenCastle() can't castle queenside with queen blocking")
	}
	legal = chessVerify("d8", "d7", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move queen")
	}
	legal = chessVerify("d1", "d2", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move queen")
	}
	legal = chessVerify("b8", "c6", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move knight")
	}
	legal = chessVerify("e1", "c1", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't castle king")
	}
	legal = chessVerify("e8", "d8", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move king")
	}
	legal = chessVerify("c1", "b1", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move king")
	}
	legal = chessVerify("d8", "e8", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move king")
	}
	legal = chessVerify("d1", "e1", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestQueenCastle() can't move rook")
	}
	legal = chessVerify("e8", "c8", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestQueenCastle() can't castle queenside as king already moved")
	}
}
