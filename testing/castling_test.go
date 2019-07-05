package testing

import (
	"testing"

	"github.com/notnil/chess"
)

// Tests king castle
func TestKingCastle(t *testing.T) {

	validator := chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}))

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
		err := validator.MoveStr(c.src + c.dst)
		if (c.legal && err == nil) || (!c.legal && err != nil) {
			t.Error(c.msg, err)
		}
	}
}

func TestQueenCastle(t *testing.T) {

	validator := chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}))

	cases := []struct {
		src   string
		dst   string
		legal bool
		msg   string
	}{
		{"d2", "d4", false, "d2 e4 can't move two squares white"},
		{"d7", "d5", false, "d7 e5 can't move two squares black"},
		{"c1", "g5", false, "c1 g5 can't move bishop"},
		{"c8", "g4", false, "c8 g4 can't move bishop"},
		{"b1", "c3", false, "b1 c3 can't move knight"},
		{"e8", "c8", true, "e8 c8  can't castle queenside with queen blocking"},
		{"d8", "d7", false, "d8 d7 can't move queen"},
		{"d1", "d2", false, "d1 d2 can't move queen"},
		{"b8", "c6", false, "b8 c6 can't move knight"},
		{"e1", "c1", false, "e1 c1 can't king castle"},
		{"e8", "d8", false, "e8 d8 can't move king"},
		{"c1", "b1", false, "c1 b1 can't move king"},
		{"d8", "e8", false, "d8 e8 can't move king"},
		{"d1", "e1", false, "d1 e1 can't move rook"},
		{"e8", "c8", true, "e8 c8 can't castle queenside as king already move"},
	}

	for _, c := range cases {
		err := validator.MoveStr(c.src + c.dst)
		if (c.legal && err == nil) || (!c.legal && err != nil) {
			t.Error(c.msg, err)
		}
	}
}
