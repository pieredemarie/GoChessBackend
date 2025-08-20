package chesslogic

import (
	"fmt"
)

//TODO
//Others pieces Legal Moves

func abs(x int) int {
	if x < 0 {return -x}
	return x
}

func squareToCoords(sq string) (int, int, error) {
	if len(sq) != 2 {
		return 0, 0, fmt.Errorf("invalid square: %s", sq)
	}
	file := int(sq[0] - 'a') 
    rank := int(sq[1] - '1') 
    if file < 0 || file > 7 || rank < 0 || rank > 7 {
        return 0, 0, fmt.Errorf("invalid square: %s", sq)
    }
    return file, rank, nil
}

func (b *Board) isLegalKnightMove(piece Piece, m Move) bool {
	fromFile, fromRank, err := squareToCoords(m.From)
    if err != nil {
        return false
    }
    toFile, toRank, err := squareToCoords(m.To)
    if err != nil {
        return false
    }

	pieceToEnd := b.GetPiece(m.To)
	if pieceToEnd != nil && pieceToEnd.Color == piece.Color { // Horses can't jump on their friends, otherwise you have a weird horse
		return false
	}

	rowsDiff := abs(fromRank - toRank)
	colsDiff := abs(fromFile - toFile)

	return (rowsDiff == 2 && colsDiff == 1) || (rowsDiff == 1 && colsDiff == 2)
}
