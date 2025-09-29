package chess

// Generate all rook moves including captures
func (p *Position) GenerateRookMoves() []Move {
	var moves []Move
	var rooks Bitboard
	var friendly Bitboard

	if p.SideToMove == "white" {
		rooks = p.WhiteRooks
		friendly = p.WhitePieces()
	} else {
		rooks = p.BlackRooks
		friendly = p.BlackPieces()
	}

	occupancy := p.GetOccupiedSquares()

	for bb := rooks; bb != 0; {
		from := PopLSB(&bb)
		idx := MagicIndex(from, occupancy, RookMasks[from], RookMagics[from], RookRelevantBitsMap[from])
		attacks := RookAttackTables[from][idx]
		attacks &^= friendly

		for a := attacks; a != 0; {
			to := PopLSB(&a)
			moves = append(moves, Move{From: from, To: to})
		}
	}

	return moves
}
