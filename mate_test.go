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
		t.Error("mate_test.go TestScholarMate() can't move two squares white")
	}
	legal = chessVerify("e7", "e5", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestScholarMate() can't move two squares black")
	}
	legal = chessVerify("f1", "c4", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestScholarMate() can't move bishop")
	}
	legal = chessVerify("b8", "c6", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestScholarMate() can't move knight")
	}
	legal = chessVerify("d1", "h5", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestScholarMate() can't move queen")
	}
	legal = chessVerify("g8", "f6", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestScholarMate() can't move knight")
	}
	legal = chessVerify("h5", "f7", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestScholarMate() can't move queen")
	}
	legal = chessVerify("e8", "f7", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestScholarMate() can't capture queen its mate")
	}
	legal = chessVerify("e8", "e7", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestScholarMate() can't move king its mate")
	}
	legal = chessVerify("a7", "a6", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestScholarMate() can't move pawn its mate")
	}
	legal = chessVerify("c6", "b4", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestScholarMate() can't move knight its mate")
	}
	legal = chessVerify("d8", "f7", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestScholarMate() can't move queen its mate")
	}
}

func TestSmotheredMate(t *testing.T) {

	var legal bool
	const gameID = 0
	// pass in blank string for white and black players as spectators
	//don't matter for testing
	initGame(gameID, "", "")

	legal = chessVerify("e2", "e4", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move two squares white")
	}
	legal = chessVerify("e7", "e5", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move two squares black")
	}
	legal = chessVerify("g1", "f3", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move knight")
	}
	legal = chessVerify("g8", "f6", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move knight")
	}
	legal = chessVerify("f1", "b5", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move bishop")
	}
	legal = chessVerify("f6", "g4", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move knight")
	}
	legal = chessVerify("e1", "g1", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move king")
	}
	legal = chessVerify("a7", "a6", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move pawn")
	}
	legal = chessVerify("g1", "h1", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move king")
	}
	legal = chessVerify("a6", "a5", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move pawn")
	}
	legal = chessVerify("f1", "g1", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move rook")
	}
	legal = chessVerify("g4", "f2", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move knight")
	}
	legal = chessVerify("a2", "a3", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestSmotheredMate() can't move pawn as white is in mate")
	}
	legal = chessVerify("h1", "h2", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestSmotheredMate() can't move king as pawn is there")
	}
	legal = chessVerify("b1", "c3", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestSmotheredMate() can't move knight as white is in mate")
	}

}
