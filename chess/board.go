package chess

import (
	"fmt"
)

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

	EnPassantTarget Bitboard

	// for O(1) lookups of pieces on a given square
	PieceMap [64]string
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

		EnPassantTarget: 0,

		PieceMap: [64]string{
			"R", "N", "B", "Q", "K", "B", "N", "R",
			"P", "P", "P", "P", "P", "P", "P", "P",
			"", "", "", "", "", "", "", "",
			"", "", "", "", "", "", "", "",
			"", "", "", "", "", "", "", "",
			"", "", "", "", "", "", "", "",
			"p", "p", "p", "p", "p", "p", "p", "p",
			"r", "n", "b", "q", "k", "b", "n", "r",
		},
	}
}

func (p *Position) WhitePieces() Bitboard {
	return p.WhiteBishops | p.WhiteKing | p.WhiteKnights | p.WhitePawns | p.WhiteRooks | p.WhiteQueens
}

func (p *Position) BlackPieces() Bitboard {
	return p.BlackBishops | p.BlackKing | p.BlackKnights | p.BlackPawns | p.BlackRooks | p.BlackQueens
}

func (p *Position) SetPiece(piece, square string) {
	index := RankFileToBitIndex(square[0], square[1])
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

	p.PieceMap[index] = piece
}

func (p *Position) RemovePiece(square string) {
	index := RankFileToBitIndex(square[0], square[1])
	mask := Bitboard(1) << index

	if p.WhitePawns&mask != 0 {
		p.WhitePawns &^= mask
		p.PieceMap[index] = ""
		return
	}
	if p.WhiteKnights&mask != 0 {
		p.WhiteKnights &^= mask
		p.PieceMap[index] = ""
		return
	}
	if p.WhiteBishops&mask != 0 {
		p.WhiteBishops &^= mask
		p.PieceMap[index] = ""
		return
	}
	if p.WhiteRooks&mask != 0 {
		p.WhiteRooks &^= mask
		p.PieceMap[index] = ""
		return
	}
	if p.WhiteQueens&mask != 0 {
		p.WhiteQueens &^= mask
		p.PieceMap[index] = ""
		return
	}
	if p.WhiteKing&mask != 0 {
		p.WhiteKing &^= mask
		p.PieceMap[index] = ""
		return
	}

	if p.BlackPawns&mask != 0 {
		p.BlackPawns &^= mask
		p.PieceMap[index] = ""
		return
	}
	if p.BlackKnights&mask != 0 {
		p.BlackKnights &^= mask
		p.PieceMap[index] = ""
		return
	}
	if p.BlackBishops&mask != 0 {
		p.BlackBishops &^= mask
		p.PieceMap[index] = ""
		return
	}
	if p.BlackRooks&mask != 0 {
		p.BlackRooks &^= mask
		p.PieceMap[index] = ""
		return
	}
	if p.BlackQueens&mask != 0 {
		p.BlackQueens &^= mask
		p.PieceMap[index] = ""
		return
	}
	if p.BlackKing&mask != 0 {
		p.BlackKing &^= mask
		p.PieceMap[index] = ""
		return
	}
}

func (p *Position) GetPieceOnSquare(square string) (piece string, empty bool) {
	index := RankFileToBitIndex(square[0], square[1])

	piece = p.PieceMap[index]

	if piece != "" {
		return piece, false
	}

	return "", true
}

func (p *Position) SetEnPassantTarget(square string) {
	p.EnPassantTarget = 1 << RankFileToBitIndex(square[0], square[1])
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

func RankFileToBitIndex(file byte, rank byte) int {
	return int((rank-'1')*8 + (file - 'a'))
}

func ToUCINotation(move Move) string {
	from := BitIndexToRankFile(move.From)
	to := BitIndexToRankFile(move.To)

	return fmt.Sprintf("%s%s", from, to)
}

func parseMove(move string) (from, to, promotionType int) {
	from = RankFileToBitIndex(move[0], move[1])
	to = RankFileToBitIndex(move[2], move[3])
	promotionType = -1 // figure out how to handle this later when implementing promotions
	return
}

func (p *Position) applyPawnMove(toMask, fromMask Bitboard, to, from int) {
	diff := abs(to - from)
	isPush := diff == 8 || diff == 16
	// TODO: set en passant target square
	if p.WhitePawns&fromMask != 0 {
		p.WhitePawns &^= fromMask
		p.WhitePawns |= toMask

		if !isPush {
			if toMask == p.EnPassantTarget {
				capIdx := to - 8
				p.BlackPawns &^= Bitboard(1) << capIdx
				p.PieceMap[capIdx] = ""
			} else {
				p.BlackPawns &^= toMask
			}
		}

		p.PieceMap[to] = "P"
		p.PieceMap[from] = ""
	} else if p.BlackPawns&fromMask != 0 {
		p.BlackPawns &^= fromMask
		p.BlackPawns |= toMask
		if !isPush {
			if toMask == p.EnPassantTarget {
				capIdx := to + 8
				p.WhitePawns &^= Bitboard(1) << capIdx
				p.PieceMap[capIdx] = ""
			} else {
				p.WhitePawns &^= toMask
			}
			p.WhitePawns &^= toMask
		}

		p.PieceMap[to] = "p"
		p.PieceMap[from] = ""
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
	pieceMoving := p.PieceMap[from]

	switch pieceMoving {
	case "P", "p":
		p.applyPawnMove(toMask, fromMask, to, from)
	default:
		panic("unexpected piece type")
	}

	diff := abs(to - from)
	if pieceMoving == "P" && diff == 16 { // white double push
		p.EnPassantTarget = 1 << (from + 8)
	} else if pieceMoving == "p" && diff == 16 { // black double push
		p.EnPassantTarget = 1 << (from - 8)
	} else {
		p.EnPassantTarget = 0 // clear for all other moves
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
