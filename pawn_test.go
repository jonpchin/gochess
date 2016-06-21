package gostuff

import(
	"testing"
)


func TestPawn(t *testing.T){

	var legal bool
	const gameID = 0
	InitGame(gameID)
	
	legal = ChessVerify("e2", "e4", gameID)
	if legal == false{
		t.Error("pawn_test.go TestPawn() 1 can't move two squares white")
	}
	legal = ChessVerify("e7", "e5", gameID)
	if legal == false{
		t.Error("pawn_test.go TestPawn() 2 can't move two squares black")
	}
	legal = ChessVerify("d2", "d5", gameID) 
	if legal == true{
		t.Error("pawn_test.go TestPawn() 3, illegal three square move")
	}
	legal = ChessVerify("d2", "d3", gameID)
	if legal == false{
		t.Error("pawn_test.go TestPawn() 4 can't move one square white")
	}
	legal = ChessVerify("d2", "d3", gameID)
	if legal == true{
		t.Error("connect pawn_test.go TestPawn() 5, no pawn on d2")
	}
	legal = ChessVerify("e5", "e4", gameID)
	if legal == true{
		t.Error("connect pawn_test.go TestPawn() 6 pawn is blocking square")
	}
	legal = ChessVerify("f7", "f5", gameID)
	if legal == false{
		t.Error("connect pawn_test.go TestPawn() 7 black moves pawn up two squares")
	}
	legal = ChessVerify("e4", "f6", gameID)
	if legal == true{
		t.Error("connect pawn_test.go TestPawn() 8 illegal white move")
	}
	legal = ChessVerify("e4", "f5", gameID)
	if legal == false{
		t.Error("connect pawn_test.go TestPawn() 9 white pawn captures")
	}
	legal = ChessVerify("g7", "g5", gameID)
	if legal == false{
		t.Error("connect pawn_test.go TestPawn() 10 black pawn two squares")
	}
	legal = ChessVerify("f5", "g6", gameID)
	if legal == false{
		t.Error("connect pawn_test.go TestPawn() 11 enpassent")
	}
	legal = ChessVerify("h7", "g6", gameID)
	if legal == false{
		t.Error("connect pawn_test.go TestPawn() black captures")
	}
	legal = ChessVerify("d3", "d2", gameID)
	if legal == true{
		t.Error("connect pawn_test.go TestPawn() pawn cannot go backwards")
	}
}