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
		t.Error("castle_test.go TestKingCastle() 1 can't move two squares white")
	}
	legal = chessVerify("e7", "e5", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 2 can't move two squares black")
	}
	legal = chessVerify("g1", "f3", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 3 can't move knight")
	}
	legal = chessVerify("d8", "h4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 4 can't move queen")
	}
	legal = chessVerify("f1", "b5", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't move bishop")
	}
	legal = chessVerify("h4", "e4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't move queen")
	}
	legal = chessVerify("e1", "e2", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestKingCastle() 5 king is in check and should not be able to move")
	}
	legal = chessVerify("d1", "e2", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't move bishop")
	}
	legal = chessVerify("g8", "f6", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't move knight")
	}
	legal = chessVerify("e1", "g1", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't castle")
	}
	legal = chessVerify("f8", "b4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't move bishop")
	}
	legal = chessVerify("e2", "c4", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't move queen")
	}
	legal = chessVerify("a7", "a6", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't move pawn")
	}
	legal = chessVerify("c4", "c5", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't move bishop")
	}
	legal = chessVerify("e8", "g8", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestKingCastle() 5 can't move through check")
	}
	legal = chessVerify("h8", "g8", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't move rook")
	}
	legal = chessVerify("g1", "f1", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestKingCastle() 5 rook is already there")
	}
	legal = chessVerify("g1", "h1", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't move bishop")
	}
	legal = chessVerify("g8", "h8", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't move rook")
	}
	legal = chessVerify("h1", "f3", "", gameID)
	if legal == true {
		t.Error("castle_test.go TestKingCastle() 5 illegal king move")
	}
	legal = chessVerify("c5", "c3", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't move queen")
	}
	legal = chessVerify("e8", "f8", "", gameID)
	if legal == false {
		t.Error("castle_test.go TestKingCastle() 5 can't move king")
	}
}
