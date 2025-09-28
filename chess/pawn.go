package chess

import "math/bits"

func generatePromotions(from, to int) []Move {
	return []Move{
		{From: from, To: to, Promo: 'r'}, // using lower case here to be consistent with UCI notation
		{From: from, To: to, Promo: 'q'},
		{From: from, To: to, Promo: 'k'},
		{From: from, To: to, Promo: 'b'},
	}
}

func (p *Position) GenerateWhitePawnMoves() []Move {
	var moves []Move

	emptySquares := p.GetEmptySquares()

	singlePush := (p.WhitePawns << 8) & emptySquares
	for singlePush != 0 {
		to := bits.TrailingZeros64(uint64(singlePush))
		from := to - 8
		isBackRank := to/8 == 7
		if isBackRank {
			moves = append(moves, generatePromotions(from, to)...)
		} else {
			moves = append(moves, Move{From: from, To: to})
		}
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

// func (p *Position) GenerateWhitePawnPushes() Bitboard {
// 	emptySquares := ^p.GetEmptySquares()

// 	singlePush := (p.WhitePawns << 8) & emptySquares

// 	rank2 := Bitboard(0x000000000000FF00)
// 	doublePush := ((p.WhitePawns & rank2) << 16) & emptySquares & (emptySquares << 8)

// 	return singlePush | doublePush
// }

func (p *Position) GenerateBlackPawnMoves() []Move {
	var moves []Move

	emptySquares := p.GetEmptySquares()

	singlePush := (p.BlackPawns >> 8) & emptySquares
	for singlePush != 0 {
		to := bits.TrailingZeros64(uint64(singlePush))
		from := to + 8
		isFirstRank := to/8 == 0
		if isFirstRank {
			moves = append(moves, generatePromotions(from, to)...)
		} else {
			moves = append(moves, Move{From: from, To: to})
		}
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

// func (p *Position) GenerateBlackPawnPushes() Bitboard {
// 	emptySquares := ^p.GetEmptySquares()

// 	singlePush := (p.BlackPawns >> 8) & emptySquares

// 	rank7 := Bitboard(0x00FF000000000000)
// 	doublePush := ((p.BlackPawns & rank7) >> 16) & emptySquares & (emptySquares >> 8)

// 	return singlePush | doublePush
// }

func (p *Position) GenerateWhitePawnCaptures() []Move {
	var moves []Move

	generateCaptures := func(captures Bitboard, offset int) {
		for captures != 0 {
			to := bits.TrailingZeros64(uint64(captures))
			from := to - offset
			isBackRank := to/8 == 7
			if isBackRank {
				moves = append(moves, generatePromotions(from, to)...)
			} else {
				moves = append(moves, Move{From: from, To: to})
			}
			captures &= captures - 1
		}
	}

	leftCaptures := ((p.WhitePawns &^ A_File) << 7) & (p.BlackPieces() | p.EnPassantTarget)
	rightCaptures := ((p.WhitePawns &^ H_File) << 9) & (p.BlackPieces() | p.EnPassantTarget)

	generateCaptures(leftCaptures, 7)
	generateCaptures(rightCaptures, 9)

	return moves
}

func (p *Position) GenerateBlackPawnCaptures() []Move {
	var moves []Move

	aFile := Bitboard(0x0101010101010101)
	hFile := Bitboard(0x8080808080808080)

	generateCaptures := func(captures Bitboard, offset int) {
		if captures != 0 {
			to := bits.TrailingZeros64(uint64(captures))
			isFirstRank := to/8 == 0
			from := to + offset
			if isFirstRank {
				moves = append(moves, generatePromotions(from, to)...)
			} else {
				moves = append(moves, Move{From: from, To: to})
			}
			captures &= captures - 1
		}
	}

	leftCaptures := ((p.BlackPawns &^ aFile) >> 9) & (p.WhitePieces() | p.EnPassantTarget)
	rightCaptures := ((p.BlackPawns &^ hFile) >> 7) & (p.WhitePieces() | p.EnPassantTarget)

	generateCaptures(leftCaptures, 9)
	generateCaptures(rightCaptures, 7)

	return moves
}
