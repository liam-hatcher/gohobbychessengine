package chess

import (
	"fmt"
	"testing"
)

func TestGenerateRookMoves_Corner(t *testing.T) {
	pos := Position{}
	pos.SetPiece('R', "a1")
	pos.SideToMove = "white"

	moves := pos.GenerateRookMoves()

	got := map[string]bool{}
	for _, m := range moves {
		key := fmt.Sprintf("%s->%s", BitIndexToRankFile(m.From), BitIndexToRankFile(m.To))
		got[key] = true
	}

	// Expected rook moves from a1 on an empty board
	expected := []string{
		"a1->a2", "a1->a3", "a1->a4", "a1->a5", "a1->a6", "a1->a7", "a1->a8",
		"a1->b1", "a1->c1", "a1->d1", "a1->e1", "a1->f1", "a1->g1", "a1->h1",
	}

	for _, exp := range expected {
		if !got[exp] {
			t.Errorf("Expected move %s not generated", exp)
		}
	}

	// Optional: check no extra moves were generated
	if len(got) != len(expected) {
		t.Errorf("Unexpected number of moves: got %d, expected %d", len(got), len(expected))
	}
}

func TestGenerateRookMoves_Center(t *testing.T) {
	pos := Position{}
	pos.SetPiece('r', "d4")
	pos.SideToMove = "black"

	moves := pos.GenerateRookMoves()

	// Build a map of "from->to" strings for easy checking
	got := map[string]bool{}
	for _, m := range moves {
		key := fmt.Sprintf("%s->%s", BitIndexToRankFile(m.From), BitIndexToRankFile(m.To))
		got[key] = true
	}

	// Expected rook moves from d4 on an empty board
	expected := []string{
		// Vertical moves
		"d4->d1", "d4->d2", "d4->d3", "d4->d5", "d4->d6", "d4->d7", "d4->d8",
		// Horizontal moves
		"d4->a4", "d4->b4", "d4->c4", "d4->e4", "d4->f4", "d4->g4", "d4->h4",
	}

	for _, exp := range expected {
		if !got[exp] {
			t.Errorf("Expected move %s not generated", exp)
		}
	}

	// Optional: check no extra moves were generated
	if len(got) != len(expected) {
		t.Errorf("Unexpected number of moves: got %d, expected %d", len(got), len(expected))
	}
}

func TestGenerateRookMoves_CenterWithBlockers(t *testing.T) {
	pos := Position{}
	pos.SetPiece('R', "d4")
	pos.SetPiece('p', "d6") // enemy piece
	pos.SetPiece('N', "f4") // friendly piece
	pos.SideToMove = "white"

	moves := pos.GenerateRookMoves()

	got := map[string]bool{}
	for _, m := range moves {
		key := fmt.Sprintf("%s->%s", BitIndexToRankFile(m.From), BitIndexToRankFile(m.To))
		got[key] = true
	}

	// Expected moves from d4 with blockers
	expected := []string{
		// Vertical moves
		"d4->d5",
		"d4->d6", // capture
		"d4->d3",
		"d4->d2",
		"d4->d1",
		// Horizontal moves
		"d4->c4",
		"d4->b4",
		"d4->a4",
		"d4->e4",
		// Note: f4 is blocked by friendly piece, so not included
	}

	for _, exp := range expected {
		if !got[exp] {
			t.Errorf("Expected move %s not generated", exp)
		}
	}

	if len(got) != len(expected) {
		t.Errorf("Unexpected number of moves: got %d, expected %d", len(got), len(expected))
	}
}
