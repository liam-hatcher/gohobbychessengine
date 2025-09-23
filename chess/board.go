package chess

import "fmt"

type Bitboard uint64

type Move struct {
	From int
	To   int
	// Promo rune // handle later
}

type Position struct {
	WhitePawns   Bitboard
	WhiteKnights Bitboard
	WhiteBishops Bitboard
	WhiteRooks   Bitboard
	WhiteQueens  Bitboard
	WhiteKing    Bitboard

	BlackPawns   Bitboard
	BlackKnights Bitboard
	BlackBishops Bitboard
	BlackRooks   Bitboard
	BlackQueens  Bitboard
	BlackKing    Bitboard
}

func NewPosition() *Position {
	return &Position{
		WhitePawns:   0x000000000000FF00,
		WhiteRooks:   0x0000000000000081,
		WhiteKnights: 0x0000000000000042,
		WhiteBishops: 0x0000000000000024,
		WhiteQueens:  0x0000000000000008,
		WhiteKing:    0x0000000000000010,

		BlackPawns:   0x00FF000000000000,
		BlackRooks:   0x8100000000000000,
		BlackKnights: 0x4200000000000000,
		BlackBishops: 0x2400000000000000,
		BlackQueens:  0x0800000000000000,
		BlackKing:    0x1000000000000000,
	}
}

func (p *Position) WhitePieces() Bitboard {
	return p.WhiteBishops | p.WhiteKing | p.WhiteKnights | p.WhitePawns | p.WhiteRooks | p.WhiteQueens
}

func (p *Position) BlackPieces() Bitboard {
	return p.BlackBishops | p.BlackKing | p.BlackKnights | p.BlackPawns | p.BlackRooks | p.BlackQueens
}

func (p *Position) SetPiece(piece, square string) {
	index := rankFileToBitIndex(square[0], square[1])
	mask := Bitboard(1) << index

	if mask&p.GetOccupiedSquares() != 0 {
		panic(fmt.Sprintf("%s is already occupied", square))
	}

	switch piece {
	case "P":
		p.WhitePawns |= mask
	case "R":
		p.WhiteRooks |= mask
	case "N":
		p.WhiteKnights |= mask
	case "B":
		p.WhiteBishops |= mask
	case "Q":
		p.WhiteQueens |= mask
	case "K":
		p.WhiteKing |= mask
	case "p":
		p.BlackPawns |= mask
	case "r":
		p.BlackRooks |= mask
	case "n":
		p.BlackKnights |= mask
	case "b":
		p.BlackBishops |= mask
	case "q":
		p.BlackQueens |= mask
	case "k":
		p.BlackKing |= mask
	default:
		panic("Piece type unknown!")
	}
}

func (p *Position) GetOccupiedSquares() Bitboard {
	return p.WhitePieces() | p.BlackPieces()
}

func (p *Position) GetEmptySquares() Bitboard {
	return ^p.GetOccupiedSquares()
}

func BitIndexToRankFile(index int) string {
	file := index % 8
	rank := index / 8
	return string(rune('a'+file)) + string(rune('1'+rank))
}

func rankFileToBitIndex(file byte, rank byte) int {
	return int((rank-'1')*8 + (file - 'a'))
}

func ToUCINotation(move Move) string {
	from := BitIndexToRankFile(move.From)
	to := BitIndexToRankFile(move.To)

	return fmt.Sprintf("%s%s", from, to)
}

func parseMove(move string) (from, to, promotionType int) {
	from = rankFileToBitIndex(move[0], move[1])
	to = rankFileToBitIndex(move[2], move[3])
	promotionType = -1 // figure out how to handle this later when implementing promotions
	return
}

func (p *Position) applyPawnMove(fromMask, toMask Bitboard) {
	if p.WhitePawns&fromMask != 0 {
		p.WhitePawns &^= fromMask
		p.WhitePawns |= toMask
	} else if p.BlackPawns&fromMask != 0 {
		p.BlackPawns &^= fromMask
		p.BlackPawns |= toMask
	}
}

// What the update function needs to do
// Parse the move string (e.g. "e2e4" → source square e2, destination square e4).
// Find which piece is moving by checking which bitboard has the source square bit set.
// Clear the source square bit in that piece’s bitboard.
// Set the destination square bit in the same piece’s bitboard.
// Handle captures: if the destination square was occupied by an opponent’s piece, clear that square from the opponent’s bitboard.
// Handle promotions, castling, and en passant (later).
func (p *Position) ApplyMove(move string) {
	from, to, _ := parseMove(move)
	fromMask := Bitboard(1) << from
	toMask := Bitboard(1) << to

	p.applyPawnMove(fromMask, toMask)
}

// debug purposes only
func PrintBitboard(bb Bitboard) {
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			square := rank*8 + file
			if (bb & (1 << square)) != 0 {
				fmt.Print("1 ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
