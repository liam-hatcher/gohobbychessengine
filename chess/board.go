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
	PieceMap [64]byte
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

func parseMove(move string) (from, to int, promotionType byte) {
	from = RankFileToBitIndex(move[0], move[1])
	to = RankFileToBitIndex(move[2], move[3])

	if len(move) == 5 {
		promotionType = move[4]
	}

	return
}

// for converting uci promotion chars to uppercase
func ToUpper(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 'a' + 'A'
	}
	return b
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
	diff := abs(to - from)
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

func (p *Position) ApplyMove(move string) {
	from, to, promotion := parseMove(move)
	fromMask := Bitboard(1) << from
	toMask := Bitboard(1) << to
	pieceMoving := p.PieceMap[from]

	switch pieceMoving {
	case 'P', 'p':
		p.applyPawnMove(toMask, fromMask, to, from, promotion)
	default:
		panic("unexpected piece type")
	}

	diff := abs(to - from)
	if pieceMoving == 'P' && diff == 16 { // white double push
		p.EnPassantTarget = 1 << (from + 8)
	} else if pieceMoving == 'p' && diff == 16 { // black double push
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

// for debugging
func PrintBitboard(bb uint64) {
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			sq := rank*8 + file
			mask := uint64(1) << sq
			if bb&mask != 0 {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
