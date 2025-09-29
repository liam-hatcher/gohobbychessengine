package chess

import (
	"fmt"
	"strconv"
	"testing"
)

func assertEmpty(t *testing.T, pos *Position, square string) {
	t.Helper()
	piece, empty := pos.GetPieceOnSquare(square)

	if !empty {
		t.Errorf("Expected %s to be empty, but found %c", square, piece)
	}
}

func assertHasPiece(t *testing.T, pos *Position, expected byte, square string) {
	t.Helper()

	piece, empty := pos.GetPieceOnSquare(square)
	if piece != expected {
		t.Errorf("expected to find piece %c on %s", expected, square)
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

	pieces := map[byte]*Bitboard{
		'P': &pos.WhitePawns,
		'R': &pos.WhiteRooks,
		'N': &pos.WhiteKnights,
		'B': &pos.WhiteBishops,
		'Q': &pos.WhiteQueens,
		'K': &pos.WhiteKing,
		'p': &pos.BlackPawns,
		'r': &pos.BlackRooks,
		'n': &pos.BlackKnights,
		'b': &pos.BlackBishops,
		'q': &pos.BlackQueens,
		'k': &pos.BlackKing,
	}

	for file := 'a'; file <= 'h'; file++ {
		for rank := 1; rank <= 8; rank++ {
			for piece, bb := range pieces {
				square := string(file) + strconv.Itoa(rank)
				assertEmpty(t, &pos, square)
				if *bb != 0 {
					t.Errorf("expected Bitboard for %c to be empty", piece)
				}

				pos.SetPiece(piece, square)
				assertHasPiece(t, &pos, piece, square)
				if *bb == 0 {
					t.Errorf("expected Bitboard for %c to NOT be empty", piece)
					t.Errorf("failed to place %c on %s", piece, square)
				}

				pos.RemovePiece(square)
				assertEmpty(t, &pos, square)
				if *bb != 0 {
					t.Errorf("expected Bitboard for %c to be empty after removal on %s", piece, square)
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

			pos.SetPiece('P', from)
			assertEmpty(t, &pos, to)           // target square is empty
			assertHasPiece(t, &pos, 'P', from) // origin square has a Pawn

			pos.ApplyMove(from + to)
			assertEmpty(t, &pos, from)       // origin square is now empty
			assertHasPiece(t, &pos, 'P', to) // target sqaure has a pawn
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

			pos.SetPiece('p', from)
			assertEmpty(t, &pos, to)           // target square is empty
			assertHasPiece(t, &pos, 'p', from) // origin square has a Pawn

			pos.ApplyMove(from + to)
			assertEmpty(t, &pos, from)       // origin square is now empty
			assertHasPiece(t, &pos, 'p', to) // target sqaure has a pawn
		}
	}
}

func TestPawnDoublePushes(t *testing.T) {
	for file := 'a'; file <= 'h'; file++ {
		// white pawns
		pos := Position{}
		from := string(file) + "2"
		to := string(file) + "4"
		pos.SetPiece('P', from)

		assertEmpty(t, &pos, to)
		assertHasPiece(t, &pos, 'P', from)
		pos.ApplyMove(from + to)
		assertEmpty(t, &pos, from)
		assertHasPiece(t, &pos, 'P', to)

		// black pawns
		pos = Position{}
		from = string(file) + "7"
		to = string(file) + "5"
		pos.SetPiece('p', from)
		assertEmpty(t, &pos, to)
		assertHasPiece(t, &pos, 'p', from)
		pos.ApplyMove(from + to)
		assertEmpty(t, &pos, from)
		assertHasPiece(t, &pos, 'p', to)
	}
}

func TestWhitePawnCaptures(t *testing.T) {
	pieces := []byte{'p', 'r', 'b', 'n', 'q'}

	testCaptures := func(file rune, rank int, dir int) {
		from := fmt.Sprintf("%c%d", file, rank)
		to := fmt.Sprintf("%c%d", file+rune(dir), rank+1)

		t.Run(
			fmt.Sprintf("capture move %s%s", from, to),
			func(t *testing.T) {
				for _, piece := range pieces {
					pos := Position{}
					pos.SetPiece('P', from)
					pos.SetPiece(piece, to)
					pos.ApplyMove(from + to)

					if pos.BlackPieces() != 0 {
						t.Errorf("failed to capture %c on %s", piece, to)
					}
				}
			})
	}

	// right captures
	for file := 'a'; file < 'h'; file++ {
		// will test captures on the back rank in a seperate test
		for rank := 2; rank < 7; rank++ {
			testCaptures(file, rank, 1)
		}
	}

	// left captures
	for file := 'b'; file <= 'h'; file++ {
		// will test captures on the back rank in a seperate test
		for rank := 2; rank < 7; rank++ {
			testCaptures(file, rank, -1)
		}
	}
}

func TestBlackPawnCaptures(t *testing.T) {
	pieces := []byte{'P', 'R', 'B', 'N', 'Q'}

	testCaptures := func(file rune, rank int, dir int) {
		from := fmt.Sprintf("%c%d", file, rank)
		to := fmt.Sprintf("%c%d", file+rune(dir), rank-1)

		t.Run(
			fmt.Sprintf("capture move %s%s", from, to),
			func(t *testing.T) {
				for _, piece := range pieces {
					pos := Position{}
					pos.SetPiece('p', from)
					pos.SetPiece(piece, to)
					pos.ApplyMove(from + to)

					if pos.WhitePieces() != 0 {
						t.Errorf("failed to capture %c on %s", piece, to)
					}
				}
			})
	}

	// right captures
	for file := 'a'; file < 'h'; file++ {
		// will test captures on the back rank in a seperate test
		for rank := 7; rank > 2; rank-- {
			testCaptures(file, rank, 1)
		}
	}

	// left captures
	for file := 'b'; file <= 'h'; file++ {
		// will test captures on the back rank in a seperate test
		for rank := 7; rank > 2; rank-- {
			testCaptures(file, rank, -1)
		}
	}
}

func TestWhitePawnEnPassant(t *testing.T) {
	pos := Position{}
	pos.SetPiece('P', "e5")
	pos.SetPiece('p', "d7")
	pos.ApplyMove("d7d5")

	assertEmpty(t, &pos, "d7")
	if pos.EnPassantTarget == 0 {
		t.Errorf("failed to set en passant target square")
	}

	pos.ApplyMove("e5d6")
	assertEmpty(t, &pos, "d5")
	assertHasPiece(t, &pos, 'P', "d6")
	if pos.BlackPawns != 0 {
		t.Errorf("failed to update bitboard")
	}
	if pos.EnPassantTarget != 0 {
		t.Errorf("failed to clear en passant target square")
	}
}

func TestBlackPawnEnPassant(t *testing.T) {
	pos := Position{}
	pos.SetPiece('p', "d4")
	pos.SetPiece('P', "e2")
	pos.ApplyMove("e2e4")

	assertEmpty(t, &pos, "e2")
	if pos.EnPassantTarget == 0 {
		t.Errorf("failed to set en passant target square")
	}

	pos.ApplyMove("d4e3")
	assertEmpty(t, &pos, "e4")
	assertHasPiece(t, &pos, 'p', "e3")
	if pos.WhitePawns != 0 {
		t.Errorf("failed to update bitboard")
	}
	if pos.EnPassantTarget != 0 {
		t.Errorf("failed to clear en passant target square")
	}
}

func TestPawnPromotions(t *testing.T) {
	promotions := [4]byte{'q', 'r', 'b', 'n'}

	tests := []struct {
		color     string
		pawn      byte
		startRank int
		endRank   int
	}{
		{"white", 'P', 7, 8},
		{"black", 'p', 2, 1},
	}

	for _, test := range tests {
		for file := 'a'; file <= 'h'; file++ {
			for _, promo := range promotions {
				from := fmt.Sprintf("%c%d", file, test.startRank)
				to := fmt.Sprintf("%c%d", file, test.endRank)
				move := from + to + string(promo)

				t.Run(fmt.Sprintf("%s pawn promotion %s", test.color, move), func(t *testing.T) {
					pos := Position{} // init inside the test
					pos.SetPiece(test.pawn, from)
					pos.ApplyMove(move)

					assertEmpty(t, &pos, from)

					// check pawn bitboard
					if test.pawn == 'P' && pos.WhitePawns != 0 {
						t.Errorf("expected white pawns bitboard to be empty")
					} else if test.pawn == 'p' && pos.BlackPawns != 0 {
						t.Errorf("expected black pawns bitboard to be empty")
					}

					piece := promo
					if test.pawn == 'P' {
						piece = ToUpper(promo)
					}

					assertHasPiece(t, &pos, piece, to)

					// check the right promoted piece bitboard
					var bb Bitboard
					switch promo {
					case 'q':
						if test.color == "white" {
							bb = pos.WhiteQueens
						} else {
							bb = pos.BlackQueens
						}
					case 'r':
						if test.color == "white" {
							bb = pos.WhiteRooks
						} else {
							bb = pos.BlackRooks
						}
					case 'b':
						if test.color == "white" {
							bb = pos.WhiteBishops
						} else {
							bb = pos.BlackBishops
						}
					case 'n':
						if test.color == "white" {
							bb = pos.WhiteKnights
						} else {
							bb = pos.BlackKnights
						}
					}

					if bb == 0 {
						t.Errorf("expected %c bitboard to be set", piece)
					}
				})
			}
		}
	}
}

func TestPawnCaptureAndPromote(t *testing.T) {
	pos := Position{}
	pos.SetPiece('n', "c8")
	pos.SetPiece('P', "b7")
	pos.ApplyMove("b7c8q")

	if pos.WhitePawns != 0 {
		t.Errorf("expected white pawn bitboard to be empty")
	}
	assertHasPiece(t, &pos, 'Q', "c8")
	if pos.WhiteQueens == 0 {
		t.Errorf("expected white queen bitboard to be set")
	}
}

func TestWhiteRookMoves(t *testing.T) {
	pos := Position{
		WhiteCastlingRights: CastlingRights{
			Short: true,
			Long:  true,
		},
	}
	pos.SetPiece('R', "a1")
	pos.SideToMove = "white"

	from := "a1"
	to := "a8"

	// Preconditions
	assertHasPiece(t, &pos, 'R', from)
	assertEmpty(t, &pos, to)

	// Apply rook move
	pos.ApplyMove(from + to)

	// Postconditions
	assertEmpty(t, &pos, from)
	assertHasPiece(t, &pos, 'R', to)

	// Castling rights: moving rook from a1 disables white's long castle
	if pos.WhiteCastlingRights.Long {
		t.Errorf("Expected long castling rights to be false after moving rook from a1")
	}
	if !pos.WhiteCastlingRights.Short {
		t.Errorf("Expected short castling rights to remain true after rook move from a1")
	}
}

func TestBlackRookMoves(t *testing.T) {
	pos := Position{
		BlackCastlingRights: CastlingRights{
			Short: true,
			Long:  true,
		},
	}
	pos.SetPiece('r', "h8")
	pos.SideToMove = "black"

	from := "h8"
	to := "h1"

	// Preconditions
	assertHasPiece(t, &pos, 'r', from)
	assertEmpty(t, &pos, to)

	// Apply rook move
	pos.ApplyMove(from + to)

	// Postconditions
	assertEmpty(t, &pos, from)
	assertHasPiece(t, &pos, 'r', to)

	// Castling rights: moving rook from a1 disables white's long castle
	if pos.BlackCastlingRights.Short {
		t.Errorf("Expected short castling rights to be false after moving rook from h8")
	}
	if !pos.BlackCastlingRights.Long {
		t.Errorf("Expected long castling rights to remain true after rook move from h8")
	}
}
