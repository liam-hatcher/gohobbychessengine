package chess

import (
	"math/bits"
	"math/rand"
)

const (
	A_File Bitboard = 0x0101010101010101
	H_File Bitboard = 0x8080808080808080
	Rank_1 Bitboard = 0x00000000000000FF
	Rank_8 Bitboard = 0xFF00000000000000
)

func RookRelevantMask(square int) Bitboard {
	startRank := square / 8
	startFile := square % 8

	rookMask := Bitboard(0)

	// Walk up (north) the file until you hit rank 6 (exclude rank 7, which is the edge for the rook).
	for rank := startRank + 1; rank < 7; rank++ {
		square := rank*8 + startFile
		rookMask |= Bitboard(1) << square
	}
	// Walk down (south) until rank 1 (exclude rank 0).
	for rank := startRank - 1; rank >= 1; rank-- {
		square := rank*8 + startFile
		rookMask |= Bitboard(1) << square
	}
	// Walk right (east) until file 6 (exclude file 7).
	for file := startFile + 1; file < 7; file++ {
		square := file + startRank*8
		rookMask |= Bitboard(1) << square
	}
	// Walk left (west) until file 1 (exclude file 0).
	for file := startFile - 1; file >= 1; file-- {
		square := file + startRank*8
		rookMask |= Bitboard(1) << square
	}

	return rookMask
}

func GenerateOccupancyVariations(mask Bitboard) []Bitboard {
	relevantSquares := []int{}
	variations := []Bitboard{}

	for mask != 0 {
		idx := bits.TrailingZeros64(uint64(mask))
		relevantSquares = append(relevantSquares, idx)
		mask &= mask - 1
	}

	n := len(relevantSquares)
	for i := 0; i < (1 << n); i++ {
		variation := Bitboard(0)
		for j, sq := range relevantSquares {
			if (i>>j)&1 == 1 {
				variation |= 1 << sq
			}
		}
		variations = append(variations, variation)
	}

	return variations
}

func ComputeRookAttacks(square int, occupancy Bitboard) Bitboard {
	startRank := square / 8
	startFile := square % 8

	attacks := Bitboard(0)

	// north
	for rank := startRank + 1; rank < 8; rank++ {
		attackIdx := (rank * 8) + startFile
		attacks |= Bitboard(1) << Bitboard(attackIdx)
		if occupancy&(1<<attackIdx) != 0 {
			break
		}
	}

	// south
	for rank := startRank - 1; rank >= 0; rank-- {
		attackIdx := (rank * 8) + startFile
		attacks |= Bitboard(1) << Bitboard(attackIdx)
		if occupancy&(1<<attackIdx) != 0 {
			break
		}
	}

	// west
	for file := startFile - 1; file >= 0; file-- {
		attackIdx := (startRank * 8) + file
		attacks |= Bitboard(1) << Bitboard(attackIdx)
		if occupancy&(1<<attackIdx) != 0 {
			break
		}
	}

	// east
	for file := startFile + 1; file < 8; file++ {
		attackIdx := (startRank * 8) + file
		attacks |= Bitboard(1) << Bitboard(attackIdx)
		if occupancy&(1<<attackIdx) != 0 {
			break
		}
	}

	return attacks
}

func GenerateRookAttackTable() [64][]Bitboard {
	var table [64][]Bitboard

	for sq := 0; sq < 64; sq++ {
		mask := RookRelevantMask(sq)
		variations := GenerateOccupancyVariations(mask)

		var attacksForSquare []Bitboard
		for _, occupancy := range variations {
			attacksForSquare = append(
				attacksForSquare,
				ComputeRookAttacks(sq, occupancy),
			)
		}
		table[sq] = attacksForSquare
	}

	return table
}

func sparseRandUint64() uint64 {
	return rand.Uint64() & rand.Uint64() & rand.Uint64()
}

// Take the square and its relevant mask.
// Loop over random candidate 64-bit numbers.
// For each candidate:
// Compute indices for all occupancy variations: (occupancy * candidate) >> (64 - relevantBits).
// Check if all indices are unique (no collisions).
// Return the first candidate that works.
func FindRookMagic(square int, relevantBits uint) uint64 {
	mask := RookRelevantMask(square)
	variations := GenerateOccupancyVariations(mask)

	var candidate uint64
	for attempts := 0; attempts < 10000000; attempts++ {
		candidate = sparseRandUint64()

		table := make(map[uint64]Bitboard)
		success := true

		for _, occupancy := range variations {
			index := uint64((occupancy * Bitboard(candidate)) >> (64 - relevantBits))
			attacks := ComputeRookAttacks(square, occupancy)

			if existing, ok := table[index]; ok {
				// Collision: two different occupancies map to same index
				if existing != attacks {
					success = false
					break
				}
			} else {
				table[index] = attacks
			}
		}

		if success {
			return candidate
		}
	}
	panic("failed to find magic number")
}

func popcount(bb Bitboard) int {
	count := 0
	for bb != 0 {
		bb &= bb - 1 // clear lowest set bit
		count++
	}
	return count
}

func RookRelevantBits(square int) uint {
	mask := RookRelevantMask(square)
	return uint(popcount(mask))
}

func GenerateRookMagic(reportProgress func()) ([64]uint64, [64][]Bitboard) {
	var magics [64]uint64
	var attackTables [64][]Bitboard

	for square := 0; square < 64; square++ {
		relevantBits := RookRelevantBits(square)
		magics[square] = FindRookMagic(square, relevantBits)

		size := 1 << relevantBits
		attackTables[square] = make([]Bitboard, size)

		mask := RookRelevantMask(square)
		variations := GenerateOccupancyVariations(mask)

		for _, occupancy := range variations {
			index := (uint64(occupancy) * magics[square]) >> (64 - relevantBits)
			attackTables[square][index] = ComputeRookAttacks(square, occupancy)
		}
		reportProgress()
	}

	return magics, attackTables
}
