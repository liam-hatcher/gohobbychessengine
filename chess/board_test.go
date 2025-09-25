package chess

import (
	"strconv"
	"testing"
)

func assertEmpty(t *testing.T, pos *Position, square string) {
	t.Helper()
	piece, empty := pos.GetPieceOnSquare(square)

	if !empty {
		t.Errorf("Expected %s to be empty, but found %s", square, piece)
	}
}

func assertHasPiece(t *testing.T, pos *Position, expected string, square string) {
	t.Helper()

	piece, empty := pos.GetPieceOnSquare(square)
	if piece != expected {
		t.Errorf("expected to find piece %s on %s", expected, square)
	}
	if empty {
		t.Errorf("expected square %s to be occupied", square)
	}
}

// The Position struct has two internal representations of the board state:
//  1. Twelve different bit boards, (one for each piece type per color)
//  2. A one dimensional sparse matrix that allows O(1) lookups of a
//     piece by it's index.
//
// This test ensures the two remain in sync when setting and removing
// pieces from the board
func TestBitboardPieceMapSynchronization(t *testing.T) {
	pos := Position{}

	pieces := map[string]*Bitboard{
		"P": &pos.WhitePawns,
		"R": &pos.WhiteRooks,
		"N": &pos.WhiteKnights,
		"B": &pos.WhiteBishops,
		"Q": &pos.WhiteQueens,
		"K": &pos.WhiteKing,
		"p": &pos.BlackPawns,
		"r": &pos.BlackRooks,
		"n": &pos.BlackKnights,
		"b": &pos.BlackBishops,
		"q": &pos.BlackQueens,
		"k": &pos.BlackKing,
	}

	for file := 'a'; file <= 'h'; file++ {
		for rank := 1; rank <= 8; rank++ {
			for piece, bb := range pieces {
				square := string(file) + strconv.Itoa(rank)
				assertEmpty(t, &pos, square)
				if *bb != 0 {
					t.Errorf("expected Bitboard for %s to be empty", piece)
				}

				pos.SetPiece(piece, square)
				assertHasPiece(t, &pos, piece, square)
				if *bb == 0 {
					t.Errorf("expected Bitboard for %s to NOT be empty", piece)
					t.Errorf("failed to place %s on %s", piece, square)
				}

				pos.RemovePiece(square)
				assertEmpty(t, &pos, square)
				if *bb != 0 {
					t.Errorf("expected Bitboard for %s to be empty after removal on %s", piece, square)
				}
			}
		}
	}
}

func TestWhitePawnSinglePushes(t *testing.T) {
	for file := 'a'; file <= 'h'; file++ {
		// white pawns start on rank 2. another test will cover pushing
		// from rank 7 (promotions), so we can test ranks 2-6 here.
		for rank := 2; rank <= 6; rank++ {
			pos := Position{}
			from := string(file) + strconv.Itoa(rank)
			to := string(file) + strconv.Itoa(rank+1)

			pos.SetPiece("P", from)
			assertEmpty(t, &pos, to)           // target square is empty
			assertHasPiece(t, &pos, "P", from) // origin square has a Pawn

			pos.ApplyMove(from + to)
			assertEmpty(t, &pos, from)       // origin square is now empty
			assertHasPiece(t, &pos, "P", to) // target sqaure has a pawn
		}
	}
}

func TestBlackPawnSinglePushes(t *testing.T) {
	for file := 'a'; file <= 'h'; file++ {
		// black pawns start on rank 7. another test will cover pushing
		// from rank 2 (promotions), so we can test ranks 7-3 here.
		for rank := 7; rank >= 3; rank-- {
			pos := Position{}
			from := string(file) + strconv.Itoa(rank)
			to := string(file) + strconv.Itoa(rank-1)

			pos.SetPiece("p", from)
			assertEmpty(t, &pos, to)           // target square is empty
			assertHasPiece(t, &pos, "p", from) // origin square has a Pawn

			pos.ApplyMove(from + to)
			assertEmpty(t, &pos, from)       // origin square is now empty
			assertHasPiece(t, &pos, "p", to) // target sqaure has a pawn
		}
	}
}

func TestPawnDoublePushes(t *testing.T) {
	for file := 'a'; file <= 'h'; file++ {
		// white pawns
		pos := Position{}
		from := string(file) + "2"
		to := string(file) + "4"
		pos.SetPiece("P", from)

		assertEmpty(t, &pos, to)
		assertHasPiece(t, &pos, "P", from)
		pos.ApplyMove(from + to)
		assertEmpty(t, &pos, from)
		assertHasPiece(t, &pos, "P", to)

		// black pawns
		pos = Position{
			BlackPawns: 0,
			PieceMap:   [64]string{},
		}
		from = string(file) + "7"
		to = string(file) + "5"
		pos.SetPiece("p", from)
		assertEmpty(t, &pos, to)
		assertHasPiece(t, &pos, "p", from)
		pos.ApplyMove(from + to)
		assertEmpty(t, &pos, from)
		assertHasPiece(t, &pos, "p", to)
	}
}

// TODO: add more test cases
func TestWhitePawnCaptures(t *testing.T) {
	pos := Position{}
	pos.SetPiece("P", "e4")
	pos.SetPiece("p", "d5")

	if pos.BlackPawns == 0 {
		t.Errorf("black pawn bit should be set")
	}

	pos.ApplyMove("e4d5")
	assertEmpty(t, &pos, "e4")
	assertHasPiece(t, &pos, "P", "d5")
	if pos.BlackPawns != 0 {
		t.Errorf("failed to update bitboard")
	}
}

func TestWhitePawnEnPassant(t *testing.T) {
	pos := Position{}
	pos.SetPiece("P", "e5")
	pos.SetPiece("p", "d7")
	pos.ApplyMove("d7d5")

	assertEmpty(t, &pos, "d7")
	if pos.EnPassantTarget == 0 {
		t.Errorf("failed to set en passant target square")
	}

	pos.ApplyMove("e5d6")
	assertEmpty(t, &pos, "d5")
	assertHasPiece(t, &pos, "P", "d6")
	if pos.BlackPawns != 0 {
		t.Errorf("failed to update bitboard")
	}
	if pos.EnPassantTarget != 0 {
		t.Errorf("failed to clear en passant target square")
	}
}

func TestBlackPawnEnPassant(t *testing.T) {
	pos := Position{}
	pos.SetPiece("p", "d4")
	pos.SetPiece("P", "e2")
	pos.ApplyMove("e2e4")

	assertEmpty(t, &pos, "e2")
	if pos.EnPassantTarget == 0 {
		t.Errorf("failed to set en passant target square")
	}

	pos.ApplyMove("d4e3")
	assertEmpty(t, &pos, "e4")
	assertHasPiece(t, &pos, "p", "e3")
	if pos.WhitePawns != 0 {
		t.Errorf("failed to update bitboard")
	}
	if pos.EnPassantTarget != 0 {
		t.Errorf("failed to clear en passant target square")
	}
}
