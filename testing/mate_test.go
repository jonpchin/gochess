package testing

import (
	"testing"

	"github.com/jonpchin/gochess/gostuff"
)

func TestScholarMate(t *testing.T) {

	var legal bool
	const gameID = 0
	// pass in blank string for white and black players as spectators
	//don't matter for testing
	gostuff.InitGame(gameID, "", "")

	legal = gostuff.ChessVerify("e2", "e4", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestScholarMate() can't move two squares white")
	}
	legal = gostuff.ChessVerify("e7", "e5", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestScholarMate() can't move two squares black")
	}
	legal = gostuff.ChessVerify("f1", "c4", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestScholarMate() can't move bishop")
	}
	legal = gostuff.ChessVerify("b8", "c6", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestScholarMate() can't move knight")
	}
	legal = gostuff.ChessVerify("d1", "h5", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestScholarMate() can't move queen")
	}
	legal = gostuff.ChessVerify("g8", "f6", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestScholarMate() can't move knight")
	}
	legal = gostuff.ChessVerify("h5", "f7", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestScholarMate() can't move queen")
	}
	legal = gostuff.ChessVerify("e8", "f7", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestScholarMate() can't capture queen its mate")
	}
	legal = gostuff.ChessVerify("e8", "e7", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestScholarMate() can't move king its mate")
	}
	legal = gostuff.ChessVerify("a7", "a6", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestScholarMate() can't move pawn its mate")
	}
	legal = gostuff.ChessVerify("c6", "b4", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestScholarMate() can't move knight its mate")
	}
	legal = gostuff.ChessVerify("d8", "f7", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestScholarMate() can't move queen its mate")
	}
}

func TestSmotheredMate(t *testing.T) {

	var legal bool
	const gameID = 0
	// pass in blank string for white and black players as spectators
	//don't matter for testing
	gostuff.InitGame(gameID, "", "")

	legal = gostuff.ChessVerify("e2", "e4", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move two squares white")
	}
	legal = gostuff.ChessVerify("e7", "e5", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move two squares black")
	}
	legal = gostuff.ChessVerify("g1", "f3", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move knight")
	}
	legal = gostuff.ChessVerify("g8", "f6", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move knight")
	}
	legal = gostuff.ChessVerify("f1", "b5", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move bishop")
	}
	legal = gostuff.ChessVerify("f6", "g4", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move knight")
	}
	legal = gostuff.ChessVerify("e1", "g1", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move king")
	}
	legal = gostuff.ChessVerify("a7", "a6", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move pawn")
	}
	legal = gostuff.ChessVerify("g1", "h1", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move king")
	}
	legal = gostuff.ChessVerify("a6", "a5", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move pawn")
	}
	legal = gostuff.ChessVerify("f1", "g1", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move rook")
	}
	legal = gostuff.ChessVerify("g4", "f2", "", gameID)
	if legal == false {
		t.Error("mate_test.go TestSmotheredMate() can't move knight")
	}
	legal = gostuff.ChessVerify("a2", "a3", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestSmotheredMate() can't move pawn as white is in mate")
	}
	legal = gostuff.ChessVerify("h1", "h2", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestSmotheredMate() can't move king as pawn is there")
	}
	legal = gostuff.ChessVerify("b1", "c3", "", gameID)
	if legal == true {
		t.Error("mate_test.go TestSmotheredMate() can't move knight as white is in mate")
	}
}
