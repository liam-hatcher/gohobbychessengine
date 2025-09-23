package chess

import "testing"

func assertEqualMoves(t *testing.T, moves []Move, expected map[string]bool) {
	t.Helper()
	for _, move := range moves {
		uci := ToUCINotation(move)
		if !expected[uci] {
			t.Errorf("unexpected move: %s", uci)
		}
	}

	actualLength := len(moves)
	expectedLength := len(expected)
	if len(moves) != len(expected) {
		t.Errorf("expected moves to have length %d but got %d instead", expectedLength, actualLength)
	}
}

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

	assertEqualMoves(t, moves, expectedMoves)
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

	captures := pos.GenerateWhitePawnCaptures()
	if len(captures) != 0 {
		t.Error("expected white pawn to have no captures")
	}

	pos = Position{}
	pos.SetPiece("p", "d7")
	pos.SetPiece("Q", "d6")

	moves = pos.GenerateBlackPawnMoves()

	if len(moves) != 0 {
		t.Errorf("a black pawn on d7 has no moves if blocked on d6")
	}
}

func TestWhitePawnCaptures(t *testing.T) {
	pos := Position{}
	pos.SetPiece("P", "d4")
	pos.SetPiece("p", "e5")
	pos.SetPiece("r", "c5")
	captures := pos.GenerateWhitePawnCaptures()
	expectedMoves := map[string]bool{
		"d4e5": true,
		"d4c5": true,
	}
	assertEqualMoves(t, captures, expectedMoves)

	pos = Position{}
	pos.SetPiece("P", "c3")
	pos.SetPiece("P", "d4") // make sure we can't capture our own piece
	pos.SetPiece("b", "b4")
	captures = pos.GenerateWhitePawnCaptures()
	if len(captures) != 1 {
		t.Errorf("expected white pawn on c3 to have 1 possible capture, but got %d", len(captures))
	}

	// 'a' file captures
	pos = Position{}
	pos.SetPiece("P", "a3")
	pos.SetPiece("q", "b4")
	captures = pos.GenerateWhitePawnCaptures()
	if len(captures) != 1 {
		t.Errorf("expected white pawn on a3 to have 1 possible capture, but got %d", len(captures))
	}

	// h file captures
	pos = Position{}
	pos.SetPiece("P", "h7")
	pos.SetPiece("q", "g8")
	captures = pos.GenerateWhitePawnCaptures()
	if len(captures) != 1 {
		t.Errorf("expected white pawn on h7 to have 1 possible capture, but got %d", len(captures))
	}
	expected := "h7g8"
	move := ToUCINotation(captures[0])
	if move != expected {
		t.Errorf("expected capture move %s but got %s", expected, move)
	}
}

func TestBlackPawnCaptures(t *testing.T) {
	pos := Position{}
	pos.SetPiece("p", "d5")
	pos.SetPiece("P", "e4")
	pos.SetPiece("N", "c4")
	captures := pos.GenerateBlackPawnCaptures()
	expectedMoves := map[string]bool{
		"d5e4": true,
		"d5c4": true,
	}
	assertEqualMoves(t, captures, expectedMoves)

	pos = Position{}
	pos.SetPiece("p", "c4")
	pos.SetPiece("p", "d3") // make sure we can't capture our own piece
	pos.SetPiece("B", "b3")
	captures = pos.GenerateBlackPawnCaptures()
	if len(captures) != 1 {
		t.Errorf("expected black pawn on c4 to have 1 possible capture, but got %d", len(captures))
	}

	// 'a' file captures
	pos = Position{}
	pos.SetPiece("p", "a6")
	pos.SetPiece("Q", "b5")
	captures = pos.GenerateBlackPawnCaptures()
	if len(captures) != 1 {
		t.Errorf("expected black pawn on a6 to have 1 possible capture, but got %d", len(captures))
	}

	// h file captures
	pos = Position{}
	pos.SetPiece("p", "h2")
	pos.SetPiece("Q", "g1")
	captures = pos.GenerateBlackPawnCaptures()
	if len(captures) != 1 {
		t.Errorf("expected black pawn on h2 to have 1 possible capture, but got %d", len(captures))
	}
	expected := "h2g1"
	move := ToUCINotation(captures[0])
	if move != expected {
		t.Errorf("expected capture move %s but got %s", expected, move)
	}
}
