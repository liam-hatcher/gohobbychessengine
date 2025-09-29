package chess

import (
	"fmt"
)

type Bitboard uint64

type Move struct {
	From  int
	To    int
	Promo byte
}

type CastlingRights struct {
	Short bool
	Long  bool
}

type Position struct {
	WhitePawns          Bitboard
	WhiteKnights        Bitboard
	WhiteBishops        Bitboard
	WhiteRooks          Bitboard
	WhiteQueens         Bitboard
	WhiteKing           Bitboard
	WhiteCastlingRights CastlingRights

	BlackPawns          Bitboard
	BlackKnights        Bitboard
	BlackBishops        Bitboard
	BlackRooks          Bitboard
	BlackQueens         Bitboard
	BlackKing           Bitboard
	BlackCastlingRights CastlingRights

	EnPassantTarget Bitboard

	// for O(1) lookups of pieces on a given square
	PieceMap [64]byte

	SideToMove string
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

		PieceMap: [64]byte{
			'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
			'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
			'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
		},

		SideToMove: "white",
	}
}

func (p *Position) WhitePieces() Bitboard {
	return p.WhiteBishops | p.WhiteKing | p.WhiteKnights | p.WhitePawns | p.WhiteRooks | p.WhiteQueens
}

func (p *Position) BlackPieces() Bitboard {
	return p.BlackBishops | p.BlackKing | p.BlackKnights | p.BlackPawns | p.BlackRooks | p.BlackQueens
}

func (p *Position) SetPiece(piece byte, square string) {
	index := RankFileToBitIndex(square[0], square[1])
	mask := Bitboard(1) << index

	if mask&p.GetOccupiedSquares() != 0 {
		panic(fmt.Sprintf("%s is already occupied", square))
	}

	switch piece {
	case 'P':
		p.WhitePawns |= mask
	case 'R':
		p.WhiteRooks |= mask
	case 'N':
		p.WhiteKnights |= mask
	case 'B':
		p.WhiteBishops |= mask
	case 'Q':
		p.WhiteQueens |= mask
	case 'K':
		p.WhiteKing |= mask
	case 'p':
		p.BlackPawns |= mask
	case 'r':
		p.BlackRooks |= mask
	case 'n':
		p.BlackKnights |= mask
	case 'b':
		p.BlackBishops |= mask
	case 'q':
		p.BlackQueens |= mask
	case 'k':
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
		p.PieceMap[index] = 0
		return
	}
	if p.WhiteKnights&mask != 0 {
		p.WhiteKnights &^= mask
		p.PieceMap[index] = 0
		return
	}
	if p.WhiteBishops&mask != 0 {
		p.WhiteBishops &^= mask
		p.PieceMap[index] = 0
		return
	}
	if p.WhiteRooks&mask != 0 {
		p.WhiteRooks &^= mask
		p.PieceMap[index] = 0
		return
	}
	if p.WhiteQueens&mask != 0 {
		p.WhiteQueens &^= mask
		p.PieceMap[index] = 0
		return
	}
	if p.WhiteKing&mask != 0 {
		p.WhiteKing &^= mask
		p.PieceMap[index] = 0
		return
	}

	if p.BlackPawns&mask != 0 {
		p.BlackPawns &^= mask
		p.PieceMap[index] = 0
		return
	}
	if p.BlackKnights&mask != 0 {
		p.BlackKnights &^= mask
		p.PieceMap[index] = 0
		return
	}
	if p.BlackBishops&mask != 0 {
		p.BlackBishops &^= mask
		p.PieceMap[index] = 0
		return
	}
	if p.BlackRooks&mask != 0 {
		p.BlackRooks &^= mask
		p.PieceMap[index] = 0
		return
	}
	if p.BlackQueens&mask != 0 {
		p.BlackQueens &^= mask
		p.PieceMap[index] = 0
		return
	}
	if p.BlackKing&mask != 0 {
		p.BlackKing &^= mask
		p.PieceMap[index] = 0
		return
	}
}

func (p *Position) GetPieceOnSquare(square string) (piece byte, empty bool) {
	index := RankFileToBitIndex(square[0], square[1])

	piece = p.PieceMap[index]

	if piece != 0 {
		return piece, false
	}

	return 0, true
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

func AtBounds(index int) bool {
	leftBound := index%8 == 0
	rightBound := index%8 == 7
	upperBound := index/8 == 7
	lowerBound := index/8 == 0

	return leftBound || rightBound || upperBound || lowerBound
}

func (p *Position) applyPromotion(promotion byte, to int, toMask Bitboard) {
	switch promotion {
	case 'Q':
		p.WhiteQueens |= toMask
	case 'B':
		p.WhiteBishops |= toMask
	case 'N':
		p.WhiteKnights |= toMask
	case 'R':
		p.WhiteRooks |= toMask
	case 'q':
		p.BlackQueens |= toMask
	case 'b':
		p.BlackBishops |= toMask
	case 'n':
		p.BlackKnights |= toMask
	case 'r':
		p.BlackRooks |= toMask
	}

	p.PieceMap[to] = promotion
}

func (p *Position) applyPawnMove(toMask, fromMask Bitboard, to, from int, promotion byte) {
	diff := Abs(to - from)
	isPush := diff == 8 || diff == 16

	var movingPawns, opponentPawns *Bitboard
	var pieceChar byte
	var enPassantOffset int

	if p.WhitePawns&fromMask != 0 {
		movingPawns = &p.WhitePawns
		opponentPawns = &p.BlackPawns
		pieceChar = 'P'
		enPassantOffset = -8
	} else if p.BlackPawns&fromMask != 0 {
		movingPawns = &p.BlackPawns
		opponentPawns = &p.WhitePawns
		pieceChar = 'p'
		enPassantOffset = 8
	}

	*movingPawns &^= fromMask

	if !isPush {
		if toMask == p.EnPassantTarget {
			capIdx := to + enPassantOffset
			*opponentPawns &^= Bitboard(1) << capIdx
			p.PieceMap[capIdx] = 0
		} else {
			if pieceChar == 'P' {
				p.BlackPawns &^= toMask
				p.BlackKnights &^= toMask
				p.BlackBishops &^= toMask
				p.BlackRooks &^= toMask
				p.BlackQueens &^= toMask
			} else {
				p.WhitePawns &^= toMask
				p.WhiteKnights &^= toMask
				p.WhiteBishops &^= toMask
				p.WhiteRooks &^= toMask
				p.WhiteQueens &^= toMask
			}
		}
	}

	isBackRank := to/8 == 7
	isFirstRank := to/8 == 0

	if isBackRank {
		promotion = ToUpper(promotion)
	}

	if isBackRank || isFirstRank {
		p.applyPromotion(promotion, to, toMask)
	} else {
		*movingPawns |= toMask
		p.PieceMap[to] = pieceChar
	}

	p.PieceMap[from] = 0
}

func (p *Position) changeTurn() {
	if p.SideToMove == "white" {
		p.SideToMove = "black"
	} else {
		p.SideToMove = "white"
	}
}

func (p *Position) updateEnpassantState(pieceMoving byte, from, to int) {
	diff := Abs(to - from)
	if pieceMoving == 'P' && diff == 16 { // white double push
		p.EnPassantTarget = 1 << (from + 8)
	} else if pieceMoving == 'p' && diff == 16 { // black double push
		p.EnPassantTarget = 1 << (from - 8)
	} else {
		p.EnPassantTarget = 0 // clear for all other moves
	}
}

func (p *Position) ApplyMove(move string) {
	from, to, promotion := ParseMove(move)
	fromMask := Bitboard(1) << from
	toMask := Bitboard(1) << to
	pieceMoving := p.PieceMap[from]

	switch pieceMoving {
	case 'P', 'p':
		p.applyPawnMove(toMask, fromMask, to, from, promotion)
	default:
		panic("unexpected piece type")
	}

	p.updateEnpassantState(pieceMoving, from, to)

	p.changeTurn()
}
