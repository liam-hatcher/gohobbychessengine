package chess

import "testing"

func TestGenerateOpeningPawnMoves_White(t *testing.T) {
	pos := NewPosition()

	moves := pos.GenerateWhitePawnMoves()

	expectedMoves := map[string]bool{
		"a2a3": true,
		"b2b3": true,
		"c2c3": true,
		"d2d3": true,
		"e2e3": true,
		"f2f3": true,
		"g2g3": true,
		"h2h3": true,
		"a2a4": true,
		"b2b4": true,
		"c2c4": true,
		"d2d4": true,
		"e2e4": true,
		"f2f4": true,
		"g2g4": true,
		"h2h4": true,
	}

	for _, move := range moves {
		uci := ToUCINotation(move)
		if !expectedMoves[uci] {
			t.Errorf("unexpected move: %s", uci)
		}
	}

	if len(moves) != 16 {
		t.Errorf("expected 16 moves but got %d", len(moves))
	}
}

func TestGenerateOpeningPawnMoves_Black(t *testing.T) {
	pos := NewPosition()

	moves := pos.GenerateBlackPawnMoves()

	expectedMoves := map[string]bool{
		"a7a6": true,
		"b7b6": true,
		"c7c6": true,
		"d7d6": true,
		"e7e6": true,
		"f7f6": true,
		"g7g6": true,
		"h7h6": true,
		"a7a5": true,
		"b7b5": true,
		"c7c5": true,
		"d7d5": true,
		"e7e5": true,
		"f7f5": true,
		"g7g5": true,
		"h7h5": true,
	}

	for _, move := range moves {
		uci := ToUCINotation(move)
		if !expectedMoves[uci] {
			t.Errorf("unexpected move: %s", uci)
		}
	}

	if len(moves) != 16 {
		t.Errorf("expected 16 moves but got %d", len(moves))
	}
}

func TestBlockedPawnMoves(t *testing.T) {
	pos := Position{}
	pos.SetPiece("P", "e4")
	pos.SetPiece("n", "e5")
	moves := pos.GenerateWhitePawnMoves()

	if len(moves) != 0 {
		t.Errorf("a white pawn on e4 has no moves if blocked on e5")
	}

	pos = Position{}
	pos.SetPiece("p", "d7")
	pos.SetPiece("Q", "d6")

	moves = pos.GenerateBlackPawnMoves()

	if len(moves) != 0 {
		t.Errorf("a black pawn on d7 has no moves if blocked on d6")
	}
}
