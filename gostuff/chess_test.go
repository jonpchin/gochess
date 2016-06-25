package gostuff

import(
	"testing"
)

func TestPawn(t *testing.T){

	var legal bool
	const gameID = 0
	initGame(gameID)
	
	legal = chessVerify("e2", "e4", gameID)
	if legal == false{
		t.Error("pawn_test.go TestPawn() 1 can't move two squares white")
	}
	legal = chessVerify("e7", "e5", gameID)
	if legal == false{
		t.Error("pawn_test.go TestPawn() 2 can't move two squares black")
	}
	legal = chessVerify("d2", "d5", gameID) 
	if legal == true{
		t.Error("pawn_test.go TestPawn() 3, illegal three square move")
	}
	legal = chessVerify("d2", "d3", gameID)
	if legal == false{
		t.Error("pawn_test.go TestPawn() 4 can't move one square white")
	}
	legal = chessVerify("d2", "d3", gameID)
	if legal == true{
		t.Error("pawn_test.go TestPawn() 5, no pawn on d2")
	}
	legal = chessVerify("e5", "e4", gameID)
	if legal == true{
		t.Error("pawn_test.go TestPawn() 6 pawn is blocking square")
	}
	legal = chessVerify("f7", "f5", gameID)
	if legal == false{
		t.Error("pawn_test.go TestPawn() 7 black moves pawn up two squares")
	}
	legal = chessVerify("e4", "f6", gameID)
	if legal == true{
		t.Error("pawn_test.go TestPawn() 8 illegal white move")
	}
	legal = chessVerify("e4", "f5", gameID)
	if legal == false{
		t.Error("pawn_test.go TestPawn() 9 white pawn captures")
	}
	legal = chessVerify("g7", "g5", gameID)
	if legal == false{
		t.Error("pawn_test.go TestPawn() 10 black pawn two squares")
	}
	legal = chessVerify("f5", "g6", gameID)
	if legal == false{
		t.Error("pawn_test.go TestPawn() 11 enpassent")
	}
	legal = chessVerify("h7", "g6", gameID)
	if legal == false{
		t.Error("pawn_test.go TestPawn() black captures")
	}
	legal = chessVerify("d3", "d2", gameID)
	if legal == true{
		t.Error("pawn_test.go TestPawn() pawn cannot go backwards")
	}
}

func TestKnght(t *testing.T){

	var legal bool
	const gameID = 0
	initGame(gameID)
	
	legal = chessVerify("g1", "e2", gameID)
	if legal == true{
		t.Error("knight_test.go 1 test white self capture failed")
	}
	legal = chessVerify("g1", "f3", gameID)
	if legal == false{
		t.Error("knight_test.go 2")
	}
	legal = chessVerify("g8", "h6", gameID)
	if legal == false{
		t.Error("knight_test.go 3")
	}
	legal = chessVerify("f3", "g4", gameID)
	if legal == true{
		t.Error("knight_test.go 4 testing illegal move")
	}
	legal = chessVerify("f3", "h4", gameID)
	if legal == false{
		t.Error("knight_test.go 5")
	}
	legal = chessVerify("b8", "c6", gameID)
	if legal == false{
		t.Error("knight_test.go 6")
	}
	legal = chessVerify("h4", "f5", gameID)
	if legal == false{
		t.Error("knight_test.go 7")
	}
	legal = chessVerify("h6", "f5", gameID)
	if legal == false{
		t.Error("knight_test.go 8 black captures white knight")
	}
}