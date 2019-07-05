package testing

import (
	"testing"

	"github.com/notnil/chess"
)

func TestPawn(t *testing.T) {

	validator := chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}))

	cases := []struct {
		src   string
		dst   string
		legal bool
	}{
		{"e2", "e4", false},
		{"e7", "e5", false},
		{"d2", "d5", true},
		{"d2", "d3", false},
		{"e5", "e4", true},
		{"f7", "f5", false},
		{"e4", "f6", true},
		{"e4", "f5", false},
		{"g7", "g5", false},
		{"f5", "g6", false},
		{"h7", "g6", false},
		{"d3", "d2", true},
	}

	for index, c := range cases {
		err := validator.MoveStr(c.src + c.dst)
		if (c.legal && err == nil) || (!c.legal && err != nil) {
			t.Error("testPawn", index, err)
		}
	}
}

func TestKnght(t *testing.T) {

	validator := chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}))

	cases := []struct {
		src   string
		dst   string
		legal bool
	}{
		{"g1", "e2", true},
		{"g1", "f3", false},
		{"g8", "h6", false},
		{"f3", "g4", true},
		{"f3", "h4", false},
		{"b8", "c6", false},
		{"h4", "f5", false},
		{"h6", "f5", false},
	}

	for index, c := range cases {
		err := validator.MoveStr(c.src + c.dst)
		if (c.legal && err == nil) || (!c.legal && err != nil) {
			t.Error("testKnight", index, err)
		}
	}
}
