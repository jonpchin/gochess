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
		t.Error("connect pawn_test.go TestPawn() 5, no pawn on d2")
	}
	legal = chessVerify("e5", "e4", gameID)
	if legal == true{
		t.Error("connect pawn_test.go TestPawn() 6 pawn is blocking square")
	}
	legal = chessVerify("f7", "f5", gameID)
	if legal == false{
		t.Error("connect pawn_test.go TestPawn() 7 black moves pawn up two squares")
	}
	legal = chessVerify("e4", "f6", gameID)
	if legal == true{
		t.Error("connect pawn_test.go TestPawn() 8 illegal white move")
	}
	legal = chessVerify("e4", "f5", gameID)
	if legal == false{
		t.Error("connect pawn_test.go TestPawn() 9 white pawn captures")
	}
	legal = chessVerify("g7", "g5", gameID)
	if legal == false{
		t.Error("connect pawn_test.go TestPawn() 10 black pawn two squares")
	}
	legal = chessVerify("f5", "g6", gameID)
	if legal == false{
		t.Error("connect pawn_test.go TestPawn() 11 enpassent")
	}
	legal = chessVerify("h7", "g6", gameID)
	if legal == false{
		t.Error("connect pawn_test.go TestPawn() black captures")
	}
	legal = chessVerify("d3", "d2", gameID)
	if legal == true{
		t.Error("connect pawn_test.go TestPawn() pawn cannot go backwards")
	}
}