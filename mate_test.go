package gostuff

import "testing"

func TestScholarMate(t *testing.T) {

	var legal bool
	const gameID = 0
	// pass in blank string for white and black players as spectators
	//don't matter for testing
	initGame(gameID, "", "")

	legal = chessVerify("e2", "e4", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestKingCastle() 1 can't move two squares white")
	}
	legal = chessVerify("e7", "e5", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestKingCastle() 2 can't move two squares black")
	}
	legal = chessVerify("f1", "c4", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestKingCastle() 2 can't move bishop")
	}
	legal = chessVerify("b8", "c6", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestKingCastle() 2 can't move knight")
	}
	legal = chessVerify("d1", "h5", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestKingCastle() 2 can't move queen")
	}
	legal = chessVerify("g8", "f6", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestKingCastle() 2 can't move knight")
	}
	legal = chessVerify("h5", "f7", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestKingCastle() 2 can't move queen")
	}
	legal = chessVerify("e8", "f7", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestKingCastle() 2 can't capture queen its mate")
	}
	legal = chessVerify("e8", "e7", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestKingCastle() 2 can't move king its mate")
	}

}

func TestSmotheredMate(t *testing.T) {

}
