package chess

import (
	"fmt"
	"math/bits"
	"strings"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// for debugging
func PrintBitboard(bb Bitboard) {
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			sq := rank*8 + file
			mask := Bitboard(1) << sq
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

// ToString converts a Bitboard to an 8x8 string of 1s and 0s for debugging
func ToString(bb Bitboard) string {
	var sb strings.Builder

	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			sq := rank*8 + file
			mask := Bitboard(1) << sq
			if bb&mask != 0 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}

			if file < 7 {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

// for converting uci promotion chars to uppercase
func ToUpper(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 'a' + 'A'
	}
	return b
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

func ParseMove(move string) (from, to int, promotionType byte) {
	from = RankFileToBitIndex(move[0], move[1])
	to = RankFileToBitIndex(move[2], move[3])

	if len(move) == 5 {
		promotionType = move[4]
	}

	return
}

func PopLSB(bb *Bitboard) int {
	lsb := *bb & -*bb                           // isolate least significant bit
	square := bits.TrailingZeros64(uint64(lsb)) // returns 0-63
	*bb &= *bb - 1                              // clear the LSB
	return square
}
