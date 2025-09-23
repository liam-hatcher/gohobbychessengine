package chess

import "math/bits"

func (p *Position) GenerateWhitePawnMoves() []Move {
	var moves []Move

	emptySquares := p.GetEmptySquares()

	singlePush := (p.WhitePawns << 8) & emptySquares
	for singlePush != 0 {
		to := bits.TrailingZeros64(uint64(singlePush))
		from := to - 8
		moves = append(moves, Move{From: from, To: to})
		singlePush &= singlePush - 1
	}

	rank2 := Bitboard(0x000000000000FF00)
	doublePush := ((p.WhitePawns & rank2) << 16) & emptySquares & (emptySquares << 8)
	for doublePush != 0 {
		to := bits.TrailingZeros64(uint64(doublePush))
		from := to - 16
		moves = append(moves, Move{From: from, To: to})
		doublePush &= doublePush - 1
	}

	return moves
}

func (p *Position) GenerateWhitePawnPushes() Bitboard {
	emptySquares := ^p.GetEmptySquares()

	singlePush := (p.WhitePawns << 8) & emptySquares

	rank2 := Bitboard(0x000000000000FF00)
	doublePush := ((p.WhitePawns & rank2) << 16) & emptySquares & (emptySquares << 8)

	return singlePush | doublePush
}

func (p *Position) GenerateBlackPawnMoves() []Move {
	var moves []Move

	emptySquares := p.GetEmptySquares()

	singlePush := (p.BlackPawns >> 8) & emptySquares
	for singlePush != 0 {
		to := bits.TrailingZeros64(uint64(singlePush))
		from := to + 8
		moves = append(moves, Move{From: from, To: to})
		singlePush &= singlePush - 1
	}

	rank7 := Bitboard(0x00FF000000000000)
	doublePush := ((p.BlackPawns & rank7) >> 16) & emptySquares & (emptySquares >> 8)
	for doublePush != 0 {
		to := bits.TrailingZeros64(uint64(doublePush))
		from := to + 16
		moves = append(moves, Move{From: from, To: to})
		doublePush &= doublePush - 1
	}

	return moves
}

func (p *Position) GenerateBlackPawnPushes() Bitboard {
	emptySquares := ^p.GetEmptySquares()

	singlePush := (p.BlackPawns >> 8) & emptySquares

	rank7 := Bitboard(0x00FF000000000000)
	doublePush := ((p.BlackPawns & rank7) >> 16) & emptySquares & (emptySquares >> 8)

	return singlePush | doublePush
}
