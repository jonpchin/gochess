package testing

import(
	"testing"
	"github.com/jonpchin/GoChess/gostuff"
)


func TestPawnMove(t *testing.T){

	var legal bool
	gostuff.InitGame(0)
	legal = gostuff.ChessVerify("e2", "e4", 0)
	if legal == false{
		t.Error("Illegal pawn move 1 connect pawn_test.go TestPawnMove()")
	}
	legal = gostuff.ChessVerify("e7", "e5", 0)
	if legal == false{
		t.Error("Illegal pawn move 2 connect pawn_test.go TestPawnMove()")
	}
	legal = gostuff.ChessVerify("d2", "d4", 0)
	if legal == false{
		t.Error("Illegal pawn move 3 connect pawn_test.go TestPawnMove()")
	}
}

func TestPawnCapture(t *testing.T){
	
}

func TestEnpassent(t *testing.T){
	
}